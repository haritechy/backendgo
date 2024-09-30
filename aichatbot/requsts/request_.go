package requsts

import (
	"fmt"
	"net/http"
	"slack-chatbot/models"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

func GetGeminiResponse(apiKey, prompt string) (string, error) {
	client := resty.New()
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	var geminiResp models.GeminiResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&geminiResp).
		Post(fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", apiKey))

	if err != nil {
		return "", fmt.Errorf("request error: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		body := resp.Body() // Get the response body
		return "", fmt.Errorf("failed to get a valid response from Gemini: %s, body: %s", resp.Status(), string(body))
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text), nil
	}
	body := resp.Body()
	return "", fmt.Errorf("no content received from Gemini API, body: %s", string(body))
}
