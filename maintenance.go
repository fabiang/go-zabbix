package zabbix

import (
	"errors"
	"strings"

	"github.com/fabiang/go-zabbix/types"
)

type MaintenanceType int
type TagsEvaltype int

var ErrMaintenanceHostNotFound = errors.New("Failed to find ID by host name")

const (
	withDataCollection MaintenanceType = iota
	withoutDataCollection

	and TagsEvaltype = iota * 2
	or

	Once int = iota
	EveryDay
	EveryWeek
	EveryMonth
)

type Maintenance struct {
	MaintenanceID       string                 `json:"maintenanceid"`
	Name                string                 `json:"name"`
	ActiveSince         types.ZBXUnixTimestamp `json:"active_since,string"`
	ActiveTill          types.ZBXUnixTimestamp `json:"active_till,string"`
	Description         string                 `json:"description"`
	Type                MaintenanceType        `json:"tags_evaltype,string"`
	ActionEvalTypeAndOr TagsEvaltype           `json:"maintenance_type,string"`
}

type MaintenanceGetParams struct {
	GetParameters

	// Sort the result by the given properties.
	// Possible values are: maintenanceid, name and maintenance_type.
	SortField []string `json:"sortfield,omitempty"`

	// Return the maintenance's time periods in the timeperiods property.
	SelectTimeperiods SelectQuery `json:"selectTimeperiods,omitempty"`

	// Return hosts assigned to the maintenance in the hosts property.
	SelectHosts SelectQuery `json:"selectHosts,omitempty"`

	// SelectGroups:Return host groups assigned to the maintenance in the groups property.
	// Deprecated: Deprecated since Zabbix 5.4 and removed in Zabbix 7.2
	SelectGroups SelectQuery `json:"selectGroups,omitempty"`

	// SelectHostGroups: Return host groups assigned to the maintenance in the groups property.
	SelectHostGroups SelectQuery `json:"selectHostGroups,omitempty"`

	// Return only maintenances with the given IDs.
	Maintenanceids []string `json:"maintenanceids,omitempty"`

	// Return only maintenances that are assigned to the given hosts.
	Hostids []string `json:"hostids,omitempty"`

	// Return only maintenances that are assigned to the given host groups.
	Groupids []string `json:"groupids,omitempty"`
}

type MaintenanceCreateParams struct {
	Maintenance

	Groupids []string `json:"groupids,omitempty"`
	// Hosts name
	HostNames   []string                 `json:"-"`
	HostIDs     []string                 `json:"hostids"`
	Timeperiods []MaintenanceTimeperiods `json:"timeperiods"`
	Tags        []string                 `json:"tags,omitempty"`
}

type MaintenanceTimeperiods struct {
	TimeperiodType int `json:"timeperiod_type,int"`
	Every          int `json:"every,string"`
	Dayofweek      int `json:"dayofweek,string"`
	StartTime      int `json:"start_time,string"`
	Period         int `json:"period,string"`
}

type MaintenanceCreateResponse struct {
	IDs []string `json:"maintenanceids"`
}

// GetMaintenance queries the Zabbix API for Maintenance matching the given search
// parameters.
func (s *Session) GetMaintenance(params *MaintenanceGetParams) ([]Maintenance, error) {
	maintenance := make([]Maintenance, 0)
	err := s.Get("maintenance.get", params, &maintenance)
	if err != nil {
		return nil, err
	}

	if len(maintenance) == 0 {
		return nil, ErrNotFound
	}

	return maintenance, nil
}

func (s *Session) CreateMaintenance(params *MaintenanceCreateParams) (response MaintenanceCreateResponse, err error) {
	if err = params.FillHostIDs(s); err != nil {
		return
	}

	err = s.Get("maintenance.create", params, &response)
	return
}

func (m *Maintenance) Delete(session *Session) error {
	ID := []string{m.MaintenanceID}
	response := make(map[string]interface{})
	if err := session.Get("maintenance.delete", ID, &response); err != nil {
		return err
	}
	return nil
}

func (m *MaintenanceCreateParams) FillHostIDs(session *Session) error {
	hosts, err := session.GetHosts(HostGetParams{})
	if err != nil {
		return err
	}

	err = ErrMaintenanceHostNotFound
	for _, name := range m.HostNames {
		for _, host := range hosts {
			if strings.ToUpper(strings.Trim(host.Hostname, " ")) == strings.ToUpper(strings.Trim(name, " ")) {
				m.HostIDs = append(m.HostIDs, host.HostID)

				err = nil
			}
		}
	}

	return err
}
