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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jackzampolin/keyserver/api"
	"github.com/spf13/cobra"
)

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Runs transaction calls",
}

var bankCmd = &cobra.Command{
	Use:   "bank",
	Short: "bank transactions",
}

var broadcastCmd = &cobra.Command{
	Use:   "broadcast [file]",
	Short: "broadcast a signed transaction",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txData, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal("error reading transaction file")
		}
		url := fmt.Sprintf("http://localhost:%d/tx/broadcast", server.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(txData))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code %d, error: %s", resp.StatusCode, string(out))
			return
		}
		fmt.Println(string(out))
	},
}

// versionCmd represents the version command
var sendCmd = &cobra.Command{
	Use:   "send [sender] [reciever] [amount] [chain-id] [memo] [fees] [gas-adjustment]",
	Short: "generate a send transaction",
	Args:  cobra.RangeArgs(4, 7),
	Run: func(cmd *cobra.Command, args []string) {
		var bs api.BankSendBody
		send, err := sdk.AccAddressFromBech32(args[0])
		if err != nil {
			return
		}
		recv, err := sdk.AccAddressFromBech32(args[1])
		if err != nil {
			return
		}

		switch len(args) {
		case 4:
			bs = api.BankSendBody{
				Sender:   send,
				Reciever: recv,
				Amount:   args[2],
				ChainID:  args[3],
			}
		case 5:
			bs = api.BankSendBody{
				Sender:   send,
				Reciever: recv,
				Amount:   args[2],
				ChainID:  args[3],
				Memo:     args[4],
			}
		case 6:
			bs = api.BankSendBody{
				Sender:   send,
				Reciever: recv,
				Amount:   args[2],
				ChainID:  args[3],
				Memo:     args[4],
				Fees:     args[5],
			}
		case 7:
			bs = api.BankSendBody{
				Sender:        send,
				Reciever:      recv,
				Amount:        args[2],
				ChainID:       args[3],
				Memo:          args[4],
				Fees:          args[5],
				GasAdjustment: args[6],
			}
		}
		url := fmt.Sprintf("http://localhost:%d/tx/bank/send", server.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs.Marshal()))
		if err != nil {
			log.Fatalf("error fetching %s", url)
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

// /keys GET
var txSign = &cobra.Command{
	Use:   "sign [name] [password] [chain-id] [account-number] [sequence] [tx-file]",
	Args:  cobra.ExactArgs(6),
	Short: "Sign a transaction",
	Run: func(cmd *cobra.Command, args []string) {
		txData, err := ioutil.ReadFile(args[5])
		if err != nil {
			log.Fatal("error reading transaction file")
		}

		postData := api.SignBody{
			Name:          args[0],
			Passphrase:    args[1],
			ChainID:       args[2],
			AccountNumber: args[3],
			Sequence:      args[4],
			Tx:            txData,
		}

		url := fmt.Sprintf("http://localhost:%d/tx/sign", server.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData.Marshal()))
		if err != nil {
			log.Fatalf("error fetching %s", url)
			return
		}
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading response body")
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("non 200 respose code %d, error: %s", resp.StatusCode, string(out))
			return
		}
		fmt.Println(string(out))
	},
}

func init() {
	txCmd.AddCommand(txSign)
	txCmd.AddCommand(bankCmd)
	txCmd.AddCommand(broadcastCmd)
	bankCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(txCmd)
}
