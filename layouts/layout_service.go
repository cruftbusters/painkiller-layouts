package layouts

import (
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
	uuidService UUIDService

	layout map[string]Layout
}

func NewLayoutService(uuidService UUIDService) LayoutService {
	return &DefaultLayoutService{
		uuidService: uuidService,
		layout:      make(map[string]Layout),
	}
}

var ErrLayoutNotFound = errors.New("layout not found")

func (service *DefaultLayoutService) Create(requestLayout Layout) Layout {
	id := service.uuidService.NewUUID()
	layout := &requestLayout
	layout.Id = id
	service.layout[id] = *layout
	return *layout
}

func (service *DefaultLayoutService) Get(id string) (Layout, error) {
	if layout, hasKey := service.layout[id]; hasKey {
		return layout, nil
	} else {
		return Layout{}, ErrLayoutNotFound
	}
}

func (service *DefaultLayoutService) GetAll(excludeMapsWithHeightmap bool) []Layout {
	layouts := make([]Layout, 0, len(service.layout))
	for _, layout := range service.layout {
		if !excludeMapsWithHeightmap || layout.HeightmapURL == "" {
			layouts = append(layouts, layout)
		}
	}
	return layouts
}

func (service *DefaultLayoutService) Patch(id string, patch Layout) (Layout, error) {
	if oldLayout, ok := service.layout[id]; ok {
		newLayout := &oldLayout
		newLayout.HeightmapURL = patch.HeightmapURL
		service.layout[id] = *newLayout
		return *newLayout, nil
	} else {
		return Layout{}, ErrLayoutNotFound
	}
}

func (service *DefaultLayoutService) Delete(id string) error {
	delete(service.layout, id)
	return nil
}
