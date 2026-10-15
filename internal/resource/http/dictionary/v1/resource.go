package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/web-rabis/reader-workstation-api/internal/manager/auth"
	"github.com/web-rabis/reader-workstation-api/internal/manager/dictionary"
)

type Resource struct {
	path          string
	authMan       *auth.Manager
	dictionaryMan dictionary.IManager
}

func NewAdminDictionaryResource(basePath string, authMan *auth.Manager, dictionaryMan dictionary.IManager) *Resource {
	return &Resource{
		path:          basePath,
		authMan:       authMan,
		dictionaryMan: dictionaryMan,
	}
}

func (res *Resource) Path() string {
	return res.path
}

func (res *Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(res.authMan.JWTAuth()))
		r.Use(auth.NewUserAccessCtx(res.authMan.JWTKey()).ChiMiddleware)
		r.Get("/type-card/{id}/backgrounds", res.backgrounds)
		r.Get("/type-card", res.getTypeCard)
		r.Get("/academic-degree", res.getAcademicDegree)
		r.Get("/education", res.getEducation)
		r.Get("/faculty", res.getFaculty)
		r.Get("/institution", res.getInstitution)
		r.Get("/nationality", res.getNationality)
		r.Get("/sex", res.getSex)
		r.Get("/social-status", res.getSocialStatus)
		r.Get("/speciality", res.getSpeciality)
	})

	return r
}
