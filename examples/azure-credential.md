### Create credential file
The file must be named as `azure.json`, you can create the file by:
```
cat <<-EOF > /local/path/to/azure.json
{
  "tenantId": "your-azure-tenant-id",
  "subscriptionId": "your-azure-subscription-id",
  "resourceGroup": "your-azure-resource-group-name",
  "useManagedIdentityExtension": true
}
EOF
```

### Create secret from file
Use the `azure.json` file to create Kubernetes secret:

```shell
kubectl create secret generic azure-credential --namespace demo --from-file /local/path/to/azure.json
```

The secret must be in the same namespace as the External-DNS