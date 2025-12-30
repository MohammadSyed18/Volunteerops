package volunteerops

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func (c *Coordinator) Save(path string) error {
	raw, err := json.MarshalIndent(c.Snapshot(), "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".volunteerops-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err = file.Write(raw); err != nil {
		return err
	}
	if err = file.Sync(); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

func Load(path string) (*Coordinator, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Bound accidental or malicious snapshot imports before decoding.
	raw, err := io.ReadAll(io.LimitReader(file, 16*1024*1024+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > 16*1024*1024 {
		return nil, errors.New("snapshot exceeds 16 MiB")
	}
	var snapshot Snapshot
	if err = json.Unmarshal(raw, &snapshot); err != nil {
		return nil, err
	}
	return Restore(snapshot)
}
