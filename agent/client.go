package main

import (
	"encoding/json"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshhttp "github.com/cloudfoundry/bosh-utils/httpclient"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/cppforlife/turbulence/tasks"
)

type Client struct {
	clientRequest clientRequest
}

func NewClient(endpoint string, httpClient boshhttp.HTTPClient, logger boshlog.Logger) Client {
	clientRequest := clientRequest{
		endpoint:   endpoint,
		httpClient: httpClient,
		logger:     logger,
	}

	return Client{clientRequest: clientRequest}
}

func (c Client) FetchTasks(agentID string) ([]tasks.Task, error) {
	var resp []tasks.Task

	// todo use query string
	// todo rename to agent_tasks
	path := fmt.Sprintf("/api/v1/agents/%s/tasks", agentID)

	err := c.clientRequest.Post(path, nil, &resp)
	if err != nil {
		return resp, bosherr.WrapErrorf(err, "Fetching tasks '%s'", agentID)
	}

	return resp, nil
}

func (c Client) FetchTaskState(taskID string) (tasks.StateResponse, error) {
	var resp tasks.StateResponse

	path := fmt.Sprintf("/api/v1/agent_tasks/%s/state", taskID)

	err := c.clientRequest.Get(path, &resp)
	if err != nil {
		return resp, bosherr.WrapErrorf(err, "Fetching task '%s' state", taskID)
	}

	return resp, nil
}

func (c Client) RecordTaskResult(taskID string, err error) error {
	var resp interface{}

	path := fmt.Sprintf("/api/v1/agent_tasks/%s", taskID)
	req := tasks.ResultRequest{}

	if err != nil {
		req.Error = err.Error()
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling task")
	}

	err = c.clientRequest.Post(path, bytes, &resp)
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating task '%s'", taskID)
	}

	return nil
}
