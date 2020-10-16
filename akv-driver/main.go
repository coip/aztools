//this is just a sample driver to run Get on a secret,
//
//	put together w RBAC SP scoped to a keyvault, utilizing the "Key Vault Secrets User (preview)" role
//	creating sp:  https://docs.microsoft.com/en-us/cli/azure/ad/sp?view=azure-cli-latest#az_ad_sp_create_for_rbac
//	role details: https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#key-vault-secrets-user-preview
//
package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/coip/aztools/akv"

	_ "github.com/joho/godotenv/autoload"
)

const (
	vaultenv  = "KVAULT_NAME"
	secretenv = "KVAULT_SECRET_NAME"
)

// NewAuthorizerFromEnvironment @ L30 depends on some env as well: https://github.com/Azure/azure-sdk-for-go#more-authentication-details

var (
	vaultName  = os.Getenv(vaultenv)
	secretName = os.Getenv(secretenv)
)

func main() {

	v := akv.NewVault(vaultName)

	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("unable to create vault authorizer: %v\n", err)
		os.Exit(1)
	}

	basicClient := keyvault.New()
	basicClient.Authorizer = authorizer

	var (
		secret *string
		err    error
	)

	if secretName = "" {
		fmt.Println(secretenv + " not set.\n")
	}


	if secret, err = v.GetSecret(basicClient, secretName); err != nil {
		panic(err)
	}

	fmt.Print(*secret)

}
