/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	config "bitbucket.com/metamorph/pkg/config"
	ctrlgRPCServer "bitbucket.com/metamorph/pkg/controller/grpc"
	ctrlgFSMServer "bitbucket.com/metamorph/pkg/controller/statemachine"
)

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Start MetaMorph controller",
	Long: `This will start MetaMorph Controller's gRPC server`,
	Run: func(cmd *cobra.Command, args []string) {
		config.SetLoggerConfig("logger.controllerpath")
		fmt.Println("Starting metamorph FSM")
     	go ctrlgFSMServer.StartMetamorphFSM(false) //check if this causes crash in case DB is not setup yet.
		fmt.Println("controller called")
		ctrlgRPCServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controllerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controllerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
