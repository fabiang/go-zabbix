package zabbix_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/fabiang/go-zabbix"
	"github.com/fabiang/go-zabbix/types"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalEvent(t *testing.T) {
	data := []byte(`{
	"eventid": "12345",
	"acknowledged": "1",
	"source": "1",
	"object": "2",
	"objectid": "678",
	"value": "3",
	"value_changed": "true",
	"hosts": [
		{
			"hostid": "234324"
		}
	],
	"clock": "1234567890",
	"ns": "98765"
}`)
	var event zabbix.Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "12345", event.EventID)
	assert.True(t, bool(event.Acknowledged))
	assert.True(t, bool(event.ValueChanged))
	assert.Equal(t, 1, event.Source)
	assert.Equal(t, 2, event.ObjectType)
	assert.Equal(t, 678, event.ObjectID)
	assert.Equal(t, 3, event.Value)
	assert.Equal(t, int64(1234567890), time.Time(event.Timestamp).Unix())
	assert.Equal(t, int(98765), time.Time(event.Timestamp).Nanosecond())
	assert.Equal(t, 1, len(event.Hosts))
	assert.Equal(t, "234324", event.Hosts[0].HostID)
}

func TestMarshalEvent(t *testing.T) {
	expected := []byte(`{
"eventid": "12345",
"acknowledged": "0",
"source": "0",
"object": "0",
"objectid": "0",
"value": "0",
"value_changed": "0",
"hosts": null,
"clock": "1234567890",
"ns": "98765"
}`)

	timestamp := time.Unix(1234567890, 98765)
	var event zabbix.Event = zabbix.Event{
		EventID:   "12345",
		Timestamp: types.ZBXUnixTimestamp(timestamp),
	}
	data, err := json.MarshalIndent(&event, "", "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(expected), string(data))
}
