package imglatex

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const PROMPTJSON = `{
  "messages": [
    {
      "role": "user",
      "content": [
        {
          "type": "text",
          "text": "Convert the given mathematical input from the image into latex output. Do not attempt to solve the problem, only convert the problem into latex. Only output latex, nothing else, including any consequent suggestions."
        },
        {
          "type": "image_url",
          "image_url": {
            "url": "data:image/jpeg;base64,%s"
          }
        }
      ]
    }
  ],
  "model": "llama-3.2-11b-vision-preview",
  "temperature": 1,
  "max_completion_tokens": 1024,
  "top_p": 1,
  "stream": false,
  "stop": null
}`

type groqJson struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		QueueTime        float64 `json:"queue_time"`
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
	XGroq             struct {
		ID string `json:"id"`
	} `json:"x_groq"`
}

type ImgLatex struct {
	apiKey string
}

func NewImgLatex(apiKey string) *ImgLatex {

	return &ImgLatex{apiKey: apiKey}

}

func (il *ImgLatex) ImageToLatex(image io.Reader) (latex string, err error) {

	imageContent, err := io.ReadAll(image)
	if err != nil {

		return "", err

	}

	base64ImgContent := make([]byte, len(imageContent)*2)
	base64.StdEncoding.Encode(base64ImgContent, imageContent)

	prompt := fmt.Sprintf(PROMPTJSON, string(base64ImgContent))

	postBuffer := bytes.NewBuffer([]byte(prompt))

	request, err := http.NewRequest(http.MethodPost, "https://api.groq.com/openai/v1/chat/completions", postBuffer)
	client := &http.Client{}

	resp, err := client.Do(request)

	if err != nil {

		return "", err

	}

	defer resp.Body.Close()

	gq := groqJson{}
	if err := json.NewDecoder(resp.Body).Decode(&gq); err != nil {

		return "", err

	}

	return gq.Choices[0].Message.Content, nil
}
