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
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	pb "github.com/fiveateooate/deployinator/deployproto"
	deployers "github.com/fiveateooate/deployinator/internal/deployers"
	"github.com/fiveateooate/deployinator/internal/pubsubclient"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func deployinateMessageHandler(ctx context.Context, msg *pubsub.Message) {
	var deploymessage pb.DeployMessage
	var response pb.DeployStatusMessage
	err := proto.Unmarshal(msg.Data, &deploymessage)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	msg.Ack()
	response.MsgID = msg.ID
	topicName := fmt.Sprintf("%s-%s-deploystatus", viper.GetString("cenv"), viper.GetString("cid"))
	pscli := pubsubclient.PubSubClient{ProjectID: viper.GetString("projectID"), TopicName: topicName}
	pscli.NewClient()
	pscli.SetTopic()
	log.Printf("Connected to topic %s\n", pscli.TopicName)
	// add some case here for different deployers

	switch deploymessage.Deployertype {
	case "helm":
		helmdeployer := deployers.NewHelmDeployer(deploymessage.Slug, deploymessage.Namespace, deploymessage.Version, viper.GetString("helmrepo"))
		log.Printf("hi: %v\n", helmdeployer)
		response.Success = true
		response.Status = fmt.Sprintf("Deploying %s to namespace %s\n", deploymessage.Slug, deploymessage.Namespace)
		err = helmdeployer.HelmDeploy(&deploymessage)
		if err != nil {
			response.Success = false
		}
		response.Status += helmdeployer.DeployResponse
	case "vaultpolicy":
		vaultDeployer := deployers.NewVaultDeployer()
		if err := vaultDeployer.Deploy(); err != nil {
			response.Status = "Fail"
			response.Success = false
		} else {
			response.Status = "Success"
			response.Success = true
		}
	default:
		response.Status = "Unknown Deployer Type"
		response.Success = false
	}
	pscli.PublishResponse(&response)
	pscli.Stop()
	return
}

func deployinateCleanup(cli *pubsubclient.PubSubClient) {
	log.Println("Stop Deplopyinating")
	cli.Disconnect()
}

func deployinate() {
	c := make(chan os.Signal, 1)
	topicName := fmt.Sprintf("%s-%s-deploy", viper.GetString("cenv"), viper.GetString("cid"))
	pscli := pubsubclient.PubSubClient{ProjectID: viper.GetString("projectID"), TopicName: topicName}
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		deployinateCleanup(&pscli)
	}()
	log.Printf("Start Deployinating\n")
	log.Printf("Listening for events on topic: %s in project: %s", topicName, viper.GetString("cenv"))
	pscli.NewClient()
	pscli.SetTopic()
	pscli.Subscribe()
	pscli.GetMessage(deployinateMessageHandler)
}

// deployinateCmd represents the deployinate command
var agentCmd = &cobra.Command{
	Use:   "agent [options]",
	Short: "run the deployinator agent",
	Long:  `Deploys stuff in a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		deployinate()
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().String("helmrepo", "stable", "helm repo to use for helm stuff")
	viper.BindPFlag("helmrepo", agentCmd.Flags().Lookup("helmrepo"))
}
