package cmd

import (
	"fmt"
	"os"

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

func runDeploymentScripts() {
	relativePath, absolutePath := fixDots(DeploymentScriptDir)
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
	fmt.Printf("\trun main: %v\n", RunMain)
	fmt.Printf("\trun cleanup: %v\n", RunCleanup)
	fmt.Println("###")

	for _, scriptList := range *deploymentScripts {
		for scriptIndex, script := range scriptList.DeploymentScripts {
			script.RunCleanup = RunCleanup
			script.RunMain = RunMain
			scriptList.DeploymentScripts[scriptIndex] = script
		}
	}
	results := deploymentScripts.Execute()

	results.Directory = fmt.Sprintf("%s/deployed-at-%s", absolutePath, commander.Timestamp())

	err = os.MkdirAll(results.Directory, os.ModePerm)
	if err != nil {
		os.Exit(1)
	}
	_, err = results.Publish(relativePath, *deploymentScripts)
	if err != nil {
		util.GetLogger().Error("loadDeploymentScripts", zap.Error(err))
		os.Exit(1)
	}

}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&DeploymentScriptDir, "scripts", "s", "./scripts/bluesage-dlp", "Directory containing deployment scripts")
	deployCmd.Flags().BoolVarP(&RunCleanup, "cleanup", "c", false, "Run the cleanup script")
	deployCmd.Flags().BoolVarP(&RunMain, "main", "m", false, "Run the main script")
}
