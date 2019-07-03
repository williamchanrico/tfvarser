package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/williamchanrico/tfvarser/cmd/tfvarser"
	"github.com/williamchanrico/tfvarser/log"
)

var (
	version            string
	showVersionAndExit bool
	verbose            bool
)

func main() {
	appFlags := &tfvarser.Flags{}
	flag.StringVar(&appFlags.Provider, "provider", "ali", "Cloud provider to do a terraform import from")
	flag.StringVar(&appFlags.ProviderObj, "obj", "ess", "Object in cloud provider to import")
	flag.StringVar(&appFlags.LimitNames, "limit-names", "", "Limit generation of objects with the exact names (separated by comma or space)")
	flag.StringVar(&appFlags.LimitIDs, "limit-ids", "", "Limit generation of objects with the exact IDs (separated by comma or space)")

	flag.BoolVar(&verbose, "verbose", false, "Show verbose output")
	flag.BoolVar(&showVersionAndExit, "version", false, "Show version and exit")
	flag.Parse()

	if showVersionAndExit {
		fmt.Printf("Tfvars %v\n", version)
		os.Exit(0)
	}

	if verbose {
		log.SetLevelString("debug")
	}

	var cfg tfvarser.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	rc, err := tfvarser.Run(appFlags, cfg)
	if err != nil {
		log.Error("tfvarser exited with error: ", err)
	}

	log.Debugf("tfvarser exited with return code %d\n", rc)
	os.Exit(rc)
}
