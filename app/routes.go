package app

import (
	"github.com/AdrianOrlow/files-api/app/handler"
	"github.com/gorilla/mux"
	"net/http"
)

// setRouters sets the all required routers
func (a *App) setRouters() {
	v1 := a.Router.PathPrefix("/v1").Subrouter()

	// Routing for handling the authentication
	a.Get(v1, "/oauth/google/login", a.handleRequest(handler.HandleGoogleLogin))
	a.Get(v1, "/oauth/google/callback", a.handleRequest(handler.HandleGoogleCallback))

	// Routing for handling the catalogs
	a.Get(v1, "/catalogs", a.handleRequest(handler.GetRootCatalogs))
	a.Get(v1, "/catalogs/{hashId}", a.handleRequest(handler.GetCatalog))
	a.Get(v1, "/catalogs/{hashId}/path", a.handleRequest(handler.GetCatalogPath))
	a.Get(v1, "/catalogs/{hashId}/files", a.handleRequest(handler.GetCatalogFiles))
	a.Get(v1, "/catalogs/{hashId}/catalogs", a.handleRequest(handler.GetCatalogCatalogs))
	a.Post(v1, "/catalogs", a.adminOnly(handler.CreateCatalog))
	a.Put(v1, "/catalogs/{hashId}", a.adminOnly(handler.UpdateCatalog))
	a.Delete(v1, "/catalogs/{hashId}", a.adminOnly(handler.DeleteCatalog))

	// Routing for handling the files
	a.Get("/files/{hashId}", a.handleRequest(handler.GetFile))
}

// Get wraps the router for GET method
func (a *App) Get(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("DELETE")
}
