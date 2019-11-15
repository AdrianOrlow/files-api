package app

import (
	"fmt"
	"github.com/AdrianOrlow/files-api/app/model"
	"github.com/AdrianOrlow/files-api/app/utils"
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
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	a.DB = model.DBMigrate(db)

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
	log.Print("Listening on " + host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func (a *App) handleRequest(h RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(a.DB, w, r)
	}
}
