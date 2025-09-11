package repository

import (
	"bytes"
	"encoding/json"
	"integration/internal/model"
	"log"
	"net/http"
)

type ExternalSystem interface {
	SendPayment(p *model.Payment) error
}

type externalSystem struct {
	url string
}

func NewExternalSystem(url string) ExternalSystem {
	return &externalSystem{url: url}
}

func (s *externalSystem) SendPayment(p *model.Payment) error {
	body, _ := json.Marshal(p)
	resp, err := http.Post(s.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("❌ External system request failed:", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
