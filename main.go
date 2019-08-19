package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func printHelp() {
	fmt.Println("usage: ", os.Args[0], "<filename>")
}

var atoms = [...]string{
	"ftyp", "pdin", "moov", "mvhd", "trak", "tkhd", "tref", "edts",
	"elst", "mdia", "mdhd", "hdlr", "minf", "vmhd", "smhd", "hmhd",
	"nmhd", "dinf", "dref", "stbl", "stsd", "stts", "ctts", "stsc",
	"stsz", "stz2", "stco", "co64", "stss", "stsh", "padb", "stdp",
	"sdtp", "sbgp", "sgpb", "subs", "mvex", "mehd", "trex", "ipmc",
	"moof", "mfhd", "traf", "tfhd", "trun", "sdtp", "sbgp", "subs",
	"mfra", "tfra", "mfro", "mdat", "free", "skip", "udta", "cprt",
	"meta", "hdlr", "dinf", "dref", "ipmc", "iloc", "ipro", "sinf",
	"frma", "imif", "schm", "schi", "iinf", "xml", "bxml", "pitm"}

func readAtom(file *os.File, level int, pos int64, size int64) {
	data := make([]byte, 8)
	var innerPos int64
	for innerPos+8 <= size {
		file.Seek(pos+innerPos, 0)
		_, err := file.Read(data)
		if err != nil {
			panic(err)
		}
		s := int64(binary.BigEndian.Uint32(data[:4]))
		atom := string(data[4:])
		known := false
		for _, a := range atoms {
			if a == atom {
				known = true
				break
			}
		}
		if s == 0 {
			s = size - innerPos
		} else if s == 1 {
			_, err = file.Read(data)
			if err != nil {
				panic(err)
			}
			s = int64(binary.BigEndian.Uint64(data))
		} else if s < 8 {
			return
		}
		if innerPos+s > size {
			return
		}
		if known {
			for index := 0; index < level; index++ {
				fmt.Print("    ")
			}
			fmt.Println(s, atom)
			if atom != "mdat" && atom != "free" {
				readAtom(file, level+1, pos+innerPos+8, s-8)
			}
		}
		innerPos += s
	}
}

func main() {
	if len(os.Args) != 2 {
		printHelp()
		return
	}
	filename := os.Args[1]
	fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()
	readAtom(file, 0, 0, stat.Size())
}
