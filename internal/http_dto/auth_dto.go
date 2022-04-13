package http_dto

type LoginType string

type LoginRequest struct {
	WalletAddress *string `json:"walletAddress"`
	DeviceId      *string `json:"deviceId"`
}

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	// either wallet address or device id needs to be defined.
	WalletAddress *string `json:"walletAddress"`
	DeviceId      *string `json:"deviceId"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
