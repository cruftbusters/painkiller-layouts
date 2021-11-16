package layouts

import (
	"database/sql"
	"errors"

	. "github.com/cruftbusters/painkiller-layouts/types"
)

type LayoutService interface {
	Create(layout Layout) Layout
	Get(id string) (Layout, error)
	GetAll(excludeMapsWithHeightmap bool) []Layout
	Patch(id string, layout Layout) (Layout, error)
	Delete(id string) error
}

type DefaultLayoutService struct {
	db          *sql.DB
	uuidService UUIDService
}

func NewLayoutService(db *sql.DB, uuidService UUIDService) LayoutService {
	return &DefaultLayoutService{
		db:          db,
		uuidService: uuidService,
	}
}

var ErrLayoutNotFound = errors.New("layout not found")

func (service *DefaultLayoutService) Create(requestLayout Layout) Layout {
	id := service.uuidService.NewUUID()
	statement, err := service.db.Prepare(`
insert into layouts (
	id,
	size_width, size_height,
	bounds_left, bounds_top, bounds_right, bounds_bottom,
	heightmap_url, hillshade_url
) values(?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	layout := &requestLayout
	layout.Id = id
	statement.Exec(
		id,
		requestLayout.Size.Width,
		requestLayout.Size.Height,
		requestLayout.Bounds.Left,
		requestLayout.Bounds.Top,
		requestLayout.Bounds.Right,
		requestLayout.Bounds.Bottom,
		requestLayout.HeightmapURL,
		requestLayout.HillshadeURL,
	)
	return *layout
}

func (service *DefaultLayoutService) Get(id string) (Layout, error) {
	statement, err := service.db.Prepare(`
select id,
size_width, size_height,
bounds_left, bounds_top, bounds_right, bounds_bottom,
heightmap_url, hillshade_url
from layouts where id = ?
`)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	layout := Layout{}
	err = statement.QueryRow(id).Scan(
		&layout.Id,
		&layout.Size.Width,
		&layout.Size.Height,
		&layout.Bounds.Left,
		&layout.Bounds.Top,
		&layout.Bounds.Right,
		&layout.Bounds.Bottom,
		&layout.HeightmapURL,
		&layout.HillshadeURL,
	)
	switch err {
	case sql.ErrNoRows:
		return layout, ErrLayoutNotFound
	case nil:
		return layout, nil
	default:
		panic(err)
	}
}

func (service *DefaultLayoutService) GetAll(excludeMapsWithHeightmap bool) []Layout {
	layouts := []Layout{}
	statement, err := service.db.Prepare(`
select id,
size_width, size_height,
bounds_left, bounds_top, bounds_right, bounds_bottom,
heightmap_url, hillshade_url
from layouts`)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	rows, err := statement.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	layout := Layout{}
	for rows.Next() {
		if err = rows.Scan(
			&layout.Id,
			&layout.Size.Width,
			&layout.Size.Height,
			&layout.Bounds.Left,
			&layout.Bounds.Top,
			&layout.Bounds.Right,
			&layout.Bounds.Bottom,
			&layout.HeightmapURL,
			&layout.HillshadeURL,
		); err != nil {
			panic(err)
		}
		if !excludeMapsWithHeightmap || layout.HeightmapURL == "" {
			layouts = append(layouts, layout)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return layouts
}

func (service *DefaultLayoutService) Patch(id string, patch Layout) (Layout, error) {
	oldLayout, err := service.Get(id)
	if err == ErrLayoutNotFound {
		return Layout{}, ErrLayoutNotFound
	} else if err != nil {
		panic(err)
	}

	statement, err := service.db.Prepare("update layouts set heightmap_url = ? where id = ?")
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	if _, err = statement.Exec(patch.HeightmapURL, id); err != nil {
		panic(err)
	}

	newLayout := &oldLayout
	newLayout.HeightmapURL = patch.HeightmapURL
	return *newLayout, nil
}

func (service *DefaultLayoutService) Delete(id string) error {
	statement, err := service.db.Prepare("delete from layouts where id = ?")
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	if _, err = statement.Exec(id); err != nil {
		panic(err)
	}
	return nil
}
