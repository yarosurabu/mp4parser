package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func printHelp() {
	fmt.Println("usage: ", os.Args[0], "<filename>")
}

var atomStruct = []byte(`
{
    "ftyp":{},
    "pdin":{},
    "moov":{
        "mvhd":{},
        "trak":{
            "tkhd":{},
            "tref":{},
            "edts":{
                "elst":{}
            },
            "mdia":{
                "mdhd":{},
                "hdlr":{},
                "minf":{
                    "vmhd":{},
                    "smhd":{},
                    "hmhd":{},
                    "nmhd":{},
                    "dinf":{
                        "dref":{}
                    },
                    "stbl":{
                        "stsd":{},
                        "stts":{},
                        "ctts":{},
                        "stsc":{},
                        "stsz":{},
                        "stz2":{},
                        "stco":{},
                        "co64":{},
                        "stss":{},
                        "stsh":{},
                        "padb":{},
                        "stdp":{},
                        "sdtp":{},
                        "sbpg":{},
                        "sgpd":{},
                        "subs":{}
                    }
                }
            }
        },
        "mvex":{
            "mehd":{},
            "trex":{}
        },
        "ipmc":{}
    },
    "moof":{
        "mfhd":{},
        "traf":{
            "tfhd":{},
            "trun":{},
            "sdtp":{},
            "sdgp":{},
            "subs":{}
        }
    },
    "mfra":{
        "tfra":{},
        "mfro":{}
    },
    "mdat":{},
    "free":{},
    "skip":{
        "udta":{
            "cprt":{}
        }
    },
    "meta":{
        "hdlr":{},
        "dinf":{
            "dref":{}
        },
        "ipmc":{},
        "iloc":{},
        "ipro":{
            "sinf":{
                "frma":{},
                "imif":{},
                "schm":{},
                "schi":{}
            }
        },
        "iinf":{},
        "xml":{},
        "bxml":{},
        "pitm":{}
    }
}`)

func readAtom(file *os.File, r *bufio.Reader, curAtom map[string]interface{}, level int, pos int64, size int64) {
	data := make([]byte, 8)
	innerPos := int64(0)
	for innerPos+8 <= size {
		_, err := file.Seek(pos+innerPos, io.SeekStart)
		if err != nil {
			panic(err)
		}
		r.Reset(file)
		_, err = r.Read(data)
		if err != nil {
			panic(err)
		}
		s := int64(binary.BigEndian.Uint32(data[:4]))
		atom := string(data[4:])
		var offset = int64(8)
		known := false
		var subAtom map[string]interface{}
		for key, val := range curAtom {
			if key == atom {
				known = true
				subAtom = val.(map[string]interface{})
				break
			}
		}
		if s == 0 {
			s = size - innerPos
		} else if s == 1 {
			_, err = r.Read(data)
			if err != nil {
				panic(err)
			}
			s = int64(binary.BigEndian.Uint64(data))
			offset += 8
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
			if len(subAtom) > 0 {
				readAtom(file, r, subAtom, level+1, pos+innerPos+offset, s-offset)
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

	buf := bufio.NewReader(file)

	var root map[string]interface{}
	err = json.Unmarshal(atomStruct, &root)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()
	start := time.Now()
	readAtom(file, buf, root, 0, 0, stat.Size())
	t := time.Now()
	fmt.Println(t.Sub(start))
}
