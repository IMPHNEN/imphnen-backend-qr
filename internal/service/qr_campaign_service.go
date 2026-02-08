package service

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"sync"
	"time"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/domain"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	ErrCampaignNotFound  = errors.New("campaign not found")
	ErrNoActiveCampaign  = errors.New("no active campaign")
	ErrInvalidImage      = errors.New("invalid image format, only PNG and JPEG are supported")
)

type QRCampaignService struct {
	repo              domain.QRCampaignRepository
	cacheMu           sync.RWMutex
	cachedQR          []byte
	cachedCampaignID  string
}

type CreateCampaignInput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewQRCampaignService(repo domain.QRCampaignRepository) *QRCampaignService {
	return &QRCampaignService{repo: repo}
}

func (s *QRCampaignService) CreateCampaign(input CreateCampaignInput, createdBy string) (*domain.QRCampaign, error) {
	// Generate QR PNG bytes (256x256, medium recovery)
	qrBytes, err := qrcode.Encode(input.URL, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	campaign := &domain.QRCampaign{
		Name:       input.Name,
		URL:        input.URL,
		QRCodeData: qrBytes,
		IsActive:   false,
		CreatedBy:  createdBy,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour),
	}

	// Create campaign first (inactive to avoid unique index conflict)
	if err := s.repo.Create(campaign); err != nil {
		return nil, err
	}

	// Then activate it (deactivates all others in a transaction)
	if err := s.repo.SetActive(campaign.ID); err != nil {
		return nil, err
	}
	campaign.IsActive = true

	// Update cache
	s.cacheMu.Lock()
	s.cachedQR = qrBytes
	s.cachedCampaignID = campaign.ID
	s.cacheMu.Unlock()

	return campaign, nil
}

func (s *QRCampaignService) GetActiveCampaign() (*domain.QRCampaign, error) {
	// Check cache first
	s.cacheMu.RLock()
	cachedID := s.cachedCampaignID
	s.cacheMu.RUnlock()

	if cachedID != "" {
		campaign, err := s.repo.FindByID(cachedID)
		if err != nil {
			return nil, err
		}
		if campaign != nil && campaign.IsActive {
			return campaign, nil
		}
		// Cache is stale, clear it
		s.cacheMu.Lock()
		s.cachedQR = nil
		s.cachedCampaignID = ""
		s.cacheMu.Unlock()
	}

	// Fallback to DB
	campaign, err := s.repo.FindActive()
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, ErrNoActiveCampaign
	}

	// Update cache
	s.cacheMu.Lock()
	s.cachedQR = campaign.QRCodeData
	s.cachedCampaignID = campaign.ID
	s.cacheMu.Unlock()

	return campaign, nil
}

func (s *QRCampaignService) GetAllCampaigns() ([]*domain.QRCampaign, error) {
	return s.repo.FindAll()
}

func (s *QRCampaignService) SetActiveCampaign(id string) error {
	if err := s.repo.SetActive(id); err != nil {
		return ErrCampaignNotFound
	}

	// Refresh cache
	campaign, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	s.cacheMu.Lock()
	if campaign != nil {
		s.cachedQR = campaign.QRCodeData
		s.cachedCampaignID = campaign.ID
	} else {
		s.cachedQR = nil
		s.cachedCampaignID = ""
	}
	s.cacheMu.Unlock()

	return nil
}

func (s *QRCampaignService) DeleteCampaign(id string) error {
	campaign, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if campaign == nil {
		return ErrCampaignNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache if deleted campaign was the cached one
	s.cacheMu.Lock()
	if s.cachedCampaignID == id {
		s.cachedQR = nil
		s.cachedCampaignID = ""
	}
	s.cacheMu.Unlock()

	return nil
}

func (s *QRCampaignService) ProcessImage(uploadedImage io.Reader) ([]byte, error) {
	// Get active campaign QR from cache
	s.cacheMu.RLock()
	qrData := s.cachedQR
	s.cacheMu.RUnlock()

	if qrData == nil {
		// Try loading from DB
		campaign, err := s.GetActiveCampaign()
		if err != nil {
			return nil, err
		}
		qrData = campaign.QRCodeData
	}

	// Decode uploaded image (PNG or JPEG)
	srcImg, _, err := image.Decode(uploadedImage)
	if err != nil {
		return nil, ErrInvalidImage
	}

	// Decode QR PNG bytes to image
	qrImg, err := png.Decode(bytes.NewReader(qrData))
	if err != nil {
		return nil, err
	}

	srcBounds := srcImg.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// Calculate QR size: ~1/5 of smallest dimension, minimum 100px
	minDim := srcW
	if srcH < minDim {
		minDim = srcH
	}
	qrSize := minDim / 5
	if qrSize < 100 {
		qrSize = 100
	}

	// Resize QR if needed (simple nearest-neighbor scaling)
	qrResized := resizeImage(qrImg, qrSize, qrSize)

	// Create output canvas
	padding := 10
	canvas := image.NewRGBA(srcBounds)
	draw.Draw(canvas, srcBounds, srcImg, srcBounds.Min, draw.Src)

	// Draw QR at bottom-right with padding
	qrRect := image.Rect(
		srcW-qrSize-padding,
		srcH-qrSize-padding,
		srcW-padding,
		srcH-padding,
	)
	draw.Draw(canvas, qrRect, qrResized, qrResized.Bounds().Min, draw.Over)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// resizeImage performs a simple nearest-neighbor resize
func resizeImage(src image.Image, width, height int) image.Image {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := srcBounds.Min.X + x*srcW/width
			srcY := srcBounds.Min.Y + y*srcH/height
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

// Register JPEG decoder for image.Decode
func init() {
	// image/jpeg and image/png decoders are registered by importing the packages
	_ = jpeg.Decode
}
