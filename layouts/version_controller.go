package layouts

import (
	"net/http"

	"github.com/cruftbusters/painkiller-gallery/types"
	"github.com/julienschmidt/httprouter"
)

type VersionController struct{}

func (c VersionController) AddRoutes(router *httprouter.Router) {
	router.GET("/version", c.GetVersion)
}

func (c VersionController) GetVersion(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	types.EncodeVersion(response, types.Version{Version: "1"})
}
