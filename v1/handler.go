package v1

import (
	"database/sql"
	"log"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Handler(router *httprouter.Router, sqlite3Connection, baseURL string) {
	db, err := sql.Open("sqlite3", sqlite3Connection)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	layoutService := NewLayoutService(db, &DefaultUUIDService{})
	LayoutController{layoutService}.AddRoutes(router)
	LayerController{NewLayerService(baseURL, db, layoutService)}.AddRoutes(router)
	(&PendingRendersController{time.Second * 5}).AddRoutes(router)
	VersionController{}.AddRoutes(router)
}