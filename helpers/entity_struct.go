package helpers

import ()

type Entity struct {
	PartitionKey            string
	RowKey                  string
	AppRegistrationName     string
	AppRegistrationObjectId string
	KeyVaultName            string
	GithubRepository        string
	SubscriptionID          string
	Owner                   string
	NotificationEmail       string
	SecretName              string
	Reprocess               bool
}
