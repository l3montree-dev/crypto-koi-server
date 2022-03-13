package service

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/cryptokoi"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/pkg/leader"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type CryptogotchiSvc interface {
	repositories.CryptogotchiRepository
	GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error)
	GenerateWithFixedTokenId(user *models.User, id uuid.UUID) (models.Cryptogotchi, error)
	MarkAsNft(crypt *models.Cryptogotchi) error
	GetNotificationListener() leader.Listener
	UpdateRanks() error
}

type CryptogotchiService struct {
	repositories.CryptogotchiRepository
	userRep                  repositories.UserRepository
	logger                   *logrus.Entry
	timeBetweenNotifications time.Duration
	notificationSvc          NotificationService
	notifications            config.PreloadedNotifications
}

func NewCryptogotchiService(rep repositories.CryptogotchiRepository, userRep repositories.UserRepository, notificationSvc NotificationService) CryptogotchiSvc {
	logger := orchardclient.Logger.WithField("component", "CryptogotchiService")
	notifications := config.GetNotifications()
	return &CryptogotchiService{
		CryptogotchiRepository:   rep,
		logger:                   logger,
		timeBetweenNotifications: 1 * time.Minute,
		notificationSvc:          notificationSvc,
		notifications:            notifications,
	}
}

func (svc *CryptogotchiService) UpdateRanks() error {
	elements, err := svc.GetLeaderboard()
	if err != nil {
		return err
	}
	for rank, el := range elements {
		el.Rank = rank + 1
		err = svc.Save(&el)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *CryptogotchiService) GenerateWithFixedTokenId(user *models.User, id uuid.UUID) (models.Cryptogotchi, error) {
	foodValue := config.DEFAULT_FOOD_VALUE
	foodDrainValue := config.DEFAULT_FOOD_DRAIN
	now := time.Now()

	tokenId, err := util.UuidToUint256(id.String())

	if err != nil {
		return models.Cryptogotchi{}, err
	}

	koi := cryptokoi.NewKoi(tokenId.String())
	name := strings.Title((koi.GetAttributes().KoiType))

	newCrypt := models.Cryptogotchi{
		// TODO: generate a random name
		Base: models.Base{
			Id: id,
		},
		Name:               util.Str(name),
		OwnerId:            user.Id,
		Food:               foodValue,
		FoodDrain:          foodDrainValue,
		PredictedDeathDate: now.Add(time.Duration(foodValue/foodDrainValue) * time.Minute),
		SnapshotValid:      now,
	}
	err = svc.Create(&newCrypt)
	return newCrypt, err
}

func (svc *CryptogotchiService) GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return models.Cryptogotchi{}, err
	}
	return svc.GenerateWithFixedTokenId(user, id)
}

func (svc *CryptogotchiService) MarkAsNft(crypt *models.Cryptogotchi) error {
	crypt.IsValidNft = true
	return svc.Save(crypt)
}

func (svc *CryptogotchiService) getUserForNotifications(phase string) ([]models.User, error) {
	duration := time.Duration(config.GetNotifications()[phase].HoursBeforeDeath) * time.Hour
	startTime := time.Now().Add(duration)
	endTime := startTime.Add(svc.timeBetweenNotifications)
	cryptogotchies, err := svc.GetCryptogotchiesWithPredictedDeathDateBetween(startTime, endTime)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(cryptogotchies))

	for i, crypt := range cryptogotchies {
		user, err := svc.userRep.GetById(crypt.OwnerId.String())
		if err != nil {
			return nil, err
		}
		users[i] = user
	}
	return users, nil
}
func (svc *CryptogotchiService) sendPhaseNotification(phase string) error {
	users, err := svc.getUserForNotifications(phase)
	if err != nil {
		return err
	}

	orchardclient.Logger.Info("Sending:", len(users), "notifications for phase:", phase)

	wg := sync.WaitGroup{}
	wg.Add(len(users))
	for _, user := range users {
		go func(u models.User) {
			defer wg.Done()
			// get the notification
			r := rand.Intn(len(svc.notifications[phase].Notifications))
			notification := svc.notifications[phase].Notifications[r]
			err = svc.notificationSvc.SendNotification(&u, notification.Title, notification.Body, nil)
			orchardclient.Logger.Error(err)
		}(user)
	}
	wg.Wait()
	return nil
}

func (svc *CryptogotchiService) GetNotificationListener() leader.Listener {
	return leader.NewListener(func(cancelChan <-chan struct{}) {
		for {
			select {
			case <-cancelChan:
				// just stop on canceling.
				return
			case <-time.After(svc.timeBetweenNotifications):
				now := time.Now()
				wg := sync.WaitGroup{}
				wg.Add(len(svc.notifications))
				for phase := range svc.notifications {
					go func(p string) {
						defer wg.Done()
						err := svc.sendPhaseNotification(p)
						if err != nil {
							orchardclient.Logger.Error(err)
						}
					}(phase)
				}
				wg.Wait()
				duration := time.Since(now)
				if duration > svc.timeBetweenNotifications {
					orchardclient.Logger.Errorf("Notification listener took too long: %s", duration)
				} else {
					svc.logger.WithField("took", duration).Info("finished notification listener")
				}
			}
		}
	})
}
