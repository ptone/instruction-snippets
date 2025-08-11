package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rs/cors"
	"google.golang.org/genai"
)

// App holds application dependencies
type App struct {
	firestoreClient *firestore.Client
	genaiClient     *genai.Client
}

// ProcessRequest defines the structure for the incoming request
type ProcessRequest struct {
	Content        string `json:"content,omitempty"`
	URL            string `json:"url,omitempty"`
	Key            string `json:"key,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	SubmitterID    string `json:"submitterId,omitempty"`
	SubmitterEmail string `json:"submitterEmail,omitempty"`
}

// Source defines the structure for the sources collection
type Source struct {
	Content        string    `firestore:"content"`
	URL            string    `firestore:"url,omitempty"`
	LastRefreshed  time.Time `firestore:"last_refreshed"`
	Type           string    `firestore:"type"`
	Status         string    `firestore:"status"`
	Key            string    `firestore:"key"`
	SubmitterID    string    `firestore:"submitterId"`
	SubmitterEmail string    `firestore:"submitterEmail"`
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
		Location: "global",
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}

	app := &App{
		firestoreClient: firestoreClient,
		genaiClient:     genaiClient,
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(http.DefaultServeMux)

	http.HandleFunc("/process", app.processHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (app *App) processSnippetsAsync(ctx context.Context, content string, sourceRef *firestore.DocumentRef, limit int) {
	log.Println("Starting snippet processing...")
	// Generate snippets from the markdown content
	snippets, err := app.generateSnippets(ctx, content, limit)
	if err != nil {
		log.Printf("Failed to generate snippets: %v", err)
		// Update the source document with an error status
		_, updateErr := sourceRef.Set(ctx, map[string]interface{}{
			"status": "error",
		}, firestore.MergeAll)
		if updateErr != nil {
			log.Printf("Failed to update source status: %v", updateErr)
		}
		return
	}

	log.Printf("Generated %d snippets", len(snippets))

	// Process and store snippets
	for i, snippetText := range snippets {
		log.Printf("Processing snippet %d/%d: %s", i+1, len(snippets), snippetText)

		labels, err := app.generateLabels(ctx, snippetText)
		if err != nil {
			log.Printf("Failed to generate labels for snippet %d: %v", i+1, err)
			continue
		}
		log.Printf("Generated labels for snippet %d: %v", i+1, labels)

		embedding, err := app.generateEmbedding(ctx, snippetText)
		if err != nil {
			log.Printf("Failed to generate embedding for snippet %d: %v", i+1, err)
			continue
		}

		newSnippet := Snippet{
			Content:    snippetText,
			Labels:     labels,
			Source:     sourceRef,
			ThumbsUp:   0,
			ThumbsDown: 0,
			CreatedAt:  time.Now(),
			Embedding:  embedding,
		}

		_, _, err = app.firestoreClient.Collection("snippets").Add(ctx, newSnippet)
		if err != nil {
			log.Printf("Failed to store snippet %d: %v", i+1, err)
			continue
		}
		log.Printf("Successfully stored snippet %d", i+1)
	}

	log.Println("Snippet processing complete.")

	// Update the source document to indicate processing is complete
	_, err = sourceRef.Set(ctx, map[string]interface{}{
		"status":         "processed",
		"last_refreshed": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Printf("Failed to update source status: %v", err)
	}
	log.Println("Source document status updated to 'processed'")
}

func (app *App) processHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /process")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	var req ProcessRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		log.Printf("Failed to read request body: %v", err)
		return
	}
	log.Printf("Request body: %s", string(body))

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		log.Printf("Failed to decode request: %v", err)
		return
	}

	if req.Content == "" && req.URL == "" {
		http.Error(w, "Request must contain either 'content' or 'url'", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var content string

	// If URL is provided, fetch content from it
	if req.URL != "" {
		url := req.URL
		if strings.Contains(url, "github.com") && strings.Contains(url, "/blob/") {
			url = strings.Replace(url, "github.com", "raw.githubusercontent.com", 1)
			url = strings.Replace(url, "/blob/", "/", 1)
		}

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
			log.Printf("Failed to fetch URL %s: %v", req.URL, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("Failed to fetch URL: status code %d", resp.StatusCode), http.StatusInternalServerError)
			log.Printf("Failed to fetch URL %s: status code %d", req.URL, resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			log.Printf("Failed to read response body from %s: %v", req.URL, err)
			return
		}
		content = string(body)
	} else {
		content = req.Content
	}

	// Use the provided key or default to the URL
	key := req.Key
	if key == "" {
		key = req.URL
	}
	if key == "" {
		http.Error(w, "A 'key' or 'url' must be provided for the source", http.StatusBadRequest)
		return
	}

	// Check if a source with the given key already exists
	iter := app.firestoreClient.Collection("sources").Where("key", "==", key).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		if err.Error() == "no more items in iterator" {
			doc = nil // Source doesn't exist yet
		} else {
			http.Error(w, "Failed to query for existing source", http.StatusInternalServerError)
			log.Printf("Failed to query for existing source: %v", err)
			return
		}
	}

	var sourceRef *firestore.DocumentRef
	if doc != nil {
		// Source exists, delete old snippets and update
		sourceRef = doc.Ref
		log.Printf("Source with key '%s' found, reprocessing...", key)

		_, err := sourceRef.Set(ctx, map[string]interface{}{"status": "processing"}, firestore.MergeAll)
		if err != nil {
			http.Error(w, "Failed to update source status", http.StatusInternalServerError)
			log.Printf("Failed to update source status: %v", err)
			return
		}

		if err := app.deleteSnippetsBySource(ctx, sourceRef); err != nil {
			http.Error(w, "Failed to delete old snippets", http.StatusInternalServerError)
			log.Printf("Failed to delete old snippets: %v", err)
			return
		}

		updateData := map[string]interface{}{
			"content":        content,
			"last_refreshed": time.Now(),
		}
		if req.URL != "" {
			updateData["url"] = req.URL
		}

		_, err = sourceRef.Set(ctx, updateData, firestore.MergeAll)
		if err != nil {
			http.Error(w, "Failed to update source", http.StatusInternalServerError)
			log.Printf("Failed to update source: %v", err)
			return
		}
	} else {
		// No existing source, create a new one
		sourceType := "file"
		if req.URL != "" {
			sourceType = "url"
		}
		source := Source{
			Content:        content,
			URL:            req.URL,
			LastRefreshed:  time.Now(),
			Type:           sourceType,
			Status:         "processing",
			Key:            key,
			SubmitterID:    req.SubmitterID,
			SubmitterEmail: req.SubmitterEmail,
		}
		sourceRef, _, err = app.firestoreClient.Collection("sources").Add(ctx, source)
		if err != nil {
			http.Error(w, "Failed to store source", http.StatusInternalServerError)
			log.Printf("Failed to add source: %v", err)
			return
		}
	}

	go app.processSnippetsAsync(context.Background(), content, sourceRef, req.Limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"documentId": sourceRef.ID})
}

func (app *App) deleteSnippetsBySource(ctx context.Context, sourceRef *firestore.DocumentRef) error {
	iter := app.firestoreClient.Collection("snippets").Where("source", "==", sourceRef).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err.Error() == "no more items in iterator" {
				break // All done
			}
			return fmt.Errorf("failed to iterate snippets: %v", err)
		}
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete snippet %s: %v", doc.Ref.ID, err)
			// Decide if you want to continue or return an error
		}
	}
	log.Printf("Deleted all snippets for source %s", sourceRef.ID)
	return nil
}

func (app *App) generateSnippets(ctx context.Context, content string, limit int) ([]string, error) {
	snippetSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"snippets": {
				Type:        genai.TypeArray,
				Description: "List of instruction snippets.",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"snippets"},
	}

	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:                 "extractSnippets",
					Description:          "Extracts discrete, standalone instruction snippets from a markdown document.",
					ParametersJsonSchema: snippetSchema,
				},
			},
		},
	}

	prompt := "Break down the following markdown into discrete, standalone instruction snippets, preserving the original markdown formatting and carriage returns. Each snippet should be a self-contained piece of instruction roughly a paragraph or so in size."
	if limit > 0 {
		prompt = fmt.Sprintf("%s Please provide no more than %d snippets.", prompt, limit)
	}
	prompt = prompt + " Markdown: " + content

	config := &genai.GenerateContentConfig{Tools: tools}
	resp, err := app.genaiClient.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if fc := part.FunctionCall; fc != nil {
			if fc.Name == "extractSnippets" {
				if snippets, ok := fc.Args["snippets"].([]interface{}); ok {
					var result []string
					for _, s := range snippets {
						if str, ok := s.(string); ok {
							result = append(result, str)
						}
					}
					return result, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("unexpected response format or empty response")
}

func (app *App) generateLabels(ctx context.Context, snippet string) ([]string, error) {
	labelSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"labels": {
				Type:        genai.TypeArray,
				Description: "List of relevant labels for a code snippet.",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"labels"},
	}

	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:                 "extractLabels",
					Description:          "Extracts relevant labels from a code snippet.",
					ParametersJsonSchema: labelSchema,
				},
			},
		},
	}

	prompt := "Generate a list of relevant topic labels for the following snippet. Snippet: " + snippet
	config := &genai.GenerateContentConfig{Tools: tools}
	resp, err := app.genaiClient.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if fc := part.FunctionCall; fc != nil {
			if fc.Name == "extractLabels" {
				if labels, ok := fc.Args["labels"].([]interface{}); ok {
					var result []string
					for _, l := range labels {
						if str, ok := l.(string); ok {
							result = append(result, str)
						}
					}
					return result, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("unexpected response format or empty response")
}

func (app *App) generateEmbedding(ctx context.Context, snippet string) ([]float32, error) {
	contents := []*genai.Content{genai.NewContentFromText(snippet, genai.RoleUser)}
	config := &genai.EmbedContentConfig{TaskType: "RETRIEVAL_DOCUMENT"}
	result, err := app.genaiClient.Models.EmbedContent(ctx, "gemini-embedding-001", contents, config)
	if err != nil {
		return nil, err
	}
	if len(result.Embeddings) == 0 || len(result.Embeddings[0].Values) == 0 {
		return nil, fmt.Errorf("empty embedding returned")
	}
	return result.Embeddings[0].Values, nil
}
