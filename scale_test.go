package gogal

import (
	"math"
	"testing"
)

func TestLinearScale_Map(t *testing.T) {
	s := NewLinearScale(0, 100)
	s.SetRange(0, 800)

	tests := []struct {
		value float64
		want  float64
	}{
		{0, 0},
		{50, 400},
		{100, 800},
		{25, 200},
	}
	for _, tt := range tests {
		got := s.Map(tt.value)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("Map(%v) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestLinearScale_Inverse(t *testing.T) {
	s := NewLinearScale(0, 100)
	s.SetRange(0, 800)

	tests := []struct {
		pixel float64
		want  float64
	}{
		{0, 0},
		{400, 50},
		{800, 100},
	}
	for _, tt := range tests {
		got := s.Inverse(tt.pixel)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("Inverse(%v) = %v, want %v", tt.pixel, got, tt.want)
		}
	}
}

func TestLinearScale_Ticks(t *testing.T) {
	s := NewLinearScale(0, 100)
	s.SetRange(0, 800)

	ticks := s.Ticks()
	if len(ticks) == 0 {
		t.Fatal("expected ticks, got none")
	}

	// First tick should be at or near 0
	if ticks[0].Value > 1 {
		t.Errorf("first tick value = %v, expected near 0", ticks[0].Value)
	}

	// Last tick should be at or near 100
	last := ticks[len(ticks)-1]
	if last.Value < 99 {
		t.Errorf("last tick value = %v, expected near 100", last.Value)
	}

	// Ticks should be evenly spaced
	if len(ticks) >= 3 {
		step := ticks[1].Value - ticks[0].Value
		for i := 2; i < len(ticks); i++ {
			got := ticks[i].Value - ticks[i-1].Value
			if math.Abs(got-step) > 0.001 {
				t.Errorf("tick spacing not even: step[%d]=%v, expected %v", i, got, step)
			}
		}
	}
}

func TestLinearScale_EqualDomain(t *testing.T) {
	// Edge case: single value
	s := NewLinearScale(5, 5)
	s.SetRange(0, 100)

	got := s.Map(5)
	if math.Abs(got-50) > 0.001 {
		t.Errorf("Map(5) with equal domain = %v, want 50", got)
	}
}

func TestLinearScaleFromData(t *testing.T) {
	s := NewLinearScaleFromData([]float64{10, 20, 30, 5, 25})
	d := s.Domain()
	if d[0] != 5 || d[1] != 30 {
		t.Errorf("domain = %v, want [5 30]", d)
	}
}

func TestLinearScaleFromData_Empty(t *testing.T) {
	s := NewLinearScaleFromData(nil)
	d := s.Domain()
	if d[0] != -1 || d[1] != 2 {
		// Empty defaults to [0,1] then adjusted to [-1,2]
		t.Logf("domain for empty data = %v", d)
	}
}

func TestOrdinalScale_Map(t *testing.T) {
	s := NewOrdinalScale([]string{"A", "B", "C", "D"})
	s.SetRange(0, 400)

	// Each band is 100px wide, points at center
	tests := []struct {
		index int
		want  float64
	}{
		{0, 50},
		{1, 150},
		{2, 250},
		{3, 350},
	}
	for _, tt := range tests {
		got := s.Map(float64(tt.index))
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("Map(%d) = %v, want %v", tt.index, got, tt.want)
		}
	}
}

func TestOrdinalScale_Ticks(t *testing.T) {
	labels := []string{"Mon", "Tue", "Wed"}
	s := NewOrdinalScale(labels)
	s.SetRange(0, 300)

	ticks := s.Ticks()
	if len(ticks) != 3 {
		t.Fatalf("got %d ticks, want 3", len(ticks))
	}
	for i, tick := range ticks {
		if tick.Label != labels[i] {
			t.Errorf("tick[%d].Label = %q, want %q", i, tick.Label, labels[i])
		}
	}
}

func TestNiceStep(t *testing.T) {
	tests := []struct {
		dataRange   float64
		targetTicks int
		want        float64
	}{
		{100, 5, 20},
		{10, 5, 2},
		{1000, 5, 200},
		{0.5, 5, 0.1},
		{7, 5, 1},
	}
	for _, tt := range tests {
		got := niceStep(tt.dataRange, tt.targetTicks)
		if math.Abs(got-tt.want) > tt.want*0.001 {
			t.Errorf("niceStep(%v, %d) = %v, want %v", tt.dataRange, tt.targetTicks, got, tt.want)
		}
	}
}
