package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func fixDots(path string) (string, string) {
	if strings.HasPrefix(path, "./") {
		raw := strings.Replace(path, "./", "", 1)
		pwd, err := os.Getwd()
		if err != nil {
			panic(path)
		}
		path = strings.Replace(path, "./", pwd+"/", 1)
		return raw, path
	}
	return path, path
}

var CurlUrl string
var CurlPayloadFile string
var CurlHeadersFile string
var FileToConvert string
var CurlMethod string
var FilterDir string
var DeploymentScriptDir string
var RunCleanup bool
var RunMain bool
var GenerateTestDir string
var FilterDirActual string
var TranslationDir string
var DataDir string
var ConfigFile string
var DefaultYaml string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "courier",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.courier.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
