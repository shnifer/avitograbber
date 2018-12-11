package main

import (
	"encoding/json"
	"github.com/peterbourgon/diskv"
	"log"
)

func init() {
	disk = diskv.New(diskv.Options{
		CacheSizeMax: 1024 * 1024,
		BasePath:     "storage",
	})

	siteParts.Parts = make(map[string][]string)
	siteParts.Names = make(map[string]map[string]string)

	buf, err := disk.Read("parts")
	if err != nil {
		log.Println("parts error ", err)
		return
	}
	err = json.Unmarshal(buf, &siteParts)
	if err != nil {
		log.Println("site unmarshal error ", err)
		return
	}

	initAsks()
}
