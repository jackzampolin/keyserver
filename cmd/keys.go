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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jackzampolin/keyserver/api"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Runs keys calls",
}

// /keys GET
var keysGet = &cobra.Command{
	Use:   "get",
	Short: "Fetch all keys managed by the keyserver",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys", server.Port)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code")
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(string(out))
	},
}

// /keys POST
var keysPost = &cobra.Command{
	Use:   "post [name] [password] [mnemonic]",
	Args:  cobra.RangeArgs(2, 3),
	Short: "Add a new key to the keyserver, optionally pass a mnemonic to restore the key",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys", server.Port)
		var addNP api.AddNewKey
		if len(args) == 2 {
			addNP = api.AddNewKey{Name: args[0], Password: args[1]}
		} else if len(args) == 3 {
			addNP = api.AddNewKey{Name: args[0], Password: args[1], Mnemonic: args[2]}
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(addNP.Marshal()))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code")
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(string(out))
	},
}

// /keys/{name} GET
var keyGet = &cobra.Command{
	Use:   "show [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Fetch details for one key",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys/%s", server.Port, args[0])
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code")
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(string(out))
	},
}

// /keys/{name} PUT
var keyPut = &cobra.Command{
	Use:   "put [name] [oldpass] [newpass]",
	Args:  cobra.ExactArgs(3),
	Short: "Update the password on a key",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys/%s", server.Port, args[0])
		kb := api.UpdateKeyBody{OldPassword: args[1], NewPassword: args[2]}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(kb.Marshal()))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 204 {
			log.Fatalf("non 204 respose code %d", resp.StatusCode)
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(out)
	},
}

// /keys/{name} DELETE
var keyDelete = &cobra.Command{
	Use:   "delete [name] [password]",
	Args:  cobra.ExactArgs(2),
	Short: "Delete a key",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost:%d/keys/%s", server.Port, args[0])
		kb := api.DeleteKeyBody{Password: args[1]}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(kb.Marshal()))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code %d", resp.StatusCode)
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		fmt.Println(out)
	},
}

func init() {
	keysCmd.AddCommand(keysGet)
	keysCmd.AddCommand(keysPost)
	keysCmd.AddCommand(keyGet)
	keysCmd.AddCommand(keyPut)
	keysCmd.AddCommand(keyDelete)
	rootCmd.AddCommand(keysCmd)
}
