package main

import (
	"context"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// OpenAI API kalitini o'rnating
	apiKey := os.Getenv("OPENAI_API_KEY") // Environment variable dan olish
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY muhit o'zgaruvchisi topilmadi")
	}

	// OpenAI clientini yaratish
	client := openai.NewClient(apiKey)

	// Embedding generatsiya qilish uchun matn
	text := "Go tilida AI dasturlash"

	// Embedding so'rovini yuborish
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2, // Modelni tanlash
	}
	resp, err := client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		log.Fatalf("Embedding generatsiya qilishda xato: %v", err)
	}

	// Natijani ko'rsatish
	embedding := resp.Data[0].Embedding // Birinchi matn uchun embedding
	fmt.Printf("Matn: %s\n", text)
	fmt.Printf("Embedding o'lchami: %d\n", len(embedding))
	fmt.Printf("Dastlabki 3 qiymat: %v\n", embedding[:3])
}