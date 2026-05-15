# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `kubeops.dev/external-dns-operator` — a Kubernetes operator that wraps [external-dns](https://github.com/kubernetes-sigs/external-dns) behind a CRD (`ExternalDNS`) so cluster admins can declare a DNS sync as a Kubernetes object instead of a long flag list on a Deployment. The controller resolves the spec into the corresponding external-dns args, gathers the right cloud credentials from referenced `Secret`s, and computes a plan against the target DNS provider.

The produced binary is `external-dns-operator`. The README is a Kubebuilder scaffold stub; treat this file as the source of truth.

## Architecture

- `cmd/external-dns-operator/` — entry point.
- `pkg/cmds/` — Cobra commands (root, run).
- `apis/external/v1alpha1/` — Kubebuilder API types. Single CRD: `ExternalDNS` (`external-dns.appscode.com/v1alpha1`).
  - `register.go`, `install/`, `fuzzer/`, `*_types.go`, generated `zz_generated.deepcopy.go`.
- `client/` — generated typed clientset.
- `crds/` — generated CRD YAMLs.
- `pkg/controllers/external-dns/` — the `ExternalDNS` reconciler.
- `pkg/credentials/` — **per-cloud credential resolution** from referenced Secrets:
  - `aws.go`, `azure.go`, `cloudflare.go`, `google.go` — one file per provider.
  - `secret.go` — generic Secret-to-config helpers.
- `pkg/plan/plan.go` — wraps external-dns's planner so the operator can preview changes.
- `pkg/informers/` — shared informer setup.
- `pkg/external/` — embeds/integrates with the upstream external-dns library.
- `examples/` — sample `ExternalDNS` manifests per provider.
- `PROJECT` — Kubebuilder metadata. Domain `appscode.com`, multigroup.
- `Dockerfile.in` (PROD, distroless), `Dockerfile.dbg` (debian), `Dockerfile.ubi` (Red Hat certified) — three image variants.
- `hack/`, `Makefile` — AppsCode build harness.
- `vendor/` — checked-in deps.

CRD API group is `external-dns.appscode.com/v1alpha1`.

## Common commands

All Make targets run inside `ghcr.io/appscode/golang-dev` — Docker must be running.

- `make ci` — CI pipeline.
- `make build` / `make all-build` — build host or all-platform binaries.
- `make gen` — regenerate clientset + manifests. Run after any change to `apis/external/v1alpha1/*_types.go`.
- `make manifests` — regenerate CRDs only.
- `make clientset` — regenerate `client/` only.
- `make fmt`, `make lint`, `make unit-tests` / `make test` — standard.
- `make verify` — `verify-gen verify-modules`; `go mod tidy && go mod vendor` must leave the tree clean.
- `make container` — build PROD, DBG, and UBI images.
- `make push` — push all three; `make docker-manifest` writes multi-arch manifests; `make release` is the full publish flow.
- `make push-to-kind` / `make deploy-to-kind` — load into Kind and Helm-install.
- `make install` / `make uninstall` / `make purge` — Helm install lifecycle.
- `make add-license` / `make check-license` — manage license headers.

Run a single Go test (requires a local Go toolchain):

```
go test ./pkg/controllers/external-dns/... -run TestName -v
```

## Conventions

- Module path is `kubeops.dev/external-dns-operator` (vanity URL). Imports must use that.
- License: `LICENSE` (Apache-2.0); new files need the standard AppsCode header (`make add-license`).
- Sign off commits (`git commit -s`); contributions follow the DCO (`DCO`).
- Vendor directory is checked in — `go mod tidy && go mod vendor` must leave the tree clean (enforced by `verify-modules`).
- Adding a new DNS provider: implement credential resolution under `pkg/credentials/<provider>.go` and wire it from the reconciler. Don't sprinkle provider-specific code across `pkg/controllers/external-dns/`.
- Do not hand-edit `zz_generated.*.go`, anything under `client/`, or `crds/` — change `apis/external/v1alpha1/*_types.go` and re-run `make gen`.
- The operator owns the relationship to upstream `external-dns`. Pin its dep deliberately; field reshuffles upstream propagate into `pkg/plan/` and `pkg/external/`.
- Three Dockerfiles, one binary — keep `Dockerfile.in`, `Dockerfile.dbg`, and `Dockerfile.ubi` in sync.
- This is a **Kubebuilder multigroup project** — use `kubebuilder` to scaffold new APIs.
