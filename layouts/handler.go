package layouts

import (
	"database/sql"
	"log"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Handler(sqlite3Connection, baseURL string) *httprouter.Router {
	db, err := sql.Open("sqlite3", sqlite3Connection)
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	layoutService := NewLayoutService(db, &DefaultUUIDService{})
	LayoutController{layoutService}.AddRoutes(router)
	LayerController{NewLayerService(baseURL, db, layoutService)}.AddRoutes(router)
	VersionController{}.AddRoutes(router)
	return router
}
