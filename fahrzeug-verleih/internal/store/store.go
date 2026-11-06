package store

import (
	"fmt"
	"rental/internal/domain"
	"strings"
	"sync"
	"time"
)

type Store struct {
	mu      sync.Mutex
	items   []domain.Item
	people  []domain.Person
	rentals []domain.Rental
}
type Snapshot struct {
	Items   []domain.Item   `json:"items"`
	People  []domain.Person `json:"people"`
	Rentals []domain.Rental `json:"rentals"`
}

func New() *Store {
	s := &Store{items: []domain.Item{}, people: []domain.Person{}, rentals: []domain.Rental{}}
	s.addItem(domain.Item{Name: "VW Golf", Details: "Auto, W-12345", HourlyCents: 1800})
	s.addItem(domain.Item{Name: "Cube Touring", Details: "Fahrrad, Rahmen M", HourlyCents: 500})
	s.addItem(domain.Item{Name: "Ford Transit", Details: "Transporter, W-98765", HourlyCents: 2500})
	return s
}
func (s *Store) Snapshot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Snapshot{append([]domain.Item{}, s.items...), append([]domain.Person{}, s.people...), append([]domain.Rental{}, s.rentals...)}
}
func (s *Store) addItem(i domain.Item) domain.Item {
	i.ID = len(s.items) + 1
	s.items = append(s.items, i)
	return i
}
func (s *Store) SaveItem(i domain.Item) (domain.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i.Name = strings.TrimSpace(i.Name)
	i.Details = strings.TrimSpace(i.Details)
	if i.Name == "" || len(i.Name) > 120 || len(i.Details) > 300 || i.HourlyCents <= 0 || i.HourlyCents > 1000000 {
		return i, fmt.Errorf("Name (max. 120 Zeichen), Details (max. 300) und Preis (0,01–10.000 €) prüfen")
	}
	if i.ID == 0 {
		return s.addItem(i), nil
	}
	for n, old := range s.items {
		if old.ID == i.ID {
			s.items[n] = i
			return i, nil
		}
	}
	return i, fmt.Errorf("Artikel nicht gefunden")
}
func (s *Store) AddPerson(p domain.Person) (domain.Person, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.Name = strings.TrimSpace(p.Name)
	p.DocumentType = strings.TrimSpace(p.DocumentType)
	p.DocumentNumber = strings.TrimSpace(p.DocumentNumber)
	if p.Name == "" || p.DocumentType == "" || p.DocumentNumber == "" || len(p.Name) > 120 || len(p.DocumentType) > 80 || len(p.DocumentNumber) > 80 {
		return p, fmt.Errorf("Name, Ausweistyp und Ausweisnummer sind erforderlich (max. 120/80/80 Zeichen)")
	}
	for _, old := range s.people {
		if strings.EqualFold(old.DocumentType, p.DocumentType) && strings.EqualFold(old.DocumentNumber, p.DocumentNumber) {
			return p, fmt.Errorf("Ausweis bereits registriert")
		}
	}
	p.ID = len(s.people) + 1
	s.people = append(s.people, p)
	return p, nil
}
func (s *Store) Reserve(r domain.Rental) (domain.Rental, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.Start.Before(time.Now().Add(-time.Minute)) {
		return r, fmt.Errorf("Beginn darf nicht in der Vergangenheit liegen")
	}
	var item *domain.Item
	for n := range s.items {
		if s.items[n].ID == r.ItemID {
			item = &s.items[n]
		}
	}
	if item == nil {
		return r, fmt.Errorf("Artikel nicht gefunden")
	}
	found := false
	for _, p := range s.people {
		if p.ID == r.PersonID {
			found = true
		}
	}
	if !found {
		return r, fmt.Errorf("Person nicht gefunden")
	}
	if r.End.Sub(r.Start) > 366*24*time.Hour {
		return r, fmt.Errorf("Zeitraum darf maximal 366 Tage betragen")
	}
	price, err := (domain.HourlyPrice{Cents: item.HourlyCents}).Price(r.Start, r.End)
	if err != nil {
		return r, fmt.Errorf("Reservierung: %w", err)
	}
	for _, old := range s.rentals {
		if old.ItemID == r.ItemID && (old.Overlaps(r.Start, r.End) || old.Status == "active") {
			return r, fmt.Errorf("Artikel ist im Zeitraum reserviert oder noch ausgeliehen")
		}
	}
	r.ID = len(s.rentals) + 1
	r.Status = "reserved"
	r.HourlyCents = item.HourlyCents
	r.PriceCents = price
	r.ReturnedAt = nil
	s.rentals = append(s.rentals, r)
	return r, nil
}
func (s *Store) Transition(id int, action string) (domain.Rental, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for n := range s.rentals {
		r := &s.rentals[n]
		if r.ID != id {
			continue
		}
		switch action {
		case "checkout":
			if r.Status != "reserved" || now.Before(r.Start) || !now.Before(r.End) {
				return *r, fmt.Errorf("Ausgabe nur innerhalb des reservierten Zeitraums möglich")
			}
			for _, other := range s.rentals {
				if other.ItemID == r.ItemID && other.Status == "active" {
					return *r, fmt.Errorf("Artikel noch ausgeliehen")
				}
			}
			r.Status = "active"
		case "return":
			if r.Status != "active" {
				return *r, fmt.Errorf("Nur ausgeliehene Artikel können zurückgegeben werden")
			}
			price, err := (domain.HourlyPrice{Cents: r.HourlyCents}).Price(r.Start, now)
			if err != nil {
				return *r, fmt.Errorf("Rückgabe: %w", err)
			}

			r.PriceCents = price
			r.ReturnedAt = &now
			r.Status = "returned"
		case "cancel":
			if r.Status != "reserved" {
				return *r, fmt.Errorf("Nur Reservierungen können storniert werden")
			}
			r.Status = "cancelled"
		default:
			return *r, fmt.Errorf("Unbekannte Aktion")
		}
		return *r, nil
	}
	return domain.Rental{}, fmt.Errorf("Reservierung nicht gefunden")
}
