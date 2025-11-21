package main

import (
	"strings"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

func parse(in []string, out []string, defs gomn.Map, sub bool) []string {
	log.Debug("parse()")

	for _, chunk := range in {
		log.Debug("splitting to check for sub function")
		subFunc := strings.FieldsFunc(chunk, subFuncSplitter)

		if len(subFunc) > 1 {
			log.Debug("detected as sub function")

			if subDefs, ok := defs[subFunc[0]].(gomn.Map); !ok {
				log.Debug("failed type assertion on sub function defs, appending chunk anyways")
				out = append(out, chunk)
			} else {
				log.Debug("success with type assertion on sub function defs")
				for _, subChunk := range subFunc {                   //
					subChunkSplit := strings.Split(subChunk, "(")      //
					newChunk, ok := subDefs[subChunkSplit[0]].(string) //
					newChunkSlice := []string{newChunk}                //
					if len(subChunkSplit) > 1 {                        //This is probably
						subChunkSplit[1] = "("+subChunkSplit[1]          //  causing a bug
						newChunkSlice = append(                          //
							parse(subChunkSplit[1:],                       //
								newChunkSlice, defsGlob, false))             //
					}                                                  //
					out = appOut(out, ok, strings.Join(newChunkSlice, ""), subChunk)
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

func subFuncSplitter(r rune) bool {
	if r == '"' {
		isString = !isString
		return false
	} else if isString {
		return false
	}
	return r == '.'
}
