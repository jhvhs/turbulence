package incident

import (
	"fmt"
	"strings"
	"time"

	"github.com/jhvhs/turbulence/incident/reporter"
	"github.com/jhvhs/turbulence/incident/selector"
	"github.com/jhvhs/turbulence/tasks"
)

type Request struct {
	Tasks    tasks.OptionsSlice
	Selector selector.Request
}

type Response struct {
	incident Incident

	ID string

	Tasks    tasks.OptionsSlice
	Selector selector.Request

	ExecutionStartedAt   string
	ExecutionCompletedAt string

	Events []reporter.EventResponse

	description string
}

type IncidentsResp []Response

func NewResponses(incidents []Incident) IncidentsResp {
	resp := []Response{}

	for _, incid := range incidents {
		resp = append(resp, NewResponse(incid))
	}

	return resp
}

func NewResponse(incident Incident) Response {
	var eventResps []reporter.EventResponse

	for _, event := range incident.Events().Events() {
		eventResps = append(eventResps, reporter.NewEventResponse(event))
	}

	var completedAt string

	if (incident.ExecutionCompletedAt() != time.Time{}) {
		completedAt = incident.ExecutionCompletedAt().Format(time.RFC3339)
	}

	return Response{
		incident: incident,

		ID: incident.ID(),

		Tasks:    incident.Tasks,
		Selector: incident.Selector,

		ExecutionStartedAt:   incident.ExecutionStartedAt().Format(time.RFC3339),
		ExecutionCompletedAt: completedAt,

		Events: eventResps,
	}
}

func (r Response) URL() string { return fmt.Sprintf("/incidents/%s", r.ID) }

func (r Response) TaskTypes() string { return strings.Join(r.incident.TaskTypes(), ", ") }

func (r Response) Description() (string, error) { return r.incident.Description() }

func (r Response) HasEventErrors() bool {
	for _, eventResp := range r.Events {
		if len(eventResp.Error) > 0 {
			return true
		}
	}

	return false
}
