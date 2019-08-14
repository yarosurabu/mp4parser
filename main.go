package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func printHelp() {
	fmt.Println("usage: ", os.Args[0], "<filename>")
}

func readAtom(file *os.File, level int, pos int64, size int64) {
	if size < 8 {
		return
	}
	file.Seek(pos, 0)
	data := make([]byte, 4)
	var innerPos int64
	for innerPos < size {
		_, err := file.Read(data)
		if err != nil {
			fmt.Println("failed to read size")
			return
		}
		s := int64(binary.BigEndian.Uint32(data))
		if s < 8 {
			//fmt.Println("size invalid")
			break
		}
		if innerPos+s > size {
			return
		}
		_, err = file.Read(data)
		if err != nil {
			fmt.Println("failed to read atom")
			return
		}
		for index := 0; index < level; index++ {
			fmt.Print(" ")
		}
		fmt.Println(s, string(data))
		if string(data) != "mdat" && string(data) != "free" {
			readAtom(file, level+1, pos+innerPos+8, s-8)
		}
		innerPos += s
		file.Seek(pos+innerPos, 0)
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
		fmt.Println("failed to open file: ", filename)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("failed to get file stat")
		return
	}

	readAtom(file, 0, 0, stat.Size())
}
