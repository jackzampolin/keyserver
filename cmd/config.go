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
	"os"

	"github.com/jackzampolin/keyserver/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// versionCmd represents the version command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Sets a default config file",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("Error finding homedir:", err)
			return
		}

		s := api.Server{
			Port:   3000,
			KeyDir: fmt.Sprintf("%s/.keyserver", home),
			Node:   "http://localhost:26657",
		}

		if _, err := os.Stat(s.KeyDir); os.IsNotExist(err) {
			err := os.MkdirAll(s.KeyDir, 0777)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		}

		conf := fmt.Sprintf("%s/config.yaml", s.KeyDir)
		if _, err := os.Stat(conf); os.IsNotExist(err) {
			out, err := yaml.Marshal(s)
			if err != nil {
				fmt.Println("Error marshaling config:", err)
				return
			}
			file, err := os.Create(conf)
			if err != nil {
				fmt.Println("Error creating config file:", err)
				return
			}
			defer file.Close()
			fmt.Fprintf(file, string(out))
		} else {
			fmt.Println("Config file already exists, skipping...")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
