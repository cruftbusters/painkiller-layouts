package layouts

import (
	"database/sql"
	"log"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Handler(sqlite3Connection, baseURL string) *httprouter.Router {
	db, err := sql.Open("sqlite3", sqlite3Connection)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	router := httprouter.New()
	layoutChannel := make(chan types.Layout)
	layoutService := NewLayoutService(db, layoutChannel, &DefaultUUIDService{})
	(&DispatchController{layoutChannel, time.Second * 3}).AddRoutes(router)
	LayoutController{layoutService}.AddRoutes(router)
	LayerController{NewLayerService(baseURL, db, layoutService)}.AddRoutes(router)
	VersionController{}.AddRoutes(router)
	return router
}
