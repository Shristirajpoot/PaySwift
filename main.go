package main

import (
	"fmt"
	"log"
	"paystack-payment/config"
	"paystack-payment/handlers"
	"paystack-payment/services"
	"paystack-payment/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := utils.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	logger := utils.NewLogger()

	paystackService := services.NewPaystackService(cfg.PaystackSecretKey)
	courierService := services.NewCourierService(cfg.CourierAPIKey)
	arkeselService := services.NewArkeselService(cfg.ArkeselAPIKey)

	r := gin.Default()

	userHandler := handlers.NewUserHandler(db, logger)
	billHandler := handlers.NewBillHandler(db, logger, courierService, arkeselService, paystackService)
	paymentHandler := handlers.NewPaymentHandler(db, logger, paystackService)
	webhookHandler := handlers.NewWebhookHandler(db, logger, paystackService)

	// Public routes
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.PUT("/users/:id", userHandler.UpdateUser)

	r.POST("/bills", billHandler.CreateBill)
	r.GET("/bills/:id", billHandler.GetBill)
	r.GET("/bills", billHandler.ListBills)
	r.PUT("/bills/:id", billHandler.UpdateBill)
	r.DELETE("/bills/:id", billHandler.DeleteBill)
	r.POST("/bills/:id/send", billHandler.SendBill)

	r.POST("/payments", paymentHandler.InitiatePayment)
	r.GET("/payments/:id", paymentHandler.GetPaymentInfo)
	r.GET("/payments", paymentHandler.ListPayments)
	r.GET("/payments/verify", paymentHandler.VerifyTransaction)

	r.POST("/webhook", webhookHandler.HandleWebhook)

	r.Run(fmt.Sprintf(":%s", cfg.ServerPort))
}
