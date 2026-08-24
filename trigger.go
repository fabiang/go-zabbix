package zabbix

import (
	"github.com/fabiang/go-zabbix/types"
)

const (
	// TriggerAlarmStateOK means a normal trigger state. Called FALSE in older Zabbix versions.
	TriggerAlarmStateOK = iota

	// TriggerAlarmStateProblem normally means that something happened.
	TriggerAlarmStateProblem
)

const (
	// TriggerStateNormal means normal trigger state
	TriggerStateNormal = iota

	// TriggerStateUnknown means unknown trigger state
	TriggerStateUnknown
)

const (
	// TriggerSeverityNotClassified is Not classified severity
	TriggerSeverityNotClassified = iota

	// TriggerSeverityInformation is Information severity
	TriggerSeverityInformation

	// TriggerSeverityWarning is Warning severity
	TriggerSeverityWarning

	// TriggerSeverityAverage is Average severity
	TriggerSeverityAverage

	// TriggerSeverityHigh is High severity
	TriggerSeverityHigh

	// TriggerSeverityDisaster is Disaster severity
	TriggerSeverityDisaster
)

// Trigger represents a Zabbix Trigger returned from the Zabbix API.
//
// See: https://www.zabbix.com/documentation/3.4/manual/config/triggers
type Trigger struct {
	// TriggerID is the ID of the Trigger.
	TriggerID string `json:"triggerid"`

	// AlarmState shows whether the trigger is in OK or problem state.
	//
	// AlarmState must be one of the TriggerAlarmState constants.
	AlarmState int `json:"value,string"`

	// Description is the name of the trigger.
	Description string `json:"description"`

	// Enabled shows whether the trigger is enabled or disabled.
	Enabled types.ZBXBoolean `json:"status"`

	// Expression is the trigger expression
	Expression string `json:"expression"`

	// Hosts is an array of Hosts that the trigger belongs.
	//
	// Hosts is only populated if TriggerGetParams.SelectHosts is given in the
	// query parameters that returned this Trigger.
	Hosts []Host `json:"hosts"`

	// Groups is an array of Hostgroups that the trigger belongs.
	//
	// Groups is only populated if TriggerGetParams.SelectGroups is given in the
	// query parameters that returned this Trigger.
	Groups []Hostgroup `json:"groups"`

	// LastChange is the time when the trigger last changed its state.
	LastChange int `json:"lastchange,string"`

	// Severity of the trigger.
	//
	// Severity must be one of the TriggerSeverity constants.
	Severity int `json:"priority,string"`

	// State of the trigger.
	//
	// State must be one of the TriggerState constants.
	State int `json:"state,string"`

	// Tags is an array of trigger tags
	//
	// Tags is only populated if TriggerGetParams.SelectTags is given in the
	// query parameters that returned this Trigger.
	Tags []TriggerTag `json:"tags"`

	// LastEvent is the latest event for the trigger
	//
	// LastEvent is only populated if TriggerGetParams.SelectLastEvent is set
	LastEvent *Event `json:"lastEvent"`

	// URL is a link to the trigger graph in Zabbix
	URL string `json:"url"`
}

// TriggerTag is trigger tag
type TriggerTag struct {
	Name  string `json:"tag"`
	Value string `json:"value"`
}

// TriggerGetParams is params for trigger.get query
type TriggerGetParams struct {
	GetParameters

	// TriggerIDs filters search results to Triggers that matched the given Trigger
	// IDs.
	TriggerIDs []string `json:"triggerids,omitempty"`

	// GroupIDs filters search results to triggers for hosts that are members of
	// the given Group IDs.
	GroupIDs []string `json:"groupids,omitempty"`

	TemplateIDs []string `json:"templateids,omitempty"`

	// HostIDs filters search results to triggers for hosts that matched the given
	// Host IDs.
	HostIDs []string `json:"hostids,omitempty"`

	ItemIDs []string `json:"itemids,omitempty"`

	ApplicationIDs []string `json:"applicationids,omitempty"`

	Functions []string `json:"functions,omitempty"`

	Group string `json:"group,omitempty"`

	Host string `json:"host,omitempty"`

	// InheritedOnly filters search results to triggers which have been
	// inherited from a template.
	InheritedOnly types.ZBXBoolean `json:"inherited,omitempty"`

	TemplatedOnly types.ZBXBoolean `json:"templated,omitempty"`

	MonitoredOnly types.ZBXBoolean `json:"monitored,omitempty"`

	ActiveOnly types.ZBXBoolean `json:"active,omitempty"`

	MaintenanceOnly types.ZBXBoolean `json:"maintenance,omitempty"`

	WithUnacknowledgedEventsOnly types.ZBXBoolean `json:"withUnacknowledgedEvents,omitempty"`

	WithAcknowledgedEventsOnly types.ZBXBoolean `json:"withAcknowledgedEvents,omitempty"`

	WithLastEventUnacknowledgedOnly types.ZBXBoolean `json:"withLastEventUnacknowledged,omitempty"`

	SkipDependent types.ZBXBoolean `json:"skipDependent,omitempty"`

	// LastChangeSince timestamp `json:"lastChangeSince,omitempty"`

	// LastChangeTill timestamp `json:"lastChangeTill,omitempty"`

	RecentProblemOnly types.ZBXBoolean `json:"only_true,omitempty"`

	MinSeverity int `json:"min_severity,omitempty"`

	ExpandComment types.ZBXBoolean `json:"expandComment,omitempty"`

	ExpandDescription types.ZBXBoolean `json:"expandDescription,omitempty"`

	ExpandExpression types.ZBXBoolean `json:"expandExpression,omitempty"`

	// SelectGroups causes all Hostgroups which contain the object that caused each
	// Trigger to be attached in the search results.
	//
	// Deprecated: Depcreated since Zabbix 6.4 and removed in Zabbix 7.2
	SelectGroups SelectQuery `json:"selectGroups,omitempty"`

	// SelectHosts causes all Hosts which contain the object that caused each
	// Trigger to be attached in the search results.
	SelectHosts SelectQuery `json:"selectHosts,omitempty"`

	SelectItems SelectQuery `json:"selectItems,omitempty"`

	SelectFunctions SelectQuery `json:"selectFunctions,omitempty"`

	SelectDependencies SelectQuery `json:"selectDependencies,omitempty"`

	SelectDiscoveryRule SelectQuery `json:"selectDiscoveryRule,omitempty"`

	SelectLastEvent SelectQuery `json:"selectLastEvent,omitempty"`

	SelectTags SelectQuery `json:"selectTags,omitempty"`
}

// GetTriggers queries the Zabbix API for Triggers matching the given search
// parameters.
//
// ErrTriggerNotFound is returned if the search result set is empty.
// An error is returned if a transport, parsing or API error occurs.
func (c *Session) GetTriggers(params TriggerGetParams) ([]Trigger, error) {
	triggers := make([]Trigger, 0)
	err := c.Get("trigger.get", params, &triggers)
	if err != nil {
		return nil, err
	}

	if len(triggers) == 0 {
		return nil, ErrNotFound
	}

	return triggers, nil
}
