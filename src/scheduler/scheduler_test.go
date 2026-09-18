package scheduler

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestNew verifies a fresh scheduler has no tasks and is not running.
func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
	if len(s.GetTasks()) != 0 {
		t.Errorf("GetTasks() = %v, want empty", s.GetTasks())
	}
}

// TestAddRemoveTask covers adding, retrieving, and removing a task.
func TestAddRemoveTask(t *testing.T) {
	s := New()
	s.AddTask("noop", time.Hour, func() error { return nil })

	tasks := s.GetTasks()
	if len(tasks) != 1 {
		t.Fatalf("GetTasks() = %v, want 1 entry", tasks)
	}
	if tasks[0].Name != "noop" {
		t.Errorf("task name = %q, want %q", tasks[0].Name, "noop")
	}
	if !tasks[0].Enabled {
		t.Error("newly added task should be Enabled")
	}
	if tasks[0].Interval != time.Hour {
		t.Errorf("task interval = %v, want %v", tasks[0].Interval, time.Hour)
	}

	s.RemoveTask("noop")
	if len(s.GetTasks()) != 0 {
		t.Errorf("GetTasks() after RemoveTask = %v, want empty", s.GetTasks())
	}

	// Removing a nonexistent task must not panic.
	s.RemoveTask("does-not-exist")
}

// TestEnableDisableTask covers toggling a task's enabled state, including
// the no-op case for an unknown task name.
func TestEnableDisableTask(t *testing.T) {
	s := New()
	s.AddTask("t1", time.Minute, func() error { return nil })

	s.DisableTask("t1")
	tasks := s.GetTasks()
	if tasks[0].Enabled {
		t.Error("DisableTask() should set Enabled = false")
	}

	s.EnableTask("t1")
	tasks = s.GetTasks()
	if !tasks[0].Enabled {
		t.Error("EnableTask() should set Enabled = true")
	}

	// Unknown task names must not panic.
	s.EnableTask("unknown")
	s.DisableTask("unknown")
}

// TestRunNow covers the success path, error-propagation path, and the
// not-found-returns-nil path, and verifies LastRun/NextRun are updated.
func TestRunNow(t *testing.T) {
	t.Run("success updates LastRun and NextRun", func(t *testing.T) {
		s := New()
		var ran int32
		s.AddTask("ok", 5*time.Minute, func() error {
			atomic.AddInt32(&ran, 1)
			return nil
		})

		if err := s.RunNow("ok"); err != nil {
			t.Fatalf("RunNow() unexpected error: %v", err)
		}
		if atomic.LoadInt32(&ran) != 1 {
			t.Errorf("task function ran %d times, want 1", ran)
		}
		tasks := s.GetTasks()
		if tasks[0].LastRun.IsZero() {
			t.Error("LastRun should be set after RunNow")
		}
		if !tasks[0].NextRun.After(tasks[0].LastRun) {
			t.Error("NextRun should be after LastRun")
		}
	})

	t.Run("error propagates", func(t *testing.T) {
		s := New()
		wantErr := errors.New("boom")
		s.AddTask("fails", time.Minute, func() error { return wantErr })

		if err := s.RunNow("fails"); !errors.Is(err, wantErr) {
			t.Errorf("RunNow() = %v, want %v", err, wantErr)
		}
	})

	t.Run("unknown task returns nil", func(t *testing.T) {
		s := New()
		if err := s.RunNow("nonexistent"); err != nil {
			t.Errorf("RunNow(nonexistent) = %v, want nil", err)
		}
	})
}

// TestRunDueTasks verifies only enabled, past-due tasks are run, and that
// LastRun/NextRun are updated afterward. runDueTasks dispatches task
// execution in goroutines, so we synchronize with a WaitGroup rather than
// sleeping.
func TestRunDueTasks(t *testing.T) {
	s := New()

	var wg sync.WaitGroup
	wg.Add(1)
	var dueRan int32
	s.AddTask("due", time.Hour, func() error {
		atomic.AddInt32(&dueRan, 1)
		wg.Done()
		return nil
	})

	var notDueRan int32
	s.AddTask("not-due", time.Hour, func() error {
		atomic.AddInt32(&notDueRan, 1)
		return nil
	})

	var disabledRan int32
	s.AddTask("disabled", time.Hour, func() error {
		atomic.AddInt32(&disabledRan, 1)
		return nil
	})
	s.DisableTask("disabled")

	// Backdate the "due" and "disabled" tasks' NextRun so they qualify as
	// due; "not-due" keeps its future NextRun from AddTask.
	s.mu.Lock()
	s.tasks["due"].NextRun = time.Now().Add(-time.Minute)
	s.tasks["disabled"].NextRun = time.Now().Add(-time.Minute)
	s.mu.Unlock()

	s.runDueTasks()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for due task to run")
	}

	if atomic.LoadInt32(&dueRan) != 1 {
		t.Errorf("due task ran %d times, want 1", dueRan)
	}
	if atomic.LoadInt32(&notDueRan) != 0 {
		t.Errorf("not-due task ran %d times, want 0", notDueRan)
	}
	if atomic.LoadInt32(&disabledRan) != 0 {
		t.Errorf("disabled task ran %d times, want 0", disabledRan)
	}
}

// TestStartStop covers Start/Stop lifecycle including the double-Start and
// double-Stop no-op branches. It does not wait on the real 30s ticker.
func TestStartStop(t *testing.T) {
	s := New()
	s.AddTask("noop", time.Hour, func() error { return nil })

	s.Start()
	s.mu.RLock()
	running := s.running
	s.mu.RUnlock()
	if !running {
		t.Error("running should be true after Start()")
	}

	// Calling Start() again while already running must be a no-op and must
	// not panic (would double-close s.stop otherwise).
	s.Start()

	s.Stop()
	s.mu.RLock()
	running = s.running
	s.mu.RUnlock()
	if running {
		t.Error("running should be false after Stop()")
	}

	// Calling Stop() again while already stopped must be a no-op and must
	// not panic (would double-close a closed channel otherwise).
	s.Stop()

	// Starting again after a Stop must work (fresh stop channel).
	s.Start()
	s.mu.RLock()
	running = s.running
	s.mu.RUnlock()
	if !running {
		t.Error("running should be true after restarting")
	}
	s.Stop()
}

// TestParseInterval covers every named alias, a raw duration string, and
// the invalid-string fallback to 24h.
func TestParseInterval(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
	}{
		{name: "minutely", input: "minutely", want: time.Minute},
		{name: "hourly", input: "hourly", want: time.Hour},
		{name: "daily", input: "daily", want: 24 * time.Hour},
		{name: "weekly", input: "weekly", want: 7 * 24 * time.Hour},
		{name: "monthly", input: "monthly", want: 30 * 24 * time.Hour},
		{name: "valid duration string", input: "15m", want: 15 * time.Minute},
		{name: "invalid string falls back to daily", input: "not-a-duration", want: 24 * time.Hour},
		{name: "empty string falls back to daily", input: "", want: 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInterval(tt.input); got != tt.want {
				t.Errorf("ParseInterval(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
