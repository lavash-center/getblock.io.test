package resources

import (
	"github.com/go-chi/chi/v5"
	"github.com/lavash-center/getblock.io.test/src/managers"
)

type Resource struct {
	path      string
	blocksMan managers.BlocksManager
}

func NewResource(path string, blocksMan managers.BlocksManager) *Resource {
	return &Resource{
		path:      path,
		blocksMan: blocksMan,
	}
}

func (res Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route(res.path, func(r chi.Router) {
		r.Get("/address", res.getBlockAddr)
	})

	return r
}
