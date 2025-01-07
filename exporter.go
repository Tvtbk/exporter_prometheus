package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"time"
)

var (
	saveRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "save_requests_total",
			Help: "Total number of Save requests",
		},
		[]string{"method"},
	)

	cpuUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percentage",
			Help: "Current CPU usage percentage",
		},
	)

	memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_percentage",
			Help: "Current memory usage percentage",
		},
	)
)

func init() {
	// Регистрация метрик в Prometheus
	prometheus.MustRegister(saveRequests)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsage)
}

func recordMetrics() {
	go func() {
		for {
			// Обновление метрик о потреблении ресурсов
			cpuPercent, _ := cpu.Percent(0, false)
			memStats, _ := mem.VirtualMemory()

			cpuUsage.Set(cpuPercent[0])
			memoryUsage.Set(memStats.UsedPercent)

			time.Sleep(5 * time.Second)
		}
	}()
}

//todo: Implement me
//func recordSaveRequest() {
//	saveRequests.WithLabelValues("POST").Inc()
//}

func main() {
	recordMetrics()

	http.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/url/save" {
			//recordSaveRequest()
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Save request recorded"))
			if err != nil {
				log.Fatal(fmt.Sprintf("Error writing response: %s", err))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":9101", nil))
}
