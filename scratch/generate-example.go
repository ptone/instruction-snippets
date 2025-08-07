package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/genai"
)

// Your GCP project
const project = "your-project"

// A GCP location like "us-central1"
const location = "some-gcp-location"

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  project,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Call the GenerateContent method.
	result, err := client.Models.GenerateContent(ctx,
		"gemini-2.0-flash",
		genai.Text("What is the sum of the first 50 prime numbers? Generate and run code for the calculation, and make sure you get all 50."),
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{CodeExecution: &genai.ToolCodeExecution{}},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	debugPrint(result)
}

func debugPrint[T any](r *T) {

	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(response))
}