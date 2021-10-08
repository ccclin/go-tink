package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/google/tink/go/keyset"
)

const (
	// Change this. AWS KMS, Google Cloud KMS and HashiCorp Vault are supported out of the box.
	// like "gcp-kms://projects/xxxxxx/locations/us-central1/keyRings/xxxxxx/cryptoKeys/xxxxxx"
	keyURI          = "gcp-kms://projects/...."
	credentialsPath = "credentials.json"
)

func main() {
	aad := []byte("Hello world")
	kh, _ := keyset.NewHandle(aead.AES128GCMKeyTemplate())
	// kh := getFromKMS("./key.json")
	createKeyJson("./key.json", kh)

	encode("file.tar", aad, kh)
	decode("file.tar", aad, kh)
}

func createKeyJson(path string, kh *keyset.Handle) {
	// Fetch the master key from a KMS.
	gcpClient, err := gcpkms.NewClientWithCredentials(keyURI, credentialsPath)
	if err != nil {
		log.Fatal(err)
	}
	registry.RegisterKMSClient(gcpClient)
	masterKey, err := gcpClient.GetAEAD(keyURI)
	if err != nil {
		log.Fatal(err)
	}
	f, _ := os.Create(path)
	writer := keyset.NewJSONWriter(f)
	defer f.Close()

	err = kh.Write(writer, masterKey)
	if err != nil {
		log.Fatal(err)
	}
}

func getFromKMS(path string) *keyset.Handle {
	// Fetch the master key from a KMS.
	gcpClient, err := gcpkms.NewClientWithCredentials(keyURI, credentialsPath)
	if err != nil {
		log.Fatal(err)
	}
	registry.RegisterKMSClient(gcpClient)
	masterKey, err := gcpClient.GetAEAD(keyURI)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Open(path)

	reader := keyset.NewJSONReader(f)

	kh, _ := keyset.Read(reader, masterKey)
	return kh
}

func encode(filePath string, aad []byte, kh *keyset.Handle) {
	a, err := aead.New(kh)
	if err != nil {
		log.Fatal(err)
	}

	fileBytes, _ := ioutil.ReadFile(filePath)
	ct, err := a.Encrypt(fileBytes, aad)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(fmt.Sprintf("%s.enc", filePath), ct, 0644)
}

func decode(filePath string, aad []byte, kh *keyset.Handle) {
	a, err := aead.New(kh)
	if err != nil {
		log.Fatal(err)
	}
	newFileBytes, _ := ioutil.ReadFile(fmt.Sprintf("%s.enc", filePath))

	pt, err := a.Decrypt(newFileBytes, aad)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(filePath, pt, 0644)
}
