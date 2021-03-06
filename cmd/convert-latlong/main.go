// convert-latlong is the small tool for converting latitude and longitude coordinates to z, x, y.
package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/bbox"
	"github.com/tierpod/metatiles-cacher/pkg/latlong"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

// Default flags values.
const (
	defaultPrefix = "/var/lib/mod_tile/style/"
	defaultExt    = "png"
)

var version string

// Pairs of integers flags with min and max values
type intPair struct {
	min, max int
}

func (i *intPair) String() string {
	return fmt.Sprintf("Int min: %v, max: %v", i.min, i.max)
}

func (i *intPair) Set(value string) error {
	values := strings.Split(value, "-")
	if len(values) != 2 {
		return errors.New("Wrong int range: need 2 integers, separated by '-'")
	}

	v1, err := strconv.Atoi(values[0])
	if err != nil {
		return err
	}

	v2, err := strconv.Atoi(values[1])
	if err != nil {
		return err
	}

	if v1 > v2 {
		i.min = v2
		i.max = v1
	} else {
		i.min = v1
		i.max = v2
	}

	return nil
}

// Pairs of float64 values wit min and max values
type float64Pair struct {
	min, max float64
}

func (f *float64Pair) String() string {
	return fmt.Sprintf("Float pair: min: %v, max: %v", f.min, f.max)
}

func (f *float64Pair) Set(value string) error {
	values := strings.Split(value, "-")
	if len(values) != 2 {
		return errors.New("Wrong float64 range")
	}

	v1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return err
	}

	v2, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return err
	}

	f.min = math.Min(v1, v2)
	f.max = math.Max(v1, v2)

	return nil
}

func main() {
	// Command line flags
	var (
		flagLat     float64Pair
		flagLong    float64Pair
		flagZooms   intPair
		flagPrefix  string
		flagExt     string
		flagMeta    bool
		flagVersion bool
	)

	flag.Var(&flagLat, "lat", "Latitude coordinates `range`, separated by '-'")
	flag.Var(&flagLong, "long", "Longitude coordinates `range`, separated by '-'")
	flag.Var(&flagZooms, "zooms", "Zooms `range`, separated by '-' (example: 10-12)")
	flag.StringVar(&flagPrefix, "prefix", defaultPrefix, "Output string `prefix`")
	flag.StringVar(&flagExt, "ext", defaultExt, "Output `extension` for tile (metatile always has 'meta' ext)")
	flag.BoolVar(&flagMeta, "meta", false, "Convert output to metatiles format?")
	flag.BoolVar(&flagVersion, "v", false, "Show version and exit")
	flag.Parse()

	/*f, err := os.Create("./pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()*/

	if flagVersion {
		fmt.Printf("Version: %v\n", version)
		os.Exit(0)
	}

	if flagZooms.min == 0 && flagZooms.max == 0 {
		fmt.Println("[ERROR] -zooms flag is not set or set zero values")
		os.Exit(1)
	}

	if flagZooms.min < 1 {
		fmt.Printf("[ERROR] Got wrong minimum zoom level: %v < 1\n", flagZooms.min)
		os.Exit(1)
	}

	if flagZooms.max > 18 {
		fmt.Printf("[ERROR] Got wrong maximum zoom level: %v > 18\n", flagZooms.max)
		os.Exit(1)
	}

	if flagLat.min == 0 && flagLat.max == 0 {
		fmt.Println("[ERROR] -lat flag is not set or set zero values")
		os.Exit(1)
	}

	if flagLong.min == 0 && flagLong.max == 0 {
		fmt.Println("[ERROR] -long flag is not set or set zero values")
		os.Exit(1)
	}

	top := latlong.LatLong{Lat: flagLat.max, Long: flagLong.min}
	bottom := latlong.LatLong{Lat: flagLat.min, Long: flagLong.max}
	zooms := util.MakeIntSlice(flagZooms.min, flagZooms.max+1)

	tiles := bbox.NewFromLatLong(zooms, top, bottom, flagExt)

	filepathPrev := ""
	for t := range tiles {
		if flagMeta {
			mt := metatile.NewFromTile(t)
			filepath := mt.Filepath(flagPrefix)
			// skip filepath if same as previous
			if filepath == filepathPrev {
				continue
			}
			fmt.Println(filepathPrev)
			filepathPrev = filepath
		} else {
			fmt.Println(t.Filepath(flagPrefix))
		}
	}
}
