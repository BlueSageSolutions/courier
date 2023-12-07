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
		DeploymentScriptDir = fixDots(DeploymentScriptDir)
		runDeploymentScripts()
	},
}

func loadDeploymentScripts() (*commander.DeploymentScripts, error) {
	fileInfo, err := os.Stat(DeploymentScriptDir)
	if err != nil {
		util.GetLogger().Error("Not a directory", zap.Error(err))
		return nil, err
	}
	path := DeploymentScriptDir
	if fileInfo.IsDir() {
		path = fmt.Sprintf("%s/.", DeploymentScriptDir)
	}
	deploymentScripts, err := commander.LoadDeploymentScripts(path)
	if err != nil {
		util.GetLogger().Error("LoadDeploymentScripts", zap.Error(err))
		return nil, err
	}
	return deploymentScripts, nil
}

func runDeploymentScripts() {
	deploymentScripts, err := loadDeploymentScripts()
	if err != nil {
		util.GetLogger().Error("loadDeploymentScripts", zap.Error(err))
		os.Exit(1)
	}
	os.Setenv("PROCESS", "transformer")
	os.Setenv("SOURCE_TYPE", "local")
	os.Setenv("ENABLE_CLOUD_WATCH", "false")
	os.Setenv("ENABLE_HEC", "false")

	fmt.Printf("directory: %s\n", DeploymentScriptDir)
	results := deploymentScripts.Execute()
	_, err = results.Publish(DeploymentScriptDir, *deploymentScripts)
	if err != nil {
		util.GetLogger().Error("loadDeploymentScripts", zap.Error(err))
		os.Exit(1)
	}

}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&DeploymentScriptDir, "deploymentScripts", "s", "./scripts/exp", "Directory containing deployment scriptss")
}
