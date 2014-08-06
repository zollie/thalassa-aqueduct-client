package aqueduct

import (
	"fmt"
	"testing"
)

const apiUrl = "http://localhost:10000"

// Test Constructor
func TestNewClient(t *testing.T) {
	_, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

// Test PUT /backends/key
func TestPutBackend(t *testing.T) {
	client, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	mem1 := Member{
		Id:   "memberID",
		Host: "192.168.0.13",
		Port: 8081,
	}

	back := &Backend{
		Type:    "static",
		Mode:    "tcp",
		Members: []Member{mem1},
	}

	re := client.PutBackend("test-back", back)
	if re != nil {
		t.Error(re)
		t.Fail()
	}
}

// Test GET /frontends
func TestGetBackEnds(t *testing.T) {
	client, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	backs, re := client.GetBackends()
	if re != nil {
		t.Error(re)
		t.Fail()
	}

	for _, b := range backs {
		fmt.Printf("%v", b)
	}
}

// Test PUT /frontends/key
func TestPutFrontend(t *testing.T) {
	client, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	front := &Frontend{
		Type:    "",
		Bind:    "*:4001",
		Backend: "testbackname",
		Mode:    "tcp",
	}

	re := client.PutFrontend("testfront", front)
	if re != nil {
		t.Error(re)
		t.Fail()
	}
}

// Test GET /frontends
func TestGetFrontEnds(t *testing.T) {
	client, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	fronts, re := client.GetFrontends()
	if re != nil {
		t.Error(re)
		t.Fail()
	}

	for _, f := range fronts {
		fmt.Printf("%v", f)
	}
}

// Test DELETE /backends/key
func TestDeleteBackend(t *testing.T) {
	client, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	re := client.DeleteBackend("test-back")
	if re != nil {
		t.Error(re)
		t.Fail()
	}
}

// Test DELETE /frontends/key
func TestDeleteFrontend(t *testing.T) {
	client, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	re := client.DeleteFrontend("testfront")
	if re != nil {
		t.Error(re)
		t.Fail()
	}
}

// Test GET /haproxy/config
func TestGetHAProxyConfig(t *testing.T) {
	client, err := NewClient(apiUrl)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	config, re := client.GetHAProxyConfig()
	if re != nil {
		t.Error(re)
		t.Fail()
	}

	fmt.Printf("Config is:\n%s", config)
}
