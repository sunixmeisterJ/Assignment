package domain

import (
	"testing"
	"time"
)

func TestHourlyPrice(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		duration time.Duration
		want     int64
	}{{time.Second, 1200}, {time.Hour, 1200}, {time.Hour + time.Second, 2400}} {
		got, err := (HourlyPrice{1200}).Price(start, start.Add(tc.duration))
		if err != nil || got != tc.want {
			t.Fatalf("duration %s: got %d, %v; want %d", tc.duration, got, err, tc.want)
		}
	}
	if _, err := (HourlyPrice{1200}).Price(start, start); err == nil {
		t.Fatal("empty period accepted")
	}
}
func TestAdjacentReservations(t *testing.T) {
	start := time.Now()
	r := Rental{Start: start, End: start.Add(time.Hour), Status: "reserved"}
	if r.Overlaps(r.End, r.End.Add(time.Hour)) {
		t.Fatal("adjacent bookings overlap")
	}
	if !r.Overlaps(start.Add(time.Minute), r.End) {
		t.Fatal("overlap missed")
	}
}
