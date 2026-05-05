package autoscaler

import "testing"

func TestScalingController(t *testing.T) {
	controller := NewScalingController()

	tests := []struct {
		name       string
		cpu        float64
		traffic    float64
		expected   int
	}{
		{"Low load stable traffic", 15, 120, 0},
		{"High load stable traffic", 85, 200, 4},
		{"High load explosive traffic", 90, 900, 4},
		{"Medium load explosive traffic", 50, 800, 4},
		{"Low load explosive traffic", 20, 700, 4},
		{"Negative traffic", 40, -30, 0},
		{"Above max CPU", 130, 500, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := controller.ComputeScalingAction(tt.cpu, tt.traffic)
			if result != tt.expected {
				t.Errorf("cpu=%.1f traffic=%.1f expected=%d got=%d", tt.cpu, tt.traffic, tt.expected, result)
			}
		})
	}
}
