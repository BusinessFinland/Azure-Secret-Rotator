package helpers

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	aztables "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"log"
)

// Shared credential that can be used across the application.
var sharedCredential *azidentity.DefaultAzureCredential

func init() {
	var err error
	sharedCredential, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Printf("Error getting Azure credential: %v\n", err)
		// Handle the error appropriately here. Perhaps panic or set a package level error.
	}
}

// GraphHelper struct to be used for Microsoft Graph API, appClient has Client methods
type GraphHelper struct {
	clientSecretCredential *azidentity.DefaultAzureCredential
	appClient              *msgraphsdk.GraphServiceClient
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}

func (g *GraphHelper) InitializeGraphForAppAuth() error {
	// Create a new Azure clientSecretCredential
	g.clientSecretCredential = sharedCredential

	// Create an auth provider using the credential
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(g.clientSecretCredential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return err
	}

	// Create a request adapter using the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return err
	}

	// Create a Graph client using requestW adapter
	client := msgraphsdk.NewGraphServiceClient(adapter)
	log.Println("Graph client initialized")
	g.appClient = client

	return nil
}

func KeyVaultClient(vaultURI string) (keyVaultClient *azsecrets.Client, err error) {
	// Create KeyVault Client using KeyVault API
	keyVaultClient, err = azsecrets.NewClient(vaultURI, sharedCredential, nil)
	if err != nil {
		log.Println("Error creating KeyVault client: %v", err)
	}
	log.Println("KeyVault client initialized")
	return keyVaultClient, err

}

func AzTablesClient(serviceUrl string) (tablesClient *aztables.Client) {
	// Create  Azure Tables client
	tablesClient, err := aztables.NewClient(serviceUrl, sharedCredential, nil)
	if err != nil {
		log.Printf("Error creating Azure Tables client: %v", err)
	}
	log.Println("Azure Tables client initialized")
	return tablesClient

}
