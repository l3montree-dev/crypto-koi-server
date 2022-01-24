package http_dto

type LoginType string

const (
	LoginTypeDeviceId      LoginType = "deviceId"
	LoginTypeWalletAddress LoginType = "walletAddress"
)

type LoginRequest struct {
	Type          LoginType `json:"type"`
	DeviceId      string    `json:"deviceId"`
	WalletAddress string    `json:"walletAddress"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
