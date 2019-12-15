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

	// Routing for handling the folders
	a.Get(v1, "/folders/public", a.handleRequest(handler.GetRootPublicFolder))
	a.Get(v1, "/folders/{hashId}", a.handleRequest(handler.GetFolder))
	a.Get(v1, "/folders/{hashId}/path", a.handleRequest(handler.GetFolderPath))
	a.Get(v1, "/folders/{hashId}/files", a.handleRequest(handler.GetFolderFiles))
	a.Get(v1, "/folders/{hashId}/folders", a.handleRequest(handler.GetFolderFolders))
	a.Post(v1, "/folders", a.adminOnly(handler.CreateFolder))
	a.Put(v1, "/folders/{hashId}", a.adminOnly(handler.UpdateFolder))
	a.Delete(v1, "/folders/{hashId}", a.adminOnly(handler.DeleteFolder))

	// Routing for handling the files
	a.Get(v1, "/files/{hashId}", a.handleRequest(handler.GetFile))
	a.Get(v1, "/files/{hashId}/download", a.handleRequest(handler.ServeFile))
	a.Post(v1, "/files", a.adminOnly(handler.CreateFile))
	a.Delete(v1, "/files/{hashId}", a.adminOnly(handler.DeleteFile))
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
