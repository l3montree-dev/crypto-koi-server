package service

import (
	"github.com/appleboy/go-fcm"
	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type NotificationSvc struct {
	fcmClient *fcm.Client
	logger    *logrus.Entry
}

type NotificationService interface {
	SendNotifications(to []*models.User, title, body string, data map[string]interface{}) map[*models.User]error
	SendNotification(to *models.User, title, body string, data map[string]interface{}) error
}

func NewNotificationSvc(apiKey string) *NotificationSvc {

	logger := orchardclient.Logger.WithField("component", "NotificationSvc")
	fcmService, err := fcm.NewClient(apiKey)

	if err != nil {
		logger.Fatal(err)
	}
	return &NotificationSvc{
		fcmClient: fcmService,
		logger:    logger,
	}
}

func (svc *NotificationSvc) SendNotifications(to []*models.User, title, body string, data map[string]interface{}) map[*models.User]error {
	var errors = make(map[*models.User]error)
	for _, user := range to {
		err := svc.SendNotification(user, title, body, data)
		if err != nil {
			errors[user] = err
		}
	}
	return errors
}

func (svc *NotificationSvc) SendNotification(to *models.User, title, body string, data map[string]interface{}) error {
	if to.PushNotificationToken == nil {
		// swallow the error.
		return nil
	}
	_, err := svc.fcmClient.Send(&fcm.Message{
		To:   *to.PushNotificationToken,
		Data: data,
		Notification: &fcm.Notification{
			Title: title,
			Body:  body,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
