package world

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// Atlas is the actual 'map' of the world
type Atlas struct {
	// TODO : Investigate if a 2d array or flat array will have better performance.
	Tiles    [][]*TileInstance
	TileDict map[string]TileGeneric
}

// TileInstance is a TileGeneric with unique properties on top
type TileInstance struct {
	// Unique properties
	x int
	y int

	// Generic properties
	generic *TileGeneric
}

// TileGeneric is generic data on a tile, used for instancing
type TileGeneric struct {
	refID string
	icon  rune
	name  string
	info  string
	clip  bool
	color string
}

// NewAtlas reads in tile dictionary and layout, and returns an Atlas
func NewAtlas() *Atlas {
	atlas := new(Atlas)
	atlas.loadTileDict()
	atlas.loadAtlas()

	return atlas
}

// loadTileDict reads external JSON and populates dictionary of tiles
func (a *Atlas) loadTileDict() {
	dat, datErr := ioutil.ReadFile("static/tiledict.json")
	if datErr != nil {
		log.Errorln(fmt.Sprint("Could not read static/tiledict.json ", datErr))
		return
	}

	jsonErr := json.Unmarshal(dat, &a.TileDict)
	if jsonErr != nil {
		log.Errorln(fmt.Sprint("Malformed json of static/tiledict.json ", jsonErr))
		return
	}

	log.WithFields(log.Fields{
		"count": len(a.TileDict),
	}).Infoln("Read static/tiledict.json")

	// Ensure default always available
	def, ok := a.TileDict["_DEFAULT"]
	if !ok {
		log.Errorln("Could not find expected _DEFAULT in tiledict.json")
	}
	// Empty tile maps to _DEFAULT
	a.TileDict[""] = def
}

// loadAtlas populates the actual tiles from file
func (a *Atlas) loadAtlas() {
	file, err := os.Open("static/atlas.csv")
	defer file.Close()
	if err != nil {
		log.Errorln(fmt.Sprint("Could not read static/tiledict.json ", err))
		return
	}

	r := csv.NewReader(file)
	y := 0

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorln(fmt.Sprint("Error in reading static/atlas.csv", err))
		}

		// New row
		row := make([]*TileInstance, len(line))
		for x, ref := range line {
			row[x] = spawnTile(a, ref, x, y)
		}
		y++

		a.Tiles = append(a.Tiles, row)
	}
	log.WithFields(log.Fields{
		"rows": y,
	}).Infoln("Read static/atlas.csv")
}

// spawnTile creates and populates a tile for an Atlas
// populates both unique and generic values
func spawnTile(a *Atlas, ref string, x int, y int) *TileInstance {
	t := new(TileInstance)

	// Uniques
	t.x = x
	t.y = y

	// Generics
	inst, ok := a.TileDict[ref]
	if ok {
		t.generic = &inst
	} else {
		// Use default if not found
		log.WithFields(log.Fields{
			"ref": ref,
			"x":   x,
			"y":   y,
		}).Errorln("Could not find tile reference")
		inst = a.TileDict["_DEFAULT"]
		t.generic = &inst
	}
	return t
}
