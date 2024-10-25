package helpers

import (
	"context"
	"github.com/google/uuid"
	graphapplications "github.com/microsoftgraph/msgraph-sdk-go/applications"
	models "github.com/microsoftgraph/msgraph-sdk-go/models"
	"log"
	"time"
)

func (g *GraphHelper) GetApplicationById(objectId string, keyVaultName string) (passwordCredentials []models.PasswordCredentialable, credentialCount int, err error) {
	// Get application information by objectId

	application, err := g.appClient.Applications().ByApplicationId(objectId).Get(context.Background(), nil)
	if err != nil {
		log.Printf("Error getting application information: %v\n", err)
	}

	passwordCredentials = application.GetPasswordCredentials()
	credentialCount = len(passwordCredentials)

	return passwordCredentials, credentialCount, err
}

func (g *GraphHelper) GetCredentialExpiration(passwordCredentials []models.PasswordCredentialable, objectId string, keyVaultName string, mySecretName string) (secretValue string, err error) {
	// Get
	secretEndDate := (passwordCredentials[0].GetEndDateTime())
	now := time.Now().UTC()
	timeFromNow := now.AddDate(0, 0, 30)
	endDate := *secretEndDate

	var displayName string
	if endDate.Before(timeFromNow) {
		if mySecretName == "" {
			displayName = *passwordCredentials[0].GetDisplayName()
		} else {
			displayName = mySecretName
		}
		secretValue, err := g.createNewSecret(objectId, displayName)
		if err != nil {
			log.Printf("Error creating new secret: %v\n", err)
		}
		return secretValue, nil
	}

	return "", err
}

func (g *GraphHelper) createNewSecret(objectId string, displayName string) (secretValue string, err error) {
	// Create a new secret for the application, if existing one is expiring

	requestBody := graphapplications.NewItemAddPasswordPostRequestBody()
	passwordCredential := models.NewPasswordCredential()
	passwordCredential.SetDisplayName(&displayName)
	requestBody.SetPasswordCredential(passwordCredential)
	addPassword, err := g.appClient.Applications().ByApplicationId(objectId).AddPassword().Post(context.Background(), requestBody, nil)
	if err != nil {
		log.Printf("Error creating new secret: %v\n", err)
		return "", err
	}
	secretValue = *addPassword.GetSecretText()

	return secretValue, nil
}

func (g *GraphHelper) DeleteSecretKeyVaultFail(passwordCredentials []models.PasswordCredentialable, objectId string) error {
	// Credential Indexing changes, new secret is always created at index 0.
	// Needs to be cleaned up!!

	application, err := g.appClient.Applications().ByApplicationId(objectId).Get(context.Background(), nil)
	credentialIndex := 0
	if err != nil {
		log.Printf("Error getting application information: %v\n", err)
	}
	passwordCredentials = application.GetPasswordCredentials()

	secret := passwordCredentials[credentialIndex]
	secretId := secret.GetKeyId()
	err = g.deleteSecret(objectId, *secretId)
	if err != nil {
		log.Printf("Error deleting secret: %v\n", err)
	}
	return nil
}

func (g *GraphHelper) GetApplicationId(objectId string) (string, error) {
	// Get application ID by objectId

	application, err := g.appClient.Applications().ByApplicationId(objectId).Get(context.Background(), nil)
	if err != nil {
		log.Printf("Error getting application information: %v\n", err)
	}
	appID := application.GetAppId()
	if appID == nil {
		log.Printf("Error getting application ID")
	}

	return *appID, nil
}

func (g *GraphHelper) GetExpirationAndDelete(passwordCredentials []models.PasswordCredentialable, objectId string) error {
	now := time.Now().UTC()
	endOfDay := now.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
	// If there are 2 credentials, and one of them has expired, delete it
	for _, secret := range passwordCredentials {
		secretEndDate := (secret.GetEndDateTime())

		// Check if secret has expired or will expire within the next 24 hours
		if secretEndDate.Before(now) || secretEndDate.After(now) && secretEndDate.Before(endOfDay) {
			log.Println("Deleting secret: ", secret, objectId)
			secretId := secret.GetKeyId()
			err := g.deleteSecret(objectId, *secretId)
			if err != nil {
				log.Printf("Error deleting secret: %v\n", err)
			}
		}
	}

	return nil
}

func (g *GraphHelper) deleteSecret(objectId string, secretId uuid.UUID) error {
	requestBody := graphapplications.NewItemRemovePasswordPostRequestBody()
	keyId := secretId
	requestBody.SetKeyId(&keyId)

	err := g.appClient.Applications().ByApplicationId(objectId).RemovePassword().Post(context.Background(), requestBody, nil)

	if err != nil {
		log.Printf("Error deleting secret: %v\n", err)
	}

	return nil
}
