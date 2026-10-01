package main

type webhookPayload struct {
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

type registerDevicePayload struct {
	InstallationID string `json:"installationId"`
}
