// utils/metrics.go

package utils

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	queueLength = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "task_queue_length",
		Help: "Current length of the task queue",
	})

	tasksProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tasks_processed_total",
		Help: "Total number of processed tasks",
	})
)

func init() {
	prometheus.MustRegister(queueLength)
	prometheus.MustRegister(tasksProcessed)
}

func IncrementQueueLength() {
	queueLength.Inc()
}

func DecrementQueueLength() {
	queueLength.Dec()
}

func SetQueueLength(value float64) {
	queueLength.Set(value)
}

func IncrementTasksProcessed() {
	tasksProcessed.Inc()
}
