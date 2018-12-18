// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger

	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "storage",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	}
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		// log

		// Print the error stack/trace to the console so that it is readable and formatted
		log.Printf("%v\n", err)
		os.Exit(9)
	}
}

func init() {
	var err error

	if os.Getenv("ENV") == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			// log

			// Print the error stack/trace to the console so that it is readable and formatted
			log.Printf("%v\n", err)
			os.Exit(9)
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			// log

			// Print the error stack/trace to the console so that it is readable and formatted
			log.Printf("%v\n", err)
			os.Exit(9)
		}
	}

	RootCmd.PersistentFlags().String("servicename", "", "Name of the service that is running")
	RootCmd.PersistentFlags().String("server-ip", "0.0.0.0", "IP the server will listen on")
	RootCmd.PersistentFlags().String("rest-port", "", "Port the rest server will listen on")
	RootCmd.PersistentFlags().String("rpc-addr", "", "Port the rpc server will listen on")
	RootCmd.PersistentFlags().String("rpc-port", "", "Port the rpc server will listen on")
	RootCmd.PersistentFlags().String("scheme", "http", "Scheme that the server will use; http for local, https in k8s")
	RootCmd.PersistentFlags().String("tls-certificate", "", "TLS certificate the server will use")
	RootCmd.PersistentFlags().String("tls-key", "", "TLS key the server will use")

	_ = viper.BindPFlag("servicename", RootCmd.PersistentFlags().Lookup("servicename"))
	_ = viper.BindPFlag("server-ip", RootCmd.PersistentFlags().Lookup("server-ip"))
	_ = viper.BindPFlag("rest-port", RootCmd.PersistentFlags().Lookup("rest-port"))
	_ = viper.BindPFlag("rpc-addr", RootCmd.PersistentFlags().Lookup("rpc-port"))
	_ = viper.BindPFlag("rpc-port", RootCmd.PersistentFlags().Lookup("rpc-port"))
	_ = viper.BindPFlag("scheme", RootCmd.PersistentFlags().Lookup("scheme"))
	_ = viper.BindPFlag("tls-certificate", RootCmd.PersistentFlags().Lookup("tls-certificate"))
	_ = viper.BindPFlag("tls-key", RootCmd.PersistentFlags().Lookup("tls-key"))
}
