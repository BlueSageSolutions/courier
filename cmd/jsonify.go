/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// jsonifyCmd represents the jsonify command
var jsonifyCmd = &cobra.Command{
	Use:   "jsonify",
	Short: "Turns a BlueSage config file into JSON",
	Long:  `Convert config to JSON.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, err := os.Open(FileToConvert)
		if err != nil {
			panic(err)
		}
		defer inputFile.Close()

		// Create a map to hold the key-value pairs
		result := make(map[string]interface{})

		// Read the file line by line
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			line := scanner.Text()

			// Skip lines that start with '#'
			if strings.HasPrefix(line, "#") {
				continue
			}

			// Match the key and value using regex
			re := regexp.MustCompile(`^(\w+)=(?:"(.*?)"|\((.*?)\)|([^()]*))$`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 0 {
				key := matches[1]
				var value string
				if matches[2] != "" {
					// Directly use the value if it's a quoted string
					value = matches[2]
				} else if matches[3] != "" {
					// Process as a list if enclosed in parentheses
					elements := strings.Fields(strings.ReplaceAll(matches[3], "\"", ""))
					result[key] = elements
				} else {
					// Use the single value
					value = strings.ReplaceAll(matches[4], "\"", "")
				}

				// Assign the value to the key if it's not an array
				if value != "" {
					result[key] = value
				}
			}
		}

		// Check for errors during file reading
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		// Convert the map to JSON
		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			panic(err)
		}

		// Print the JSON result
		fmt.Println(string(jsonResult))
	},
}

func Save(results json.RawMessage) {
	// Create and open the output file
	outputFile, err := os.Create(FileToConvert + ".json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Write the JSON to the file
	_, err = outputFile.Write(results)
	if err != nil {
		panic(err)
	}

}

func init() {
	rootCmd.AddCommand(jsonifyCmd)
	jsonifyCmd.Flags().StringVarP(&FileToConvert, "file", "f", "", "File to jsonify")

}
