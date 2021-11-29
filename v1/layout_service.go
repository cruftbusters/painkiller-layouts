package v1

import (
	"database/sql"
	"errors"

	. "github.com/cruftbusters/painkiller-layouts/types"
)

type LayoutService interface {
	Create(layout Layout) Layout
	Get(id string) (Layout, error)
	GetAll() []Layout
	GetAllWithNoHeightmap() []Layout
	GetAllWithHeightmapWithoutHillshade() []Layout
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
	scale,
	size_width, size_height,
	bounds_left, bounds_top, bounds_right, bounds_bottom,
	heightmap_url, hillshade_url
) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	layout := &requestLayout
	layout.Id = id
	statement.Exec(
		id,
		requestLayout.Scale,
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
scale,
size_width, size_height,
bounds_left, bounds_top, bounds_right, bounds_bottom,
heightmap_url, hillshade_url
from layouts where id = ?
`)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	layout, err := scan(statement.QueryRow(id).Scan)
	switch err {
	case sql.ErrNoRows:
		return layout, ErrLayoutNotFound
	case nil:
		return layout, nil
	default:
		panic(err)
	}
}

func (service *DefaultLayoutService) GetAll() []Layout {
	layouts, err := service.getAllSQL(`
select id,
scale,
size_width, size_height,
bounds_left, bounds_top, bounds_right, bounds_bottom,
heightmap_url, hillshade_url
from layouts`)
	if err != nil {
		panic(err)
	}
	return layouts
}

func (service *DefaultLayoutService) GetAllWithNoHeightmap() []Layout {
	layouts, err := service.getAllSQL(`
select id,
scale,
size_width, size_height,
bounds_left, bounds_top, bounds_right, bounds_bottom,
heightmap_url, hillshade_url
from layouts
where heightmap_url == ''`)
	if err != nil {
		panic(err)
	}
	return layouts
}

func (service *DefaultLayoutService) GetAllWithHeightmapWithoutHillshade() []Layout {
	layouts, err := service.getAllSQL(`
select id,
scale,
size_width, size_height,
bounds_left, bounds_top, bounds_right, bounds_bottom,
heightmap_url, hillshade_url
from layouts
where heightmap_url != ''
and hillshade_url == ''`)
	if err != nil {
		panic(err)
	}
	return layouts
}

func (service *DefaultLayoutService) getAllSQL(sql string) ([]Layout, error) {
	rows, err := service.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	layouts := []Layout{}
	for rows.Next() {
		layout, err := scan(rows.Scan)
		if err != nil {
			return nil, err
		}
		layouts = append(layouts, layout)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return layouts, nil
}

func (service *DefaultLayoutService) Patch(id string, patch Layout) (Layout, error) {
	oldLayout, err := service.Get(id)
	if err == ErrLayoutNotFound {
		return Layout{}, ErrLayoutNotFound
	} else if err != nil {
		panic(err)
	}

	newLayout := &oldLayout

	if patch.HeightmapURL != "" {
		statement, err := service.db.Prepare("update layouts set heightmap_url = ? where id = ?")
		if err != nil {
			panic(err)
		}
		defer statement.Close()
		if _, err = statement.Exec(patch.HeightmapURL, id); err != nil {
			panic(err)
		}
		newLayout.HeightmapURL = patch.HeightmapURL
	}

	if patch.HillshadeURL != "" {
		statement, err := service.db.Prepare("update layouts set hillshade_url = ? where id = ?")
		if err != nil {
			panic(err)
		}
		defer statement.Close()
		if _, err = statement.Exec(patch.HillshadeURL, id); err != nil {
			panic(err)
		}
		newLayout.HillshadeURL = patch.HillshadeURL
	}

	if patch.Scale != 0 {
		statement, err := service.db.Prepare("update layouts set scale = ? where id = ?")
		if err != nil {
			panic(err)
		}
		defer statement.Close()
		if _, err = statement.Exec(patch.Scale, id); err != nil {
			panic(err)
		}
		newLayout.Scale = patch.Scale
	}

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

func scan(scan func(dest ...interface{}) error) (Layout, error) {
	layout := Layout{}
	err := scan(
		&layout.Id,
		&layout.Scale,
		&layout.Size.Width,
		&layout.Size.Height,
		&layout.Bounds.Left,
		&layout.Bounds.Top,
		&layout.Bounds.Right,
		&layout.Bounds.Bottom,
		&layout.HeightmapURL,
		&layout.HillshadeURL,
	)
	return layout, err
}
