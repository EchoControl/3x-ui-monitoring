package metrics

import (
	"3x-ui-monitoring/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

type Metrics struct {
	config *config.Config
}

func NewMetrics(cfg *config.Config) *Metrics {
	return &Metrics{config: cfg}
}

var (
	customMetric1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ram_usage",
		Help: "RAM usage of server ${host}",
	})
	customMetric2 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "online_users",
		Help: "Current online users of server ${host}",
	})
	customMetric3 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "CPU usage of server ${host}",
	})
)

func init() {
	prometheus.MustRegister(customMetric1)
	prometheus.MustRegister(customMetric2)
	prometheus.MustRegister(customMetric3)
}

func (m *Metrics) loginToAPI(*http.Client) (*http.Client, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	data := url.Values{
		"username": {m.config.Username},
		"password": {m.config.Password},
	}
	// Send POST request to login

	loginUrl := fmt.Sprintf("https://%s:%s/%s/login", m.config.Host, m.config.Port, m.config.Basepath)

	resp, err := client.PostForm(loginUrl, data)
	if err != nil {
		return nil, fmt.Errorf("error logging in: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}
	// Store cookies from the response
	client.Jar, err = cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %v", err)
	}
	for _, cookie := range resp.Cookies() {
		client.Jar.SetCookies(resp.Request.URL, []*http.Cookie{cookie})
	}
	// Cookies from the login will be stored in the client's jar for reuse
	return client, nil
}

func (m *Metrics) getOnlineCount(client *http.Client) (int, error) {
	if client.Jar == nil {
		return 0, fmt.Errorf("cookie jar is not initialized")
	}
	onlineUrl := fmt.Sprintf("https://%s:%s/%s/%s", m.config.Host, m.config.Port, m.config.Basepath, m.config.Online_path)
	//log.Printf("Fetching onlines from: %s", onlineUrl)

	resp, err := client.PostForm(onlineUrl, url.Values{})
	if err != nil {
		return 0, fmt.Errorf("error fetching onlines: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("fetching onlines failed with status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("error decoding onlines response: %v", err)
	}

	if success, ok := data["success"].(bool); !ok || !success {
		return 0, fmt.Errorf("fetching onlines was not successful")
	}

	obj, ok := data["obj"].([]interface{})
	if !ok || obj == nil {
		return 0, nil
	}

	return len(obj), nil
}

func ramUsage() (float64, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	ramUsage := float64(memStats.Alloc) / (1024 * 1024) // Convert bytes to MB
	return ramUsage, nil
}

func cpuUsage() (float64, error) {
	// Get the percentage of CPU usage over a short interval (e.g., 1 second)
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil // Return the first CPU usage value (overall usage)
	}
	return 0, fmt.Errorf("unable to retrieve CPU usage")
}

func (m *Metrics) updateMetrics() {
	client, err := m.loginToAPI(&http.Client{})
	if err != nil {
		log.Printf("error logging in: %v", err)
		return
	}

	cpuUsage, err := cpuUsage()
	if err != nil {
		log.Printf("error fetching CPU usage: %v", err)
		return
	}
	customMetric3.Set(cpuUsage)

	onlineUsers, err := m.getOnlineCount(client)
	if err != nil {
		log.Printf("error fetching onlines: %v", err)
		return
	}
	customMetric2.Set(float64(onlineUsers))

	ram, err := ramUsage()
	if err != nil {
		log.Printf("error fetching RAM usage: %v", err)
		return
	}
	customMetric1.Set(ram)
}

func (m *Metrics) StartPolling(internal time.Duration) {
	ticker := time.NewTicker(internal)
	defer ticker.Stop()
	log.Println("Starting metrics polling...")
	for range ticker.C {
		// log.Println("Updating metrics...")
		m.updateMetrics()
	}
}
