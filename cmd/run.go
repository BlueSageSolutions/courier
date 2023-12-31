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

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a series of deployment scripts.",
	Long:  `Run a series of deployment scripts`,
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

	if len(os.Getenv("COURIER_CLIENT")) > 0 {
		Client = os.Getenv("COURIER_CLIENT")
	}

	if len(os.Getenv("COURIER_ENVIRONMENT")) > 0 {
		Environment = os.Getenv("COURIER_ENVIRONMENT")
	}

	if strings.ToUpper(os.Getenv("COURIER_DESTROY")) == "TRUE" {
		Destroy = true
	}

	if strings.ToUpper(os.Getenv("COURIER_CLEAN_TMP")) == "TRUE" {
		CleanTmp = true
	}

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
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&Client, "client", "c", "", "Name of client")
	runCmd.Flags().StringVarP(&Environment, "environment", "e", "", "Type of environment")
	runCmd.Flags().StringVarP(&DeploymentScripts, "directory", "s", "./scripts/bluesage-dlp/databases/create/create-users.yaml", "Either the location of a single deployment script;\nor, the directory containing deployment scripts")
	runCmd.Flags().StringVarP(&Reports, "reports", "r", "./reports", "Directory where the reports will be written to")
	runCmd.Flags().BoolVarP(&Destroy, "destroy", "z", false, "Run the cleanup script")
	runCmd.Flags().BoolVarP(&CleanTmp, "cleantmp", "t", false, "Clean /tmp after execution")
	runCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "Dry run")
}
