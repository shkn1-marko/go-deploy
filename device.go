package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"
)

type deviceFile struct {
	InstallationID string `json:"installationID"`
}

type DeviceStore struct {
	mu   sync.RWMutex
	path string
	id   string
}

func NewDeviceStore(path string) (*DeviceStore, error) {
	s := &DeviceStore{path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var f deviceFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	s.id = f.InstallationID
	return s, nil
}

func (s *DeviceStore) Register(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("empty installation ID")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.writeAtomic(deviceFile{InstallationID: id}); err != nil {
		return err
	}
	s.id = id
	return nil
}

func (s *DeviceStore) Current() (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.id, s.id != ""
}

func (s *DeviceStore) writeAtomic(f deviceFile) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
