package volunteerops

import (
	"testing"
	"time"
)

func TestPromotionRechecksConflicts(t *testing.T) {
	c := New()
	start := time.Date(2025, 12, 1, 9, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	for _, id := range []string{"a", "b"} {
		if err := c.AddVolunteer(Volunteer{ID: id, AvailableFrom: start, AvailableTo: end}); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"x", "y"} {
		if err := c.AddShift(Shift{ID: id, Start: start, End: end, Capacity: 1}); err != nil {
			t.Fatal(err)
		}
	}
	c.Request("a", "x")
	c.Request("b", "x")
	c.Request("b", "y")
	if err := c.Cancel("a", "x"); err != nil {
		t.Fatal(err)
	}
	shift, _ := c.Shift("x")
	if len(shift.Assigned) != 0 || len(shift.Waitlist) != 1 {
		t.Fatal("promoted a conflicting volunteer")
	}
	if err := c.Cancel("b", "x"); err != nil {
		t.Fatal(err)
	}
}
func TestCancelShiftIsIdempotent(t *testing.T) {
	c := New()
	now := time.Now()
	c.AddVolunteer(Volunteer{ID: "a", AvailableFrom: now, AvailableTo: now.Add(time.Hour)})
	c.AddShift(Shift{ID: "x", Start: now, End: now.Add(time.Hour), Capacity: 1})
	c.Request("a", "x")
	c.CancelShift("x")
	c.CancelShift("x")
	if c.volunteers["a"].Assignments != 0 {
		t.Fatal("incorrect assignment count")
	}
}
