package main

import "github.com/charmbracelet/log"

var (
	isImp bool
)

func getHeader(in []string) []string {
	log.Debug("getting header")
	var header []string
	for i, val := range in {
		if val == headEnd {
			log.Debug("header ending matched")
			header = in[:i]
			input = in[i+1:]
		}
	}

	log.Debug("returning header")
	return header
}
	
func parseHeader(in []string, out []string) []string {
	log.Debug("parsing header")
	for _, chunk := range in {
		if newChunk, ok := headDefs[chunk].(string); ok {
			log.Debug("matched header definition")
			out = append(out, newChunk)
		} else if impDefs, ok := importsMap[chunk].(string); ok {
			log.Debug("matched as import")
			log.Debug("import bool flipped")
			isImp = !isImp
			out = append(out, impDefs)
		} else if imp, ok := importDefs[chunk].(string); ok && isImp {
			log.Debug("matched import definition")
			out = append(out, imp)
		} else {
			log.Debug("chunk unmatched, adding unchanged")
			out = append(out, chunk)
		}
	}

	log.Debug("returning header")
	return out
}
