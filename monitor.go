package unittest

import (
	"io/ioutil"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

// Command handler
var CommandHandler func(string) bool

func init() {
	go func() {
		for {
			if input, err := ioutil.ReadFile("unittest.cmd"); err == nil && len(input) > 0 {
				ioutil.WriteFile("unittest.cmd", []byte(""), 0744)

				cmd := strings.Trim(string(input), " \n\r\t")

				var (
					profile  *pprof.Profile
					filename string
				)

				switch cmd {
				case "lookup goroutine":
					profile = pprof.Lookup("goroutine")
					filename = "unittest.goroutine"
				case "lookup heap":
					profile = pprof.Lookup("heap")
					filename = "unittest.heap"
				case "lookup threadcreate":
					profile = pprof.Lookup("threadcreate")
					filename = "unittest.thread"
				default:
					if CommandHandler == nil || !CommandHandler(cmd) {
						println("unknow command: '" + cmd + "'")
					}
				}

				if profile != nil {
					file, err := os.Create(filename)
					if err != nil {
						println("couldn't create " + filename)
					} else {
						profile.WriteTo(file, 2)
					}
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
