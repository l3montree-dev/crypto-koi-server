package models

// the owner of a cryptogotchi.
// there is a variety of possible authentication methods.
// for example, a user can be authenticated by their wallet and the device id
type User struct {
	Base
	// never return the device id.
	// this is a rather private information
	DeviceId string `json:"-" gorm:"type:varchar(255),unique"`
	// never return the wallet address of the user.
	WalletAddress string `json:"-" gorm:"type:varchar(255),unique"`
}
