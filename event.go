package zabbix

import (
	"encoding/json"
	"time"

	"github.com/fabiang/go-zabbix/types"
)

const (
	// EventSourceTrigger indicates that an Event was created by a Trigger.
	EventSourceTrigger = iota

	// EventSourceDiscoveryRule indicates than an Event was created by a
	// Discovery Rule.
	EventSourceDiscoveryRule

	// EventSourceAutoRegistration indicates that an Event was created by an
	// active Host registration rule.
	EventSourceAutoRegistration

	// EventSourceInternal indicates that an Event was created by an Internal
	// Event.
	EventSourceInternal
)

const (
	// EventObjectTypeTrigger indicates that an Event with Source type
	// EventSourceTrigger or EventSourceInternal is related to a Trigger.
	EventObjectTypeTrigger = iota

	// EventObjectTypeDiscoveredHost indicates that an Event with Source type
	// EventSourceDiscoveryRule is related to a discovered Host.
	EventObjectTypeDiscoveredHost

	// EventObjectTypeDiscoveredService indicates that an Event with Source type
	// EventSourceDiscoveryRule is related to a discovered Service.
	EventObjectTypeDiscoveredService

	// EventObjectTypeAutoRegisteredHost indicates that an Event with Source
	// type EventSourceAutoRegistration is related to an auto-registered Host.
	EventObjectTypeAutoRegisteredHost

	// EventObjectTypeItem indicates that an Event with Source type
	// EventSourceInternal is related to an Item.
	EventObjectTypeItem

	// EventObjectTypeLLDRule indicates that an Event with Source type
	// EventSourceInternal is related to a low-level Discovery Rule.
	EventObjectTypeLLDRule
)

const (
	// TriggerEventValueOK indicates that the Object related to an Event with
	// Source type EventSourceTrigger is in an "OK" state.
	TriggerEventValueOK = iota

	// TriggerEventValueProblem indicates that the Object related to an Event with
	// Source type EventSourceTrigger is in a "Problem" state.
	TriggerEventValueProblem
)

const (
	// DiscoveryEventValueUp indicates that the Host or Service related to an
	// Event with Source type EventSourceDiscoveryRule is in an "Up" state.
	DiscoveryEventValueUp = iota

	// DiscoveryEventValueDown indicates that the Host or Service related to an
	// Event with Source type EventSourceDiscoveryRule is in a "Down" state.
	DiscoveryEventValueDown

	// DiscoveryEventValueDiscovered indicates that the Host or Service related
	// to an Event with Source type EventSourceDiscoveryRule is in a
	// "Discovered" state.
	DiscoveryEventValueDiscovered

	// DiscoveryEventValueLost indicates that the Host or Service related to an
	// Event with Source type EventSourceDiscoveryRule is in a "Lost" state.
	DiscoveryEventValueLost
)

const (
	// InternalEventValueNormal indicates that the Object related to an Event
	// with Source type EventSourceInternal is in a "Normal" state.
	InternalEventValueNormal = iota

	// InternalEventValueNotSupported indicates that the Object related to an
	// Event with Source type EventSourceInternal is in an "Unknown" or
	// "Not supported" state.
	InternalEventValueNotSupported
)

// Event represents a Zabbix Event returned from the Zabbix API. Events are
// readonly as they may only be created by the Zabbix server.
//
// See: https://www.zabbix.com/documentation/2.2/manual/config/events
type Event struct {
	// EventID is the ID of the Event.
	EventID string `json:"eventid"`

	// Acknowledged indicates if the Event has been acknowledged by an operator.
	Acknowledged types.ZBXBoolean `json:"acknowledged"`

	// Source is the type of the Event source.
	//
	// Source must be one of the EventSource constants.
	Source int `json:"source,string"`

	// ObjectType is the type of the Object that is related to the Event.
	// ObjectType must be one of the EventObjectType constants.
	ObjectType int `json:"object,string"`

	// ObjectID is the unique identifier of the Object that caused this Event.
	ObjectID int `json:"objectid,string"`

	// Value is the state of the related Object.
	//
	// Value must be one of the EventValue constants, according to the Event's
	// Source type.
	Value int `json:"value,string"`

	// ValueChanges indicates if the state of the related Object has changed
	// since the previous Event.
	ValueChanged types.ZBXBoolean `json:"value_changed"`

	// Hosts is an array of Host which contained the Object which created this
	// Event.
	//
	// Hosts is only populated if EventGetParams.SelectHosts is given in the
	// query parameters that returned this Event and the Event Source is one of
	// EventSourceTrigger or EventSourceDiscoveryRule.
	Hosts []Host `json:"hosts"`

	// Mapped Clock/Nanoseconds from JSON response
	Timestamp types.ZBXUnixTimestamp `json:"-"`
}

