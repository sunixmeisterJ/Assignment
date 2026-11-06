package store

import (
	"rental/internal/domain"
	"sync"
	"testing"
	"time"
)

func TestConcurrentReservationAndLifecycle(t *testing.T) {
	s := New()
	p, err := s.AddPerson(domain.Person{Name: "Test", DocumentType: "Pass", DocumentNumber: "DEMO"})
	if err != nil {
		t.Fatal(err)
	}
	r := domain.Rental{ItemID: 1, PersonID: p.ID, Start: time.Now().Add(-time.Second), End: time.Now().Add(time.Hour)}
	var wg sync.WaitGroup
	var mu sync.Mutex
	success := 0
	var booked domain.Rental
	for n := 0; n < 20; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, e := s.Reserve(r)
			if e == nil {
				mu.Lock()
				success++
				booked = v
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if success != 1 {
		t.Fatalf("expected one successful booking, got %d", success)
	}
	item := s.Snapshot().Items[0]
	original := item.HourlyCents
	item.HourlyCents *= 2
	if _, err := s.SaveItem(item); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Transition(booked.ID, "checkout"); err != nil {
		t.Fatal(err)
	}
	returned, err := s.Transition(booked.ID, "return")
	if err != nil {
		t.Fatal(err)
	}
	if returned.PriceCents != original || returned.ReturnedAt == nil {
		t.Fatalf("tariff snapshot not retained: %+v", returned)
	}
	if _, err := s.Transition(booked.ID, "return"); err == nil {
		t.Fatal("double return accepted")
	}
	if _, err := s.Reserve(r); err != nil {
		t.Fatalf("returned item blocked: %v", err)
	}
}
func TestValidationAndCancellation(t *testing.T) {
	s := New()
	if _, err := s.AddPerson(domain.Person{}); err == nil {
		t.Fatal("empty person accepted")
	}
	p, _ := s.AddPerson(domain.Person{Name: "Test", DocumentType: "Pass", DocumentNumber: "X"})
	if _, err := s.AddPerson(p); err == nil {
		t.Fatal("duplicate document accepted")
	}
	r := domain.Rental{ItemID: 1, PersonID: p.ID, Start: time.Now().Add(time.Hour), End: time.Now().Add(2 * time.Hour)}
	b, err := s.Reserve(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Transition(b.ID, "checkout"); err == nil {
		t.Fatal("early checkout accepted")
	}
	if _, err := s.Transition(b.ID, "cancel"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Reserve(r); err != nil {
		t.Fatal(err)
	}
}
