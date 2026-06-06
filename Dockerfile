# Copyright AppsCode Inc. and Contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Unified Dockerfile proposal — replaces Dockerfile.in / Dockerfile.dbg
# / Dockerfile.ubi with a single multi-stage build selectable via
# `docker build --target {prod|dbg|ubi}`.
#
# Build args (all required):
#   BIN  — binary name (must match bin/${BIN}-${TARGETOS}-${TARGETARCH})
#   BASE_PROD — base image for the prod target (e.g. gcr.io/distroless/...)
#   BASE_DBG  — base image for the dbg target  (debian-based, runs dlv)
#   BASE_UBI  — base image for the ubi target  (Red Hat UBI minimal)
#   TAG       — release tag, only used for UBI label metadata
#
# Each target produces an image that contains the prebuilt binary
# bin/${BIN}-${TARGETOS}-${TARGETARCH} (built outside this Dockerfile,
# same as the current setup).

ARG BIN
ARG BASE_PROD
ARG BASE_DBG
ARG BASE_UBI
ARG TAG=latest

# ---------- prod (distroless / alpine) ----------------------------------
FROM ${BASE_PROD} AS prod
ARG BIN
ARG TARGETOS
ARG TARGETARCH

LABEL org.opencontainers.image.source="https://github.com/kubeops/external-dns-operator"

RUN set -x \
  && apk add --update --upgrade --no-cache pcre2 ca-certificates tzdata \
  && echo 'Etc/UTC' > /etc/timezone

ADD bin/${BIN}-${TARGETOS}-${TARGETARCH} /${BIN}
USER 65534
ENTRYPOINT ["/${BIN}"]

# ---------- dbg (debian + dlv) ------------------------------------------
FROM ghcr.io/appscode/dlv:1.25 AS dlv

FROM ${BASE_DBG} AS dbg
ARG BIN
ARG TARGETOS
ARG TARGETARCH

LABEL org.opencontainers.image.source="https://github.com/kubeops/external-dns-operator"

RUN set -x \
  && apt-get update \
  && apt-get upgrade -y \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/* \
  && echo 'Etc/UTC' > /etc/timezone

ADD bin/${BIN}-${TARGETOS}-${TARGETARCH} /${BIN}
COPY --from=dlv /usr/local/bin/dlv /bin/dlv

EXPOSE 40000
ENTRYPOINT ["/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/${BIN}", "--"]

# ---------- ubi (Red Hat UBI minimal) -----------------------------------
FROM ${BASE_UBI} AS ubi
ARG BIN
ARG TARGETOS
ARG TARGETARCH
ARG TAG

LABEL org.opencontainers.image.source="https://github.com/kubeops/external-dns-operator" \
      name="External DNS Operator" \
      maintainer=AppsCode \
      vendor=AppsCode \
      version=${TAG} \
      release=${TAG} \
      summary="External DNS Operator configures external DNS servers dynamically from Kubernetes resources" \
      description="External DNS Operator configures external DNS servers dynamically from Kubernetes resources."

RUN mkdir -p /licenses
COPY LICENSE /licenses/

RUN set -x \
  && microdnf update -y \
  && microdnf install -y tzdata \
  && microdnf clean all

ENV TZ=Etc/UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ADD bin/${BIN}-${TARGETOS}-${TARGETARCH} /${BIN}
USER 65534
ENTRYPOINT ["/${BIN}"]
