package domain

import "time"

type QRCampaign struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	QRCodeData []byte    `json:"-"`
	IsActive   bool      `json:"is_active"`
	CreatedBy  string    `json:"created_by"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type QRCampaignRepository interface {
	Create(campaign *QRCampaign) error
	FindByID(id string) (*QRCampaign, error)
	FindActive() (*QRCampaign, error)
	FindAll() ([]*QRCampaign, error)
	SetActive(id string) error
	Delete(id string) error
}
