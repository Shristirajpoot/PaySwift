package services

import (
	"paystack-payment/config"

	"github.com/go-resty/resty/v2"
)

type CourierService struct {
	apiKey string
	client *resty.Client
}

func NewCourierService(apiKey string) *CourierService {
	cfg, err := config.Load()

	if err != nil {
		return nil
	}

	return &CourierService{
		apiKey: cfg.CourierAPIKey,
		client: resty.New(),
	}
}

func (s *CourierService) SendEmail(email, subject, message string) error {
	_, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.apiKey).
		SetBody(map[string]interface{}{
			"message": map[string]interface{}{
				"to": map[string]string{
					"email": email,
				},
				"content": map[string]string{
					"title": subject,
					"body":  message,
				},
			},
		}).
		Post("https://api.courier.com/send")

	return err
}
