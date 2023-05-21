package srtm

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	SRTMGL1 = 1
	SRTMGL3 = 3
)

var (
	srtmBaseUrl    string
	srtmResolution int
)

func Init(baseUrl string, resolution int) {
	srtmBaseUrl = baseUrl
	srtmResolution = resolution
}

func LoadTile(lat, lon float64) ([][]float64, error) {

	tileName := GetSrtmTileName(lat, lon)

	data, err := loadFile(fmt.Sprintf("%s%s.SRTMGL%d.hgt.zip", srtmBaseUrl, tileName, srtmResolution))

	if err != nil {
		log.Fatal(err)
	}

	bytes, err := unpackSrtm(data)

	if err != nil {
		log.Fatal(err)
	}

	tileName = tileName[0:7]

	points, _ := parse(tileName, bytes)

	return points, nil
}

func GetSrtmTileName(lat, lon float64) string {
	tileLat := math.Floor(math.Abs(lat))
	if lat < 0 {
		tileLat = tileLat*-1 - 1
	}

	tileLon := math.Floor(math.Abs(lon))
	if lon < 0 {
		tileLon = tileLon*-1 - 1
	}

	ns := 'N'
	if tileLat < 0 {
		ns = 'S'
	}
	we := 'E'
	if tileLon < 0 {
		we = 'W'
	}

	filename := fmt.Sprintf("%s%02d%s%03d",
		string(ns), int(math.Abs(tileLat)),
		string(we), int(math.Abs(tileLon)))

	return filename
}

func parse(tileName string, bytes []byte) ([][]float64, error) {

	ns := tileName[0]
	we := tileName[3]

	var err error

	var tlat, tlon int

	if tlat, err = strconv.Atoi(tileName[1:3]); err != nil {
		return nil, err
	} else if ns == 'S' {
		tlat = -tlat
	}

	if tlon, err = strconv.Atoi(tileName[4:7]); err != nil {
		return nil, err
	} else if we == 'W' {
		tlon = -tlon
	}

	var resolution int

	if len(bytes) == 3601*3601*2 {
		resolution = 3601
	} else if len(bytes) == 1201*1201*2 {
		resolution = 1201
	} else {
		log.Fatal("unknown SRTM data format")
	}

	var points [][]float64

	for i := 0; i < len(bytes); i += 2 {

		if bytes[i] == 0x80 {
			continue
		}

		index := i / 2
		row := index / resolution
		col := index % resolution

		lat := float64(tlat) + float64(row)/float64(resolution)
		lon := float64(tlon) + float64(col)/float64(resolution)

		elevation := binary.BigEndian.Uint16([]byte{bytes[i], bytes[i+1]})

		points = append(points, []float64{lat, lon, float64(elevation)})
	}

	return points, nil
}

func unpackSrtm(data []byte) ([]byte, error) {

	arc, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	if err != nil {
		return nil, err
	}

	file, err := arc.File[0].Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func loadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
