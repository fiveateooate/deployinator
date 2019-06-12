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
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	pb "github.com/fiveateooate/deployinator/deployproto"
	"github.com/fiveateooate/deployinator/internal/pubsubclient"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func processMessage(ctx context.Context, msg *pubsub.Message) {
	var message pb.DeployMessage
	var response pb.DeployStatusMessage
	err := proto.Unmarshal(msg.Data, &message)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf("got messageid %s", msg.ID)
	topicName := fmt.Sprintf("%s-%s-deploystatus", viper.GetString("cenv"), viper.GetString("cid"))
	pscli := pubsubclient.PubSubClient{ProjectID: viper.GetString("cenv"), TopicName: topicName}
	pscli.Connect()
	response.Status = fmt.Sprintf("Deploying %s to namespace  %s.\n", message.Name, message.Namespace)
	response.MsgID = msg.ID
	pscli.PublishResponse(&response)
	msg.Ack()
	return
}

// deployService - deploy a service and stream messages
// publish status messages to deploystatus topic
func deployService(host string) error {
	service := pb.DeployMessage{Name: "MyService", Namespace: "MyNamespace", Cid: viper.GetString("cid"), Cenv: viper.GetString("cenv")}
	log.Printf("Triggering a deploy of %v", service)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDeployinatorClient(conn)
	stream, err := c.TriggerDeploy(context.Background(), &service)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.TriggerDeploy(_) = _, %v", c, err)
		}
		log.Println(resp)
	}
	return nil
}

func deployStatus(host string) {
	log.Println("deploy status")
	service := pb.DeployMessage{Name: "MyService", Cid: viper.GetString("cid"), Cenv: viper.GetString("cenv")}
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDeployinatorClient(conn)
	resp, err := c.DeployStatus(context.Background(), &service)
	log.Printf("error: %v, resp: %v\n", err, resp)
}

func clientCleanup(cli *pubsubclient.PubSubClient) {
	log.Println("Stopping Deplopyinator Client")
	cli.Disconnect()
}

func runClient() {
	c := make(chan os.Signal, 1)
	topicName := fmt.Sprintf("%s-%s-deploy", viper.GetString("cenv"), viper.GetString("cid"))
	pscli := pubsubclient.PubSubClient{ProjectID: viper.GetString("cenv"), TopicName: topicName}
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		clientCleanup(&pscli)
	}()
	log.Printf("Starting Deployinator Client\n")
	log.Printf("Listening for events on topic: %s in project: %s", topicName, viper.GetString("cenv"))
	pscli.Connect()
	pscli.Subscribe()
	pscli.GetMessage(processMessage)
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
		switch viper.GetString("cmd") {
		case "trigger":
			host := fmt.Sprintf("%s:%s", viper.GetString("serverAddr"), viper.GetString("serverPort"))
			deployService(host)
		case "status":
			host := fmt.Sprintf("%s:%s", viper.GetString("serverAddr"), viper.GetString("serverPort"))
			deployStatus(host)
		case "deployinator":
			runClient()
		default:
			log.Printf("unknown command %s\n", viper.GetString("cmd"))
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().String("server-addr", "127.0.0.1", "server address")
	viper.BindPFlag("serverAddr", clientCmd.Flags().Lookup("server-addr"))
	clientCmd.Flags().Int("server-port", 9091, "server port")
	viper.BindPFlag("serverPort", clientCmd.Flags().Lookup("server-port"))
	clientCmd.Flags().String("cmd", "", "run client command (deployinator|trigger)")
	viper.BindPFlag("cmd", clientCmd.Flags().Lookup("cmd"))
}
