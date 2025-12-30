package volunteerops

import (
	"errors"
	"sort"
	"sync"
	"time"
)

type Volunteer struct {
	ID                         string   `json:"id"`
	Skills                     []string `json:"skills"`
	AvailableFrom, AvailableTo time.Time
	Assignments                int `json:"assignments"`
}
type Shift struct {
	ID             string `json:"id"`
	Start, End     time.Time
	Capacity       int      `json:"capacity"`
	RequiredSkills []string `json:"requiredSkills"`
	Assigned       []string `json:"assigned"`
	Waitlist       []string `json:"waitlist"`
	Cancelled      bool     `json:"cancelled"`
}
type Coordinator struct {
	mu         sync.Mutex
	volunteers map[string]*Volunteer
	shifts     map[string]*Shift
}

func New() *Coordinator {
	return &Coordinator{volunteers: map[string]*Volunteer{}, shifts: map[string]*Shift{}}
}
func (c *Coordinator) AddVolunteer(v Volunteer) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v.ID == "" || !v.AvailableTo.After(v.AvailableFrom) {
		return errors.New("invalid volunteer")
	}
	if _, ok := c.volunteers[v.ID]; ok {
		return errors.New("duplicate volunteer")
	}
	x := v
	x.Skills = append([]string(nil), v.Skills...)
	c.volunteers[v.ID] = &x
	return nil
}
func (c *Coordinator) AddShift(s Shift) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s.ID == "" || s.Capacity < 1 || !s.End.After(s.Start) {
		return errors.New("invalid shift")
	}
	if _, ok := c.shifts[s.ID]; ok {
		return errors.New("duplicate shift")
	}
	if len(s.Assigned) > 0 || len(s.Waitlist) > 0 {
		return errors.New("new shift cannot contain assignments")
	}
	x := s
	x.RequiredSkills = append([]string(nil), s.RequiredSkills...)
	c.shifts[s.ID] = &x
	return nil
}
func contains(all []string, want string) bool {
	for _, x := range all {
		if x == want {
			return true
		}
	}
	return false
}
func (c *Coordinator) eligible(v *Volunteer, s *Shift) error {
	if s.Cancelled {
		return errors.New("shift cancelled")
	}
	if v.AvailableFrom.After(s.Start) || v.AvailableTo.Before(s.End) {
		return errors.New("outside availability")
	}
	for _, skill := range s.RequiredSkills {
		if !contains(v.Skills, skill) {
			return errors.New("missing required skill")
		}
	}
	for _, other := range c.shifts {
		if other.ID != s.ID && contains(other.Assigned, v.ID) && s.Start.Before(other.End) && other.Start.Before(s.End) {
			return errors.New("shift conflict")
		}
	}
	return nil
}
func (c *Coordinator) Request(volunteerID, shiftID string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, s := c.volunteers[volunteerID], c.shifts[shiftID]
	if v == nil || s == nil {
		return "", errors.New("unknown volunteer or shift")
	}
	if contains(s.Assigned, volunteerID) || contains(s.Waitlist, volunteerID) {
		return "", errors.New("already requested")
	}
	if err := c.eligible(v, s); err != nil {
		return "", err
	}
	if len(s.Assigned) < s.Capacity {
		s.Assigned = append(s.Assigned, volunteerID)
		v.Assignments++
		return "assigned", nil
	}
	s.Waitlist = append(s.Waitlist, volunteerID)
	sort.SliceStable(s.Waitlist, func(i, j int) bool {
		return c.volunteers[s.Waitlist[i]].Assignments < c.volunteers[s.Waitlist[j]].Assignments
	})
	return "waitlisted", nil
}

func (c *Coordinator) Shift(id string) (Shift, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := c.shifts[id]
	if s == nil {
		return Shift{}, false
	}
	copy := *s
	copy.RequiredSkills = append([]string(nil), s.RequiredSkills...)
	copy.Assigned = append([]string{}, s.Assigned...)
	copy.Waitlist = append([]string{}, s.Waitlist...)
	return copy, true
}
