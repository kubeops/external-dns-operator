### Create static credentials
* Create a Google Service Account(GSA) that has access to the CloudDNS Zone
```bash
GKE_PROJECT_ID="your-organization-project-id"
DNS_SA_NAME="external-dns-sa"
DNS_SA_EMAIL="$DNS_SA_NAME@${GKE_PROJECT_ID}.iam.gserviceaccount.com"

# create GSA used to access the Cloud DNS zone
gcloud iam service-accounts create $DNS_SA_NAME --display-name $DNS_SA_NAME

# assign google service account to dns.admin role in cloud-dns project
gcloud projects add-iam-policy-binding $DNS_PROJECT_ID --member serviceAccount:$DNS_SA_EMAIL --role "roles/dns.admin"
```
* Generate static credential from the ExternalDNS GSA
```bash
# download static credentials
gcloud iam service-accounts keys create /local/path/to/credentials.json --iam-account $DNS_SA_EMAIL
```
### Create Secret from
Create a Kubernetes secret with the credentials in the `same namespace of External-DNS` operator. 
```shell
kubectl create secret generic google-credential --namespace demo --from-file /local/path/to/credentials.json
```
The name and key of this secret will be used in `spec.google.secretRef.name` and `spec.google.secretRef.credentialKey`
