package bitdb

import (
	"bufio"
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Connection contains the base url and BitDB object
type Connection struct {
	BaseURL string
	BitDB   BitDB
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// NewSocket builds and returns a Connection
func NewSocket(version int, url string) *Connection {
	request := new(Connection)
	request.BaseURL = url
	request.BitDB.Version = version
	return request
}

// Stream opens a connection to a bitsocket url and recieves a stream of server sent events
func (bs *Connection) Stream(query interface{}, jq string) (events chan *Event, err error) {
	bs.BitDB.Query.Find = query
	bs.BitDB.Response.Function = jq
	j, err := json.Marshal(bs.BitDB)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	b64Query := b64.StdEncoding.EncodeToString([]byte(j))
	events, err = stream(bs.BaseURL + b64Query)
	return events, err
}

// RawBitsocket opens a connection to a bitsocket url and recieves a stream of server sent events
func RawBitsocket(bitquery []byte, url string) (events chan *Event, err error) {
	b64Query := b64.StdEncoding.EncodeToString(bitquery)
	events, err = stream(url + b64Query)
	return events, err

}

func stream(url string) (events chan *Event, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("got response status code %d", resp.StatusCode)
	}
	events = make(chan *Event)
	var buf bytes.Buffer
	go func() {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "error during resp.Body read:%s\n", err)
				close(events)
			}
			switch {
			// event data
			case bytes.HasPrefix(line, []byte("data:")):
				buf.Write(line[6:])
			// end of event
			case bytes.Equal(line, []byte("\n")):
				b := buf.Bytes()
				if bytes.HasPrefix(b, []byte("{")) {
					buf.Reset()
					ev := Event{}
					err = json.Unmarshal(b, &ev)
					if err != nil {
						fmt.Println("JSON ERROR: ", err)
					}
					if ev.Type == "mempool" {
						events <- &ev
					}
				}
			default:
				fmt.Fprintf(os.Stderr, "Error: len:%d\n%s", len(line), line)
				close(events)
			}
		}
	}()
	return events, nil
}
