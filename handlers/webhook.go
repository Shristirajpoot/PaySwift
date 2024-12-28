package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"paystack-payment/services"
	"paystack-payment/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebhookHandler struct {
	db              *gorm.DB
	logger          *utils.Logger
	paystackService *services.PaystackService
}

func NewWebhookHandler(db *gorm.DB, logger *utils.Logger, paystackService *services.PaystackService) *WebhookHandler {
	return &WebhookHandler{db: db, logger: logger, paystackService: paystackService}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("Failed to read request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	signature := c.GetHeader("X-Paystack-Signature")
	if signature == "" {
		h.logger.Error("Missing Paystack signature")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Paystack signature"})
		return
	}

	if !h.paystackService.VerifyWebhookSignature(signature, body) {
		h.logger.Error("Invalid Paystack signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.Error("Failed to parse webhook payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	event, ok := payload["event"].(string)
	if !ok || event != "charge.success" {
		c.JSON(http.StatusOK, gin.H{"message": "Unhandled event"})
		return
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		h.logger.Error("Invalid data in webhook payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload data"})
		return
	}

	reference, ok := data["reference"].(string)
	if !ok {
		h.logger.Error("Invalid reference in webhook payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload data"})
		return
	}

	// Here, we're just returning the reference for the frontend to use
	c.JSON(http.StatusOK, gin.H{"reference": reference})
}
