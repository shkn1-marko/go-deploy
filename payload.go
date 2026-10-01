package main

type webhookPayload struct {
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}
