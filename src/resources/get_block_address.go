package resources

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func (res Resource) getBlockAddr(w http.ResponseWriter, r *http.Request) {
	addr, err := res.blocksMan.GetBlockAddress()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	render.JSON(w, r, addr)
}
