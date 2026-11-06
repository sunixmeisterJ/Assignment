package domain

import (
	"fmt"
	"time"
)

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Details     string `json:"details"`
	HourlyCents int64  `json:"hourlyCents"`
}

// PriceRule is the extension point for alternative tariffs.
type PriceRule interface {
	Price(time.Time, time.Time) (int64, error)
}
type HourlyPrice struct{ Cents int64 }

func (p HourlyPrice) Price(start, end time.Time) (int64, error) {
	if p.Cents <= 0 || !end.After(start) {
		return 0, fmt.Errorf("Preis und Zeitraum müssen positiv sein")
	}
	if end.Sub(start) > 366*24*time.Hour {
		return 0, fmt.Errorf("Zeitraum darf maximal 366 Tage betragen")
	}
	hours := int64((end.Sub(start) + time.Hour - 1) / time.Hour)
	return hours * p.Cents, nil
}

type Person struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DocumentType   string `json:"documentType"`
	DocumentNumber string `json:"documentNumber"`
}
type Rental struct {
	ID          int        `json:"id"`
	ItemID      int        `json:"itemId"`
	PersonID    int        `json:"personId"`
	Start       time.Time  `json:"start"`
	End         time.Time  `json:"end"`
	Status      string     `json:"status"`
	HourlyCents int64      `json:"hourlyCents"`
	PriceCents  int64      `json:"priceCents"`
	ReturnedAt  *time.Time `json:"returnedAt,omitempty"`
}

func (r Rental) Overlaps(start, end time.Time) bool {
	return (r.Status == "reserved" || r.Status == "active") && start.Before(r.End) && r.Start.Before(end)
}
