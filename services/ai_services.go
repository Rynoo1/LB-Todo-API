package services

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AiService struct {
	client *openai.Client
}

func NewAiService(apiKey string) *AiService {
	client := openai.NewClient(option.WithAPIKey(apiKey))

	return &AiService{
		client: &client,
	}
}

func (s *AiService) AiServiceCall(title, description, action string) (string, error) {
	file := "services/templates/" + action + ".txt"

	systemPromptBytes, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	systemPrompt := string(systemPromptBytes)

	// var promptId = "pmpt_69303e98ab8881908aaa5cfe81563b330e77cb5d01fc7e1e"

	sysmsg := openai.ChatCompletionMessageParamUnion{
		OfSystem: &openai.ChatCompletionSystemMessageParam{
			Content: openai.ChatCompletionSystemMessageParamContentUnion{
				OfString: openai.String(systemPrompt),
			},
		},
	}

	userContent := fmt.Sprintf(
		"Title: %s\n\nDescription: %s\n\nRespond ONLY with JSON",
		title, description,
	)

	userMsg := openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(userContent),
			},
		},
	}

	messages := []openai.ChatCompletionMessageParamUnion{sysmsg, userMsg}

	resp, err := s.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model:       "gpt-4.1-nano",
		Temperature: openai.Float(0.0),
		Messages:    messages,
	})

	if err != nil {
		return "", err
	}

	// fmt.Println(resp.OutputText())

	output := resp.Choices[0].Message.Content
	return output, nil
}
