package gogue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

type LayerType string

type Layer struct {
	Height, Width int
	X, Y          int
	Data          []int
	Name          string
	Type          LayerType
	Visible       bool
}

type Properties map[string]string

type Tileset struct {
	FirstGID int `json:"firstgid"`
	Properties
	TileProperties map[string]Properties `json:"tileproperties"`
}

type Map struct {
	Height int
	Width  int
	Layers []Layer
	Properties
	Tilesets []Tileset
}

func NewMap(path string) *Map {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("File error: %v\n", e)
	}

	m := new(Map)
	json.Unmarshal(file, m)
	return m
}

func mergeProperties(ps ...Properties) Properties {
	properties := make(Properties)
	for _, p := range ps {
		for k, v := range p {
			properties[k] = v
		}
	}
	return properties
}

func (m Map) tileProperties(tileID int) Properties {
	var ts Tileset
	for _, t := range m.Tilesets {
		if tileID >= t.FirstGID {
			ts = t
		}
	}
	properties := ts.TileProperties[strconv.Itoa(tileID-ts.FirstGID)]
	return mergeProperties(ts.Properties, properties)
}

func (m Map) layerByName(name string) (l Layer, err error) {
	for _, layer := range m.Layers {
		if layer.Name == name {
			l = layer
			return
		}
	}
	err = errors.New("layer not found")
	return
}

func (m Map) groundLayer(z int) (Layer, error) {
	return m.layerByName(fmt.Sprintf("%d Ground", z))
}

func (m Map) wallsLayer(z int) (Layer, error) {
	return m.layerByName(fmt.Sprintf("%d Walls", z))
}

func (m Map) Get(pos Position) (t Tile, err error) {
	if pos.Y >= m.Height || pos.Y < 0 || pos.X >= m.Width || pos.X < 0 {
		err = errors.New("coordinate out of bounds")
		return
	}
	groundLayer, err := m.groundLayer(pos.Z)
	if err != nil {
		return
	}
	wallsLayer, err := m.wallsLayer(pos.Z)
	if err != nil {
		return
	}
	groundTile := groundLayer.Data[pos.Y*m.Height+pos.X]
	wallTile := wallsLayer.Data[pos.Y*m.Height+pos.X]

	groundTileProperties := m.tileProperties(groundTile)
	wallTileProperties := m.tileProperties(wallTile)
	t = Tile{
		Tiles:      []int{groundTile, wallTile},
		Properties: mergeProperties(groundTileProperties, wallTileProperties),
	}
	return
}
