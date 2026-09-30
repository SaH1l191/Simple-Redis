package main

import (
	"fmt"
	"simple-redis/store"
	"time"
)

type MetricsMiddleware struct {
	store store.Storer

	getCalls    int
	setCalls    int
	deleteCalls int
	lenCalls    int

	getMisses       int
	totalGetLatency time.Duration
	totalSetLatency time.Duration
}

func NewMetricsMiddleware(store store.Storer) *MetricsMiddleware {
	return &MetricsMiddleware{
		store: store,
	}
}

func (m *MetricsMiddleware) Get(key string) (string, error) {
	start := time.Now()
	val, err := m.store.Get(key)
	if err != nil {
		m.getMisses++
	}
	m.getCalls++

	duration := time.Since(start)
	m.totalGetLatency += duration
	return val, err
}

func (m *MetricsMiddleware) Set(key, value string) error {

	start := time.Now()

	err := m.store.Set(key, value)
	m.setCalls++

	timeSinceStart := time.Since(start)
	m.totalSetLatency += timeSinceStart
	return err
}

func (m *MetricsMiddleware) Delete(key string) {
	m.store.Delete(key)
	m.deleteCalls++
}

func (m *MetricsMiddleware) Len() int {
	m.lenCalls++
	return m.store.Len()
}

func (m *MetricsMiddleware) Keys() []string {
	return m.store.Keys()
}

func (m *MetricsMiddleware) Report() {
	fmt.Printf("GET calls: %d (missed %d)\n", m.getCalls, m.getMisses)
	fmt.Printf("SET calls: %d\n", m.setCalls)
	fmt.Printf("DEL calls: %d\n", m.deleteCalls)
	fmt.Printf("LEN calls: %d\n", m.lenCalls)

	if m.getCalls > 0 {
		avgGetLatency := float64(m.totalGetLatency) / float64(m.getCalls)
		fmt.Printf("Average GET latency: %.2f ms\n", avgGetLatency)
	}

	if m.setCalls > 0 {
		avgSetLatency := float64(m.totalSetLatency) / float64(m.setCalls)
		fmt.Printf("Average SET latency: %.2f ms\n", avgSetLatency)
	}

}
