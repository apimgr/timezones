package timezones

import (
	"encoding/json"
	"reflect"
	"testing"
)

// fixtureJSON is a small, hand-crafted timezone dataset covering the cases
// exercised below: multiple UTC identifiers, DST/non-DST entries, duplicate
// abbreviations, and a negative offset.
const fixtureJSON = `[
	{
		"value": "Eastern Standard Time",
		"abbr": "EST",
		"offset": -5,
		"isdst": false,
		"text": "(UTC-05:00) Eastern Time (US & Canada)",
		"utc": ["America/New_York", "America/Detroit"]
	},
	{
		"value": "Pacific Standard Time",
		"abbr": "PST",
		"offset": -8,
		"isdst": true,
		"text": "(UTC-08:00) Pacific Time (US & Canada)",
		"utc": ["America/Los_Angeles"]
	},
	{
		"value": "Central Standard Time",
		"abbr": "PST",
		"offset": -6,
		"isdst": true,
		"text": "(UTC-06:00) Central Time (US & Canada)",
		"utc": ["America/Chicago"]
	}
]`

func newTestService(t *testing.T) *Service {
	t.Helper()
	svc, err := NewService([]byte(fixtureJSON))
	if err != nil {
		t.Fatalf("NewService() unexpected error: %v", err)
	}
	return svc
}

// TestNewService covers the happy path and the malformed-JSON error path.
func TestNewService(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		svc, err := NewService([]byte(fixtureJSON))
		if err != nil {
			t.Fatalf("NewService() unexpected error: %v", err)
		}
		if svc.Count() != 3 {
			t.Errorf("Count() = %d, want 3", svc.Count())
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		_, err := NewService([]byte(`{not valid json`))
		if err == nil {
			t.Fatal("NewService() expected error for malformed JSON, got nil")
		}
	})

	t.Run("empty array", func(t *testing.T) {
		svc, err := NewService([]byte(`[]`))
		if err != nil {
			t.Fatalf("NewService() unexpected error: %v", err)
		}
		if svc.Count() != 0 {
			t.Errorf("Count() = %d, want 0", svc.Count())
		}
	})
}

// TestCount verifies the count matches the fixture size.
func TestCount(t *testing.T) {
	svc := newTestService(t)
	if got := svc.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}
}

// TestGetRawJSON verifies the original bytes are returned unmodified.
func TestGetRawJSON(t *testing.T) {
	svc := newTestService(t)
	if string(svc.GetRawJSON()) != fixtureJSON {
		t.Error("GetRawJSON() did not return the original bytes")
	}
}

// TestGetAll verifies every fixture entry round-trips correctly.
func TestGetAll(t *testing.T) {
	svc := newTestService(t)
	all := svc.GetAll()
	if len(all) != 3 {
		t.Fatalf("GetAll() returned %d entries, want 3", len(all))
	}

	var want []Timezone
	if err := json.Unmarshal([]byte(fixtureJSON), &want); err != nil {
		t.Fatalf("failed to unmarshal fixture for comparison: %v", err)
	}
	for i := range want {
		if !reflect.DeepEqual(all[i], want[i]) {
			t.Errorf("GetAll()[%d] = %+v, want %+v", i, all[i], want[i])
		}
	}
}

// TestSearch covers matches on value, abbr, text and UTC identifier, case
// insensitivity, and the no-match empty-result case.
func TestSearch(t *testing.T) {
	svc := newTestService(t)

	tests := []struct {
		name      string
		query     string
		wantCount int
	}{
		{name: "match value", query: "eastern", wantCount: 1},
		{name: "match abbr case-insensitive", query: "pst", wantCount: 2},
		{name: "match text substring", query: "us & canada", wantCount: 3},
		{name: "match utc identifier", query: "america/chicago", wantCount: 1},
		{name: "no match", query: "nonexistent-zone", wantCount: 0},
		{name: "empty query matches all", query: "", wantCount: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Search(tt.query)
			if len(got) != tt.wantCount {
				t.Errorf("Search(%q) returned %d results, want %d", tt.query, len(got), tt.wantCount)
			}
		})
	}
}

// TestGetByOffset covers a matching offset and a non-matching offset.
func TestGetByOffset(t *testing.T) {
	svc := newTestService(t)

	got := svc.GetByOffset(-5)
	if len(got) != 1 || got[0].Value != "Eastern Standard Time" {
		t.Errorf("GetByOffset(-5) = %+v, want single Eastern entry", got)
	}

	if got := svc.GetByOffset(999); len(got) != 0 {
		t.Errorf("GetByOffset(999) = %+v, want empty", got)
	}
}

// TestGetByAbbr covers a matching abbreviation (including duplicates and
// case-insensitivity) and a non-matching abbreviation.
func TestGetByAbbr(t *testing.T) {
	svc := newTestService(t)

	got := svc.GetByAbbr("pst")
	if len(got) != 2 {
		t.Errorf("GetByAbbr(\"pst\") returned %d results, want 2", len(got))
	}

	got = svc.GetByAbbr("EST")
	if len(got) != 1 || got[0].Value != "Eastern Standard Time" {
		t.Errorf("GetByAbbr(\"EST\") = %+v, want single Eastern entry", got)
	}

	if got := svc.GetByAbbr("ZZZ"); len(got) != 0 {
		t.Errorf("GetByAbbr(\"ZZZ\") = %+v, want empty", got)
	}
}

// TestGetByUTC covers a matching UTC identifier (case-insensitive) and the
// nil-on-no-match case.
func TestGetByUTC(t *testing.T) {
	svc := newTestService(t)

	got := svc.GetByUTC("america/new_york")
	if got == nil || got.Value != "Eastern Standard Time" {
		t.Errorf("GetByUTC(\"america/new_york\") = %+v, want Eastern Standard Time", got)
	}

	if got := svc.GetByUTC("Nowhere/Nothing"); got != nil {
		t.Errorf("GetByUTC(\"Nowhere/Nothing\") = %+v, want nil", got)
	}
}

// TestGetByValue covers a matching value (case-insensitive) and the
// nil-on-no-match case.
func TestGetByValue(t *testing.T) {
	svc := newTestService(t)

	got := svc.GetByValue("pacific standard time")
	if got == nil || got.Abbr != "PST" || got.Offset != -8 {
		t.Errorf("GetByValue(\"pacific standard time\") = %+v, want Pacific entry", got)
	}

	if got := svc.GetByValue("does not exist"); got != nil {
		t.Errorf("GetByValue(\"does not exist\") = %+v, want nil", got)
	}
}

// TestGetStats verifies every stat key against the known fixture.
func TestGetStats(t *testing.T) {
	svc := newTestService(t)
	stats := svc.GetStats()

	if got := stats["total_timezones"]; got != 3 {
		t.Errorf("total_timezones = %v, want 3", got)
	}
	// UTC entries: 2 + 1 + 1 = 4
	if got := stats["total_utc_entries"]; got != 4 {
		t.Errorf("total_utc_entries = %v, want 4", got)
	}
	// DST entries: PST + Central = 2
	if got := stats["dst_timezones"]; got != 2 {
		t.Errorf("dst_timezones = %v, want 2", got)
	}
	if got := stats["non_dst_timezones"]; got != 1 {
		t.Errorf("non_dst_timezones = %v, want 1", got)
	}
}
