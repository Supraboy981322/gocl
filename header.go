package main

func getHeader() []string {
	var header []string
	for i, val := range input {
		if val == headEnd {
			header = input[:i]
			input = input[i:]
		}
	}

	return header
}
	
func parseHeader(in []string, out []string) []string {
	for _, chunk := range in {
		if newChunk, ok := headDefs[chunk].(string); ok {
			subChunk, _ := headDefs[""].(string)
			out = appOut(out, false, subChunk, newChunk)
		} else {
			out = append(out, chunk)
		}
	}

	return out
}
