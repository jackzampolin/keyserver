// Copyright © 2018 Jack Zampolin <jack@blockstack.com>
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
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the server",
	Run: func(cmd *cobra.Command, args []string) {

		log.Println(fmt.Sprintf("Listening on port ':%v'...", server.Port))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", server.Port), handlers.LoggingHandler(os.Stdout, server.Router())))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
