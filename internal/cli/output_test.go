package cli

import (
	"testing"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:  "empty input",
			input: []string{},
			want:  map[string]interface{}{},
		},
		{
			name:  "string value",
			input: []string{"format=png"},
			want:  map[string]interface{}{"format": "png"},
		},
		{
			name:  "integer value",
			input: []string{"vw=1920"},
			want:  map[string]interface{}{"vw": 1920},
		},
		{
			name:  "float value",
			input: []string{"deviceScale=1.5"},
			want:  map[string]interface{}{"deviceScale": 1.5},
		},
		{
			name:  "boolean true",
			input: []string{"fullPage=true"},
			want:  map[string]interface{}{"fullPage": true},
		},
		{
			name:  "boolean false",
			input: []string{"darkMode=false"},
			want:  map[string]interface{}{"darkMode": false},
		},
		{
			name:  "multiple options",
			input: []string{"vw=1920", "vh=1080", "fullPage=true", "format=webp"},
			want: map[string]interface{}{
				"vw":       1920,
				"vh":       1080,
				"fullPage": true,
				"format":   "webp",
			},
		},
		{
			name:  "value with equals sign",
			input: []string{"selector=div.class=value"},
			want:  map[string]interface{}{"selector": "div.class=value"},
		},
		{
			name:    "invalid format no equals",
			input:   []string{"invalid"},
			wantErr: true,
		},
		{
			name:    "invalid format empty key",
			input:   []string{"=value"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOptions(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("parseOptions() got %d options, want %d", len(got), len(tt.want))
				return
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("parseOptions()[%s] = %v (%T), want %v (%T)", k, got[k], got[k], v, v)
				}
			}
		})
	}
}
