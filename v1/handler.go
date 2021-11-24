package v1

import (
	"database/sql"
	"log"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Handler(router *httprouter.Router, sqlite3Connection, baseURL string) {
	db, err := sql.Open("sqlite3", sqlite3Connection)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	awaitingHeightmap := NewAwaitingLayerService(2)

	layoutService := NewLayoutService(db, &DefaultUUIDService{})
	LayoutController{
		NewLayoutAwaitingLayerWire(
			layoutService,
			awaitingHeightmap,
		),
	}.AddRoutes(router)
	LayerController{NewLayerService(baseURL, db, layoutService)}.AddRoutes(router)

	(&AwaitingLayersController{
		awaitingHeightmap,
		NewAwaitingLayerService(2),
	}).AddRoutes(router)

	VersionController{}.AddRoutes(router)
}
