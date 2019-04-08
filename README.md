# Keyserver

This is a basic `keyserver` for cosmos-sdk applications. It contains the following routes:

```
GET     /version
GET     /keys
POST    /keys
GET     /keys/{name}?bech=acc
PUT     /keys/{name}
DELETE  /keys/{name}
POST    /tx/sign
```

First, build and start the server:

```bash
> make install
> keyserver config
> keyserver serve
```

Then you can use the included CLI to create keys, use the mnemonics to create them in `gaiacli` as well:

```bash
# Create a new key with generated mnemonic
> keyserver keys post jack foobarbaz | jq

# Save the mnemonic from the above command and add it to gaiacli
> gaiacli keys add jack --recover

# Next create a single node testnet
> gaiad init testing --chain-id testing
> gaiacli config chain-id testing
> gaiad add-genesis-account jack 1000000000stake
> gaiad gentx --name jack
> gaiad collect-gentxs
> gaiad start
```

In another window, generate the transaction to sign, sign it and broadcast:
```bash
> mkdir -p test_data
> gaiacli keys add jill
> gaiacli tx send $(gaiacli keys show jill -a) 10000stake --memo "sending the things" --gas auto --fees 100stake --from jack --generate-only > test_data/unsigned.json
> keyserver tx sign jack foobarbaz testing 0 1 test_data/unsigned.json > test_data/signed.json
> gaiacli tx broadcast test_data/signed.json
```
