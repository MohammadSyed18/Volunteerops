package volunteerops

import "errors"

func (c *Coordinator) promote(s *Shift) {
	remaining := make([]string, 0, len(s.Waitlist))
	for _, id := range s.Waitlist {
		volunteer := c.volunteers[id]
		if volunteer == nil {
			continue
		}
		if len(s.Assigned) < s.Capacity && c.eligible(volunteer, s) == nil {
			s.Assigned = append(s.Assigned, id)
			volunteer.Assignments++
		} else {
			remaining = append(remaining, id)
		}
	}
	s.Waitlist = remaining
}

func (c *Coordinator) Cancel(volunteerID, shiftID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	shift := c.shifts[shiftID]
	if shift == nil {
		return errors.New("unknown shift")
	}
	for i, id := range shift.Waitlist {
		if id == volunteerID {
			shift.Waitlist = append(shift.Waitlist[:i], shift.Waitlist[i+1:]...)
			return nil
		}
	}
	for i, id := range shift.Assigned {
		if id == volunteerID {
			shift.Assigned = append(shift.Assigned[:i], shift.Assigned[i+1:]...)
			c.volunteers[id].Assignments--
			c.promote(shift)
			return nil
		}
	}
	return errors.New("assignment not found")
}

func (c *Coordinator) CancelShift(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	shift := c.shifts[id]
	if shift == nil {
		return errors.New("unknown shift")
	}
	if shift.Cancelled {
		return nil
	}
	for _, volunteer := range shift.Assigned {
		c.volunteers[volunteer].Assignments--
	}
	shift.Assigned = nil
	shift.Waitlist = nil
	shift.Cancelled = true
	return nil
}
