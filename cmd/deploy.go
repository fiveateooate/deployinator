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

	"github.com/fiveateooate/deployinator/internal/envfilehandler"
	"github.com/wsxiaoys/terminal/color"

	pb "github.com/fiveateooate/deployinator/deployproto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// deployService - deploy a service and stream messages
// publish status messages to deploystatus topic
func deployService(host string) error {
	var envyml envfilehandler.Envfile
	envyml.LoadEnvfile(viper.GetString("deploydescription"))
	service := pb.DeployMessage{Slug: envyml.Slug, Namespace: envyml.Domain, Cid: viper.GetString("cid"), Cenv: viper.GetString("cenv")}
	service.Version = "v2.8.1-6"
	log.Printf("Service: %v\n", service)
	log.Printf("Triggering a deploy of %s", service.Slug)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDeployinatorClient(conn)
	resp, err := c.TriggerDeploy(context.Background(), &service)
	if err != nil {
		color.Printf("@r%s", err)
		return err
	}
	log.Println(resp.Status)
	return nil
}

// deployCmd represents the client command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "trigger a deploy of a service",
	Long:  `deploy and stuff`,
	Run: func(cmd *cobra.Command, args []string) {
		host := fmt.Sprintf("%s:%s", viper.GetString("serverAddr"), viper.GetString("serverPort"))
		deployService(host)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
