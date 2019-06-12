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
	"net"
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

type deployinatorServer struct{}
type contextKey string

func cleanup(server *grpc.Server) {
	log.Println("Stopping Deployinator Server")
	server.Stop()
}

func newDeployinatorServer() pb.DeployinatorServer {
	return &deployinatorServer{}
}

func processServerMessage(ctx context.Context, msg *pubsub.Message) {
	var message pb.DeployStatusMessage
	stream := ctx.Value("stream")
	err := proto.Unmarshal(msg.Data, &message)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	message.Status = fmt.Sprintf("Status: %s, MsgID: %s.\n", message.Status, message.MsgID)
	stream.Send(message)
	msg.Ack()
	return
}

// send a message to pubsub
// sub to deploy status and stream messages back
func (ds *deployinatorServer) TriggerDeploy(in *pb.DeployMessage, stream pb.Deployinator_TriggerDeployServer) error {
	// log.Printf("triggering a deploy\n")
	response := new(pb.DeployStatusMessage)
	response.Status = "Starting deploy"
	stream.Send(response)
	// log.Println(in)
	topicName := fmt.Sprintf("%s-%s-deploy", in.Cenv, in.Cid)
	response.Status = fmt.Sprintf("Connecting to topic %s", topicName)
	stream.Send(response)
	cli := pubsubclient.PubSubClient{ProjectID: in.Cenv, TopicName: topicName}
	cli.Connect()
	response.Status = fmt.Sprintf("Connected to topic %s", topicName)
	stream.Send(response)
	msgid, err := cli.Publish(in)

	if err != nil {
		log.Println(err)
		return err
	}
	response.Status = fmt.Sprintf("Published %v to %s", in, topicName)
	stream.Send(response)
	cli.TopicName = fmt.Sprintf("%s-%s-deploystatus", in.Cenv, in.Cid)
	if err := cli.Subscribe(); err != nil {
		log.Println(err)
		return err
	}
	response.Status = fmt.Sprintf("Subscribed to %s", cli.TopicName)
	stream.Send(response)

	if msgid == response.MsgID {
		stream.Send(response)
	}

	k := contextKey("stream")
	cli.CTX = context.WithValue(context.Background(), k, &stream)
	// now if pub successful stream deploystatus messages back
	cli.GetMessage(processServerMessage)
	stream.Send(response)
	// is this needed? - cli.Stop()
	return nil
}

func (ds *deployinatorServer) DeployStatus(ctx context.Context, in *pb.DeployMessage) (*pb.DeployStatusMessage, error) {
	log.Printf("DeployStatus\n")
	response := new(pb.DeployStatusMessage)
	response.Status = "deploy status"
	return response, nil
}

func runServer(listentAddr string, port string) {
	c := make(chan os.Signal, 1)
	server := grpc.NewServer()
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(server)
	}()
	log.Printf("Starting Deployinator Server: %s:%s\n", listentAddr, port)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", listentAddr, port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	pb.RegisterDeployinatorServer(server, newDeployinatorServer())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	log.Println("Goodbye from Deployinator Server")
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run the deployinator server",
	Long: `Deployinator Server

	listens for deploy events
	sends events to pubsub
	listens for responses
	streams them back
	
	future - states in places?`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer(viper.GetString("listenAddr"), viper.GetString("port"))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("listen-addr", "0.0.0.0", "listen address")
	viper.BindPFlag("listenAddr", serverCmd.Flags().Lookup("listen-addr"))
	serverCmd.Flags().Int("port", 9091, "listen port")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
}