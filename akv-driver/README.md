# ~infra prereq: 

1) Keyvault resource, using the RBAC permissions model
2) some secret already present in the vault

# prep: making SP

### auth into a user via az cli, provision *service principal* scoped to KV w/ access 

``` bash
docker run -it --rm mcr.microsoft.com/azure-cli:2.9.1 bash

SERVICE_PRINCIPAL_NAME=secretReaderSP
AKV_ID=keyvaultresourcename

#role ID from https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#key-vault-secrets-user-preview
SP_PASSWD=$(az ad sp create-for-rbac --name http://$SERVICE_PRINCIPAL_NAME --scopes $AKV_ID --role 4633458b-17de-408a-b874-0445c86b69e6 --query password --output tsv)
SP_APP_ID=$(az ad sp show --id http://$SERVICE_PRINCIPAL_NAME --query appId --output tsv)

echo "ClientID: [$SP_APP_ID]"
echo "ClientSecret: [$SP_PASSWD]"
```

# running example:

1) `$ mv .exampleenv .env`
2) update `.env` w/ your values
3) `KVAULT_SECRET_NAME=$SECRETNAME go run main.go` 

should print the secret value.