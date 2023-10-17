package http

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/siemasusel/companiesapi/internal/api"
	"github.com/siemasusel/companiesapi/internal/command"
	"github.com/siemasusel/companiesapi/internal/domain"
	"github.com/siemasusel/companiesapi/internal/query"
	"github.com/siemasusel/companiesapi/internal/utils/apperr"
	"github.com/siemasusel/companiesapi/internal/utils/httperr"
)

type server struct {
	commands *command.Commands
	queries  *query.Queries
}

func newServer(commands *command.Commands, queries *query.Queries) *server {
	return &server{
		commands: commands,
		queries:  queries,
	}
}

func (s *server) CreateCompany(w http.ResponseWriter, r *http.Request) {
	payload := api.CreateCompany{}
	if err := render.Decode(r, &payload); err != nil {
		httperr.RespondWithError(apperr.NewBadRequest(err, "invalid request"), w, r)
		return
	}

	copmanyType, err := api.UnmarshalCompanyType(api.CompanyType(payload.Type))
	if err != nil {
		httperr.RespondWithError(apperr.NewInvalidArgument(err, "invalid company type '%s'", payload.Type), w, r)
		return
	}

	var description string
	if payload.Description != nil {
		description = *payload.Description
	}

	id, err := s.commands.CreateCompany(r.Context(), payload.Name, description, payload.EmployeesCount, payload.Registered, copmanyType)
	if err != nil {
		httperr.RespondWithError(err, w, r)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, api.IdResponse{Id: id.String()})
}

func (s *server) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	payload := api.UpdateCompany{}
	if err := render.Decode(r, &payload); err != nil {
		httperr.RespondWithError(apperr.NewBadRequest(err, "invalid request"), w, r)
		return
	}

	id, err := getCompanyIDFromRequest(r)
	if err != nil {
		httperr.RespondWithError(apperr.NewBadRequest(err, "invalid company id"), w, r)
		return
	}

	var companyType *domain.CompanyType
	if payload.Type != nil {
		ct, err := api.UnmarshalCompanyType(api.CompanyType(*payload.Type))
		if err != nil {
			httperr.RespondWithError(apperr.NewInvalidArgument(err, "invalid company type '%s'", *payload.Type), w, r)
			return
		}

		companyType = &ct
	}

	err = s.commands.UpdateCompany(r.Context(), id, payload.Name, payload.Description, payload.EmployeesCount, payload.Registered, companyType)
	if err != nil {
		httperr.RespondWithError(err, w, r)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (s *server) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, err := getCompanyIDFromRequest(r)
	if err != nil {
		httperr.RespondWithError(apperr.NewBadRequest(err, "invalid company id"), w, r)
		return
	}

	err = s.commands.DeleteCompany(r.Context(), id)
	if err != nil {
		httperr.RespondWithError(err, w, r)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (s *server) GetCompany(w http.ResponseWriter, r *http.Request) {
	id, err := getCompanyIDFromRequest(r)
	if err != nil {
		httperr.RespondWithError(apperr.NewBadRequest(err, "invalid company id"), w, r)
		return
	}

	company, err := s.queries.ReadCompany(r.Context(), id)
	if err != nil {
		httperr.RespondWithError(err, w, r)
		return
	}

	render.Respond(w, r, company)
}
