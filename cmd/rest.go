// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

	"github.com/scottshotgg/storage/server"
	"github.com/spf13/cobra"
)

// restCmd represents the REST server command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Storage REST server",
	Long:  "The server command that starts the Storage REST server",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.RunREST()
		if err != nil {
			// log

			// Print the error stack/trace to the console so that it is readable and formatted
			log.Printf("%v\n", err)
			os.Exit(9)
		}
	},
}

func init() {
	RootCmd.AddCommand(restCmd)
}
