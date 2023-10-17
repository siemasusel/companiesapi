package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/command"
	"github.com/siemasusel/companiesapi/internal/query"
	"github.com/siemasusel/companiesapi/internal/utils/chilog"
	"golang.org/x/exp/slog"
)

func NewRouter(logger slog.Handler, commands *command.Commands, queries *query.Queries) http.Handler {
	srv := newServer(commands, queries)
	r := chi.NewRouter()

	r.Use(chilog.NewStructuredLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)

	r.Route("/companies", func(r chi.Router) {
		r.Get("/{company_id}", srv.GetCompany)
		r.Post("/", srv.CreateCompany)
		r.Put("/{company_id}", srv.UpdateCompany)
		r.Delete("/{company_id}", srv.DeleteCompany)
	})

	return r
}

func getCompanyIDFromRequest(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, "company_id"))
}
