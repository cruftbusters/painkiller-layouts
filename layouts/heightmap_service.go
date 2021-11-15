package layouts

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/cruftbusters/painkiller-layouts/types"
)

type HeightmapService interface {
	Put(id string, heightmap []byte) error
	Get(id string) ([]byte, string, error)
}

func NewHeightmapService(
	baseURL string,
	db *sql.DB,
	layoutService LayoutService,
) HeightmapService {
	if _, err := db.Exec(`
create table if not exists heightmaps(
	id text primary key,
	heightmap blob
)`); err != nil {
		panic(err)
	}
	return &DefaultHeightmapService{
		baseURL,
		db,
		layoutService,
	}
}

type DefaultHeightmapService struct {
	baseURL       string
	db            *sql.DB
	layoutService LayoutService
}

var ErrHeightmapNotFound = errors.New("heightmap not found")

func (s *DefaultHeightmapService) Put(id string, heightmap []byte) error {
	_, err := s.layoutService.Get(id)
	if err != nil {
		return err
	}
	statement, err := s.db.Prepare("insert into heightmaps (id, heightmap) values(?, ?)")
	if err != nil {
		panic(err)
	}
	if _, err = statement.Exec(id, heightmap); err != nil {
		panic(err)
	}
	heightmapURL := fmt.Sprintf("%s/v1/layouts/%s/heightmap.jpg", s.baseURL, id)
	_, err = s.layoutService.Patch(id, types.Layout{HeightmapURL: heightmapURL})
	return err
}

func (s *DefaultHeightmapService) Get(id string) ([]byte, string, error) {
	if _, err := s.layoutService.Get(id); err != nil {
		return nil, "", err
	}
	statement, err := s.db.Prepare("select heightmap from heightmaps where id = ?")
	if err != nil {
		panic(err)
	}
	var heightmap []byte
	err = statement.QueryRow(id).Scan(&heightmap)
	switch err {
	case sql.ErrNoRows:
		return nil, "", ErrHeightmapNotFound
	case nil:
		return heightmap, "image/jpeg", nil
	default:
		panic(err)
	}
}
