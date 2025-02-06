package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	apiKey = "AIzaSyDdsVVKIqddYVkS8FxrsfCsFx6TIXNvsA4"
)

func main() {
	ctx := context.Background()

	// Create a client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Choose a model
	model := client.EmbeddingModel("embedding-001")

	text := "Go tilida AI dasturlash"
	resp, err := model.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		log.Fatalf("Embedding xatosi: %v", err)
	}

	fmt.Printf("Matn: %s\n", text)
	fmt.Printf("Embedding o'lchami: %d\n", len(resp.Embedding.Values))
	fmt.Printf("Dastlabki 3 qiymat: %v\n", resp.Embedding.Values[:3])
}
