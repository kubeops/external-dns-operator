### Create configuration file
* Create `Resource Group` and `DNS Zone`
```bash
az group create --name "MyDnsResourceGroup" --location "eastus"
az network dns zone create --resource-group "MyDnsResourceGroup" --name "example.com"
```
* Create a `Service Principal`
```bash
EXTERNALDNS_NEW_SP_NAME="ExternalDnsServicePrincipal" # name of the service principal
AZURE_DNS_ZONE_RESOURCE_GROUP="MyDnsResourceGroup" # name of resource group where dns zone is hosted
AZURE_DNS_ZONE="example.com" # DNS zone name like example.com or sub.example.com

# Create the service principal
DNS_SP=$(az ad sp create-for-rbac --name $EXTERNALDNS_NEW_SP_NAME)
EXTERNALDNS_SP_APP_ID=$(echo $DNS_SP | jq -r '.appId')
EXTERNALDNS_SP_PASSWORD=$(echo $DNS_SP | jq -r '.password')
```
* Grant access to Azure DNS zone for the service principal.
```bash
# fetch DNS id used to grant access to the service principal
DNS_ID=$(az network dns zone show --name $AZURE_DNS_ZONE \
 --resource-group $AZURE_DNS_ZONE_RESOURCE_GROUP --query "id" --output tsv)

# 1. as a reader to the resource group
# az role assignment create --role "Reader" --assignee $EXTERNALDNS_SP_APP_ID --scope $DNS_ID

# 2. as a contributor to DNS Zone itself
az role assignment create --role "Contributor" --assignee $EXTERNALDNS_SP_APP_ID --scope $DNS_ID
```
* Write the credentials to a local path
```bash
cat <<-EOF > /local/path/to/azure.json
{
  "tenantId": "$(az account show --query tenantId -o tsv)",
  "subscriptionId": "$(az account show --query id -o tsv)",
  "resourceGroup": "$AZURE_DNS_ZONE_RESOURCE_GROUP",
  "aadClientId": "$EXTERNALDNS_SP_APP_ID",
  "aadClientSecret": "$EXTERNALDNS_SP_PASSWORD"
}
EOF
```
* Once you have completed all the process you will have a json file in `/local/path/to/azure.json` path
```json
{
  "tenantId": "your-azure-tenant-id",
  "subscriptionId": "your-azure-subscription-id",
  "resourceGroup": "your-azure-resource-group-name",
  "aadClientId": "your-azure-client-id",
  "aadClientSecret": "your-azure-client-password"
}
```

### Create secret from file
Use the JSON file to create a secret.

```shell
kubectl create secret generic azure-credential --namespace demo --from-file /local/path/to/azure.json
```

The secret must be in the same namespace as the External-DNS

This secret name and secret key will be used in `spec.azure.secretRef.name` and `spec.azure.secretRef.credentialKey`
