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
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"

	_ "github.com/joho/godotenv/autoload"
)

type Vault string

func (v Vault) String() string {
	return string(v)
}

func NewVault(vaultname string) Vault {
	return Vault(vaultname)
}

const vaultURLfmt = "https://%s.vault.azure.net"

func (v Vault) URL() string {
	return fmt.Sprintf(vaultURLfmt, v)
}

func (v Vault) ListSecrets(basicClient keyvault.BaseClient) {
	secretList, err := basicClient.GetSecrets(context.Background(), v.URL(), nil)
	if err != nil {
		fmt.Printf("unable to get list of secrets: %v\n", err)
		os.Exit(1)
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
}

func (v Vault) GetSecret(basicClient keyvault.BaseClient, secname string) (*string, error) {
	secretResp, err := basicClient.GetSecret(context.Background(), v.URL(), secname, "")
	if err != nil {
		return nil, fmt.Errorf("unable to get value for secret: %v\n", err)
	}
	return secretResp.Value, nil
}

//CreateUpdateSecret will create or update a secret, and return the ID. on error, ID will be nil
func (v Vault) CreateUpdateSecret(basicClient keyvault.BaseClient, secname, secvalue string) (*string, error) {
	var secParams keyvault.SecretSetParameters
	secParams.Value = &secvalue
	newBundle, err := basicClient.SetSecret(context.Background(), v.URL(), secname, secParams)
	if err != nil {
		return nil, fmt.Errorf("unable to add/update secret: %v\n", err)
	}
	return newBundle.ID, nil
}

func (v Vault) DeleteSecret(basicClient keyvault.BaseClient, secname string) error {
	_, err := basicClient.DeleteSecret(context.Background(), v.URL(), secname)
	if err != nil {
		return fmt.Errorf("error deleting secret: %v\n", err)
	}
	return nil
}
