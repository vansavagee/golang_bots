package gpt

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func GetResponseFromGPT(s string) (string, error) {

	client := resty.New()
	response, err := client.R().
		SetAuthToken(os.Getenv("GptKey")).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo-0613",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": s}},
			"max_tokens": 3000,
		}).
		Post(os.Getenv("GptEndpoint"))

	if err != nil {
		return "", fmt.Errorf("error while sending send the request: %v", err)
	}
	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("error while decoding JSON response: %v", err)
	}

	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content, nil
}
