package main

import (
	"os"
	"slices"
	"github.com/Supraboy981322/gomn"
	"github.com/charmbracelet/log"
)

func checkArgs() {
	log.Debug("checking args")
	var taken []int
	for i, arg := range args {
		if !slices.Contains(taken, i) {
			log.Debug("arg not yet used")
			var isArg bool
			for k, a := range arg {
				if a == '-' && !isArg {
					log.Debug("start of arg")
					isArg = true
				} else if isArg {
					log.Debug("checking arg char")
					switch a {
       		case 'i':
						log.Debug("matched char as input arg")
						inputFile = args[i+1]
						log.Debugf("input file:  %s", inputFile)
						taken = append(taken, i+1)
					case 'o':
						log.Debug("matched char as output arg")
						outputFile = args[i+1]
						log.Debugf("output file:  %s", outputFile)
						taken = append(taken, i+1)
					case '-':
						log.Debug("matched char as second dash") 
						log.Debug("ignoring remaining arg chars")
					  isArg = false
						log.Debug("passing full arg to traditional parsing func")
						taken = append(taken, checkFullArg(i, arg)...)
					default:
						log.Debug("arg not matched")
						log.Debug("assumed invalid")
						invArg(k, a, arg)
					}
				}
			}
		} else { log.Debug("arg already used") }
	}

	if inputFile == "" {
		log.Debug("input file not set")
		log.Fatal("no file input")
	}
}

func checkFullArg(current int, arg string) []int {
	log.Debug("checking full arg")
	var taken []int
	switch arg {
	case "--input", "--i", "--in", "--source", "--s":
		log.Debug("matched as input arg")
		inputFile = args[current+1]
		log.Debugf("input file:  %s", inputFile)
		taken = append(taken, current+1)
	case "--output", "--o", "--out", "--new":
		log.Debug("matched as output arg")
		outputFile = args[current+1]
		log.Debugf("output file:  %s", outputFile)
		taken = append(taken, current+1)
	default:
		log.Debug("arg not matched")
		log.Debug("assumed invalid")
		er.Fatal("invalid arg:  \033[1;31m" + arg + "\033[0m")
	}
	return taken
}

func invArg(index int, which rune, arg string) {
	log.Debug("printing invalid arg")
	var pointer string
	for i := 0; i < index+19; i++ {
		pointer += " "
	}
	log.Debug("constructing pointer")
	pointer += "\033[1;31m^\033[0m"
	log.Debug("constructing arg")
	arg = arg[:index] + "\033[1;31m" +
			string(arg[index]) + "\033[0m" +
			arg[index+1:]
	log.Debug("printing message")
	er.Fatal("invalid arg:  " + arg + "\n" +
			pointer + "\n")
}

func readConf() {
	log.Debug("reading config")
	var ok bool
	var defsFile []byte
	var err error

	if defsFile, err = os.ReadFile("defs.gomn"); err != nil {
		log.Fatalf("failed to read defs.gomn:\n  %v", err)
	} else { log.Debug("defs.gomn found") }

	if defsGlob, err = gomn.Parse(string(defsFile)); err != nil {
		log.Fatal(err)
	} else { log.Debug("parsed config") ; log.Debug("global defs set") }

	if rcDefs, ok = defsGlob[0].(gomn.Map); !ok {
		log.Warn("rc definitions not found, may produce odd results")
		kilOcont("continuing anyways")
	} else {
		log.Debug("found rc defs")

		if killOnWarn, _ = rcDefs["kill on warn"].(bool); killOnWarn {
			log.Info("configured to kill on warn")
		} else { log.Debug("not configured to kill on warn") }

		if headEnd, ok = rcDefs["head end"].(string); !ok {
			kilOcont("\"head end\" not defined in rc definitions, this will probably cause problems")
		} else { log.Debug("end of header string def found") }

		if writeFile, ok = rcDefs["write to file"].(bool); !ok {
			log.Info("configured to not write to file")
		} else { log.Debug("configured to write output to file") }

		if printOut, ok = rcDefs["print output"].(bool); !ok {
			log.Info("configured to not print output")
			if !writeFile {
				kilOcont("no output configured")
			} else { log.Debug("output to file only") }
		} else { log.Debug("configured to print output to stdout") } 

		if fileExt, ok = rcDefs["output file extension"].(string); !ok {
			kilOcont("no file extention configured")
		} else { log.Debug("file extension def found") }

		if headDefs, ok = rcDefs["head defs"].(gomn.Map); !ok {
			kilOcont("head defs not defined")
		} else { log.Debug("head defs found") }

		if importsMap, ok = headDefs["imports"].(gomn.Map); !ok {
			kilOcont("imports map not found, could be a non-problem")
		} else if importDefs, ok = importsMap["defs"].(gomn.Map); !ok {
			kilOcont("imports map found, but no definitions found") 
		} else { log.Debug("import map and defs found") }

		if debug, ok = rcDefs["debug"].(bool); debug {
			log.SetLevel(log.DebugLevel)
			log.Info("debug mode enabled")
		} else { log.Debug("debug mode not enabled") }
	}
}

func appOut(old []string, cond bool, newVal string, oldVal string) []string {
	log.Debug("appOut()")

	if cond {
		log.Debug("condition true, returning with new value appended")
		return append(old, newVal)
	} else if !cont { //wanted to put a funny fatal message
		log.Debug("condition false, returning with old value appended") 
		return append(old, oldVal)
	} else { log.Fatal("shrodinger's boolean") }

	log.Debug("returning blank string slice")
	return []string{}
}

func whitespaceSplitter(r rune) bool {
	if r == '"' {
		log.Debug("flipping isString bool")
		isString = !isString
		log.Debug("returning as false")
		return false
	} else if isString {
		log.Debug("isString == true, returning false")
		return false
	} else {
		log.Debug("not a string")
		switch r {
		case '\n', ' ':
			log.Debug("detected as newline or space")
			log.Debug("appending string-ed rune")
			splitters = append(splitters, string(r))
			log.Debug("returning true")
			return true
			break
		case '.':
			log.Debug("detected as dot")
			log.Debug("appending dot")
			splitters = append(splitters, ".")
			log.Debug("returning false")
			return false
		default:
			break
		}
	}
	log.Debug("not matched as whitespace or dot, returning false")
	return false
}

func kilOcont(str string) {
	log.Debug("kilOcont()")
	log.Warn(str)
	if killOnWarn {
		log.Debug("exiting...")
		os.Exit(1)
	} else {
		log.Debug("killOnWarn == false")
		log.Info("continuing anyways")
	}
}
