package runner

import (
	"testing"

	"github.com/darkLord19/foglet/internal/state"
)

func TestNotificationsEnabledDefaultsFalse(t *testing.T) {
	r := &Runner{}
	if r.notificationsEnabled() {
		t.Fatal("expected notifications to be disabled by default")
	}
}

func TestNotificationsEnabledReadsSetting(t *testing.T) {
	dir := t.TempDir()
	store, err := state.NewStore(dir)
	if err != nil {
		t.Fatalf("new state store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	r := &Runner{state: store}
	if r.notificationsEnabled() {
		t.Fatal("expected notifications to be disabled when setting is missing")
	}

	if err := store.SetSetting("default_notify", "true"); err != nil {
		t.Fatalf("set default_notify: %v", err)
	}
	if !r.notificationsEnabled() {
		t.Fatal("expected notifications to be enabled when default_notify is true")
	}

	if err := store.SetSetting("default_notify", "FALSE"); err != nil {
		t.Fatalf("set default_notify: %v", err)
	}
	if r.notificationsEnabled() {
		t.Fatal("expected notifications to be disabled when default_notify is false")
	}
}

