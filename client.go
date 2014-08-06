package aqueduct

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

var (
	ErrNotFound = errors.New("Not found")
	ErrKeyReq   = errors.New("Key is required")
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
}

func NewClient(host string) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	httpClient := newHTTPClient(u)
	return &Client{u, httpClient}, nil
}

func (client *Client) GetFrontends() ([]Frontend, error) {
	data, err := client.doRequest("GET", "/frontends", nil)
	if err != nil {
		return nil, err
	}
	ret := []Frontend{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (client *Client) GetFrontendByKey(key string) (*Frontend, error) {
	if key == "" {
		return nil, ErrKeyReq
	}

	data, err := client.doRequest("GET", "/frontends/"+key, nil)
	if err != nil {
		return nil, err
	}
	ret := &Frontend{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (client *Client) GetBackends() ([]Backend, error) {
	data, err := client.doRequest("GET", "/backends", nil)
	if err != nil {
		return nil, err
	}
	ret := []Backend{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (client *Client) GetBackendByKey(key string) (*Backend, error) {
	if key == "" {
		return nil, ErrKeyReq
	}

	data, err := client.doRequest("GET", "/backends/"+key, nil)
	if err != nil {
		return nil, err
	}
	ret := &Backend{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (client *Client) PutBackend(key string, back *Backend) error {
	if key == "" {
		return ErrKeyReq
	}

	b, merr := json.Marshal(back)
	if merr != nil {
		return merr
	}

	_, rerr := client.doRequest("PUT", "/backends/"+key, b)
	if rerr != nil {
		return rerr
	}

	return nil
}

func (client *Client) UpdateBackend(key string, back *Backend) error {
	if key == "" {
		return ErrKeyReq
	}

	b, merr := json.Marshal(back)
	if merr != nil {
		return merr
	}

	_, rerr := client.doRequest("POST", "/backends/"+key, b)
	if rerr != nil {
		return rerr
	}

	return nil
}

func (client *Client) DeleteBackend(key string) error {
	if key == "" {
		return ErrKeyReq
	}

	_, err := client.doRequest("DELETE", "/backends/"+key, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) PutFrontend(key string, front *Frontend) error {
	if key == "" {
		return ErrKeyReq
	}

	b, merr := json.Marshal(front)
	if merr != nil {
		return merr
	}

	_, rerr := client.doRequest("PUT", "/frontends/"+key, b)
	if rerr != nil {
		return rerr
	}

	return nil
}

func (client *Client) DeleteFrontend(key string) error {
	if key == "" {
		return ErrKeyReq
	}

	_, err := client.doRequest("DELETE", "/frontends/"+key, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) GetHAProxyConfig() (string, error) {
	path := "/haproxy/config"

	data, err := client.doRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	ret := string(data)
	return ret, nil
}

// Lifted from DockerClient
func (client *Client) doRequest(method string, path string, body []byte) ([]byte, error) {
	log.Printf("Sending request with body: %s", body)

	b := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, client.URL.String()+path, b)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, ErrNotFound
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %s", resp.Status, data)
	}
	return data, nil
}

// Lifted from DockerClient
func newHTTPClient(u *url.URL) *http.Client {
	httpTransport := &http.Transport{}
	if u.Scheme == "unix" {
		socketPath := u.Path
		unixDial := func(proto string, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		}
		httpTransport.Dial = unixDial
		// Override the main URL object so the HTTP lib won't complain
		u.Scheme = "http"
		u.Host = "unix.sock"
	}
	u.Path = ""
	return &http.Client{Transport: httpTransport}
}
