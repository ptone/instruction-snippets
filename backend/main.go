package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/vertexai/genai"
)

const (
	projectID = "new-test-297222"
	location  = "us-central1"
)

// Snippet represents the data structure for a snippet in Firestore.

type Snippet struct {
	Text      string   `firestore:"text"`
	Labels    []string `firestore:"labels"`
	CreatedAt int64    `firestore:"created_at"`
	ThumbsUp  int      `firestore:"thumbs_up"`
	ThumbsDown int      `firestore:"thumbs_down"`
}

func main() {
	http.HandleFunc("/ingest", ingestHandler)

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ingestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	markdownContent := string(body)

	snippets, err := processMarkdown(markdownContent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing markdown: %v", err), http.StatusInternalServerError)
		return
	}

	if err := saveSnippets(snippets); err != nil {
		http.Error(w, fmt.Sprintf("Error saving snippets: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(snippets)
}

func processMarkdown(content string) ([]Snippet, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("error creating genai client: %v", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel("gemini-pro")

	// First, get the snippets
	prompt := fmt.Sprintf("Extract the instruction snippets from the following markdown document. Each snippet should be a standalone instruction. Return the snippets as a JSON array of strings. Markdown: %s", content)
	resp, err := gemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating content: %v", err)
	}

	// Assuming the response is a JSON array of strings
	var snippetTexts []string
	if err := json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &snippetTexts); err != nil {
		return nil, fmt.Errorf("error unmarshalling snippets: %v", err)
	}

	var snippets []Snippet
	for _, text := range snippetTexts {
		// Then, get the labels for each snippet
		labelPrompt := fmt.Sprintf("Generate a list of relevant labels for the following instruction snippet. Return the labels as a JSON array of strings. Snippet: %s", text)
		labelResp, err := gemini.GenerateContent(ctx, genai.Text(labelPrompt))
		if err != nil {
			return nil, fmt.Errorf("error generating labels: %v", err)
		}

		var labels []string
		if err := json.Unmarshal([]byte(fmt.Sprintf("%s", labelResp.Candidates[0].Content.Parts[0])), &labels); err != nil {
			return nil, fmt.Errorf("error unmarshalling labels: %v", err)
		}

		snippets = append(snippets, Snippet{
			Text:   text,
			Labels: labels,
		})
	}

	return snippets, nil
}

func saveSnippets(snippets []Snippet) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("error creating firestore client: %v", err)
	}
	defer client.Close()

	for _, snippet := range snippets {
		_, _, err := client.Collection("snippets").Add(ctx, snippet)
		if err != nil {
			return fmt.Errorf("error adding snippet: %v", err)
		}
	}

	return nil
}