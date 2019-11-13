package app

import (
	"github.com/AdrianOrlow/files-api/app/handler"
	"net/http"
)

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the authentication
	a.Get("/oauth/google/login", a.handleRequest(handler.HandleGoogleLogin))
	a.Get("/oauth/google/callback", a.handleRequest(handler.HandleGoogleCallback))

	// Routing for handling the catalogs
	a.Get("/catalogs", a.handleRequest(handler.GetPublicCatalogs))
	a.Get("/catalogs/{hashId}", a.handleRequest(handler.GetCatalog))
	a.Get("/catalogs/{hashId}/files", a.handleRequest(handler.GetCatalogFiles))
	a.Get("/catalogs/{hashId}/catalogs", a.handleRequest(handler.GetCatalogCatalogs))
	a.Post("/catalogs", a.adminOnly(handler.CreateCatalog))
	a.Put("/catalogs/{hashId}", a.adminOnly(handler.UpdateCatalog))
	a.Delete("/catalogs/{hashId}", a.adminOnly(handler.DeleteCatalog))

	// Routing for handling the files
	a.Get("/files/{hashId}", a.handleRequest(handler.GetFile))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}
