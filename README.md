# go-bitdb
Using go-bitdb you can query a bitdb node and connect to a bitsocket node in golang.

## BitDB


[BitDb](https://bitdb.network/) is an autonomous database that continuously synchronizes itself with Bitcoin.
BitDB stores every bitcoin transaction in a structured document format that can be queried against like a regular database.
With a simple MongoDB query, anyone can easily query, filter, and build powerful decentralized applications on Bitcoin.

## Bitsocket

[Bitsocket](https://bitsocket.org/) is a bitcoin notification service. By constructing a [bitquery](https://docs.bitdb.network/docs/query_v3), you can subscribe to the bitcoin blockchain and recieve transactions that meet a certain criteria.

## BitDb Usage
The following code creates a BitDb instance and queries a bitdb node for a transaction by its TxID.
The jq string is the processing function which will return an object containing the amount in the first output of the transaction

``` go
version := 3 // API version
bitdbURL := "https://bitdb.network/q/" // API url
apiKey := "qq54zc33pttdp6l8ycnnj99ahan8a2hfrygqyz0fc3"
BitDb := bitdb.New(version, bitdbURL, apiKey) // Create new instance

type TxHash struct {
  Hash string `json:"tx.h"`
}
txHash := TxHash{Hash: "ffff5c6d0660068381b26fe3546eb2a51faf1a0a1a707db1ca32a5b168a7301b"}
jq := ".[] | .out[0] | {amount: .e.v}"
response, err := BitDb.Request(txHash, jq)
if err != nil {
  fmt.Println(err)
}
fmt.Println(string(response))
```
Output:
``` json
{"u":[],"c":{"amount":7051}}
```
Rather than using structs you can query with a json string:
```go
var bitquery = []byte(`
{
  "v": 3,
  "q": {
    "find": {
      "tx.h": "ffff5c6d0660068381b26fe3546eb2a51faf1a0a1a707db1ca32a5b168a7301b"
    }
  },
  "r": {
    "f": ".[] | .out[0] | {amount: .e.v}"
  }
}
`)
response, err := BitDb.RawRequest(bitquery)
if err != nil {
  fmt.Println(err)
}
fmt.Println(string(response))
```
## Bitsocket Usage
