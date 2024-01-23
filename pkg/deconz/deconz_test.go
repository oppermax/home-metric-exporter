package deconz

import "testing"

func Test_convertLevel(t *testing.T) {
	tests := []struct {
		name  string
		state State
		want  float64
	}{
		{name: "lower than fifty",
			state: State{On: true, Bri: 44},
			want:  17},
		{name: "zero",
			state: State{On: false, Bri: 127},
			want:  0},
		{name: "one hundred",
			state: State{On: true, Bri: 254},
			want:  100},
		{name: "fifty",
			state: State{On: true, Bri: 127},
			want:  50},
		{name: "more than one hundred",
			state: State{On: true, Bri: 3000},
			want:  1181},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertLevel(tt.state); got != tt.want {
				t.Errorf("convertLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
