// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fiveateooate/deployinator/internal/pubsubclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deployService - deploy a service and stream messages
// publish status messages to deploystatus topic
func deployService() error {
	return nil
}

func clientCleanup(cli *pubsubclient.PubSubClient) {
	log.Println("Stopping Deplopyinator Client")
	cli.Disconnect()
}

func runClient(host string) {
	c := make(chan os.Signal, 1)
	cli := pubsubclient.PubSubClient{ProjectID: viper.GetString("projectID"), TopicName: viper.GetString("topicName")}
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		clientCleanup(&cli)
	}()
	log.Printf("Starting Deployinator Client\n")
	log.Printf("topic: %s project: %s", viper.GetString("topicName"), viper.GetString("projectID"))
	cli.Connect()
	// probably need to pass a handle func
	cli.Subscribe()
	cli.GetMessage()
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "run the deployinator client",
	Long: `twoish modes, probably dumb

	1. subscribe to a deploy queue for specific cluster and deploy stuff
		in this mode runs in whatever/every/all/your base k8s(kates ! k-eights) clusters
	2. send deploy trigger to deployinator server and wait for responses
		in this mode it should run in same k8s(kates ! k-eights) cluster as server`,
	Run: func(cmd *cobra.Command, args []string) {
		host := fmt.Sprintf("%s:%s", viper.GetString("serverAddr"), viper.GetString("serverPort"))
		runClient(host)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().String("server-addr", "127.0.0.1", "server address")
	viper.BindPFlag("serverAddr", clientCmd.Flags().Lookup("server-addr"))
	clientCmd.Flags().Int("server-port", 9091, "server port")
	viper.BindPFlag("serverPort", clientCmd.Flags().Lookup("server-port"))
	clientCmd.Flags().String("cmd", "", "add an alert")
	viper.BindPFlag("cmd", clientCmd.Flags().Lookup("cmd"))
}
