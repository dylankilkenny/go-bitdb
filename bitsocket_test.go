package bitdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSocket(t *testing.T) {
	bitsocket := NewSocket(3, "https://bitsocket.org/s/")
	assert.Equal(t, bitsocket.BitDB.Version, 3)
	assert.Equal(t, bitsocket.BaseURL, "https://bitsocket.org/s/")
}

func TestStream(t *testing.T) {
	bitsocket := NewSocket(3, "https://bitsocket.org/s/")
	events, err := bitsocket.Stream("", ".[] | .out[0] | {amount: .e.v}")
	if err != nil {
		t.Error(err)
	}
	for event := range events {
		_, ok := event.Data.(map[string]interface{})
		if !ok {
			t.Error("Response assertion error")
		}
		break
	}
}

func TestRawBitsocket(t *testing.T) {
	var bitquery = []byte(`
		{
			"v": 3,
			"q": {
				"find": {}
			},
			"r": {
				"f": ".[] | .out[0] | {amount: .e.v}"
			}
		}
	`)
	events, err := RawBitsocket(bitquery, "https://bitsocket.org/s/")
	if err != nil {
		t.Error(err)
	}
	for event := range events {
		_, ok := event.Data.(map[string]interface{})
		if !ok {
			t.Error("Response assertion error")
		}
		break
	}
}
