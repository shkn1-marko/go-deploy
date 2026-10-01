package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type server struct {
	store *DeviceStore
}

func Start() {
	store, err := NewDeviceStore(os.Getenv("GDEP_DEVICE_FILE"))
	if err != nil {
		log.Fatalln("--", err)
	}

	s := &server{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.handleWebhook)
	mux.HandleFunc("/register-device", s.handleRegisterDevice)

	log.Println("--", ":9091")
	if err := http.ListenAndServe(":9091", mux); err != nil {
		log.Fatalln("--")
	}
}

func (s *server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !validSignature(
		os.Getenv("GDEP_WEBHOOK_SECRET"),
		r.Header.Get("X-Hub-Signature-256"),
		body,
	) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("X-GitHub-Event") != "push" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var payload webhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	device, _ := s.store.Current()
	go RunDeploy(device, payload.Repository.Name)
}

func (s *server) handleRegisterDevice(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !validSignature(
		os.Getenv("GDEP_REGISTER_SECRET"),
		r.Header.Get("X-GDEP-Signature-256"),
		body,
	) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var payload registerDevicePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	if err := s.store.Register(payload.InstallationID); err != nil {
		http.Error(w, "invalid installation ID", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("--")
}

func validSignature(secret, header string, body []byte) bool {
	const prefix = "sha256="
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(header[len(prefix):]), []byte(expected))
}
