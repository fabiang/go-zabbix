# go-zabbix

Go bindings for the Zabbix API

[![go report card](https://goreportcard.com/badge/github.com/fabiang/go-zabbix "go report card")](https://goreportcard.com/report/github.com/fabiang/go-zabbix)
[![GPL license](https://img.shields.io/badge/license-GPL-brightgreen.svg)](https://opensource.org/licenses/gpl-license)
[![GoDoc](https://godoc.org/github.com/fabiang/go-zabbix?status.svg)](https://godoc.org/github.com/fabiang/go-zabbix)

## Overview

This project provides bindings to interoperate between programs written in Go
language and the Zabbix monitoring API.

A number of Zabbix API bindings already exist for Go with varying levels of
maturity. This project aims to provide an alternative implementation which is
stable, fast, and allows for loose typing (using types such as`interface{}` or
`map[string]interface{}`) as well as strong types (such as `Host` or `Event`).

The package aims to have comprehensive coverage of Zabbix API methods from v1.8
through to v7.4 without introducing limitations to the native API methods.

## Fork

Currently maintained fork of https://github.com/zabbix-tools/go-zabbix

New Features:

* Support for Zabbix JSONRPC API 4.0 - 7.4
* Support for host interfaces
* More info on hosts
* Support for proxies
* Allow executing scripts on Zabbix Server

## Getting started

Get the package:

```
go get "github.com/fabiang/go-zabbix"
```

```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/fabiang/go-zabbix"
)

func main() {
	// Default approach - without session caching
	session, err := zabbix.NewSession("http://zabbix/api_jsonrpc.php", "Admin", "zabbix")
	if err != nil {
		panic(err)
	}

	version, err := session.GetVersion()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to Zabbix API v%s", version)
}
```

### Use session builder with caching.

You can use own cache by implementing SessionAbstractCache interface.
Optionally an http.Client can be passed to the builder, allowing to skip TLS verification, pass proxy settings, etc.

```go
func main() {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true
			}
		}
	}

	cache := zabbix.NewSessionFileCache().SetFilePath("./zabbix_session")
	session, err := zabbix.CreateClient("http://zabbix/api_jsonrpc.php").
		WithCache(cache).
		WithHTTPClient(client).
		WithCredentials("Admin", "zabbix").
		Connect()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	version, err := session.GetVersion()

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("Connected to Zabbix API v%s", version)
}
```

## Running the tests

### Unit tests
Running the unit tests:

```bash
go test -v -short "./..."
# or:
make unittests
```

### Integration tests

To run the integration tests against a specific Zabbix Server version, you'll need Docker. Then start the containers:

```bash
export ZBX_VERSION=6.4
docker compose up -d
# server should be running in a minute
# run tests:
go test -v -run Integration "./..."
# or:
make integration
```

## License

Released under the [GNU GPL License](https://github.com/fabiang/go-zabbix/blob/main/LICENSE)
