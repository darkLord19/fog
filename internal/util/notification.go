package util

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Notify sends a desktop notification on macOS.
func Notify(title, message, sessionID string) {
	if runtime.GOOS != "darwin" {
		return
	}

	title = strings.TrimSpace(title)
	message = strings.TrimSpace(message)
	sessionID = strings.TrimSpace(sessionID)

	if title == "" || message == "" {
		return
	}

	subtitle := ""
	if sessionID != "" {
		subtitle = sessionID
		if len(subtitle) > 8 {
			subtitle = subtitle[:8]
		}
		subtitle = "Session " + subtitle
	}

	script := ""
	if subtitle != "" {
		script = fmt.Sprintf("display notification %q with title %q subtitle %q", message, title, subtitle)
	} else {
		script = fmt.Sprintf("display notification %q with title %q", message, title)
	}

	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Start(); err != nil {
		return
	}
	go func() { _ = cmd.Wait() }()
}
