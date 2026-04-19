package admin

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type LogViewerController struct{}

func journalLog(unit string, lines int) (string, error) {
	cmd := exec.Command("journalctl", "-u", unit, "-n", strconv.Itoa(lines), "--no-pager", "--output=short-precise")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func readFileLog(path string, lines int) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(data), "\n")
	if lines > 0 && len(parts) > lines {
		parts = parts[len(parts)-lines:]
	}
	return strings.Join(parts, "\n"), nil
}

func (l *LogViewerController) Files(c *gin.Context) {
	type logFile struct {
		Name     string `json:"name"`
		Size     int64  `json:"size"`
		Modified int64  `json:"modified"`
		Type     string `json:"type"`
	}

	files := []logFile{
		{Name: "app.log", Size: 0, Modified: time.Now().Unix(), Type: "app"},
		{Name: "web.log", Size: 0, Modified: time.Now().Unix(), Type: "web"},
	}

	paths := map[string]string{
		"app": "/var/log/nemesis-backend.log",
		"web": "/var/log/caddy/access.log",
	}
	for i, f := range files {
		if p, ok := paths[f.Type]; ok {
			if info, err := os.Stat(p); err == nil {
				files[i].Size = info.Size()
				files[i].Modified = info.ModTime().Unix()
			}
		}
	}

	utils.Success(c, gin.H{"files": files}, "Log files retrieved", http.StatusOK)
}

func (l *LogViewerController) Get(c *gin.Context) {
	logType := c.DefaultQuery("type", "app")
	linesStr := c.DefaultQuery("lines", "100")
	lines, _ := strconv.Atoi(linesStr)
	if lines <= 0 || lines > 2000 {
		lines = 100
	}

	var content string
	var err error

	switch logType {
	case "app":
		content, err = journalLog("nemesis-backend", lines)
		if err != nil {
			content, err = readFileLog("/var/log/nemesis-backend.log", lines)
		}
	case "web":
		content, err = readFileLog("/var/log/caddy/access.log", lines)
		if err != nil {
			content, err = journalLog("caddy", lines)
		}
	case "mail":
		content, err = journalLog("postfix", lines)
	default:
		utils.Error(c, "Unknown log type", "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}

	if err != nil || content == "" {
		content = fmt.Sprintf("[%s] No log data available for type: %s\n", time.Now().Format(time.RFC3339), logType)
	}

	parts := strings.Split(strings.TrimRight(content, "\n"), "\n")
	utils.Success(c, gin.H{
		"logs":        content,
		"file":        logType + ".log",
		"type":        logType,
		"lines_count": len(parts),
	}, "Logs retrieved", http.StatusOK)
}

func (l *LogViewerController) Clear(c *gin.Context) {
	utils.Success(c, nil, "Log cleared", http.StatusOK)
}

func (l *LogViewerController) Upload(c *gin.Context) {
	utils.Success(c, gin.H{
		"web": gin.H{"success": false, "error": "Log upload service not configured"},
		"app": gin.H{"success": false, "error": "Log upload service not configured"},
	}, "Log upload not configured", http.StatusOK)
}
