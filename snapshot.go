package volunteerops

import (
	"errors"
	"fmt"
	"sort"
)

// Snapshot is a complete, versioned coordinator state. Assignment counters are
// checked against actual shifts when restoring, rather than trusted as input.
type Snapshot struct {
	Version    int         `json:"version"`
	Volunteers []Volunteer `json:"volunteers"`
	Shifts     []Shift     `json:"shifts"`
}

func (c *Coordinator) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := Snapshot{Version: 1, Volunteers: []Volunteer{}, Shifts: []Shift{}}
	for _, original := range c.volunteers {
		v := *original
		v.Skills = append([]string{}, original.Skills...)
		out.Volunteers = append(out.Volunteers, v)
	}
	for _, original := range c.shifts {
		s := *original
		s.RequiredSkills = append([]string{}, original.RequiredSkills...)
		s.Assigned = append([]string{}, original.Assigned...)
		s.Waitlist = append([]string{}, original.Waitlist...)
		out.Shifts = append(out.Shifts, s)
	}
	sort.Slice(out.Volunteers, func(i, j int) bool { return out.Volunteers[i].ID < out.Volunteers[j].ID })
	sort.Slice(out.Shifts, func(i, j int) bool { return out.Shifts[i].ID < out.Shifts[j].ID })
	return out
}

func Restore(snapshot Snapshot) (*Coordinator, error) {
	if snapshot.Version != 1 {
		return nil, errors.New("unsupported snapshot version")
	}
	candidate := New()
	expected := map[string]int{}
	for _, v := range snapshot.Volunteers {
		expected[v.ID] = v.Assignments
		v.Assignments = 0
		if err := candidate.AddVolunteer(v); err != nil {
			return nil, fmt.Errorf("volunteer %s: %w", v.ID, err)
		}
	}
	for _, original := range snapshot.Shifts {
		s := original
		s.Assigned = nil
		s.Waitlist = nil
		if original.Cancelled && (len(original.Assigned) > 0 || len(original.Waitlist) > 0) {
			return nil, errors.New("cancelled shift contains requests")
		}
		if len(original.Assigned) > original.Capacity {
			return nil, errors.New("shift exceeds capacity")
		}
		if err := candidate.AddShift(s); err != nil {
			return nil, err
		}
	}
	for _, s := range snapshot.Shifts {
		for _, id := range s.Assigned {
			state, err := candidate.Request(id, s.ID)
			if err != nil {
				return nil, fmt.Errorf("invalid assignment %s/%s: %w", id, s.ID, err)
			}
			if state != "assigned" {
				return nil, errors.New("assignment would exceed capacity")
			}
		}
		seen := map[string]bool{}
		for _, id := range s.Waitlist {
			if candidate.volunteers[id] == nil || seen[id] || contains(s.Assigned, id) {
				return nil, errors.New("invalid waitlist entry")
			}
			seen[id] = true
		}
		// An existing waitlist entry may currently conflict with another shift.
		// Promotion will recheck eligibility at the moment a place opens.
		candidate.shifts[s.ID].Waitlist = append([]string{}, s.Waitlist...)
	}
	for id, count := range expected {
		if count != candidate.volunteers[id].Assignments {
			return nil, errors.New("assignment counter mismatch")
		}
	}
	return candidate, nil
}
