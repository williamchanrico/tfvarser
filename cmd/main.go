package main

import (
	"flag"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/williamchanrico/tfvarser/cmd/tfvarser"
)

func main() {
	appFlags := &tfvarser.Flags{}
	flag.StringVar(&appFlags.Provider, "provider", "ali", "Cloud provider to do a terraform import from")
	flag.StringVar(&appFlags.ProviderObj, "obj", "ess", "Object in cloud provider to import")
	flag.StringVar(&appFlags.LimitNames, "limit-names", "", "Limit generation of objects containing the specified names (separated by comma or space)")
	flag.StringVar(&appFlags.LimitIDs, "limit-ids", "", "Limit generation of objects with the exact IDs (separated by comma or space)")
	flag.Parse()

	var cfg tfvarser.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	rc, err := tfvarser.Run(appFlags, cfg)
	if err != nil {
		log.Println("tfvarser exited with error: ", err)
	}

	log.Printf("tfvarser exited with return code %d\n", rc)
	os.Exit(rc)
}
