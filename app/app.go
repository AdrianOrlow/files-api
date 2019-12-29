package app

import (
	"github.com/AdrianOrlow/files-api/app/utils"
	"github.com/gorilla/handlers"
	"log"
	"net/http"

	"github.com/AdrianOrlow/files-api/app/handler"
	"github.com/AdrianOrlow/files-api/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	err := a.InitializeDatabase(config)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.Initialize(config)
	if err != nil {
		log.Fatal(err)
	}

	handler.InitializeAuth(config)

	a.Router = mux.NewRouter()
	a.setRouters()
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

// Run the app on it's router
func (a *App) Run(host string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Print("Listening on " + host)
	log.Fatal(http.ListenAndServe(host, handlers.CORS(headersOk, originsOk, methodsOk)(a.Router)))
}

func (a *App) handleRequest(h RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(a.DB, w, r)
	}
}
