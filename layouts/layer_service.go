package layouts

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/cruftbusters/painkiller-layouts/types"
)

type LayerService interface {
	Put(id string, layer []byte) error
	Get(id string) ([]byte, string, error)
}

func NewLayerService(
	baseURL string,
	db *sql.DB,
	layoutService LayoutService,
) LayerService {
	if _, err := db.Exec(`
create table if not exists heightmaps(
	id text primary key on conflict replace,
	heightmap blob
)`); err != nil {
		panic(err)
	}
	return &DefaultLayerService{
		baseURL,
		db,
		layoutService,
	}
}

type DefaultLayerService struct {
	baseURL       string
	db            *sql.DB
	layoutService LayoutService
}

var ErrLayerNotFound = errors.New("layer not found")

func (s *DefaultLayerService) Put(id string, layer []byte) error {
	_, err := s.layoutService.Get(id)
	if err != nil {
		return err
	}
	statement, err := s.db.Prepare("insert into heightmaps (id, heightmap) values(?, ?)")
	if err != nil {
		panic(err)
	}
	if _, err = statement.Exec(id, layer); err != nil {
		panic(err)
	}
	layerURL := fmt.Sprintf("%s/v1/layouts/%s/heightmap.jpg", s.baseURL, id)
	_, err = s.layoutService.Patch(id, types.Layout{HeightmapURL: layerURL})
	return err
}

func (s *DefaultLayerService) Get(id string) ([]byte, string, error) {
	if _, err := s.layoutService.Get(id); err != nil {
		return nil, "", err
	}
	statement, err := s.db.Prepare("select heightmap from heightmaps where id = ?")
	if err != nil {
		panic(err)
	}
	var layer []byte
	err = statement.QueryRow(id).Scan(&layer)
	switch err {
	case sql.ErrNoRows:
		return nil, "", ErrLayerNotFound
	case nil:
		return layer, "image/jpeg", nil
	default:
		panic(err)
	}
}
