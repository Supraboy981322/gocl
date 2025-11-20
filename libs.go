package main

import (
	"os"
	"slices"
	"github.com/Supraboy981322/gomn"
	"github.com/charmbracelet/log"
)

func checkArgs() {
	var taken []int
	for i, arg := range args {
		if !slices.Contains(taken, i) {
			var isArg bool
			for k, a := range arg {
				if a == '-' && !isArg {
					isArg = true
				} else if isArg {
					switch a {
       		case 'i':
						inputFile = args[i+1]
						taken = append(taken, i+1)
					case 'o':
						outputFile = args[i+1]
						taken = append(taken, i+1)
					case '-':
					  isArg = false
						taken = append(taken, checkFullArg(i, arg)...)
					default:
						invArg(k, a, arg)
					}
				}
			}
		}
	}

	if inputFile == "" {
		log.Fatal("no file input")
	}
}

func checkFullArg(current int, arg string) []int {
	var taken []int
	switch arg {
	case "--input", "--i", "--in", "--source", "--s":
		inputFile = args[current+1]
		taken = append(taken, current+1)
	case "--output", "--o", "--out":
		outputFile = args[current+1]
		taken = append(taken, current+1)
	default:
		er.Fatal("invalid arg:  \033[1;31m" + arg + "\033[0m")
	}
	return taken
}

func invArg(index int, which rune, arg string) { 
	var pointer string
	for i := 0; i < index+19; i++ {
		pointer += " "
	}
	pointer += "\033[1;31m^\033[0m"
	arg = arg[:index] + "\033[1;31m" +
			string(arg[index]) + "\033[0m" +
			arg[index+1:]
	er.Fatal("invalid arg:  " + arg + "\n" +
			pointer + "\n")
}

func readConf() {
	var ok bool
	var defsFile []byte
	var err error

	if defsFile, err = os.ReadFile("defs.gomn"); err != nil {
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
			kilOcont("\"head end\" not defined in rc definitions, this will probably cause problems")
		}
		if writeFile, ok = rcDefs["write to file"].(bool); !ok {
			log.Info("configured to not write to file")
		}
		if printOut, ok = rcDefs["print output"].(bool); !ok {
			log.Info("configured to not print output")
			if !writeFile {
				kilOcont("no output configured")
			}
		}
		if fileExt, ok = rcDefs["output file extension"].(string); !ok {
			kilOcont("no file extention configured")
		}
		if headDefs, ok = rcDefs["head defs"].(gomn.Map); !ok {
			kilOcont("head defs not defined")
		}
		if importsMap, ok = headDefs["imports"].(gomn.Map); !ok {
			kilOcont("imports map not found, could be a non-problem")
		} else if importDefs, ok = importsMap["defs"].(gomn.Map); !ok {
			kilOcont("imports map found, but no definitions found") 
		}
		if debug, ok = rcDefs["debug"].(bool); ok {
			log.Info("debug mode enabled")
		}
	}
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
	} else if isString {
		return false
	} else {
		switch r {
		case '\n', ' ':
			splitters = append(splitters, string(r))
			return true
			break
		case '.':
			splitters = append(splitters, ".")
			return false
		default:
			break
		}
	}
	return false
}

func kilOcont(str string) {
	log.Warn(str)
	if killOnWarn {
		os.Exit(1)
	} else {
		log.Info("continuing anyways")
	}
}
