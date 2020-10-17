//this is just a sample driver to run Get on a secret,
//
//	put together w RBAC SP scoped to a keyvault, utilizing the "Key Vault Secrets User (preview)" role
//	creating sp:  https://docs.microsoft.com/en-us/cli/azure/ad/sp?view=azure-cli-latest#az_ad_sp_create_for_rbac
//	role details: https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#key-vault-secrets-user-preview
//
package main

import (
	"log"
	"os"

	"github.com/coip/aztools/akv"
)

const (
	vaultenv  = "KVAULT_NAME"
	secretenv = "KVAULT_SECRET_NAME"
)

// NewAuthorizerFromEnvironment depends on some env as well: https://github.com/Azure/azure-sdk-for-go#more-authentication-details

var (
	vaultName  = os.Getenv(vaultenv)
	secretName = os.Getenv(secretenv)
)

func main() {

	var (
		v = akv.NewVault(vaultName)

		secret *string
		err    error
	)

	if secretName == "" {
		log.Fatal(secretenv + " not set.\n")
	}

	if secret, err := v.GetSecret(secretName); err != nil {
		panic(err)
	} else {
		log.Print(*secret)
	}

}
