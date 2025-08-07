package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/genai"
)

// App holds application dependencies
type App struct {
	firestoreClient *firestore.Client
	genaiClient     *genai.Client
}

// ProcessRequest defines the structure for the incoming request
type ProcessRequest struct {
	Content string `json:"content"`
}

// Source defines the structure for the sources collection
type Source struct {
	Content       string    `firestore:"content"`
	LastRefreshed time.Time `firestore:"last_refreshed"`
	Type          string    `firestore:"type"`
}

// Snippet defines the structure for the snippets collection
type Snippet struct {
	Content    string                 `firestore:"content"`
	Labels     []string               `firestore:"labels"`
	Source     *firestore.DocumentRef `firestore:"source"`
	ThumbsUp   int                    `firestore:"thumbs_up"`
	ThumbsDown int                    `firestore:"thumbs_down"`
	CreatedAt  time.Time              `firestore:"created_at"`
	Embedding  []float32              `firestore:"embedding"`
}

func main() {
	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, "new-test-297222")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer firestoreClient.Close()

	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  "new-test-297222",
		Location: "us-central1",
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}

	app := &App{
		firestoreClient: firestoreClient,
		genaiClient:     genaiClient,
	}

	http.HandleFunc("/process", app.processHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (app *App) processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	source := Source{
		Content:       req.Content,
		LastRefreshed: time.Now(),
		Type:          "file", // Assuming file upload for now
	}

	ref, _, err := app.firestoreClient.Collection("sources").Add(ctx, source)
	if err != nil {
		http.Error(w, "Failed to store source", http.StatusInternalServerError)
		log.Printf("Failed to add source: %v", err)
		return
	}

	log.Printf("Stored source with ID: %s", ref.ID)

	// Generate snippets from the markdown content
	snippets, err := app.generateSnippets(ctx, req.Content)
	if err != nil {
		http.Error(w, "Failed to generate snippets", http.StatusInternalServerError)
		log.Printf("Failed to generate snippets: %v", err)
		return
	}

	// Process and store snippets
	for _, snippetText := range snippets {
		log.Printf("Generated snippet: %s", snippetText)

		labels, err := app.generateLabels(ctx, snippetText)
		if err != nil {
			log.Printf("Failed to generate labels for snippet: %v", err)
			continue
		}
		log.Printf("Generated labels: %v", labels)

		embedding, err := app.generateEmbedding(ctx, snippetText)
		if err != nil {
			log.Printf("Failed to generate embedding for snippet: %v", err)
			continue
		}

		newSnippet := Snippet{
			Content:    snippetText,
			Labels:     labels,
			Source:     ref,
			ThumbsUp:   0,
			ThumbsDown: 0,
			CreatedAt:  time.Now(),
			Embedding:  embedding,
		}

		_, _, err = app.firestoreClient.Collection("snippets").Add(ctx, newSnippet)
		if err != nil {
			log.Printf("Failed to store snippet: %v", err)
			continue
		}
	}

	// Update the last refreshed timestamp on the source document
	_, err = ref.Set(ctx, map[string]interface{}{
		"last_refreshed": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Printf("Failed to update source last_refreshed: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing started"))
}

func (app *App) generateSnippets(ctx context.Context, content string) ([]string, error) {
	prompt := "Break down the following markdown into discrete, standalone instruction snippets. Each snippet should be a self-contained piece of instruction. Return the snippets as a JSON array of strings. Markdown: " + content
	resp, err := app.genaiClient.Models.GenerateContent(ctx, "gemini-pro", genai.Text(prompt), nil)
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		text := fmt.Sprint(part)
		var snippets []string
		if err := json.Unmarshal([]byte(text), &snippets); err != nil {
			return nil, fmt.Errorf("failed to unmarshal snippets: %w", err)
		}
		return snippets, nil
	}
	return nil, fmt.Errorf("unexpected response format or empty response")
}

func (app *App) generateLabels(ctx context.Context, snippet string) ([]string, error) {
	prompt := "Generate a list of relevant labels for the following snippet. Return the labels as a JSON array of strings. Snippet: " + snippet
	resp, err := app.genaiClient.Models.GenerateContent(ctx, "gemini-pro", genai.Text(prompt), nil)
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		text := fmt.Sprint(part)
		var labels []string
		if err := json.Unmarshal([]byte(text), &labels); err != nil {
			return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
		}
		return labels, nil
	}
	return nil, fmt.Errorf("unexpected response format or empty response")
}

func (app *App) generateEmbedding(ctx context.Context, snippet string) ([]float32, error) {
	contents := []*genai.Content{genai.NewContentFromText(snippet, genai.RoleUser)}
	result, err := app.genaiClient.Models.EmbedContent(ctx, "embedding-001", contents, nil)
	if err != nil {
		return nil, err
	}
	if len(result.Embeddings) == 0 || len(result.Embeddings[0].Values) == 0 {
		return nil, fmt.Errorf("empty embedding returned")
	}
	return result.Embeddings[0].Values, nil
}