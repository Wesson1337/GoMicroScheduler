package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"gomicrosched/internal/domain"
	"gomicrosched/pkg/kafka"
	"gomicrosched/pkg/smtp"
	"log"
)

func processMessage(smtpClient smtp.SMTPClient) kafka.ProcessFunc {
	return func(ctx context.Context, w kafka.MessageWriter, msg kafka.Message) {
		var notification domain.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("Error unmarshalling message: %v, err - %v", msg, err)
			return
		}

		emailMsg := getEmailMessage(notification)
		err = smtpClient.SendMail([]string{notification.Email}, []byte(emailMsg))
		if err != nil {
			log.Printf("Error sending message to %v: %v, err - %v", notification.Email, emailMsg, err)
			return
		}
		log.Printf("Successfully sent email to %v: %v", notification.Email, emailMsg)
	}
}

func getEmailMessage(notification domain.Notification) string {
	switch notification.Error {
	case "":
		return fmt.Sprintf(
			"Your job with ID %d has been successfully done. Result: %v",
			notification.JobID,
			notification.Result,
		)
	default:
		return fmt.Sprintf(
			"Your job with ID %d has been done with error: %v",
			notification.JobID,
			notification.Error,
		)
	}
}
