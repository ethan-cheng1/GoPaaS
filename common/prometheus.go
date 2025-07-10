package common

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"strconv"
)

func PrometheusBoot(port int)  {
	http.Handle("/metrics",promhttp.Handler())
	// Start web service
	go func() {
		err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port),nil)
		if err !=nil {
			log.Fatal("Start failed")
		}
		log.Info("Monitoring started, port: " + strconv.Itoa(port))
	}()
}