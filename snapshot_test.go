package volunteerops

import (
	"path/filepath"
	"sync"
	"testing"
)

func TestSnapshotRoundTripAndIsolation(t *testing.T) {
	c := fixture()
	c.Request("a", "s")
	c.Request("b", "s")
	original := c.Snapshot()
	path := filepath.Join(t.TempDir(), "state.json")
	if err := c.Save(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = loaded.Cancel("a", "s"); err != nil {
		t.Fatal(err)
	}
	shift, _ := loaded.Shift("s")
	if shift.Assigned[0] != "b" {
		t.Fatal("promotion was not restored")
	}
	original.Shifts[0].Assigned[0] = "changed"
	unchanged, _ := c.Shift("s")
	if unchanged.Assigned[0] != "a" {
		t.Fatal("snapshot aliased coordinator")
	}
}
func TestMalformedSnapshotDoesNotRestore(t *testing.T) {
	c := fixture()
	c.Request("a", "s")
	state := c.Snapshot()
	state.Volunteers[0].Assignments = 99
	if _, err := Restore(state); err == nil {
		t.Fatal("accepted counter mismatch")
	}
	state = c.Snapshot()
	state.Shifts[0].Waitlist = []string{"unknown"}
	if _, err := Restore(state); err == nil {
		t.Fatal("accepted missing waitlist volunteer")
	}
	state = c.Snapshot()
	state.Shifts[0].Assigned = []string{"a", "a"}
	if _, err := Restore(state); err == nil {
		t.Fatal("accepted excess assignment")
	}
}
func TestSnapshotsAreConsistentDuringRequests(t *testing.T) {
	c := fixture()
	var group sync.WaitGroup
	for _, id := range []string{"a", "b", "c"} {
		group.Add(1)
		go func(id string) { defer group.Done(); c.Request(id, "s") }(id)
	}
	for i := 0; i < 20; i++ {
		if _, err := Restore(c.Snapshot()); err != nil {
			t.Fatal(err)
		}
	}
	group.Wait()
}
