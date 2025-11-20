package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	input []string
	splitters []string
	defsGlob gomn.Map
)

func init() {
	defsFile, err := os.ReadFile("defs.gomn")
	if err != nil {
		log.Fatalf("failed to read defs.gomn:\n  %v", err)
	}
	
	defsGlob, err = gomn.Parse(string(defsFile))
	if err != nil {
		log.Fatal(err)
	}

	inFile, err := os.ReadFile("foo.gocl")
	if err != nil {
		log.Fatal(err)
	}

	strFile := string(inFile)
	var trimmedFile []string
	for _, line := range strings.Split(strFile, "\n") {
		trimmedFile = append(trimmedFile, strings.TrimSpace(line))
		fmt.Println(line)
	}
	input = append(input, strings.Join(trimmedFile, "\n"))
}

func main() {
	input := strings.FieldsFunc(input[0], whitespaceSplitter)
	var output []string
	output = parse(input, output, defsGlob, false)

	for i, chunk := range output {
		fmt.Print(chunk + splitters[i])
	}//	fmt.Printf("new: %#v\n", output)
}

func parse(in []string, out []string, defs gomn.Map, sub bool) []string {
	for _, chunk := range in {
		subFunc := strings.FieldsFunc(chunk, subFuncSplitter)
		if len(subFunc) > 1 {
			subDefs, ok := defs[subFunc[0]].(gomn.Map)
			if !ok {
				out = append(out, chunk)
			} else {
				for _, subChunk := range subFunc {
					newChunk, ok := subDefs[subChunk].(string)
					out = appOut(out, ok, newChunk, subChunk)
				}
			}
		} else {
			if newChunk, ok := defs[chunk].(string); ok || sub {
				subChunk, _ := defs[""].(string)
				out = appOut(out, sub, subChunk, newChunk)
			} else {
				out = append(out, chunk)
			}
		}
	}
	return out
}

func appOut(old []string, cond bool, newVal string, oldVal string) []string {
	if cond {
		return append(old, newVal) 
	} else {
		return append(old, oldVal)
	}
	return []string{}
}

func whitespaceSplitter(r rune) bool {
	switch r {
	case '\n':
		splitters = append(splitters, "\n")
		return true
	case ' ':
		splitters = append(splitters, " ")
		return true
	case '.':
		splitters = append(splitters, ".")
	}
	return false
}

func subFuncSplitter(r rune) bool {
	return r == '.'
}
