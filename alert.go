package zabbix

import (
	"github.com/fabiang/go-zabbix/types"
)

const (
	// AlertTypeMessage indicates that an Alert is a notification message.
	AlertTypeMessage = iota

	// AlertTypeRemoteCommand indicates that an Alert is a remote command call.
	AlertTypeRemoteCommand
)

const (
	// AlertMessageStatusNotSent indicates that an Alert of type
	// AlertTypeMessage has not been sent yet.
	AlertMessageStatusNotSent = iota

	// AlertMessageStatusSent indicates that an Alert of type AlertTypeMessage
	// has been sent successfully.
	AlertMessageStatusSent

	// AlertMessageStatusFailed indicates that an Alert of type AlertTypeMessage
	// failed to send.
	AlertMessageStatusFailed
)

const (
	// AlertCommandStatusRun indicates that an Alert of type
	// AlertTypeRemoteCommand has been run.
	AlertCommandStatusRun = 1 + iota

	// AlertCommandStatusAgentUnavailable indicates that an Alert of type
	// AlertTypeRemoteCommand failed to run as the Zabbix Agent was unavailable.
	AlertCommandStatusAgentUnavailable
)

// Alert represents a Zabbix Alert returned from the Zabbix API.
//
// See: https://www.zabbix.com/documentation/2.2/manual/config/notifications
type Alert struct {
	// AlertID is the unique ID of the Alert.
	AlertID string `json:"alertid"`

	// ActionID is the unique ID of the Action that generated this Alert.
	ActionID string `json:"actionid"`

	// AlertType is the type of the Alert.
	// AlertType must be one of the AlertType constants.
	AlertType int `json:"alerttype,string"`

	// Timestamp is the UTC timestamp at which the Alert was generated.
	Timestamp types.ZBXUnixTimestamp `json:"clock"`

	// ErrorText is the error message if there was a problem sending a message
	// or running a remote command.
	ErrorText string `json:"error"`

	// EscalationStep is the escalation step during which the Alert was
	// generated.
	EscalationStep int `json:"esc_step,string"`

	// EventID is the unique ID of the Event that triggered this Action that
	// generated this Alert.
	EventID string `json:"eventid"`

	// MediaTypeID is the unique ID of the Media Type that was used to send this
	// Alert if the AlertType is AlertTypeMessage.
	MediaTypeID string `json:"mediatypeid"`

	// Message is the Alert message body if AlertType is AlertTypeMessage.
	Message string `json:"message"`

	// RetryCount is the number of times Zabbix tried to send a message.
	RetryCount int `json:"retries,string"`

	// Recipient is the end point address of a message if AlertType is
	// AlertTypeMessage.
	Recipient string `json:"sendto"`

	// Status indicates the outcome of executing the Alert.
	//
	// If AlertType is AlertTypeMessage, Status must be one of the
	// AlertMessageStatus constants.
	//
	// If AlertType is AlertTypeRemoteCommand, Status must be one of the
	// AlertCommandStatus constants.
	Status int `json:"status,string"`

	// Subject is the Alert message subject if AlertType is AlertTypeMessage.
	Subject string `json:"subject"`

	// UserID is the unique ID of the User the Alert message was sent to.
	UserID string `json:"userid"`

	// Hosts is an array of Hosts that triggered this Alert.
	//
	// Hosts is only populated if AlertGetParams.SelectHosts is given in the
	// query parameters that returned this Alert.
	Hosts []Host `json:"hosts"`
}

// AlertGetParams is query params for alert.get call
type AlertGetParams struct {
	GetParameters

	// SelectHosts causes all Hosts which triggered the Alert to be attached in
	// the search results.
	SelectHosts SelectQuery `json:"selectHosts,omitempty"`

	// SelectMediaTypes causes the Media Types used for the Alert to be attached
	// in the search results.
	SelectMediaTypes SelectQuery `json:"selectMediatypes,omitempty"`

	// SelectUsers causes all Users to which the Alert was addressed to be
	// attached in the search results.
	SelectUsers SelectQuery `json:"selectUsers,omitempty"`
}

// GetAlerts queries the Zabbix API for Alerts matching the given search
// parameters.
//
// ErrNotFound is returned if the search result set is empty.
// An error is returned if a transport, parsing or API error occurs.
func (c *Session) GetAlerts(params AlertGetParams) ([]Alert, error) {
	alerts := make([]Alert, 0)
	err := c.Get("alert.get", params, &alerts)
	if err != nil {
		return nil, err
	}

	if len(alerts) == 0 {
		return nil, ErrNotFound
	}

	return alerts, nil
}
