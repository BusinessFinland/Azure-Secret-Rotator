package main

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/go-github/v62/github"
	"log"
	"os"
	apis "secretRotator/apis"
	"secretRotator/helpers"
	graphhelper "secretRotator/helpers"
)

func main() {
	log.Println("Starting Secret Rotator Orhcestrator")

	helpers.GetEnvironmentVariablesFromKeyVault()

	tenantId := helpers.SecretsMap["TENANT-ID"]

	graphHelper, azTablesClient, githubClient := initializeClients()
	entities := graphhelper.GetEntitiesFromStorage(azTablesClient)

	for _, entity := range entities {
		processEntity(entity, graphHelper, *azTablesClient, githubClient, tenantId)
	}

}

func initializeClients() (graphHelper *graphhelper.GraphHelper, azTablesClient *aztables.Client, githubClient *github.Client) {
	// Initialize Graph Client
	graphHelper = graphhelper.NewGraphHelper()
	if err := graphHelper.InitializeGraphForAppAuth(); err != nil {
		log.Fatal("Error initializing Graph Client")
	}

	storageAccountName := helpers.SecretsMap["AZURE-STORAGE-ACCOUNT"]
	tableName := helpers.SecretsMap["AZURE-TABLE-NAME"]
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", storageAccountName, tableName)
	azTablesClient = helpers.AzTablesClient(serviceURL)

	// Check if GitHub private key exists, if not don't try to innitialize GitHub client
	if _, err := os.Stat("azuresecretrotator.private-key.pem"); err == nil {
		githubClient = apis.GithubClient()
		log.Println("GitHub client initialized")
		return graphHelper, azTablesClient, githubClient

	} else {
		return graphHelper, azTablesClient, nil
	}
}

func processEntity(entity graphhelper.Entity, graphHelper *graphhelper.GraphHelper, azTablesClient aztables.Client, githubClient *github.Client, tenantId string) {
	// Handle entity...

	var mySecretName string

	mySecretName = entity.SecretName

	if entity.KeyVaultName == "" {
		log.Printf("KeyVault is mandatory field")
		log.Printf("KeyVaultName is empty for entity: %s", entity.AppRegistrationName)
		return
	}

	if entity.PartitionKey == "0" {
		return
	}

	if entity.Reprocess == true {
		log.Println("Reprocessing entity...")
		vaultUri := fmt.Sprintf("https://%s.vault.azure.net/", entity.KeyVaultName)
		keyVaultClient, err := helpers.KeyVaultClient(vaultUri)
		if err != nil {
			log.Printf("error creating KeyVault client: %v", err)
		}
		secretValue, err := helpers.GetKeyVaultSecret(mySecretName, entity.KeyVaultName, keyVaultClient)
		if githubClient != nil {
			ok := apis.CreateOrUpdateGithubSecret(githubClient, entity.GithubRepository, entity.AppRegistrationObjectId, secretValue, tenantId, entity.SubscriptionID)
			if ok == true {
				reprocess := false
				helpers.MarkEntityForReprocessing(&azTablesClient, entity, reprocess)
			}
		}
	}

	if entity.Reprocess == false {
		secretValue, err := getAndUpdateAppRegistrationAndKeyVault(entity, graphHelper, azTablesClient, mySecretName)
		if err != nil {
			log.Printf("error getting and updating app registration and key vault: %v", err)
		}
		// Make Github call
		appId, err := graphHelper.GetApplicationId(entity.AppRegistrationObjectId)
		if err != nil {
			log.Printf("error getting application ID: %v", err)
		}
		if githubClient != nil {
			ok := apis.CreateOrUpdateGithubSecret(githubClient, entity.GithubRepository, appId, secretValue, tenantId, entity.SubscriptionID)
			if ok != true {
				reprocess := true
				helpers.MarkEntityForReprocessing(&azTablesClient, entity, reprocess)
			}
		}
	}

}

func getAndUpdateAppRegistrationAndKeyVault(entity graphhelper.Entity, graphHelper *graphhelper.GraphHelper, azTablesClient aztables.Client, mySecretName string) (secretValue string, err error) {
	objectId := entity.AppRegistrationObjectId
	keyVaultName := entity.KeyVaultName

	passwordCredentials, credentialCount, err := graphHelper.GetApplicationById(objectId, keyVaultName)
	if err != nil {
		log.Printf("error retrieving application by ID: %v", err)
		return "", err
	}

	if credentialCount == 0 {
		// More functionality to create the first secret could be added later on.
		log.Printf("no credentials found for object ID: %s", objectId)
		return "", err
	}

	switch {
	case credentialCount < 2:
		secretValue, err = graphHelper.GetCredentialExpiration(passwordCredentials, objectId, keyVaultName, mySecretName)
		if err != nil {
			log.Printf("error getting credential expiration: %v", err)
			return "", err
		}
		if secretValue != "" {
			ok := helpers.UploadToKeyVault(secretValue, keyVaultName, mySecretName)
			if !ok {
				// Cleanup actions after failing to upload to KeyVault.
				err = graphHelper.DeleteSecretKeyVaultFail(passwordCredentials, objectId)
				if err != nil {
					log.Printf("failed to clean up after upload to KeyVault failed: %v", err)
				}
				log.Printf("failed to upload secret to KeyVault for object ID: %s", objectId)
				return "", err
			}
			if entity.NotificationEmail != "" {
				//	apis.SendMail(entity.NotificationEmail, entity.AppRegistrationName, keyVaultName, mySecretName)
				log.Println("Send email, replace with sendgrid in Azure")
			}
		}

	case credentialCount >= 2:
		err = graphHelper.GetExpirationAndDelete(passwordCredentials, objectId)
		if err != nil {
			log.Printf("error managing multiple credential expirations: %v", err)
			return "", err
		}
	default:
		// Optional: handle other cases if necessary.
	}

	return secretValue, nil
}
