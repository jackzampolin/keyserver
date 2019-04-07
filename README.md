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

`/tx/sign` is currently not working.


To test, setup a [single node testnet](https://cosmos.network/docs/gaia/deploy-testnet.html#single-node-local-manual-testnet), start the server and create a key using the following curl:

```bash
> make install
> keyserver config
> keyserver serve
> curl -XPOST -d "{\"name\":\"foo\",\"password\":\"foobarbaz\"}" localhost:3000/keys | jq
{
  "name": "foo",
  "type": "local",
  "address": "cosmos1k8jntawgxwsff70k9dfk8zk0p29lt3uxrwcwuy",
  "pubkey": "cosmospub1addwnpepqftzdupkclzgy3fcvrf0q9tmfm7j0dreg0w4pcz9ed056kjg2723uetjln0",
  "mnemonic": "ask year mother egg long monster bulb seminar make mother bomb gossip slab alter zoo mesh black deer property ritual own pool dinner near"
}
```

Keep the `mnemonic` to test with the CLI if necessary (restore it using `gaiacli keys add --restore`)

To test signing, modify the JSON in `keyserver_unsigned.json` for testing and post it to the server with the following curl:

```bash
> curl -XPOST -d "@./test_data/keyserver_unsigned.json" localhost:3000/tx/sign > ./test_data/keyserver_signed.json
> gaiacli tx broadcast ./test_data/keyserver_signed.json
Response:
  TxHash: 04B7B34DF69D3FF06382E9BEEBEAE63252B4B6AF8C209A3C10AE2E4661C7AF28
ERROR: {"codespace":"sdk","code":4,"message":"signature verification failed"}
```
