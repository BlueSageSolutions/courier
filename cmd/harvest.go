/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

// harvestCmd represents the harvest command
var harvestCmd = &cobra.Command{
	Use:   "harvest",
	Short: "Enumerates all template variables in a repeatable SQL script.",
	Long:  `Enumerates all template variables in a SQL script.`,
	Run: func(cmd *cobra.Command, args []string) {
		SQLTemplate = fixDots(SQLTemplate)
		file, err := os.Open(SQLTemplate)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// Regular expression to match tokens
		tokenRegex, err := regexp.Compile(`<[^>]*>`)
		if err != nil {
			fmt.Println("Error compiling regex:", err)
			return
		}
		dictionary := make(map[string]string, 0)
		// Read file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			// Find and print all tokens in the line
			matches := tokenRegex.FindAllString(line, -1)
			for _, match := range matches {
				dictionary[match] = match
			}
		}
		for key := range dictionary {
			fmt.Println(key)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(harvestCmd)
	harvestCmd.Flags().StringVarP(&SQLTemplate, "repeatable-script", "r", "", "The repeatable SQL script.")
}
