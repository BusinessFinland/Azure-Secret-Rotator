package helpers

import (
	"context"
	"fmt"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"log"
	//	"os"
	"strings"
)

var SecretsMap = make(map[string]string)

func UploadToKeyVault(secretValue string, vaultURI string, mySecretName string) bool {
	// Upload secret to keyvault with name client_secret
	vaultURI = fmt.Sprintf("https://%s.vault.azure.net", vaultURI)
	keyVaultClient, err := KeyVaultClient(vaultURI)
	if err != nil {
		log.Printf("Error creating KeyVault client: %v\n", err)
	}
	version := ""
	_, err = keyVaultClient.GetSecret(context.TODO(), mySecretName, version, nil)
	if err != nil {
		ok := setKeyVaultSecret(secretValue, mySecretName, vaultURI, keyVaultClient)
		if !ok {
			return false
		}
	}
	return true
}

func setKeyVaultSecret(secretValue string, mySecretName string, vaultURI string, keyVaultClient *azsecrets.Client) bool {
	// Upload secret to keyvault with name client_secret
	// Create slice of strings from mySecretName, remove all underscores, and join them back together as camel case
	mySecretName = trimSecretName(mySecretName)

	params := azsecrets.SetSecretParameters{Value: &secretValue}
	_, err := keyVaultClient.SetSecret(context.TODO(), mySecretName, params, nil)
	if err != nil {
		log.Printf("Error uploading secret to keyvault: %v", err)
		return false
	}
	return true
}

func GetKeyVaultSecret(mySecretName string, vaultURI string, keyVaultClient *azsecrets.Client) (string, error) {
	// Get secret from keyvault
	mySecretName = trimSecretName(mySecretName)
	version := ""
	secret, err := keyVaultClient.GetSecret(context.TODO(), mySecretName, version, nil)
	if err != nil {
		log.Printf("Error getting secret from keyvault: %v", err)
		return "", err
	}
	return *secret.Value, nil
}

func trimSecretName(mySecretName string) string {
	// Create slice of strings from secretName, remove all underscores, and join them back together as camel case
	mySecretName = strings.ToUpper(mySecretName)

	mySecretName = strings.ReplaceAll(mySecretName, "_", "-")
	return mySecretName
}

func GetEnvironmentVariablesFromKeyVault() {
	// Get Environment Variables from KeyVault
	// Use Managed Identity to get Environment variables
	// For container, use os.getEnv("KEYVAULT_NAME")

	vaultURI := fmt.Sprintf("https://%s.vault.azure.net/", "kvazsecrotator")
	fmt.Println("Remember to change to os.Getenv before moving to container")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Error getting Azure credentials: %v", err)
	}

	client, err := azsecrets.NewClient(vaultURI, cred, nil)

	// List secrets
	pager := client.NewListSecretPropertiesPager(nil)
	secretNames := []string{}

	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, secret := range page.Value {
			secretNames = append(secretNames, secret.ID.Name())
		}
		for _, secretName := range secretNames {
			secret, err := client.GetSecret(context.TODO(), secretName, "", nil)
			if err != nil {
				log.Fatalf("Error getting secret: %v", err)
			}
			SecretsMap[secretName] = *secret.Value
		}
	}
}
