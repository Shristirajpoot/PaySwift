package utils

import (
	"fmt"
	"paystack-payment/config"
	"paystack-payment/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dbURL string) (*gorm.DB, error) {
	cfg, err := config.Load()

	if err != nil {
		return nil, fmt.Errorf("could not load config items: %v", err)
	}

	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.PG_USER, cfg.PG_PASS, cfg.PG_HOST, cfg.PG_PORT, cfg.PG_NAME)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Drop existing tables
	err = db.Migrator().DropTable(&models.User{}, &models.Bill{}, &models.Payment{})
	if err != nil {
		return nil, fmt.Errorf("failed to drop existing tables: %v", err)
	}

	// Create tables
	err = db.AutoMigrate(&models.User{}, &models.Bill{}, &models.Payment{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate schema: %v", err)
	}

	// Add unique constraint to Payment.Reference
	err = db.Exec("ALTER TABLE payments ADD CONSTRAINT uni_payment_reference UNIQUE (reference)").Error
	if err != nil {
		return nil, fmt.Errorf("failed to add unique constraint to Payment.Reference: %v", err)
	}
	return db, nil
}
