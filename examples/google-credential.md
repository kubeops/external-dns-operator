### Create Secret from `credentials.json`
Create a Kubernetes secret with the credentials in the same namespace of External-DNS operator. <br>
`kubectl create secret generic "google-credential" --namespace demo --from-file /local/path/to/credentials.json` <br>
 The file containing the GKE credentials must be named as `credentials.json`