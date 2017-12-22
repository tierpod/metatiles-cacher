// Unpack tiles data from metatile(s) to directory.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/metatile"
)

const defaultExt = ".png"
const usage = `Unpack tiles data from metatiles to directory.

Usage: unpack-metatiles -dir DIR [-ext EXT] /path/to/file1 [path/to/file2]
`

func main() {
	// Command line flags
	var (
		flagDir string
		flagExt string
	)

	flag.Usage = func() {
		fmt.Printf(usage)
		flag.PrintDefaults()
	}

	flag.StringVar(&flagDir, "dir", "", "output directory")
	flag.StringVar(&flagExt, "ext", defaultExt, "append extension")
	flag.Parse()
	files := flag.Args()

	if flagDir == "" {
		fmt.Println("-dir is not set")
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("files(s) is not set")
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Printf("unpack %v to %v\n", f, flagDir)

		// read metatile file
		f, err := os.Open(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		// decode headers from reader
		ml, err := metatile.DecodeHeader(f)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(ml)

		// seek to the start of reader
		_, err = f.Seek(0, 0)
		if err != nil {
			fmt.Println(err)
			return
		}

		// decode data from reader
		data, err := metatile.DecodeData(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(len(data))
		size := int(math.Sqrt(float64(ml.Count)))
		err = writeToFile(data, flagDir, flagExt, int(ml.Z), int(ml.X), int(ml.Y), size)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("----------")
	}
}

func writeToFile(data [][]byte, dir, ext string, z, x, y, size int) error {
	// create dir using z/x format
	files := []string{}
	for xx := x; xx < x+size; xx++ {
		for yy := y; yy < y+size; yy++ {
			subDir := path.Join(dir, strconv.Itoa(z), strconv.Itoa(xx))
			err := os.MkdirAll(subDir, 0777)
			if err != nil {
				return err
			}
			fileName := path.Join(subDir, strconv.Itoa(yy)+ext)
			files = append(files, fileName)
		}
	}

	// write data to files
	for i, d := range data {
		err := ioutil.WriteFile(files[i], d, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
