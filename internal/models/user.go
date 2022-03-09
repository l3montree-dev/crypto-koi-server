package models

// the owner of a cryptogotchi.
// there is a variety of possible authentication methods.
// for example, a user can be authenticated by their wallet and the device id
type User struct {
	Base
	Cryptogotchies []Cryptogotchi `json:"cryptogotchies" gorm:"foreignKey:OwnerId;references:Id"`
	Name           *string        `json:"name" gorm:"type:varchar(255);default:null"`
	// never return the wallet address of the user.
	WalletAddress         *string `json:"-" gorm:"type:varchar(255);unique"`
	DeviceId              *string `json:"-" gorm:"type:varchar(255);unique"`
	RefreshToken          string  `json:"-" gorm:"type:varchar(255);not null;unique"`
	PushNotificationToken *string `json:"-" gorm:"type:varchar(255)"`
}
