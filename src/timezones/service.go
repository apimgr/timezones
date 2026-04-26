package timezones

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Timezone represents a single timezone entry
type Timezone struct {
	Value  string   `json:"value"`
	Abbr   string   `json:"abbr"`
	Offset float64  `json:"offset"`
	IsDST  bool     `json:"isdst"`
	Text   string   `json:"text"`
	UTC    []string `json:"utc"`
}

// Service provides timezone lookup functionality
type Service struct {
	timezones []Timezone
	rawJSON   []byte
}

// NewService creates a new timezone service from embedded JSON data
func NewService(jsonData []byte) (*Service, error) {
	var timezones []Timezone
	if err := json.Unmarshal(jsonData, &timezones); err != nil {
		return nil, fmt.Errorf("failed to parse timezones JSON: %w", err)
	}

	return &Service{
		timezones: timezones,
		rawJSON:   jsonData,
	}, nil
}

// Count returns the total number of timezones
func (s *Service) Count() int {
	return len(s.timezones)
}

// GetRawJSON returns the raw JSON data
func (s *Service) GetRawJSON() []byte {
	return s.rawJSON
}

// GetAll returns all timezones
func (s *Service) GetAll() []Timezone {
	return s.timezones
}

// Search searches timezones by query string
func (s *Service) Search(query string) []Timezone {
	query = strings.ToLower(query)
	var results []Timezone

	for _, tz := range s.timezones {
		// Search in value, abbr, text, and UTC identifiers
		if strings.Contains(strings.ToLower(tz.Value), query) ||
			strings.Contains(strings.ToLower(tz.Abbr), query) ||
			strings.Contains(strings.ToLower(tz.Text), query) {
			results = append(results, tz)
			continue
		}

		// Search in UTC identifiers
		for _, utc := range tz.UTC {
			if strings.Contains(strings.ToLower(utc), query) {
				results = append(results, tz)
				break
			}
		}
	}

	return results
}

// GetByOffset returns timezones with a specific UTC offset
func (s *Service) GetByOffset(offset float64) []Timezone {
	var results []Timezone
	for _, tz := range s.timezones {
		if tz.Offset == offset {
			results = append(results, tz)
		}
	}
	return results
}

// GetByAbbr returns timezones matching an abbreviation
func (s *Service) GetByAbbr(abbr string) []Timezone {
	abbr = strings.ToUpper(abbr)
	var results []Timezone
	for _, tz := range s.timezones {
		if strings.ToUpper(tz.Abbr) == abbr {
			results = append(results, tz)
		}
	}
	return results
}

// GetByUTC returns a timezone by UTC identifier
func (s *Service) GetByUTC(utc string) *Timezone {
	utc = strings.ToLower(utc)
	for _, tz := range s.timezones {
		for _, u := range tz.UTC {
			if strings.ToLower(u) == utc {
				return &tz
			}
		}
	}
	return nil
}

// GetByValue returns a timezone by its value name
func (s *Service) GetByValue(value string) *Timezone {
	value = strings.ToLower(value)
	for _, tz := range s.timezones {
		if strings.ToLower(tz.Value) == value {
			return &tz
		}
	}
	return nil
}

// GetStats returns statistics about the timezone data
func (s *Service) GetStats() map[string]interface{} {
	// Count unique UTC identifiers
	utcCount := 0
	for _, tz := range s.timezones {
		utcCount += len(tz.UTC)
	}

	// Count DST timezones
	dstCount := 0
	for _, tz := range s.timezones {
		if tz.IsDST {
			dstCount++
		}
	}

	return map[string]interface{}{
		"total_timezones":   len(s.timezones),
		"total_utc_entries": utcCount,
		"dst_timezones":     dstCount,
		"non_dst_timezones": len(s.timezones) - dstCount,
	}
}
