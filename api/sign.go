package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	txbldr "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

// SignBody is the body for a sign request
type SignBody struct {
	Tx            json.RawMessage `json:"tx"`
	Name          string          `json:"name"`
	Passphrase    string          `json:"passphrase"`
	ChainID       string          `json:"chain_id"`
	AccountNumber string          `json:"account_number"`
	Sequence      string          `json:"sequence"`
}

// Marshal returns the json byte representation of the sign body
func (sb SignBody) Marshal() []byte {
	out, err := json.Marshal(sb)
	if err != nil {
		panic(err)
	}
	return out
}

// StdSignMsg returns a StdSignMsg from a SignBody request
func (sb SignBody) StdSignMsg() (stdSign txbldr.StdSignMsg, stdTx auth.StdTx, err error) {
	err = cdc.UnmarshalJSON(sb.Tx, &stdTx)
	if err != nil {
		return
	}
	acc, err := strconv.ParseInt(sb.AccountNumber, 10, 64)
	if err != nil {
		return
	}

	seq, err := strconv.ParseInt(sb.Sequence, 10, 64)
	if err != nil {
		return
	}

	stdSign = txbldr.StdSignMsg{
		Memo:          stdTx.Memo,
		Msgs:          stdTx.Msgs,
		ChainID:       sb.ChainID,
		AccountNumber: uint64(acc),
		Sequence:      uint64(seq),
		Fee: auth.StdFee{
			Amount: stdTx.Fee.Amount,
			Gas:    uint64(stdTx.Fee.Gas),
		},
	}

	return
}

// Sign handles the /tx/sign route
func (s *Server) Sign(w http.ResponseWriter, r *http.Request) {
	var m SignBody

	kb, err := keys.NewKeyBaseFromDir(s.KeyDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	err = cdc.UnmarshalJSON(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	stdSign, stdTx, err := m.StdSignMsg()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	sigBytes, pubkey, err := kb.Sign(m.Name, m.Passphrase, sdk.MustSortJSON(cdc.MustMarshalJSON(stdSign)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	sigs := append(stdTx.GetSignatures(), auth.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	})

	signedStdTx := auth.NewStdTx(stdTx.GetMsgs(), stdTx.Fee, sigs, stdTx.GetMemo())
	out, err := cdc.MarshalJSON(signedStdTx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}
