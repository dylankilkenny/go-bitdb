package bitdb

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// A Request contains the apikey, url and BitDB object
type Request struct {
	apiKey  string
	BaseURL string
	BitDB   BitDB
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

type Response struct {
	Confirmed   interface{} `json:"c"`
	Unconfirmed interface{} `json:"u"`
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
	request.BaseURL = url
	request.BitDB.Version = version
	return request
}

// Request queries a bitdb node and returns the result
func (b *Request) Request(query interface{}, jq string) (*Response, error) {
	b64Query, err := buildQuery(b, query, jq)
	if err != nil {
		log.Fatal(err)
	}
	body, err := fetch(b.BaseURL+b64Query, b.apiKey)
	if err != nil {
		fmt.Println(err)
	}
	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON ERROR: ", err)
	}
	return &response, nil
}

// RawRequest accepts a bitquery object as a json string
func (b *Request) RawRequest(bitquery []byte) (*Response, error) {
	b64Query := b64.StdEncoding.EncodeToString(bitquery)
	body, err := fetch(b.BaseURL+b64Query, b.apiKey)
	if err != nil {
		fmt.Println(err)
	}
	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON ERROR: ", err)
	}
	return &response, nil
}

// TxHash returns a struct for querying a transaction by hash
func (b *Request) TxHash(prevTxID string) TxHash {
	return TxHash{Hash: prevTxID}
}

func fetch(url, apiKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if apiKey != "" {
		req.Header.Set("key", apiKey)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 304 {
		return nil, errors.New("Auth: API key required")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body, nil
}

func buildQuery(b *Request, query interface{}, jq string) (string, error) {
	b.BitDB.Query.Find = query
	b.BitDB.Response.Function = jq
	j, err := json.Marshal(b.BitDB)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return b64.StdEncoding.EncodeToString([]byte(j)), nil
}
