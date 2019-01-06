package bitdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	bitdb := New(3, "apikey", "https://bitdb.network/q/")
	assert.Equal(t, bitdb.BitDB.Version, 3)
	assert.Equal(t, bitdb.apiKey, "apikey")
	assert.Equal(t, bitdb.BaseURL, "https://bitdb.network/q/")
}

func TestBuildQuery(t *testing.T) {
	bitdb := New(3, "apikey", "https://bitdb.network/q/")
	b64Query, err := buildQuery(bitdb, "test", "test")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, b64Query, "eyJ2IjozLCJxIjp7ImZpbmQiOiJ0ZXN0In0sInIiOnsiZiI6InRlc3QifX0=")
}

func TestRequest(t *testing.T) {
	bitdb := New(3, "qrqm04uwrd0wwaguxea079h0znswwt3quuejvl6zd6", "https://bitdb.network/q/")
	resp, err := bitdb.Request(TxHash{Hash: "ff353d0f15d6c46fa1c6a884798b2947d266e2417788b5e0935b243a92cfb52b"}, ".[] | .out[0] | {amount: .e.v}")
	if err != nil {
		t.Error(err)
	}
	_, ok := resp.Confirmed.(map[string]interface{})
	if !ok {
		t.Error("Response assertion error")
	}
}

func TestRawRequest(t *testing.T) {
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
	bitdb := New(3, "qrqm04uwrd0wwaguxea079h0znswwt3quuejvl6zd6", "https://bitdb.network/q/")
	response, err := bitdb.RawRequest(bitquery)
	if err != nil {
		fmt.Println(err)
	}
	_, ok := response.Confirmed.(map[string]interface{})
	if !ok {
		t.Error("Response assertion error")
	}

}
