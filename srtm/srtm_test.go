package srtm

import "testing"

func TestGetSrtmTileName(t *testing.T) {
	cities := map[string]struct {
		coordinates [2]float64
		tile        string
	}{
		"New York":     {coordinates: [2]float64{40.7128, -74.0060}, tile: "N40W075"},
		"London":       {coordinates: [2]float64{51.5074, -0.1278}, tile: "N51W001"},
		"Berlin":       {coordinates: [2]float64{52.5200, 13.4050}, tile: "N52E013"},
		"Quito":        {coordinates: [2]float64{-0.1807, -78.4678}, tile: "S01W079"},
		"Buenos Aires": {coordinates: [2]float64{-34.6037, -58.3816}, tile: "S35W059"},
	}

	for _, city := range cities {
		if tile := GetSrtmTileName(city.coordinates[0], city.coordinates[1]); tile != city.tile {
			t.Errorf("expected %s got %s", city.tile, tile)
		}
	}
}
