package main

import (
	"context"
	"log"
	"strconv"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type StageStatus string

const (
	StageOK      StageStatus = "OK"
	StageFailed  StageStatus = "FAILED"
	StageSkipped StageStatus = "SKIPPED"
)

type DeployStatus struct {
	Name      string      `json:"name"`
	Build     StageStatus `json:"buildStatus"`
	Deploy    StageStatus `json:"deployStatus"`
	Cause     string      `json:"cause,omitempty"`
	Output    string      `json:"output,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

type FCMSender struct {
	client *messaging.Client
}

func NewFCMSender(path string) (*FCMSender, error) {
	app, err := firebase.NewApp(context.Background(), nil,
		option.WithAuthCredentialsFile(option.ServiceAccount, path))
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}
	return &FCMSender{client: client}, nil
}

func (f *FCMSender) SendDeployStatus(device string, status DeployStatus) {
	_, err := f.client.Send(context.Background(), &messaging.Message{
		Token: device,
		Data: map[string]string{
			"name":         status.Name,
			"buildStatus":  string(status.Build),
			"deployStatus": string(status.Deploy),
			"cause":        status.Cause,
			"output":       status.Output,
			"timestamp":    strconv.FormatInt(status.Timestamp, 10),
		},
	})
	if err != nil {
		log.Println("--", err)
	}
}
