package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
)

type ClaudeAgent struct {
	client *anthropic.Client
}

func NewClaudeAgent(key string) *ClaudeAgent {
	c := anthropic.NewClient(key)
	return &ClaudeAgent{
		client: c,
	}
}

func (c *ClaudeAgent) AskHint(sfen string) (string, error) {
	return c.Ask(sfen,
		"You are a shogi tutor. You will receive a SFEN string representing a shogi game. Respond ONLY with a valid JSON object that strictly adheres to the following schema: { \"next_move\": string } where the value is the suggested movement in Hodges notation. Do not include any extra text or explanation.")
}

func (c *ClaudeAgent) AskMovement(sfen string, difficulty AgentLevel) (string, error) {
	level := "begginer"
	switch difficulty {
	case Pro:
		level = "profesional"
	case Medium:
		level = "medium level"
	case Begginer:
		level = "begginer"
	default:
		level = "medium level"
	}

	system := fmt.Sprintf("You are a %s shogi player. You are going to receive a sfen string representing a shogi game and are going to respond with the next movement for the current player according to your level using Hodges notation", level)
	return c.Ask(sfen, system)
}

func (c *ClaudeAgent) Ask(message string, system string) (string, error) {
	request := anthropic.MessagesRequest{
		Model:  anthropic.ModelClaude3Dot5HaikuLatest,
		System: system,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(message),
		},
		MaxTokens: 1000,
		/*
			Tools: []anthropic.ToolDefinition{
				{
					Name:        "get_shogi_movement",
					Description: "Get the next movement in a given shogi game",
					InputSchema: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"next_move": {
								Type:        jsonschema.String,
								Description: "The suggested next shogi movement in Hodges notation",
							},
						},
						Required: []string{"next_move"},
					},
				},
			},*/
	}
	resp, err := c.client.CreateMessages(context.Background(), request)
	if err != nil {
		return "", err
	}
	/*
		request.Messages = append(request.Messages, anthropic.Message{
			Role:    anthropic.RoleAssistant,
			Content: resp.Content,
		})

		var toolUse *anthropic.MessageContentToolUse

		for _, c := range resp.Content {
			if c.Type == anthropic.MessagesContentTypeToolUse {
				toolUse = c.MessageContentToolUse
			}
		}

		if toolUse == nil {
			panic("tool use not found")
		}

		request.Messages = append(request.Messages, anthropic.NewToolResultsMessage(toolUse.ID, "7f7g", false))

		finalResp, err := c.client.CreateMessages(context.Background(), request)
		if err != nil {
			panic(err)
		}
	*/
	var finalText string
	for _, msg := range resp.Content {
		if msg.Type == anthropic.MessagesContentTypeText && msg.Text != nil {
			finalText = *msg.Text
			break
		}
	}
	if finalText == "" {
		return "", errors.New("final answer not found")
	}
	// Unmarshal the final answer.
	var result Movement
	if err := json.Unmarshal([]byte(finalText), &result); err != nil {
		return "", fmt.Errorf("failed to parse final answer as JSON: %w", err)
	}
	if result.NextMove == "" {
		return "", errors.New("final answer JSON does not contain a valid next_move")
	}
	return result.NextMove, nil
}
