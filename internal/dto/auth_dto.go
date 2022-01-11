package dto

type LoginRequest struct {
	Type          string `json:"type"`
	DeviceId      string `json:"deviceId"`
	WalletAddress string `json:"walletAddress"`
}
