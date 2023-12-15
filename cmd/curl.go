package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, headerFile := fixDots(CurlHeadersFile)
		headers, err := LoadJson(headerFile)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
		_, payloadFile := fixDots(CurlPayloadFile)
		payload, err := LoadJson(payloadFile)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
		j, _ := CurlEquivalent(CurlUrl, headers, payload, CurlMethod)
		fmt.Printf("%s", j)
	},
}

type ErrorPayload struct {
	Error string `json:"error"`
}

func LoadJson(filename string) (json.RawMessage, error) {
	var jsonBlob json.RawMessage
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonFile, &jsonBlob)

	if err != nil {
		return nil, err
	}
	return jsonBlob, nil
}

func CurlEquivalent(url string, headers json.RawMessage, postPayload json.RawMessage, method string) (json.RawMessage, error) {
	// Create a new HTTP request with the given method
	req, err := http.NewRequest(method, url, bytes.NewBuffer(postPayload))
	if err != nil {
		errorPayload := ErrorPayload{Error: err.Error()}
		errorJSON, _ := json.Marshal(errorPayload)
		return errorJSON, err
	}
	// Add headers to the request
	headersMap := map[string]string{}

	err = json.Unmarshal([]byte(headers), &headersMap)
	if err != nil {
		errorPayload := ErrorPayload{Error: err.Error()}
		errorJSON, _ := json.Marshal(errorPayload)
		return errorJSON, err
	}

	for key, value := range headersMap {
		req.Header.Add(key, value)
	}
	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorPayload := ErrorPayload{Error: err.Error()}
		errorJSON, _ := json.Marshal(errorPayload)
		return errorJSON, err
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorPayload := ErrorPayload{Error: err.Error()}
		errorJSON, _ := json.Marshal(errorPayload)
		return errorJSON, err
	}
	return body, nil
}

func init() {
	rootCmd.AddCommand(curlCmd)
	curlCmd.Flags().StringVar(&CurlUrl, "url", "", "url")
	curlCmd.Flags().StringVar(&CurlMethod, "method", "POST", "method")
	curlCmd.Flags().StringVar(&CurlHeadersFile, "headers", "./headers.json", "headers")
	curlCmd.Flags().StringVar(&CurlPayloadFile, "payload", "./payload.json", "payload")
}
