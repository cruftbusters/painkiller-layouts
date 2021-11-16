package layouts

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/cruftbusters/painkiller-layouts/types"
)

type LayerService interface {
	Put(id, name string, layer []byte) error
	Get(id, name string) ([]byte, string, error)
}

func NewLayerService(
	baseURL string,
	db *sql.DB,
	layoutService LayoutService,
) LayerService {
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

func (s *DefaultLayerService) Put(id, name string, layer []byte) error {
	_, err := s.layoutService.Get(id)
	if err != nil {
		return err
	}
	statement, err := s.db.Prepare("insert into layers (id, name, layer) values(?, ?, ?)")
	if err != nil {
		panic(err)
	}
	if _, err = statement.Exec(id, name, layer); err != nil {
		panic(err)
	}
	layerURL := fmt.Sprintf("%s/v1/layouts/%s/%s", s.baseURL, id, name)
	_, err = s.layoutService.Patch(id, types.Layout{HeightmapURL: layerURL})
	return err
}

func (s *DefaultLayerService) Get(id, name string) ([]byte, string, error) {
	if _, err := s.layoutService.Get(id); err != nil {
		return nil, "", err
	}
	statement, err := s.db.Prepare("select layer from layers where id = ? and name = ?")
	if err != nil {
		panic(err)
	}
	var layer []byte
	err = statement.QueryRow(id, name).Scan(&layer)
	switch err {
	case sql.ErrNoRows:
		return nil, "", ErrLayerNotFound
	case nil:
		return layer, "image/jpeg", nil
	default:
		panic(err)
	}
}
