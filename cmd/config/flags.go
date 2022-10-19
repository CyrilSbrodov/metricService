package config

import "flag"

var (
	flagAddress        = "ADDRESS"
	flagRestore        = "RESTORE"
	flagStoreInterval  = "STORE_INTERVAL"
	flagStoreFile      = "STORE_FILE"
	flagPollInterval   = "REPORT_INTERVAL"
	flagReportInterval = "POLL_INTERVAL"
)

func init() {
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address of service")
	flag.StringVar(&flagRestore, "r", "true", "restore from file")
	flag.StringVar(&flagStoreInterval, "i", "300", "upload interval")
	flag.StringVar(&flagStoreFile, "f", "/tmp/devops-metrics-db.json", "name of file")
	flag.StringVar(&flagPollInterval, "p", "2", "update interval")
	flag.StringVar(&flagReportInterval, "r", "10", "upload interval to server")
}
