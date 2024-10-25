package apis

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"secretRotator/helpers"
)

func SendMail(notificationEmail string, appRegistrationName string, secretName string, keyVaultName string) {
	sendGridApiKey := helpers.SecretsMap["SENDGRID-API-KEY"]

	from := mail.NewEmail("Example User", "azuresecretrotator@example.com")
	subject := "New Secret has been created"

	to := mail.NewEmail("Example User", notificationEmail)
	plainTextContent := fmt.Sprintf(textContent, keyVaultName, secretName)
	htmlContent := "Is this needed?"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendGridApiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

}

var textContent = `You can find the new secret in the Key Vault: %s with the name: %s Make sure to update your applications with the new secret within 4 weeks. After that, the old secret will be deleted. If secret is programitically used from Key Vault, you can ignore this email.`
