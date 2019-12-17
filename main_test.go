package main

import (
	"testing"
	"time"
)

func parseTime(v string) *time.Time {
	t, _ := time.Parse(time.RFC1123Z, v)
	return &t
}

func TestFormatExpirationTime(t *testing.T) {

	var testCases = []struct {
		in          *time.Time
		out         string
		description string
	}{
		{
			parseTime("Tue, 17 Dec 2019 12:25:28 -0500"),
			"Tue 12:25",
			"Test base case",
		},
		{
			parseTime("Tue, 17 Dec 2019 17:25:28 -0000"),
			"Tue 12:25",
			"Test time in UTC",
		},
		{
			parseTime("Tue, 17 Dec 2019 18:25:28 -0000"),
			"Tue 13:25",
			"Test time in after 1:00 PM UTC",
		},
		{
			nil,
			"",
			"Test when input is nil",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			actualOutput := formatExpirationTime(tt.in)
			if actualOutput != tt.out {
				t.Errorf("got %v, want %v", actualOutput, tt.out)
			}
		})
	}

}
