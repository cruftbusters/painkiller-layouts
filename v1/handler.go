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

	awaitingHeightmap := NewAwaitingLayerService(8)
	awaitingHillshade := NewAwaitingLayerService(8)
	layoutService := NewLayoutService(db, &DefaultUUIDService{})
	wire := NewLayoutAwaitingLayerWire(
		layoutService,
		awaitingHeightmap,
		awaitingHillshade,
	)

	LayoutController{wire}.AddRoutes(router)
	LayerController{NewLayerService(baseURL, db, wire)}.AddRoutes(router)

	(&AwaitingLayersController{
		awaitingHeightmap,
		awaitingHillshade,
	}).AddRoutes(router)

	VersionController{}.AddRoutes(router)
}
