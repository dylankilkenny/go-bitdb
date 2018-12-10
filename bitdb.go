package bitdb

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// A Request contains the apikey, url and BitDB object
type Request struct {
	apiKey string
	url    string
	BitDB  BitDB
}

type query struct {
	Find interface{} `json:"find"`
}

type response struct {
	Function string `json:"f"`
}

// BitDB contains the bitdb query object
//
//This object will be transformed into json, then base 64 and used
// as the bitdb endpoint url
type BitDB struct {
	Version  int      `json:"v"`
	Query    query    `json:"q"`
	Response response `json:"r"`
}

// TxHash is a sample struct which can be passed to the query field
// on the BitDB object. Instantiate TxHash with a transaction hash,
// set it on BitDB.Query, and the bitdb node will return the transaction object.
type TxHash struct {
	Hash string `json:"tx.h"`
}

// New builds and returns a Connection
func New(version int, apiKey string, url string) *Request {
	request := new(Request)
	request.apiKey = apiKey
	request.url = url
	request.BitDB.Version = version
	return request
}

// Request queries a bitdb node and returns the result
func (b *Request) Request(query interface{}, jq string) (string, error) {
	b.BitDB.Query.Find = query
	b.BitDB.Response.Function = jq
	j, err := json.Marshal(b.BitDB)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	b64Query := b64.StdEncoding.EncodeToString([]byte(j))
	req, err := http.NewRequest("GET", b.url+b64Query, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("key", b.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body), nil
}

// RawRequest accepts a bitquery object as a json string
func RawRequest(bitquery []byte, url string, apiKey string) (string, error) {
	b64Query := b64.StdEncoding.EncodeToString(bitquery)
	req, err := http.NewRequest("GET", url+b64Query, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("key", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body), nil
}
