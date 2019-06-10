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

	pb "github.com/fiveateooate/deployinator/deployproto"
	"github.com/fiveateooate/deployinator/internal/pubsubclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type deployinatorServer struct{}

func cleanup(server *grpc.Server) {
	log.Println("Stopping Deployinator Server")
	server.Stop()
}

func newDeployinatorServer() pb.DeployinatorServer {
	return &deployinatorServer{}
}

func (ds *deployinatorServer) DeployService(ctx context.Context, in *pb.DeployMessage) (*pb.DeployResponse, error) {
	response := new(pb.DeployResponse)
	response.Success = "pass"
	log.Println(in)
	cli := pubsubclient.PubSubClient{ProjectID: viper.GetString("projectID"), TopicName: viper.GetString("topicName")}
	cli.Connect()
	cli.Publish(in)
	log.Println("done")
	return response, nil
}

func runServer(listentAddr string, listenPort string) {
	c := make(chan os.Signal, 1)
	server := grpc.NewServer()
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(server)
	}()
	log.Println("Starting Deployinator Server")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", listentAddr, listenPort))
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer(viper.GetString("listenAddr"), viper.GetString("port"))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("listen-addr", "0.0.0.0", "listen address")
	viper.BindPFlag("listenAddr", serverCmd.Flags().Lookup("listen-addr"))
	serverCmd.Flags().Int("port", 9091, "listen port")
	viper.BindPFlag("listenPort", serverCmd.Flags().Lookup("port"))
}
