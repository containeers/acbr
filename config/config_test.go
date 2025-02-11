package config

import "testing"

func TestGetMaxResults(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want int32
	}{
		{
			name: "zero value returns default",
			cfg:  Config{MaxResults: 0},
			want: 50,
		},
		{
			name: "negative value returns default",
			cfg:  Config{MaxResults: -1},
			want: 50,
		},
		{
			name: "value over limit returns default",
			cfg:  Config{MaxResults: 51},
			want: 50,
		},
		{
			name: "valid value returns same",
			cfg:  Config{MaxResults: 30},
			want: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.GetMaxResults(); got != tt.want {
				t.Errorf("Config.GetMaxResults() = %v, want %v", got, tt.want)
			}
		})
	}
}
