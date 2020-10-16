package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	_ "github.com/joho/godotenv/autoload"
)

// From the Azure portal, get your Storage account's name and account key.
var (
	accountName, accountKey = os.Getenv("acctName"), os.Getenv("acctKey")
	containerName           = os.Getenv("container")
)

func main() {

	// Use your Storage account's name and key to create a credential object; this is used to access your account.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create a request pipeline that is used to process HTTP(S) requests and responses. It requires
	// your account credentials. In more advanced scenarios, you can configure telemetry, retry policies,
	// logging, and other options. Also, you can configure multiple request pipelines for different scenarios.
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your Storage account blob service URL endpoint.
	// The URL typically looks like this:
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	// Create an ServiceURL object that wraps the service URL and a request pipeline.
	serviceURL := azblob.NewServiceURL(*u, p)

	// Now, you can use the serviceURL to perform various container and blob operations.

	// All HTTP operations allow you to specify a Go context.Context object to control cancellation/timeout.
	ctx := context.Background() // This example uses a never-expiring context.

	// This example shows several common operations just to get you started.

	// Create a URL that references a to-be-created container in your Azure Storage account.
	// This returns a ContainerURL object that wraps the container's URL and a request pipeline (inherited from serviceURL)
	containerURL := serviceURL.NewContainerURL(containerName) // Container names require lowercase

	// Create the container on the service (with no metadata and no public access)
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		log.Fatal(err)
	}
	var f *os.File
	const filename = "imgexample.png"
	if f, err = os.Open(filename); err != nil {
		panic(err)
	}
	// Create a URL that references a to-be-created blob in your Azure Storage account's container.
	// This returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := containerURL.NewBlockBlobURL(filename) // Blob names can be mixed case

	// Create the blob with string (plain text) content.
	_, err = blobURL.Upload(ctx, f, azblob.BlobHTTPHeaders{ContentType: "image/png"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	// Download the blob's contents and verify that it worked correctly
	get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		log.Fatal(err)
	}

	reader := get.Body(azblob.RetryReaderOptions{})
	output, err := os.Create("testoutput.png")
	io.Copy(output, reader)
	// The client must close the response body when finished with it
	if err = reader.Close(); err != nil {
		panic(err)
	}

	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			log.Fatal(err)
		}
		// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print("Blob name: " + blobInfo.Name + "\n")
		}
	}

	// Delete the blob we created earlier.
	_, err = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	/*
		// Delete the container we created earlier.
		_, err = containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
		if err != nil {
			log.Fatal(err)
		}
	*/

}
