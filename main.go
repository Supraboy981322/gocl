package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	isString bool
	killOnWarn bool
	input []string
	splitters []string
	headEnd string
	defsGlob gomn.Map
	rcDefs gomn.Map	
	imports gomn.Map
	headDefs gomn.Map
)

func init() {
	var ok bool

	defsFile, err := os.ReadFile("defs.gomn")
	if err != nil {
		log.Fatalf("failed to read defs.gomn:\n  %v", err)
	}
	
	if defsGlob, err = gomn.Parse(string(defsFile)); err != nil {
		log.Fatal(err)
	}

	if rcDefs, ok = defsGlob[0].(gomn.Map); !ok {
		log.Warn("rc definitions not found, may produce odd results")
		kilOcont("continuing anyways")
	} else {
		if killOnWarn, _ = rcDefs["kill on warn"].(bool); killOnWarn {
			log.Info("configured to kill on warn")
		}
		if headEnd, ok = rcDefs["head end"].(string); !ok {
			log.Warn("\"head end\" not defined in rc definitions, this will probably cause problems")
			kilOcont("continuing anyways")
		}
	}
	
	inFile, err := os.ReadFile("foo.gocl")
	if err != nil {
		log.Fatal(err)
	}
	
	strFile := string(inFile)
	var trimmedFile []string
	for _, line := range strings.Split(strFile, "\n") {
		trimmedFile = append(trimmedFile, strings.TrimSpace(line))
//		fmt.Println(line)
	}
	input = append(input, strings.Join(trimmedFile, "\n"))
}

func main() {
	var ok bool

	input := strings.FieldsFunc(input[0], whitespaceSplitter)

	inputHeader := getHeader()
	input = input[len(inputHeader):]

	if headDefs, ok = rcDefs["head defs"].(gomn.Map); !ok {
		kilOcont("head defs not defined")
	}

	if imports, ok = headDefs["imports"].(gomn.Map); !ok {
		kilOcont("imports not defined, could be a non-problem")
	}

	//parse the header
	var outputHeader []string
	outputHeader = parseHeader(inputHeader, outputHeader)

	//parse the main script
	var outputMain []string
	outputMain = parse(input, outputMain, defsGlob, false)

	//combine them for output
	output := make([]string, len(outputMain)+len(outputHeader))
	copy(output, outputHeader)
	copy(output[len(outputHeader):], outputMain)

	//print it
	//  (for testing, will be changed to write to file)
	for i, chunk := range output {
		fmt.Print(chunk + splitters[i])
	}//	fmt.Printf("new: %#v\n", output)
}


func importParser(old []string, defs gomn.Map) []string {
	var out []string
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
	if r == '"' {
		isString = !isString
		return false
	} else if isString{
		return false
	} else {
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
	}
	return false
}

func subFuncSplitter(r rune) bool {
	if isString {
		return false
	} else if r == '"' {
		isString = true
		return false
	}
	return r == '.'
}

func kilOcont(str string) {
	log.Warn(str)
	if killOnWarn {
		os.Exit(1)
	} else {
		log.Info("continuing anyways")
	}
}
