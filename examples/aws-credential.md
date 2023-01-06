### Create credential file
The file must be named as `credentials`, you can create credential file from command line as,
```
cat <<-EOF > /local/path/to/credentials

[default]
aws_access_key_id = exampleawsaccesskey123
aws_secret_access_key = exampleawssecretaccesskey456
EOF
```

### Create secret from file
Use the `credentials` file to create secret:
```bash
kubectl create secret generic aws-credential --namespace demo --from-file /local/path/to/credentials
```

The secret must be in the same namespace as the External-DNS