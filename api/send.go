package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// BankSendBody contains the necessary data to make a send transaction
type BankSendBody struct {
	Sender        sdk.AccAddress `json:"sender"`
	Reciever      sdk.AccAddress `json:"reciever"`
	Amount        string         `json:"amount"`
	ChainID       string         `json:"chain-id"`
	Memo          string         `json:"memo,omitempty"`
	Fees          string         `json:"fees,omitempty"`
	GasAdjustment string         `json:"gas_adjustment,omitempty"`
}

func (sb BankSendBody) Marshal() []byte {
	out, err := json.Marshal(sb)
	if err != nil {
		panic(err)
	}
	return out
}

// BankSend handles the /tx/bank/send route
func (s *Server) BankSend(w http.ResponseWriter, r *http.Request) {
	var sb BankSendBody

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	err = cdc.UnmarshalJSON(body, &sb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	coins, err := sdk.ParseCoins(sb.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("failed to parse amount %s into sdk.Coins", sb.Amount)).marshal())
		return
	}

	var fees sdk.Coins
	if sb.Fees != "" {
		fees, err = sdk.ParseCoins(sb.Fees)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newError(fmt.Errorf("failed to parse fees %s into sdk.Coins", sb.Fees)).marshal())
			return
		}
	}

	stdTx := auth.NewStdTx(
		[]sdk.Msg{bank.NewMsgSend(sb.Sender, sb.Reciever, coins)},
		auth.NewStdFee(20000, fees),
		[]auth.StdSignature{{}},
		sb.Memo,
	)

	gas, err := s.SimulateGas(cdc.MustMarshalBinaryLengthPrefixed(stdTx))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	if gas != 0 && sb.GasAdjustment != "" {
		adj, err := strconv.ParseFloat(sb.GasAdjustment, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newError(fmt.Errorf("failed to parse gasAdjustment %d into float64", sb.GasAdjustment)).marshal())
			return
		}
		gas = uint64(adj * float64(gas))
	}

	stdTx = auth.NewStdTx(
		stdTx.Msgs,
		auth.NewStdFee(gas, stdTx.Fee.Amount),
		[]auth.StdSignature{},
		stdTx.Memo,
	)

	w.WriteHeader(http.StatusOK)
	w.Write(cdc.MustMarshalJSON(stdTx))
	return
}
