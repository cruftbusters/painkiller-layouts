package v1

import (
	"database/sql"
	"log"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Handler(router *httprouter.Router, sqlite3Connection, baseURL string) {
	db, err := sql.Open("sqlite3", sqlite3Connection)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	layoutsAwaitingHeightmap := make(chan types.Layout, 2)
	layoutsAwaitingHillshade := make(chan types.Layout, 2)

	layoutService := NewLayoutService(db, &DefaultUUIDService{})
	LayoutController{
		layoutService,
		layoutsAwaitingHeightmap,
		layoutsAwaitingHillshade,
	}.AddRoutes(router)
	LayerController{NewLayerService(baseURL, db, layoutService)}.AddRoutes(router)

	(&AwaitingLayoutController{
		time.Second * 5,
		layoutsAwaitingHeightmap,
		layoutsAwaitingHillshade,
	}).AddRoutes(router)

	VersionController{}.AddRoutes(router)
}
