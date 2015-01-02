package engi

import (
	"strconv"
)

type Tilemap struct {
	Tiles    [][]Tile
	Tilesize int
}

type Tile struct {
	Point
	Image Drawable
	Solid bool
}

func NewTilemap(mapString [][]string, sheet *Texture, tilesize int) *Tilemap {
	tilemap := Tilemap{}
	position := Point{}
	tilemap.Tilesize = tilesize

	tilemap.Tiles = make([][]Tile, len(mapString))
	for i := range tilemap.Tiles {
		tilemap.Tiles[i] = make([]Tile, len(mapString[0]))
	}

	for y, slice := range mapString {
		for x, key := range slice {
			var image Drawable
			solid := false
			index, err := strconv.Atoi(key)
			if err != nil {
				panic(err)
			}
			if index > 0 {
				image = getRegionOfSpriteSheet(sheet, tilemap.Tilesize, index)
				solid = true
			}

			tile := Tile{Point: Point{position.X + float32(x*tilemap.Tilesize), position.Y + float32(y*tilemap.Tilesize)}, Image: image, Solid: solid}
			tilemap.Tiles[y][x] = tile
		}
	}

	return &tilemap
}

func CollideTilemap(e *Entity, et *Entity, t *Tilemap) {
	var eSpace *SpaceComponent
	var tSpace *SpaceComponent

	if !e.GetComponent(&eSpace) || !et.GetComponent(&tSpace) {
		return
	}

	for _, slice := range t.Tiles {
		for _, tile := range slice {
			if tile.Solid {
				aabb := AABB{Point{tile.X + tSpace.Position.X, tile.Y + tSpace.Position.Y}, Point{tile.X + tSpace.Position.X + 16, tile.Y + tSpace.Position.Y + 16}}
				if IsIntersecting(eSpace.AABB(), aabb) {
					mtd := MinimumTranslation(eSpace.AABB(), aabb)
					eSpace.Position.X += mtd.X
					eSpace.Position.Y += mtd.Y
					Mailbox.Dispatch("CollisionMessage", CollisionMessage{e, et})
				}
			}
		}
	}
}

func getRegionOfSpriteSheet(texture *Texture, tilesize int, index int) *Region {
	width := texture.Width()
	widthInSprites := width / float32(tilesize)

	pointer := Point{}
	step := 0
	for step != (index) {
		step += 1
		if pointer.X < (widthInSprites - 1) {
			pointer.X += 1
		} else {
			pointer.X = 0
			pointer.Y += 1
		}
	}

	pointer.MultiplyScalar(float32(tilesize))

	return NewRegion(texture, int(pointer.X), int(pointer.Y), tilesize, tilesize)
}
