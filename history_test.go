package zabbix_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/fabiang/go-zabbix"
	"github.com/fabiang/go-zabbix/types"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalHistory(t *testing.T) {
	data := []byte(`{
	"itemid": "12345",
	"value": "3",
	"logeventid": "345",
	"severity": "321",
	"source": "639",
	"clock": "1234567890",
	"ns": "98765"
}`)
	var history zabbix.History
	err := json.Unmarshal(data, &history)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 12345, history.ItemID)
	assert.Equal(t, "3", history.Value)
	assert.Equal(t, int64(1234567890), time.Time(history.Timestamp).Unix())
	assert.Equal(t, int(98765), time.Time(history.Timestamp).Nanosecond())
	assert.Equal(t, 345, history.LogEventID)
	assert.Equal(t, 321, history.Severity)
	assert.Equal(t, "639", history.Source)
}

func TestMarshalHistory(t *testing.T) {
	expected := []byte(`{
"itemid": "0",
"value": "thevalue",
"clock": "1234567890",
"ns": "98765"
}`)

	timestamp := time.Unix(1234567890, 98765)
	var history zabbix.History = zabbix.History{
		Value:     "thevalue",
		Timestamp: types.ZBXUnixTimestamp(timestamp),
	}
	data, err := json.MarshalIndent(&history, "", "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(expected), string(data))
}