func (h *Event) UnmarshalJSON(data []byte) error {
	type Alias Event

	var aux struct {
		Alias
		Clock       int64 `json:"clock,string"`
		Nanoseconds int64 `json:"ns,string"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*h = Event(aux.Alias)
	ts := types.ZBXUnixTimestamp(time.Unix(aux.Clock, aux.Nanoseconds))
	h.Timestamp = ts
	return nil
}

func (h Event) MarshalJSON() ([]byte, error) {
	type Alias Event

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

// EventGetParams is query params for event.get call
type EventGetParams struct {
	GetParameters

	// EventIDs filters search results to Events that matched the given Event
	// IDs.
	EventIDs []string `json:"eventids,omitempty"`

	// GroupIDs filters search results to events for hosts that are members of
	// the given Group IDs.
	GroupIDs []string `json:"groupids,omitempty"`

	// HostIDs filters search results to events for hosts that matched the given
	// Host IDs.
	HostIDs []string `json:"hostids,omitempty"`

	// ObjectIDs filters search results to events for Objects that matched
	// the given Object IDs.
	ObjectIDs []string `json:"objectids,omitempty"`

	// ObjectType filters search results to events created by the given Object
	// Type. Must be one of the EventObjectType constants.
	//
	// Default: EventObjectTypeTrigger
	ObjectType int `json:"object"`

	// AcknowledgedOnly filters search results to event which have been
	// acknowledged.
	AcknowledgedOnly bool `json:"acknowledged"`

	// MinEventID filters search results to Events with an ID greater or equal
	// to the given ID.
	MinEventID string `json:"eventid_from,omitempty"`

	// MaxEventID filters search results to Events with an ID lesser or equal
	// to the given ID.
	MaxEventID string `json:"eventid_till,omitempty"`

	// MinTime filters search results to Events with a timestamp lesser than or
	// equal to the given timestamp.
	MinTime int64 `json:"time_from,omitempty"`

	// MaxTime filters search results to Events with a timestamp greater than or
	// equal to the given timestamp.
	MaxTime int64 `json:"time_till,omitempty"`

	// Value filters search results to Events with the given values. Each value
	// must be one of the EventValue constants for the given ObjectType.
	Value []int `json:"value,omitempty"`

	// SelectHosts causes all Hosts which contain the object that caused each
	// Event to be attached in the search results.
	SelectHosts SelectQuery `json:"selectHosts,omitempty"`

	// SelectRelatedObject causes the object which caused each Event to be
	// attached in the search results.
	SelectRelatedObject SelectQuery `json:"selectRelatedObject,omitempty"`

	// SelectAlerts causes Alerts generated by each Event to be attached in the
	// search results.
	//
	// Deprecated: Depcreated since Zabbix 6.4 and removed in Zabbix 7.2
	SelectAlerts SelectQuery `json:"select_alerts,omitempty"`

	// SelectAcknowledgements causes Acknowledgments for each Event to be
	// attached in the search results in reverse chronological order.
	//
	// Deprecated: Depcreated since Zabbix 6.4 and removed in Zabbix 7.2
	SelectAcknowledgements SelectQuery `json:"select_acknowledges,omitempty"`
}

// GetEvents queries the Zabbix API for Events matching the given search
// parameters.
//
// ErrEventNotFound is returned if the search result set is empty.
// An error is returned if a transport, parsing or API error occurs.
func (c *Session) GetEvents(params EventGetParams) ([]Event, error) {
	events := make([]Event, 0)
	err := c.Get("event.get", params, &events)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, ErrNotFound
	}

	return events, nil
}
