### Create Secret from
Create a Kubernetes secret with the credentials in the same namespace of External-DNS operator. 
```shell
kubectl create secret generic google-credential --namespace demo --from-file /local/path/to/credentials.json
```
The name and key of this secret will be used in `spec.google.secretRef.name` and `spec.google.secretRef.credentialKey`
