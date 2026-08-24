package zabbix

import (
	"bytes"
	"encoding/json"

	"github.com/fabiang/go-zabbix/types"
)

const (
	// ProxyStatusActive active proxy
	ProxyStatusActive = 5
	// ProxyStatusPassive passive proxy
	ProxyStatusPassive = 6

	// ProxyTLSConnectUnencryped connect unencrypted to or from host
	ProxyTLSConnectUnencryped = 1

	// ProxyTLSConnectPSK connect with PSK to or from host
	ProxyTLSConnectPSK = 2

	// ProxyTLSConnectCertificate connect with certificate to or from host
	ProxyTLSConnectCertificate = 4
)

// Proxy Proxy infomation returned from Zabbix API
type Proxy struct {
	ProxyID     string `json:"proxyid"`
	Host        string `json:"host"`
	Status      int    `json:"status,string"`
	Description string `json:"description"`

	// How should we connect to proxy
	TLSConnect int `json:"tls_connect,string"`

	// What type of connections we accept from proxy
	TLSAccept int `json:"tls_accept,string"`

	TLSIssuer  string `json:"tls_issuer"`
	TLSSubject string `json:"tls_subject"`

	// Note those were removed in Zabbix 5.4
	TLSPSKIdentity string `json:"tls_psk_identity"`
	TLSPSK         string `json:"tls_psk"`

	ProxyAddresses types.ZBXProxyAddresses `json:"proxy_address"`

	Interface ProxyInterface `json:"interface,omitempty"`
}

// ProxyInterface Interface information for a Zabbix proxy
// There is no way to query this directly, query for a Zabbix proxy and use SelectInterface
type ProxyInterface struct {
	// (readonly) ID of the interface.
	InterfaceID string `json:"interfaceid"`

	// ID of the host the interface belongs to.
	HostID string `json:"hostid"`

	// DNS name used by the interface.
	DNS string `json:"dns"`

	// IP address used by the interface.
	IP string `json:"ip"`

	// Whether the connection should be made via IP.
	UseIP types.ZBXBoolean `json:"useip"`

	Port int `json:"port,string"`
}

// Handle PHPs empty associative arrays are arrays when turning them into JSON
func (t *ProxyInterface) UnmarshalJSON(in []byte) error {
	if bytes.Equal(in, []byte("[]")) {
		return nil
	}

	m := (*ProxyInterface)(t)
	return json.Unmarshal(in, m)
}

type ProxyGetParams struct {
	GetParameters

	ProxyIDs []string `json:"proxyids,omitempty"`

	SelectHosts     SelectQuery `json:"selectHosts,omitempty"`
	SelectInterface SelectQuery `json:"selectInterface,omitempty"`
}

// GetProxies queries the Zabbix API for Proxies matching the given search
// parameters.
//
// ErrEventNotFound is returned if the search result set is empty.
// An error is returned if a transport, parsing or API error occurs.
func (c *Session) GetProxies(params ProxyGetParams) ([]Proxy, error) {
	proxies := make([]Proxy, 0)
	err := c.Get("proxy.get", params, &proxies)
	if err != nil {
		return nil, err
	}

	if len(proxies) == 0 {
		return nil, ErrNotFound
	}

	return proxies, nil
}
