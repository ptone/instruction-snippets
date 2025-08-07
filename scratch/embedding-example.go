import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	contents := []*genai.Content{
		genai.NewContentFromText("What is the meaning of life?"),
		genai.NewContentFromText("How does photosynthesis work?"),
		genai.NewContentFromText("Tell me about the history of the internet."),
	}
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		contents,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	embeddings, err := json.MarshalIndent(result.Embeddings, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(embeddings))
}