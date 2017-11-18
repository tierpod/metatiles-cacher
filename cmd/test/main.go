package main

import (
	"fmt"
	"log"

	"github.com/tierpod/metatiles-cacher/pkg/config"
)

func main() {
	cfg, err := config.Load("./config/dev.yaml")
	if err != nil {
		log.Fatal(err)
	}

	s, err := cfg.Source("example")
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("%+v\n", cfg)
	//fmt.Printf("%+v\n%v\n", s, s.HasRegion())
	fmt.Println(s.Zoom)
}
