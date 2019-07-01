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
	"os"

	sharedfuncs "github.com/fiveateooate/deployinator/internal/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//type messageHandler func(context.Context, pubsub.Message)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deployinator",
	Short: "Deplopyinate all the things",
	Long: `Deployinator
	deploy things to things
	
	hopefully modulular enough to allow
	different ways to deploy things
	to different things`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/deployinator/deployconfig.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		defaultconfig := "/etc/deployinator/deployconfig.yaml"
		if sharedfuncs.FileExists(defaultconfig) {
			viper.SetConfigFile(defaultconfig)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	deployCmd.Flags().String("server-addr", "127.0.0.1", "server address")
	viper.BindPFlag("serverAddr", deployCmd.Flags().Lookup("server-addr"))
	deployCmd.Flags().Int("server-port", 9091, "server port")
	viper.BindPFlag("serverPort", deployCmd.Flags().Lookup("server-port"))
	deployCmd.Flags().String("cid", "cu1", "cluster id")
	viper.BindPFlag("cid", deployCmd.Flags().Lookup("cid"))
	deployCmd.Flags().String("cenv", "local", "cluster env")
	viper.BindPFlag("cenv", deployCmd.Flags().Lookup("cenv"))
}
