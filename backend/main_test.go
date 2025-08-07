package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/genai"
)

func TestReadSampleFile(t *testing.T) {
	content, err := ioutil.ReadFile("../samples/GEMINI-brief.md")
	if err != nil {
		t.Fatalf("Failed to read sample file: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("Sample file is empty")
	}
}

func TestProcessHandler_Integration(t *testing.T) {

	// This test makes real calls to Google Cloud Vertex AI APIs.
	// Ensure you have authenticated with `gcloud auth application-default login`.
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		// Check for ADC file existence as a proxy for being logged in.
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("could not get user home directory: %v", err)
		}
		adcFile := home + "/.config/gcloud/application_default_credentials.json"
		if _, err := os.Stat(adcFile); os.IsNotExist(err) {
			t.Skip("Skipping integration test: Application Default Credentials not found. Run 'gcloud auth application-default login'.")
		}
	}

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, "new-test-297222")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer firestoreClient.Close()

	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  "new-test-297222",
		Location: "global",
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		t.Fatalf("Failed to create genai client: %v", err)
	}

	app := &App{
		firestoreClient: firestoreClient,
		genaiClient:     genaiClient,
	}

	content, err := ioutil.ReadFile("../samples/GEMINI-brief.md")
	if err != nil {
		t.Fatalf("Failed to read sample file: %v", err)
	}

	reqBody := ProcessRequest{
		Content: string(content),
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/process", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.processHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
	}

	// TODO: Add verification that the data was written to firestore.
}
