package agent

import (
	"context"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIAgent struct {
	client         *openai.Client
	schemaResponse openai.ResponseFormatJSONSchemaJSONSchemaParam
}

type Movement struct {
	NextMove string `json:"next_move" jsonschema_description:"The suggested next shogi movement"`
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var MovementResponseSchema = GenerateSchema[Movement]()

func NewOpenAIAgent(key string) *OpenAIAgent {
	return &OpenAIAgent{
		client: openai.NewClient(option.WithAPIKey(key)),
		schemaResponse: openai.ResponseFormatJSONSchemaJSONSchemaParam{
			Name:        openai.F("movement"),
			Description: openai.F("Next movement in a shogi game"),
			Schema:      openai.F(MovementResponseSchema),
			Strict:      openai.Bool(true),
		},
	}
}

func (c *OpenAIAgent) AskHint(sfen string) (string, error) {
	system := "You are a shogi tutor. You are going to receive a sfen string representing a shogi game and are meant to respond with a movement suggestion for the current player"
	return c.Ask(sfen, system)
}

func (c *OpenAIAgent) AskMovement(sfen string, difficulty AgentLevel) (string, error) {
	var level string
	switch difficulty {
	case Pro:
		level = "professional"
	case Medium:
		level = "medium level"
	case Begginer:
		level = "beginner"
	default:
		level = "medium level"
	}
	system := fmt.Sprintf("You are a %s shogi player. You are going to receive a sfen string representing a shogi game and are going to respond with the next movement for the current player according to your level", level)
	return c.Ask(sfen, system)
}

func (c *OpenAIAgent) Ask(message string, system string) (string, error) {
	var messages []openai.ChatCompletionMessageParamUnion
	if system != "" {
		messages = append(messages, openai.SystemMessage(system))
	}
	messages = append(messages, openai.UserMessage(message))
	chatCompletion, err := c.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(c.schemaResponse),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
