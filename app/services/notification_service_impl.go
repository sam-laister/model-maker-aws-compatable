package services

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/v4/messaging"
	"github.com/Soup666/diss-api/model"
	fcm "github.com/appleboy/go-fcm"
)

type NotificationServiceImpl struct {
}

func NewNotificationService() NotificationService {
	return &NotificationServiceImpl{}
}

func (s *NotificationServiceImpl) SendMessage(notification *model.Notification) (*model.Notification, error) {
	ctx := context.Background()
	client, err := fcm.NewClient(
		ctx,
		fcm.WithCredentialsFile("./service-account-key.json"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Send to a topic
	resp, err := client.Send(
		ctx,
		&messaging.Message{
			Notification: &messaging.Notification{
				Title: notification.Title,
				Body:  notification.Message,
			},
			Topic: fmt.Sprintf("%d", notification.UserID),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success:", resp.SuccessCount, "Failure:", resp.FailureCount)
	return notification, nil
}
