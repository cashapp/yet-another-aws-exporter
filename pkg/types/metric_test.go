package types

import "testing"

func TestMetric_PrefixMetricName(t *testing.T) {
	tests := []struct {
		name   string
		metric *Metric
		want   string
	}{
		{
			name: "Adds `aws_` prefix to scraper names",
			metric: &Metric{
				Name: "foo",
			},
			want: "aws_foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metric.PrefixMetricName(); got != tt.want {
				t.Errorf("Scraper.PrefixMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}
