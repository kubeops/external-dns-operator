### Create credential file
First, you have to add a credential file using AWS secrets,
```bash
cat <<-EOF > /local/path/to/credentials

[default]
aws_access_key_id = exampleawsaccesskey123
aws_secret_access_key = exampleawssecretaccesskey456
EOF
```

### Create secret from file
Use the above file to create a secret:
```bash
kubectl create secret generic aws-credential --namespace demo --from-file /local/path/to/credentials
```

The secret must be in the same **namespace** as the External-DNS

You can use this secret to create records in the AWS provider. The secret name should be used against `spec.aws.secretRef.name` field and secret key against `spec.aws.secretRef.credentialKey`
