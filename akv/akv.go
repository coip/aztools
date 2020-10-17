package akv

// adapted from https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/keyvault/examples/go-keyvault-msi-example.go

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// for reading a secret @ runtime

// You need to set four environment variables before using the app:
// AZURE_SUBSCRIPTION_ID
// AZURE_TENANT_ID
// AZURE_CLIENT_ID 		($SP_APP_ID from genSP.sh)
// AZURE_CLIENT_SECRET	($SP_PASSWD from genSP.sh)
// KVAULT_SECRET_NAME to the secret's name.

import (
	"context"
	"fmt"
	"path"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"

	_ "github.com/joho/godotenv/autoload"
)

type Vault struct {
	Name string

	client keyvault.BaseClient
}

func (v Vault) String() string {
	return v.Name
}

func NewVault(vaultname string) *Vault {
	var (
		authorizer autorest.Authorizer
		err        error
	)

	if authorizer, err = kvauth.NewAuthorizerFromEnvironment(); err != nil {
		panic(err)
	}

	var v = &Vault{Name: vaultname}

	v.client = keyvault.New()
	v.client.Authorizer = authorizer

	return v
}

const vaultURLfmt = "https://%s.vault.azure.net"

func (v Vault) URL() string {
	return fmt.Sprintf(vaultURLfmt, v)
}

func (v Vault) ListSecrets() error {
	var (
		secretList keyvault.SecretListResultPage
		err        error
	)

	if secretList, err = v.client.GetSecrets(context.Background(), v.URL(), nil); err != nil {
		return fmt.Errorf("unable to get list of secrets: %v\n", err)
	}

	// group by ContentType
	secWithType := make(map[string][]string)
	secWithoutType := make([]string, 1)
	for _, secret := range secretList.Values() {
		if secret.ContentType != nil {
			_, exists := secWithType[*secret.ContentType]
			if exists {
				secWithType[*secret.ContentType] = append(secWithType[*secret.ContentType], path.Base(*secret.ID))
			} else {
				tempSlice := make([]string, 1)
				tempSlice[0] = path.Base(*secret.ID)
				secWithType[*secret.ContentType] = tempSlice
			}
		} else {
			secWithoutType = append(secWithoutType, path.Base(*secret.ID))
		}
	}

	for k, v := range secWithType {
		fmt.Println(k)
		for _, sec := range v {
			fmt.Println(sec)
		}
	}
	for _, wov := range secWithoutType {
		fmt.Println(wov)
	}

	return nil

}

func (v Vault) GetSecret(secname string) (*string, error) {
	var (
		secret keyvault.SecretBundle

		err error
	)

	if secret, err = v.client.GetSecret(context.Background(), v.URL(), secname, ""); err != nil {
		return nil, fmt.Errorf("unable to get value for secret: %v\n", err)
	}

	return secret.Value, nil

}

//CreateUpdateSecret will create or update a secret, and return the ID. on error, ID will be nil
func (v Vault) CreateUpdateSecret(secname, secvalue string) (*string, error) {
	var (
		params = keyvault.SecretSetParameters{Value: &secvalue}
	)

	if newBundle, err := v.client.SetSecret(context.Background(), v.URL(), secname, params); err != nil {
		return nil, fmt.Errorf("unable to add/update secret: %v\n", err)
	} else {
		return newBundle.ID, nil
	}

}

func (v Vault) DeleteSecret(secname string) error {
	if _, err := v.client.DeleteSecret(context.Background(), v.URL(), secname); err != nil {
		return fmt.Errorf("error deleting secret: %v\n", err)
	}
	return nil
}
