package main

import (
	"io/ioutil"

	"github.com/pieterclaerhout/go-james"
	"github.com/pieterclaerhout/go-log"
)

func main() {

	args, err := james.ParseBuildArgs()
	log.CheckError(err)

	//log.InfoDump(args, "post_build arguments")
	log.Info("Built loup-gorou", args.Version)
	log.Info("Now copying .env...")

	data, err := ioutil.ReadFile("cmd/loup-gorou/.env")
	log.CheckError(err)
	// Write data to dst
	err = ioutil.WriteFile("build/.env", data, 0644)
	log.CheckError(err)
	log.Info(".env copied to build/")
}
