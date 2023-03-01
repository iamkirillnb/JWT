package core

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

func NewExporter(sn string) *ocagent.Exporter {
	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(os.Getenv("OC_AGENT_HOST")),
		ocagent.WithServiceName(sn))
	if err != nil {
		log.Fatalf("Failed to create ocagent-exporter: %v", err)
		return nil
	}
	return oce
}
