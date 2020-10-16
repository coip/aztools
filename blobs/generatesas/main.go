package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// This is the name of the container and blob that we're creating a SAS to.
var (
	acctKey       = os.Getenv("key")
	acctName      = os.Getenv("storageacct")
	containerName = os.Getenv("container") // Container names require lowercase
)

const (
	ttl = 1 //hrs to expiry
)

//using https://github.com/Azure/azure-storage-blob-go/blob/master/azblob/zt_examples_test.go#L350-L401
func main() {

	blobName := os.Args[1]
	if blobName == "" {
		fmt.Printf("provide blobName, eg './generatesas <BLOBNAME>'")
		os.Exit(1)
	}

	fmt.Printf("generating sas for [%s/%s]\n", containerName, blobName)

	// Use your Storage account's name and key to create a credential object; this is required to sign a SAS.
	credential, err := azblob.NewSharedKeyCredential(acctName, acctKey)
	if err != nil {
		log.Fatal(err)
	}

	// Set the desired SAS signature values and sign them with the shared key credentials to get the SAS query parameters.
	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS, // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(ttl * time.Hour),
		ContainerName: containerName,
		BlobName:      blobName,

		// To produce a container SAS (as opposed to a blob SAS), assign to Permissions using
		// ContainerSASPermissions and make sure the BlobName field is "" (the default).
		Permissions: azblob.BlobSASPermissions{Add: false, Read: true, Write: false}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		log.Fatal(err)
	}

	// Create the SAS token URL partial for the resource you wish to access
	// Since this is a blob SAS, the URL is to the Azure storage blob.
	qp := sasQueryParams.Encode()

	fmt.Printf("access blob via sas[%s]\n", qp)

}
