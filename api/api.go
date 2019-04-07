package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/codec"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/keyerror"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	txbldr "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/gorilla/mux"
)

const (
	maxValidAccountValue = int(0x80000000 - 1)
	maxValidIndexalue    = int(0x80000000 - 1)
)

var cdc *codec.Codec

func init() {
	cdc = app.MakeCodec()
}

// Server represents the API server
type Server struct {
	Port   int    `json:"port"`
	KeyDir string `json:"key_dir"`

	Version string `yaml:"version,omitempty"`
	Commit  string `yaml:"commit,omitempty"`
	Branch  string `yaml:"branch,omitempty"`
}

// Router returns the router
func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/version", s.VersionHandler).Methods("GET")
	router.HandleFunc("/keys", s.GetKeys).Methods("GET")
	router.HandleFunc("/keys", s.PostKeys).Methods("POST")
	router.HandleFunc("/keys/{name}", s.GetKey).Methods("GET")
	router.HandleFunc("/keys/{name}", s.PutKey).Methods("PUT")
	router.HandleFunc("/keys/{name}", s.DeleteKey).Methods("DELETE")
	router.HandleFunc("/tx/sign", s.Sign).Methods("POST")

	return router
}

// SignBody is the body for a sign request
type SignBody struct {
	Tx            json.RawMessage `json:"tx"`
	Name          string          `json:"name"`
	Passphrase    string          `json:"passphrase"`
	ChainID       string          `json:"chain_id"`
	AccountNumber string          `json:"account_number"`
	Sequence      string          `json:"sequence"`
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

// VersionHandler handles the /version route
func (s *Server) VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(s.newVersion().marshal())
}

type version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
}

func (s *Server) newVersion() version {
	return version{s.Version, s.Commit, s.Branch}
}

func (v version) marshal() []byte {
	out, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return out
}

type restError struct {
	Error string `json:"error"`
}

func newError(err error) restError {
	return restError{err.Error()}
}

func (e restError) marshal() []byte {
	out, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return out
}

// GetKeys is the handler for the GET /keys
func (s *Server) GetKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	kb, err := keys.NewKeyBaseFromDir(s.KeyDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	infos, err := kb.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	if len(infos) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	keysOutput, err := ckeys.Bech32KeysOutput(infos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	out, err := json.Marshal(keysOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}

// AddNewKey is the necessary data for adding a new key
type AddNewKey struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Mnemonic string `json:"mnemonic"`
	Account  int    `json:"account,string,omitempty"`
	Index    int    `json:"index,string,omitempty"`
}

func (ak AddNewKey) marshal() []byte {
	out, err := json.Marshal(ak)
	if err != nil {
		panic(err)
	}
	return out
}

// PostKeys is the handler for the POST /keys
func (s *Server) PostKeys(w http.ResponseWriter, r *http.Request) {
	var m AddNewKey

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

	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	if m.Name == "" || m.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("must include both password and name with request")).marshal())
		return
	}

	// if mnemonic is empty, generate one
	mnemonic := m.Mnemonic
	if mnemonic == "" {
		_, mnemonic, _ = ckeys.NewInMemory().CreateMnemonic("inmemorykey", ckeys.English, "123456789", ckeys.Secp256k1)
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("invalid mnemonic")).marshal())
		return
	}

	if m.Account < 0 || m.Account > maxValidAccountValue {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("invalid account number")).marshal())
		return
	}

	if m.Index < 0 || m.Index > maxValidIndexalue {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("invalid index number")).marshal())
		return
	}

	_, err = kb.Get(m.Name)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("key %s already exists", m.Name)).marshal())
		return
	}

	account := uint32(m.Account)
	index := uint32(m.Index)
	info, err := kb.CreateAccount(m.Name, mnemonic, ckeys.DefaultBIP39Passphrase, m.Password, account, index)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	keyOutput, err := ckeys.Bech32KeyOutput(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	keyOutput.Mnemonic = mnemonic

	out, err := json.Marshal(keyOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}

// GetKey is the handler for the GET /keys/{name}
func (s *Server) GetKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	kb, err := keys.NewKeyBaseFromDir(s.KeyDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]
	bechPrefix := r.URL.Query().Get("bech")

	if bechPrefix == "" {
		bechPrefix = "acc"
	}

	bechKeyOut, err := getBechKeyOut(bechPrefix)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	info, err := kb.Get(name)
	if keyerror.IsErrKeyNotFound(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(newError(err).marshal())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	keyOutput, err := bechKeyOut(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	out, err := json.Marshal(keyOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}

type bechKeyOutFn func(keyInfo ckeys.Info) (ckeys.KeyOutput, error)

func getBechKeyOut(bechPrefix string) (bechKeyOutFn, error) {
	switch bechPrefix {
	case "acc":
		return ckeys.Bech32KeyOutput, nil
	case "val":
		return ckeys.Bech32ValKeyOutput, nil
	case "cons":
		return ckeys.Bech32ConsKeyOutput, nil
	}

	return nil, fmt.Errorf("invalid Bech32 prefix encoding provided: %s", bechPrefix)
}

// UpdateKeyBody update key password request REST body
type UpdateKeyBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

func (u UpdateKeyBody) marshal() []byte {
	out, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return out
}

// PutKey is the handler for the PUT /keys/{name}
func (s *Server) PutKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var m UpdateKeyBody

	kb, err := keys.NewKeyBaseFromDir(s.KeyDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	err = kb.Update(name, m.OldPassword, func() (string, error) { return m.NewPassword, nil })
	if keyerror.IsErrKeyNotFound(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(newError(err).marshal())
		return
	} else if keyerror.IsErrWrongPassword(err) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(newError(err).marshal())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

// DeleteKeyBody request
type DeleteKeyBody struct {
	Password string `json:"password"`
}

func (u DeleteKeyBody) marshal() []byte {
	out, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return out
}

// DeleteKey is the handler for the DELETE /keys/{name}
func (s *Server) DeleteKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var m DeleteKeyBody

	kb, err := keys.NewKeyBaseFromDir(s.KeyDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	err = kb.Delete(name, m.Password, false)
	if keyerror.IsErrKeyNotFound(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(newError(err).marshal())
		return
	} else if keyerror.IsErrWrongPassword(err) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(newError(err).marshal())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
