package main

import (
	"encoding/json"
	"github.com/peterbourgon/diskv"
	"log"
)

var disk *diskv.Diskv

func usedHashes() map[[16]byte]struct{} {
	hashes, err := disk.Read("used-hashes")
	if err != nil {
		log.Println("can't read used-hashes")
		return make(map[[16]byte]struct{})
	}
	num := len(hashes) / 16
	result := make(map[[16]byte]struct{})
	for i := 0; i < num; i++ {
		var v [16]byte
		copy(v[:], hashes[i*16:])
		result[v] = struct{}{}
	}
	return result
}

func saveHashes(m map[[16]byte]struct{}) {
	hashes := make([]byte, len(m)*16)
	i := 0
	for hash := range m {
		copy(hashes[i*16:], hash[:])
		i++
	}
	disk.Write("used-hashes", hashes)
}

func getMailOptions() MailOptions {
	buf, err := disk.Read("mail-opts")
	if err != nil {
		log.Println("can't read email opts, saving default")
		def, err := json.Marshal(MailOptions{})
		if err != nil {
			log.Println("json error ", err)
			return MailOptions{}
		}
		disk.Write("mail-opts", def)
		return MailOptions{}
	}
	var res MailOptions
	err = json.Unmarshal(buf, &res)
	if err != nil {
		log.Println("unmarshal error: ", err)
	}
	return res
}
