package services

import (
	"fmt"
	"net/http"
	"net/url"
	"paystack-payment/config"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ArkeselService struct {
	apiKey string
	client *resty.Client
}

func NewArkeselService(apiKey string) *ArkeselService {
	cfg, err := config.Load()

	if err != nil {
		return nil
	}
	return &ArkeselService{
		apiKey: cfg.ArkeselAPIKey,
		client: resty.New(),
	}
}

func (s *ArkeselService) SendSMS(phoneNumbers []string, message string) error {
	params := url.Values{}
	params.Add("action", "send-sms")
	params.Add("api_key", s.apiKey)
	params.Add("to", strings.Join(phoneNumbers, ","))
	params.Add("from", "TaxFlow")
	params.Add("sms", message)

	resp, err := s.client.R().
		SetQueryString(params.Encode()).
		Get("https://sms.arkesel.com/sms/api")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("SMS API returned non-OK status: %s", resp.Status())
	}

	return nil
}
