package handlers

import (
	"net/http"
	"paystack-payment/models"
	"paystack-payment/services"
	"paystack-payment/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BillHandler struct {
	db              *gorm.DB
	logger          *utils.Logger
	courierService  *services.CourierService
	arkeselService  *services.ArkeselService
	paystackService *services.PaystackService
}

func NewBillHandler(db *gorm.DB, logger *utils.Logger, courierService *services.CourierService, arkeselService *services.ArkeselService, paystackService *services.PaystackService) *BillHandler {
	return &BillHandler{
		db:              db,
		logger:          logger,
		courierService:  courierService,
		arkeselService:  arkeselService,
		paystackService: paystackService,
	}
}

func (h *BillHandler) CreateBill(c *gin.Context) {
	var bill models.Bill
	if err := c.ShouldBindJSON(&bill); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&bill).Error; err != nil {
		h.logger.Error("Failed to create bill", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bill"})
		return
	}

	c.JSON(http.StatusCreated, bill)
}

func (h *BillHandler) GetBill(c *gin.Context) {
	id := c.Param("id")
	var bill models.Bill

	if err := h.db.First(&bill, id).Error; err != nil {
		h.logger.Error("Failed to get bill", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

func (h *BillHandler) ListBills(c *gin.Context) {
	var bills []models.Bill
	if err := h.db.Find(&bills).Error; err != nil {
		h.logger.Error("Failed to list bills", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bills"})
		return
	}

	c.JSON(http.StatusOK, bills)
}

func (h *BillHandler) UpdateBill(c *gin.Context) {
	id := c.Param("id")
	var bill models.Bill

	if err := h.db.First(&bill, id).Error; err != nil {
		h.logger.Error("Bill not found", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	if err := c.ShouldBindJSON(&bill); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&bill).Error; err != nil {
		h.logger.Error("Failed to update bill", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bill"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

func (h *BillHandler) DeleteBill(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.Bill{}, id).Error; err != nil {
		h.logger.Error("Failed to delete bill", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill deleted successfully"})
}

func (h *BillHandler) SendBill(c *gin.Context) {
	id := c.Param("id")
	var bill models.Bill

	if err := h.db.Preload("User").First(&bill, id).Error; err != nil {
		h.logger.Error("Bill not found", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	reference, authURL, err := h.paystackService.InitiateTransaction(bill.User.Email, bill.Amount)
	if err != nil {
		h.logger.Error("Failed to generate payment link", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate payment link"})
		return
	}

	message := utils.FormatBillNotification(bill, authURL)

	// Send email and SMS concurrently
	emailChan := make(chan error)
	smsChan := make(chan error)

	go func() {
		emailChan <- h.courierService.SendEmail(bill.User.Email, "New Bill Notification", message)
	}()

	go func() {
		smsChan <- h.arkeselService.SendSMS([]string{bill.User.Phone}, message)
	}()

	// Wait for both email and SMS to be sent
	emailErr := <-emailChan
	smsErr := <-smsChan

	if emailErr != nil {
		h.logger.Error("Failed to send email notification", "error", emailErr)
	}

	if smsErr != nil {
		h.logger.Error("Failed to send SMS notification", "error", smsErr)
	}

	if emailErr != nil || smsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notifications"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Bill sent successfully", "payment_link": authURL})
}