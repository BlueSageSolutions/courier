package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BlueSageSolutions/courier/pkg/commander"
	"github.com/BlueSageSolutions/courier/pkg/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var deployCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a series of deployment scriptss.",
	Long:  `Run a series of deployment scriptss`,
	Run: func(cmd *cobra.Command, args []string) {
		runDeploymentScripts()
	},
}

func loadDeploymentScripts(scriptLocation string) (*commander.DeploymentScripts, error) {
	fileInfo, err := os.Stat(scriptLocation)
	if err != nil {
		util.GetLogger().Error("Not a directory", zap.Error(err))
		return nil, err
	}
	path := scriptLocation
	if fileInfo.IsDir() {
		path = fmt.Sprintf("%s/.", scriptLocation)
	}
	deploymentScripts, err := commander.LoadDeploymentScripts(path)
	if err != nil {
		util.GetLogger().Error("LoadDeploymentScripts", zap.Error(err))
		return nil, err
	}
	return deploymentScripts, nil
}

func deleteJSONFiles(directory string) error {
	// Read the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	// Loop through the files in the directory
	for _, file := range files {
		if file.IsDir() {
			// Skip directories
			continue
		}

		// Check if the file has a ".json" extension
		if strings.HasSuffix(file.Name(), ".json") {
			// Construct the full file path
			filePath := filepath.Join(directory, file.Name())

			// Delete the JSON file
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("Error deleting %s: %v\n", filePath, err)
			} else {
				fmt.Printf("Deleted %s\n", filePath)
			}
		}
	}
	return nil
}

func runDeploymentScripts() {
	absolutePath := fixDots(DeploymentScripts)
	reportsPath := fixDots(Reports)
	deploymentScripts, err := loadDeploymentScripts(absolutePath)
	if err != nil {
		util.GetLogger().Error("loadDeploymentScripts", zap.Error(err))
		os.Exit(1)
	}
	os.Setenv("PROCESS", "transformer")
	os.Setenv("SOURCE_TYPE", "local")
	os.Setenv("ENABLE_CLOUD_WATCH", "false")
	os.Setenv("ENABLE_HEC", "false")

	fmt.Printf("### arguments:\n\tdirectory: %s\n", absolutePath)
	fmt.Printf("\tdry run: %v\n", DryRun)
	fmt.Printf("\tdestroy: %v\n", Destroy)
	fmt.Println("###")

	for _, scriptList := range *deploymentScripts {
		for scriptIndex, script := range scriptList.DeploymentScripts {
			script.Destroy = Destroy
			script.DryRun = DryRun
			scriptList.DeploymentScripts[scriptIndex] = script
		}
	}
	executionContext := &commander.ExecutionContext{Client: Client, Environment: Environment}
	results := deploymentScripts.Execute(executionContext)

	results.Directory = fmt.Sprintf("%s/deployed-at-%s", reportsPath, commander.Timestamp())

	err = os.MkdirAll(results.Directory, os.ModePerm)
	if err != nil {
		os.Exit(1)
	}
	_, err = results.Publish(reportsPath, *deploymentScripts)
	if err != nil {
		util.GetLogger().Error("loadDeploymentScripts", zap.Error(err))
		os.Exit(1)
	}
	if CleanTmp {
		deleteJSONFiles("/tmp")
	}

}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&Client, "client", "c", "testclient", "Name of client")
	deployCmd.Flags().StringVarP(&Environment, "environment", "e", "dev", "Type of environment")
	deployCmd.Flags().StringVarP(&DeploymentScripts, "directory", "s", "./scripts/bluesage-dlp/databases/create/schema.yaml", "Directory containing deployment scripts")
	deployCmd.Flags().StringVarP(&Reports, "reports", "r", "./reports", "Directory where the reports will be written to")
	deployCmd.Flags().BoolVarP(&Destroy, "destroy", "z", false, "Run the cleanup script")
	deployCmd.Flags().BoolVarP(&CleanTmp, "cleantmp", "t", false, "Clean /tmp after execution")
	deployCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "Dry run")
}
