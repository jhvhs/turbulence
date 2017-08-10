package main

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/jhvhs/turbulence/director"
	"github.com/jhvhs/turbulence/incident"
	"github.com/jhvhs/turbulence/incident/reporter"
	"github.com/jhvhs/turbulence/scheduledinc"
	"github.com/jhvhs/turbulence/tasks"
)

type Repos struct {
	incidentsRepo          incident.Repo
	scheduledIncidentsRepo scheduledinc.Repo
	tasksRepo              tasks.Repo
}

func NewRepos(
	uuidGen boshuuid.Generator,
	reporter reporter.Reporter,
	director director.Director,
	incidentNotifier incident.RepoNotifier,
	scheduledIncidentNotifier scheduledinc.RepoNotifier,
	logger boshlog.Logger,
) (Repos, error) {
	tasksRepo := tasks.NewRepo(logger)

	incidentsRepo := incident.NewRepo(
		uuidGen,
		incidentNotifier,
		reporter,
		director,
		tasksRepo,
		logger,
	)

	scheduledIncidentsRepo := scheduledinc.NewRepo(
		uuidGen,
		scheduledIncidentNotifier,
		incidentsRepo,
		logger,
	)

	return Repos{incidentsRepo, scheduledIncidentsRepo, tasksRepo}, nil
}

func (r Repos) IncidentsRepo() incident.Repo              { return r.incidentsRepo }
func (r Repos) ScheduledIncidentsRepo() scheduledinc.Repo { return r.scheduledIncidentsRepo }
func (r Repos) TasksRepo() tasks.Repo                     { return r.tasksRepo }
