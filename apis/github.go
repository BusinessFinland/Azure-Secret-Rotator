package apis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/google/go-github/v62/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
	"log"
	"os"
	"secretRotator/helpers"
	"strconv"
)

type AzureAuthCreds struct {
	ClientID       string
	ClientSecret   string
	TenantID       string
	SubscriptionID string
}

func GithubClient() (githubClient *github.Client) {
	// Intialize github client using github app and return the client for further use.

	gitHubAppId := helpers.SecretsMap["GITHUB-APP-ID"]
	InstallationIdStr := helpers.SecretsMap["GITHUB-INSTALLATION-ID"]

	privateKey, err := (os.ReadFile("./azuresecretrotator.private-key.pem"))
	if err != nil {
		log.Printf("Error reading private key: %v", err)
	}

	appIdParse, _ := strconv.ParseInt((gitHubAppId), 10, 64)
	installationIdParse, _ := strconv.ParseInt((InstallationIdStr), 10, 64)

	appTokenSource, err := githubauth.NewApplicationTokenSource(appIdParse, []byte(privateKey))
	if err != nil {
		log.Printf("Error creating application token source: %v", err)
		return
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(installationIdParse, appTokenSource)

	// oauth2.NewClient create a new http.Client that adds an Authorization header with the token.
	// Transport src use oauth2.ReuseTokenSource to reuse the token.
	// The token will be reused until it expires.
	// The token will be refreshed if it's expired.
	httpClient := oauth2.NewClient(context.Background(), installationTokenSource)

	githubClient = github.NewClient(httpClient)

	return githubClient

}

func EncodeWithPublicKey(text string, publicKey string) (string, error) {
	// Github secret API requires encrypted secret
	// Decode the public key from base64

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	// Decode the public key
	var publicKeyDecoded [32]byte
	copy(publicKeyDecoded[:], publicKeyBytes)

	// Encrypt the secret value
	encrypted, err := box.SealAnonymous(nil, []byte(text), (*[32]byte)(publicKeyBytes), rand.Reader)

	if err != nil {
		return "", err
	}
	// Encode the encrypted value in base64
	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)

	return encryptedBase64, nil
}

func CreateOrUpdateGithubSecret(githubClient *github.Client, repository string, applicationId string, secretValue, tenantId string, subscriptionId string) bool {
	// Create or update repository secret using public key generated above
	// Calls EncodeWithPublicKey to encrypt the secret using public key

	organization := helpers.SecretsMap["ORGANIZATION-NAME"]

	// Get github public key
	pubKey, _, err := githubClient.Actions.GetRepoPublicKey(context.Background(), organization, repository)
	if err != nil {
		log.Printf("Error getting installations:", err)
		return false
	}

	azureAuthCred := AzureAuthCreds{
		ClientID:       applicationId,
		ClientSecret:   secretValue,
		TenantID:       tenantId,
		SubscriptionID: subscriptionId,
	}

	azureAuthCredJSON, err := json.Marshal(azureAuthCred)
	if err != nil {
		log.Printf("Error marshalling AzureAuthCreds to JSON:", err)
		return false
	}

	// encrypt secret
	encryptedSecret, err := EncodeWithPublicKey(string(azureAuthCredJSON), *pubKey.Key)

	eSecret := &github.EncryptedSecret{
		Name:           "CLIENT_SECRET",
		KeyID:          *pubKey.KeyID,
		EncryptedValue: encryptedSecret,
	}

	// create or update repository secret using public key generated above
	_, err = githubClient.Actions.CreateOrUpdateRepoSecret(context.Background(), organization, repository, eSecret)
	if err != nil {
		log.Printf("Error creating or updating secret:", err)
		return false
	}
	return true

}
