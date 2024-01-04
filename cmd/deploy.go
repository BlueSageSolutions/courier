/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/BlueSageSolutions/courier/pkg/commander"
	"github.com/BlueSageSolutions/courier/pkg/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func loadManifest(location string) (*commander.Manifest, error) {
	fileInfo, err := os.Stat(location)
	if err != nil {
		util.GetLogger().Error("Not a directory", zap.Error(err))
		return nil, err
	}
	if fileInfo.IsDir() {
		location = fmt.Sprintf("%s/.", location)
	}
	manifest, err := commander.LoadManifest(location)
	return manifest, err
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Using a deployment manifest, execute the steps necessary to fulfill the manifest's expectations.",
	Long:  `Using a deployment manifest, execute the steps necessary to fulfill the manifest's expectations.`,
	Run: func(cmd *cobra.Command, args []string) {
		manifestPath := fixDots(Manifest)
		reportsPath := fixDots(Reports)
		manifest, err := loadManifest(manifestPath)
		if err != nil {
			util.GetLogger().Error("loadManifest", zap.Error(err))
			os.Exit(1)
		}
		err = manifest.Validate()
		if err != nil {
			util.GetLogger().Error("Validate", zap.Error(err))
			os.Exit(1)
		}
		results, err := manifest.Process()
		if err != nil {
			util.GetLogger().Error("Process", zap.Error(err))
			os.Exit(1)
		}
		err = results.Publish(manifest, reportsPath)
		if err != nil {
			util.GetLogger().Error("Publish", zap.Error(err))
			os.Exit(1)
		}

	},
}

// Validate manifest
// Iterate over environments
// create schemas and users
/*
client-profile:
  name: client1
  infrastructure-type: t2
  contacts:
    - fred@client1.com
  email-domain: client1.com
environments:
  - name: dev
    start-date: 02/01/2024
    operations:
      - create
services:
  - name: lion
    url: loan.client1.com
    operations:
      - build
      - install
      - configure
      - start
databases:
  - name: default
    operations:
      - create
*/
func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&Manifest, "manifest", "m", "./manifests/IT-1202.yaml", "The deployment manifest. Named after the Jira ticket ID.")
	deployCmd.Flags().StringVarP(&Reports, "reports", "r", "./reports", "Directory where the reports will be written to")
	deployCmd.Flags().BoolVarP(&CleanTmp, "cleantmp", "t", false, "Clean /tmp after execution")
	deployCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "Dry run")
}
