### Create credential file
The file must be named as `credentials`, you can from command line as,
```
cat <<-EOF > /local/path/to/credentials

[default]
aws_access_key_id = exampleawsaccesskey123
aws_secret_access_key = exampleawssecretaccesskey456
EOF
```

### Create secret from file
Use the `credentials` file to create secret:<br>
`kubectl create secret generic aws-credential --namespace demo --from-file /local/path/to/credentials`
<br>

The secret must be in the same namespace as the External-DNS