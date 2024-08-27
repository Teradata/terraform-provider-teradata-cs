package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetEnvironments() (*[]Environment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/environments", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	environments := []Environment{}
	err = json.Unmarshal(body, &environments)
	if err != nil {
		return nil, err
	}
	return &environments, nil
}

func (c *Client) CreateEnvironment(env EnvironmentCreateRequest) (*Environment, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/environments", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	reqbody, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(reqbody))
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	environment := Environment{}
	err = json.Unmarshal(body, &environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

func (c *Client) GetEnvironment(envName string) (*Environment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/environments/%s", c.HostURL, envName), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	environment := Environment{}
	err = json.Unmarshal(body, &environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

func (c *Client) UpdateEnvironment(envName string, operation string) (*Environment, error) {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/environments/%s", c.HostURL, envName), nil)

	if err != nil {
		return nil, err
	}
	postBody, err := json.Marshal(map[string]string{
		"operation": operation,
	})
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	environment := Environment{}
	err = json.Unmarshal(body, &environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

func (c *Client) DeleteEnvironment(envName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/environments/%s", c.HostURL, envName), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}
