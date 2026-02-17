package runner

import "strings"

func (r *Runner) notificationsEnabled() bool {
	if r == nil || r.state == nil {
		return false
	}
	notify, found, err := r.state.GetSetting("default_notify")
	if err != nil || !found {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(notify), "true")
}
