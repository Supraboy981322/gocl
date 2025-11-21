package main

import (
	"os"
	"strings"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	//sorry, but I have a lot of
	// vars to keep track of things 
	//  and for configs
	args = os.Args[1:]
	er = log.New(os.Stderr)
	isString bool
	fileExt string
	execute bool
	killOnWarn bool
	writeFile bool
	printOut bool
	debug bool
	input []string
	inputFile string
	outputFile string
	inputHeader []string
	splitters []string
	headEnd string
	splitHead bool
	splitHeadAt int
	defsGlob gomn.Map
	rcDefs gomn.Map	
	importsMap gomn.Map
	importDefs gomn.Map
	headDefs gomn.Map
)

func init() {
	readConf()
	checkArgs()

	//i get an odd bug if not defined this way
	var err error

	if inputFile == "" {
		log.Debug("input file not set")
	}
	
	inFile, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	
	strFile := string(inFile)
	var trimmedFile []string
	for _, line := range strings.Split(strFile, "\n") {
		if trimmedLine := strings.TrimSpace(line); trimmedLine != headEnd {
			trimmedFile = append(trimmedFile, trimmedLine)
		} else {
			inputHeader = strings.FieldsFunc(strings.Join(trimmedFile, "\n"), whitespaceSplitter)
			trimmedFile = []string{""}
		}
//		fmt.Println(line)
	}
	input = append(input, strings.Join(trimmedFile, "\n"))
}

func main() {

	input := strings.FieldsFunc(input[0], whitespaceSplitter)


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

	if debug { 
		log.Debugf("new: %#v\n", output)
	}

	var finalOut string
	for i, chunk := range output {
		if len(splitters)-1 < i {
			splitters = append(splitters, "\n")
		}
		finalOut += chunk + splitters[i]
	}

	if printOut {
		log.Print(finalOut)
	}
	if writeFile {
		if outputFile == "" {
			orExt := filepath.Ext(inputFile)
			orName := strings.TrimSuffix(filepath.Base(inputFile), orExt)
			outputFile = orName + "." + fileExt
			log.Warn("no output file, using input file name:  " + outputFile)
		}
		if err := os.WriteFile(outputFile, []byte(finalOut), 0644); err != nil {
			log.Fatalf("failed to write to file:  %v", err)
		}
	}
}
