package api

import (
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func (s *Server) Broadcast(w http.ResponseWriter, r *http.Request) {
	var stdTx auth.StdTx
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	err = cdc.UnmarshalJSON(body, &stdTx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	res, err := rpcclient.NewHTTP(s.Node, "/websocket").BroadcastTxAsync(txBytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(cdc.MustMarshalJSON(sdk.NewResponseFormatBroadcastTx(res)))
	return
}
