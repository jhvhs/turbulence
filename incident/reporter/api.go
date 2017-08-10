package reporter

import (
	"fmt"
	"html/template"
	"time"
)

type EventResponse struct {
	event *Event

	ID   string
	Type string

	Instance EventInstanceResp

	ExecutionStartedAt   string
	ExecutionCompletedAt string

	Error string
}

type EventInstanceResp struct {
	ID         string
	Group      string
	Deployment string
	AZ         string
}

func NewEventResponse(event *Event) EventResponse {
	var completedAt string

	if (event.ExecutionCompletedAt != time.Time{}) {
		completedAt = event.ExecutionCompletedAt.Format(time.RFC3339)
	}

	return EventResponse{
		event: event,

		ID:   event.ID,
		Type: event.Type,

		Instance: EventInstanceResp{
			ID:         event.Instance.ID,
			Group:      event.Instance.Group,
			Deployment: event.Instance.Deployment,
			AZ:         event.Instance.AZ,
		},

		ExecutionStartedAt:   event.ExecutionStartedAt.Format(time.RFC3339),
		ExecutionCompletedAt: completedAt,

		Error: event.ErrorStr(),
	}
}

func (r EventResponse) IsAction() bool { return r.event.IsAction() }

func (r EventResponse) DescriptionHTML() template.HTML {
	content := ""

	if len(r.Instance.ID) > 0 {
		content = fmt.Sprintf("<span>Instance</span> %s/%s", r.Instance.Group, r.Instance.ID)
		content += " <span>Deployment</span> " + r.Instance.Deployment

		if len(r.Instance.AZ) > 0 {
			content += " <span>AZ</span> " + r.Instance.AZ
		}
	}

	return template.HTML(content)
}
