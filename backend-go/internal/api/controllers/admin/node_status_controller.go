package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type NodeStatusController struct{}

type wingsSystemInfo struct {
	Architecture  string  `json:"architecture"`
	CPUCount      int     `json:"cpu_count"`
	KernelVersion string  `json:"kernel_version"`
	OS            string  `json:"os"`
	MemoryBytes   int64   `json:"memory_bytes"`
	DiskBytes     int64   `json:"disk_bytes"`
	Utilization   wingsUtilization `json:"utilization"`
}

type wingsUtilization struct {
	MemoryBytes    int64   `json:"memory_bytes"`
	DiskBytes      int64   `json:"disk_bytes"`
	CPUPercent     float64 `json:"cpu_percent"`
	NetworkRxBytes int64   `json:"network_rx_bytes"`
	NetworkTxBytes int64   `json:"network_tx_bytes"`
}

type nodeStatusResult struct {
	ID          uint        `json:"id"`
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	FQDN        string      `json:"fqdn"`
	LocationID  int         `json:"location_id"`
	Status      string      `json:"status"`
	Utilization interface{} `json:"utilization"`
	Error       *string     `json:"error"`
}

func queryWingsStatus(node models.Node) (wingsSystemInfo, error) {
	url := fmt.Sprintf("%s://%s:%d/api/system", node.Scheme, node.FQDN, node.DaemonListen)
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return wingsSystemInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+node.DaemonToken)
	resp, err := client.Do(req)
	if err != nil {
		return wingsSystemInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return wingsSystemInfo{}, fmt.Errorf("wings returned %d", resp.StatusCode)
	}
	var info wingsSystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return wingsSystemInfo{}, err
	}
	return info, nil
}

func (n *NodeStatusController) GlobalStatus(c *gin.Context) {
	var nodes []models.Node
	database.DB.Find(&nodes)

	type nodeUtil struct {
		MemoryTotal   int64   `json:"memory_total"`
		MemoryUsed    int64   `json:"memory_used"`
		DiskTotal     int64   `json:"disk_total"`
		DiskUsed      int64   `json:"disk_used"`
		SwapTotal     int64   `json:"swap_total"`
		SwapUsed      int64   `json:"swap_used"`
		CPUPercent    float64 `json:"cpu_percent"`
		LoadAverage1  float64 `json:"load_average1"`
		LoadAverage5  float64 `json:"load_average5"`
		LoadAverage15 float64 `json:"load_average15"`
	}

	var results []nodeStatusResult
	var totalMemory, usedMemory, totalDisk, usedDisk int64
	var totalCPU float64
	healthyCount, unhealthyCount := 0, 0

	for _, node := range nodes {
		info, err := queryWingsStatus(node)
		if err != nil {
			errStr := err.Error()
			results = append(results, nodeStatusResult{
				ID: node.ID, UUID: node.UUID, Name: node.Name, FQDN: node.FQDN,
				LocationID: node.LocationID, Status: "unhealthy", Utilization: nil, Error: &errStr,
			})
			unhealthyCount++
			continue
		}

		memUsed := info.Utilization.MemoryBytes
		diskUsed := info.Utilization.DiskBytes
		totalMemory += info.MemoryBytes
		usedMemory += memUsed
		totalDisk += info.DiskBytes
		usedDisk += diskUsed
		totalCPU += info.Utilization.CPUPercent
		healthyCount++

		util := nodeUtil{
			MemoryTotal: info.MemoryBytes,
			MemoryUsed:  memUsed,
			DiskTotal:   info.DiskBytes,
			DiskUsed:    diskUsed,
			CPUPercent:  info.Utilization.CPUPercent,
		}
		results = append(results, nodeStatusResult{
			ID: node.ID, UUID: node.UUID, Name: node.Name, FQDN: node.FQDN,
			LocationID: node.LocationID, Status: "healthy", Utilization: util, Error: nil,
		})
	}

	avgCPU := 0.0
	if healthyCount > 0 {
		avgCPU = totalCPU / float64(healthyCount)
	}

	global := gin.H{
		"total_nodes":    len(nodes),
		"healthy_nodes":  healthyCount,
		"unhealthy_nodes": unhealthyCount,
		"total_memory":   totalMemory,
		"used_memory":    usedMemory,
		"total_disk":     totalDisk,
		"used_disk":      usedDisk,
		"avg_cpu_percent": avgCPU,
	}

	if results == nil {
		results = []nodeStatusResult{}
	}
	utils.Success(c, gin.H{"global": global, "nodes": results}, "Node status retrieved", http.StatusOK)
}
