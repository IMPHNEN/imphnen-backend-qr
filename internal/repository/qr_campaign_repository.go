package repository

import (
	"database/sql"
	"time"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/domain"
	"github.com/google/uuid"
)

type qrCampaignRepository struct {
	db *sql.DB
}

func NewQRCampaignRepository(db *sql.DB) domain.QRCampaignRepository {
	return &qrCampaignRepository{db: db}
}

func (r *qrCampaignRepository) Create(campaign *domain.QRCampaign) error {
	if campaign.ID == "" {
		campaign.ID = uuid.New().String()
	}
	now := time.Now()
	campaign.CreatedAt = now
	campaign.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO qr_campaigns (id, name, url, qr_code_data, is_active, created_by, expires_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		campaign.ID, campaign.Name, campaign.URL, campaign.QRCodeData, campaign.IsActive,
		campaign.CreatedBy, campaign.ExpiresAt, campaign.CreatedAt, campaign.UpdatedAt,
	)
	return err
}

func (r *qrCampaignRepository) FindByID(id string) (*domain.QRCampaign, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, nil
	}

	campaign := &domain.QRCampaign{}
	err := r.db.QueryRow(
		`SELECT id, name, url, qr_code_data, is_active, created_by, expires_at, created_at, updated_at
		 FROM qr_campaigns WHERE id = $1::uuid`, id,
	).Scan(&campaign.ID, &campaign.Name, &campaign.URL, &campaign.QRCodeData, &campaign.IsActive,
		&campaign.CreatedBy, &campaign.ExpiresAt, &campaign.CreatedAt, &campaign.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return campaign, err
}

func (r *qrCampaignRepository) FindActive() (*domain.QRCampaign, error) {
	campaign := &domain.QRCampaign{}
	err := r.db.QueryRow(
		`SELECT id, name, url, qr_code_data, is_active, created_by, expires_at, created_at, updated_at
		 FROM qr_campaigns WHERE is_active = true LIMIT 1`,
	).Scan(&campaign.ID, &campaign.Name, &campaign.URL, &campaign.QRCodeData, &campaign.IsActive,
		&campaign.CreatedBy, &campaign.ExpiresAt, &campaign.CreatedAt, &campaign.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return campaign, err
}

func (r *qrCampaignRepository) FindAll() ([]*domain.QRCampaign, error) {
	rows, err := r.db.Query(
		`SELECT id, name, url, qr_code_data, is_active, created_by, expires_at, created_at, updated_at
		 FROM qr_campaigns ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*domain.QRCampaign
	for rows.Next() {
		campaign := &domain.QRCampaign{}
		if err := rows.Scan(&campaign.ID, &campaign.Name, &campaign.URL, &campaign.QRCodeData, &campaign.IsActive,
			&campaign.CreatedBy, &campaign.ExpiresAt, &campaign.CreatedAt, &campaign.UpdatedAt); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}
	return campaigns, rows.Err()
}

func (r *qrCampaignRepository) SetActive(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deactivate all campaigns
	if _, err := tx.Exec(`UPDATE qr_campaigns SET is_active = false, updated_at = $1`, time.Now()); err != nil {
		return err
	}

	// Activate the target campaign
	result, err := tx.Exec(`UPDATE qr_campaigns SET is_active = true, updated_at = $1 WHERE id = $2::uuid`, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (r *qrCampaignRepository) Delete(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return err
	}

	_, err := r.db.Exec(`DELETE FROM qr_campaigns WHERE id = $1::uuid`, id)
	return err
}
