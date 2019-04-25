package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/morya/utils/log"
)

var (
	flagConfig   = flag.String("config", "drone.json", "config file")
	flagLogLevel = flag.String("loglevel", "info", "[debug,info,warn,error]")
	flagListen   = flag.String("listen", ":6789", "http webhook listen address")
)

func DumpObject(obj interface{}) string {
	var buff = &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	enc.SetIndent("", "  ")
	enc.Encode(obj)
	return buff.String()
}

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevelString(*flagLogLevel)

	if err := loadConfig(*flagConfig); err != nil {
		log.InfoError(err)
		return
	}

	log.Infof("obj = %s", DumpObject(config))

	drone := newDrone()

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		drone.Stop()
	}()

	drone.Run(*flagListen)
}
