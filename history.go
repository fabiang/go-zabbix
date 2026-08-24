package zabbix

import (
	"encoding/json"
	"time"

	"github.com/fabiang/go-zabbix/types"
)

// History represents a Zabbix History returned from the Zabbix API.
//
// See: https://www.zabbix.com/documentation/4.0/manual/api/reference/history/object
type History struct {
	// ItemID is the ID of the related item.
	ItemID int `json:"itemid,string"`

	// Value is the received value.
	// Possible types: 0 - float; 1 - character; 2 - log; 3 - int; 4 - text;
	Value string `json:"value"`

	// LogEventID is the Windows event log entry ID.
	LogEventID int `json:"logeventid,string,omitempty"`

	// Severity is the Windows event log entry level.
	Severity int `json:"severity,string,omitempty"`

	// Source is the Windows event log entry source.
	Source string `json:"source,omitempty"`

	// Mapped Clock/Nanoseconds from JSON response
	Timestamp types.ZBXUnixTimestamp `json:"-"`
}

func (h *History) UnmarshalJSON(data []byte) error {
	type Alias History

	var aux struct {
		Alias
		Clock       int64 `json:"clock,string"`
		Nanoseconds int64 `json:"ns,string"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*h = History(aux.Alias)
	ts := types.ZBXUnixTimestamp(time.Unix(aux.Clock, aux.Nanoseconds))
	h.Timestamp = ts
	return nil
}

func (h History) MarshalJSON() ([]byte, error) {
	type Alias History

	var clock, nanoseconds int64

	time := time.Time(h.Timestamp)

	clock = time.Unix()
	nanoseconds = int64(time.Nanosecond())

	return json.Marshal(&struct {
		Alias
		Clock       int64 `json:"clock,string"`
		Nanoseconds int64 `json:"ns,string"`
	}{
		Alias:       Alias(h),
		Clock:       clock,
		Nanoseconds: nanoseconds,
	})
}

type HistoryGetParams struct {
	GetParameters

	// History object types to return
	// Possible values: 0 - numeric float, 1 - character, 2 - log,
	// 3 - numeric signed, 4, text
	// Default: 3
	History int `json:"history"`

	// HistoryIDs filters search results to histories with the given History ID's.
	HistoryIDs []string `json:"historyids,omitempty"`

	// ItemIDs filters search results to histories belong to the hosts
	// of the given Item ID's.
	ItemIDs []string `json:"itemids,omitempty"`

	// Return only values that have been received after or at the given time.
	TimeFrom float64 `json:"time_from,omitempty"`

	// Return only values that have been received before or at the given time.
	TimeTill float64 `json:"time_till,omitempty"`
}

// GetHistories queries the Zabbix API for Histories matching the given search
// parameters.
//
// ErrEventNotFound is returned if the search result set is empty.
// An error is returned if a transport, parsing or API error occurs.
func (c *Session) GetHistories(params HistoryGetParams) ([]History, error) {
	histories := make([]History, 0)
	err := c.Get("history.get", params, &histories)
	if err != nil {
		return nil, err
	}

	if len(histories) == 0 {
		return nil, ErrNotFound
	}

	return histories, nil
}
