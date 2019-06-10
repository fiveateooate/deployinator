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
	cli.Subscribe()
	cli.GetMessage()
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
