package handler

import (
	"log"
	"net/http"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/service"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/utils"
	"github.com/labstack/echo/v4"
)

type QRCampaignHandler struct {
	campaignService *service.QRCampaignService
}

func NewQRCampaignHandler(campaignService *service.QRCampaignService) *QRCampaignHandler {
	return &QRCampaignHandler{campaignService: campaignService}
}

func (h *QRCampaignHandler) CreateCampaign(c echo.Context) error {
	var input service.CreateCampaignInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", "bad_request")
	}

	if input.Name == "" || input.URL == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "name and url are required", "validation_error")
	}

	createdBy := c.Get("user_id").(string)

	campaign, err := h.campaignService.CreateCampaign(input, createdBy)
	if err != nil {
		log.Printf("[ERROR] CreateCampaign: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create campaign", "internal_error")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "campaign created", campaign)
}

func (h *QRCampaignHandler) GetAllCampaigns(c echo.Context) error {
	campaigns, err := h.campaignService.GetAllCampaigns()
	if err != nil {
		log.Printf("[ERROR] GetAllCampaigns: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch campaigns", "internal_error")
	}

	return utils.SuccessResponse(c, http.StatusOK, "campaigns retrieved", campaigns)
}

func (h *QRCampaignHandler) SetActiveCampaign(c echo.Context) error {
	id := c.Param("id")

	if err := h.campaignService.SetActiveCampaign(id); err != nil {
		log.Printf("[ERROR] SetActiveCampaign: %v", err)
		return utils.ErrorResponse(c, http.StatusNotFound, "campaign not found", "campaign_not_found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "campaign activated", nil)
}

func (h *QRCampaignHandler) DeleteCampaign(c echo.Context) error {
	id := c.Param("id")

	if err := h.campaignService.DeleteCampaign(id); err != nil {
		if err == service.ErrCampaignNotFound {
			return utils.ErrorResponse(c, http.StatusNotFound, "campaign not found", "campaign_not_found")
		}
		log.Printf("[ERROR] DeleteCampaign: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete campaign", "internal_error")
	}

	return utils.SuccessResponse(c, http.StatusOK, "campaign deleted", nil)
}

func (h *QRCampaignHandler) ProcessImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "image file is required", "validation_error")
	}

	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "failed to read image file", "bad_request")
	}
	defer src.Close()

	result, err := h.campaignService.ProcessImage(src)
	if err != nil {
		if err == service.ErrNoActiveCampaign {
			return utils.ErrorResponse(c, http.StatusNotFound, "no active campaign", "no_active_campaign")
		}
		if err == service.ErrInvalidImage {
			return utils.ErrorResponse(c, http.StatusBadRequest, "invalid image format, only PNG and JPEG are supported", "invalid_image")
		}
		log.Printf("[ERROR] ProcessImage: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to process image", "internal_error")
	}

	return c.Blob(http.StatusOK, "image/png", result)
}
