package http_dto

type LoginType string

type LoginRequest struct {
	WalletAddress string `json:"walletAddress"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
