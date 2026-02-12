package daemon

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const defaultHealthTimeout = 15 * time.Second

// EnsureRunning checks /health and starts fogd if needed.
func EnsureRunning(fogHome string, port int, timeout time.Duration) (string, error) {
	if timeout <= 0 {
		timeout = defaultHealthTimeout
	}

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	healthURL := baseURL + "/health"

	if isHealthy(healthURL, 2*time.Second) {
		return baseURL, nil
	}

	if err := startFogd(fogHome, port); err != nil {
		return "", err
	}

	if err := waitForHealth(healthURL, timeout); err != nil {
		return "", err
	}

	return baseURL, nil
}

func startFogd(fogHome string, port int) error {
	logsDir := filepath.Join(fogHome, "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return fmt.Errorf("create logs dir: %w", err)
	}

	logFile := filepath.Join(logsDir, "fogd.log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open fogd log file: %w", err)
	}
	defer f.Close()

	cmd := exec.Command("fogd", "--port", strconv.Itoa(port))
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.Env = os.Environ()
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start fogd: %w", err)
	}

	return nil
}

func waitForHealth(healthURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if isHealthy(healthURL, 2*time.Second) {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("fogd health check timed out: %s", healthURL)
}

func isHealthy(healthURL string, timeout time.Duration) bool {
	client := &http.Client{Timeout: timeout}
	return isHealthyWithClient(client, healthURL)
}

func isHealthyWithClient(client *http.Client, healthURL string) bool {
	resp, err := client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// OpenBrowser opens URL in the system browser (macOS/Linux).
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform for auto-open: %s", runtime.GOOS)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}
	return nil
}
