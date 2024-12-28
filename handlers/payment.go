package handlers

import (
	"net/http"
	"paystack-payment/models"
	"paystack-payment/services"
	"paystack-payment/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	db              *gorm.DB
	logger          *utils.Logger
	paystackService *services.PaystackService
}

func NewPaymentHandler(db *gorm.DB, logger *utils.Logger, paystackService *services.PaystackService) *PaymentHandler {
	return &PaymentHandler{
		db:              db,
		logger:          logger,
		paystackService: paystackService,
	}
}

func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	var req struct {
		BillID uint `json:"bill_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bill models.Bill
	if err := h.db.Preload("User").First(&bill, req.BillID).Error; err != nil {
		h.logger.Error("Bill not found", "error", err, "billID", req.BillID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	reference, authURL, err := h.paystackService.InitiateTransaction(bill.User.Email, bill.Amount)
	if err != nil {
		h.logger.Error("Failed to initiate Paystack transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate payment"})
		return
	}

	payment := models.Payment{
		BillID:    bill.ID,
		UserID:    bill.UserID,
		Amount:    bill.Amount,
		Reference: reference,
		Status:    "pending",
	}

	if err := h.db.Create(&payment).Error; err != nil {
		h.logger.Error("Failed to create payment record", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authorization_url": authURL, "reference": reference})
}

func (h *PaymentHandler) GetPaymentInfo(c *gin.Context) {
	id := c.Param("id")
	var payment models.Payment

	if err := h.db.Preload("User").Preload("Bill").First(&payment, id).Error; err != nil {
		h.logger.Error("Payment not found", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) ListPayments(c *gin.Context) {
	var payments []models.Payment
	if err := h.db.Preload("User").Preload("Bill").Find(&payments).Error; err != nil {
		h.logger.Error("Failed to list payments", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) VerifyTransaction(c *gin.Context) {
	reference := c.Query("reference")
	if reference == "" {
		h.logger.Error("Missing transaction reference")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction reference is required"})
		return
	}

	transaction, err := h.paystackService.VerifyTransaction(reference)
	if err != nil {
		h.logger.Error("Failed to verify transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify transaction"})
		return
	}

	var payment models.Payment
	if err := h.db.Where("reference = ?", reference).First(&payment).Error; err != nil {
		h.logger.Error("Payment not found", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if !transaction.Status {
		payment.Status = "failed"
	} else if transaction.Data.Status == "success" {
		payment.Status = "verified"
	}

	if err := h.db.Save(&payment).Error; err != nil {
		h.logger.Error("Failed to update payment status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": payment.Status})
}