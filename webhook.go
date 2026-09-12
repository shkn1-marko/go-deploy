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

const webhookAddr = ":9091"

type webhookPayload struct {
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

func StartWebhookServer() {
	http.HandleFunc("/webhook", handleWebhook)
	log.Println("webhook server listening on", webhookAddr)

	if err := http.ListenAndServe(webhookAddr, nil); err != nil {
		log.Fatalln("webhook server failed:", err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !validSignature(r.Header.Get("X-Hub-Signature-256"), body) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var payload webhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Println("push event received:", payload.Repository.Name)
	go Deploy(payload.Repository.Name)
}

func validSignature(header string, body []byte) bool {
	const prefix = "sha256="
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return false
	}

	secret := os.Getenv("GDEP_WEBHOOK_SECRET")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(header[len(prefix):]), []byte(expected))
}
