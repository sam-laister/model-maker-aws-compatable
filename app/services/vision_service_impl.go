package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type VisionServiceImpl struct{}

func NewVisionService() *VisionServiceImpl {
	return &VisionServiceImpl{}
}

func (s *VisionServiceImpl) AnalyseImage(imagePath string, prompt string) (string, error) {
	if prompt == "" {
		prompt = "Caption this image, give one sentence only"
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	file, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Create the request.
	req := []genai.Part{
		genai.ImageData("jpeg", imageBytes),
		genai.Text(prompt),
	}

	// Generate content.
	resp, err := model.GenerateContent(ctx, req...)
	if err != nil {
		panic(err)
	}

	var result string = extractText(resp)

	return result, nil
}

func extractText(resp *genai.GenerateContentResponse) string {
	var builder strings.Builder

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
				if text, ok := part.(genai.Text); ok {
					builder.WriteString(string(text))
				}
			}
		}
	}
	return builder.String()

}
