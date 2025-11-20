package main

import (
	"strings"
	"github.com/Supraboy981322/gomn"
)

func parse(in []string, out []string, defs gomn.Map, sub bool) []string {
	for _, chunk := range in {
		subFunc := strings.FieldsFunc(chunk, subFuncSplitter)
		if len(subFunc) > 1 {
			subDefs, ok := defs[subFunc[0]].(gomn.Map)
			if !ok {
				out = append(out, chunk)
			} else {
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
