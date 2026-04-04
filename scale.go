package gogal

import (
	"fmt"
	"math"
)

// Tick represents a tick mark on an axis.
type Tick struct {
	Value    float64 // data value
	Position float64 // pixel position within range
	Label    string
}

// Scale maps data values to pixel positions.
type Scale interface {
	// Map converts a data value to a pixel position.
	Map(value float64) float64
	// Inverse converts a pixel position back to a data value.
	Inverse(pixel float64) float64
	// Ticks returns computed tick marks.
	Ticks() []Tick
	// Domain returns [min, max] data range.
	Domain() [2]float64
	// SetRange sets the [min, max] pixel range.
	SetRange(min, max float64)
}

// LinearScale maps a continuous data domain to a pixel range.
type LinearScale struct {
	domainMin, domainMax float64
	rangeMin, rangeMax   float64
	tickCount            int
	format               string
}

// NewLinearScale creates a scale from the given data domain.
func NewLinearScale(domainMin, domainMax float64) *LinearScale {
	if domainMin == domainMax {
		domainMin -= 1
		domainMax += 1
	}
	return &LinearScale{
		domainMin: domainMin,
		domainMax: domainMax,
		tickCount: 5,
		format:    "%.4g",
	}
}

// NewLinearScaleFromData creates a scale that fits the given values.
func NewLinearScaleFromData(values []float64) *LinearScale {
	if len(values) == 0 {
		return NewLinearScale(0, 1)
	}
	min, max := values[0], values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return NewLinearScale(min, max)
}

func (s *LinearScale) SetRange(min, max float64) {
	s.rangeMin = min
	s.rangeMax = max
}

func (s *LinearScale) SetFormat(format string) {
	s.format = format
}

func (s *LinearScale) Map(value float64) float64 {
	if s.domainMax == s.domainMin {
		return (s.rangeMin + s.rangeMax) / 2
	}
	t := (value - s.domainMin) / (s.domainMax - s.domainMin)
	return s.rangeMin + t*(s.rangeMax-s.rangeMin)
}

func (s *LinearScale) Inverse(pixel float64) float64 {
	if s.rangeMax == s.rangeMin {
		return (s.domainMin + s.domainMax) / 2
	}
	t := (pixel - s.rangeMin) / (s.rangeMax - s.rangeMin)
	return s.domainMin + t*(s.domainMax-s.domainMin)
}

func (s *LinearScale) Domain() [2]float64 {
	return [2]float64{s.domainMin, s.domainMax}
}

func (s *LinearScale) Ticks() []Tick {
	step := niceStep(s.domainMax-s.domainMin, s.tickCount)
	start := math.Ceil(s.domainMin/step) * step
	var ticks []Tick
	for v := start; v <= s.domainMax+step*0.001; v += step {
		ticks = append(ticks, Tick{
			Value:    v,
			Position: s.Map(v),
			Label:    fmt.Sprintf(s.format, v),
		})
	}
	return ticks
}

// OrdinalScale maps discrete indices to equal-width bands.
type OrdinalScale struct {
	labels             []string
	rangeMin, rangeMax float64
}

// NewOrdinalScale creates a scale from the given labels.
func NewOrdinalScale(labels []string) *OrdinalScale {
	return &OrdinalScale{labels: labels}
}

func (s *OrdinalScale) SetRange(min, max float64) {
	s.rangeMin = min
	s.rangeMax = max
}

func (s *OrdinalScale) Map(value float64) float64 {
	n := float64(len(s.labels))
	if n == 0 {
		return s.rangeMin
	}
	bandWidth := (s.rangeMax - s.rangeMin) / n
	return s.rangeMin + (value+0.5)*bandWidth
}

func (s *OrdinalScale) Inverse(pixel float64) float64 {
	n := float64(len(s.labels))
	if n == 0 {
		return 0
	}
	bandWidth := (s.rangeMax - s.rangeMin) / n
	return (pixel-s.rangeMin)/bandWidth - 0.5
}

func (s *OrdinalScale) Domain() [2]float64 {
	return [2]float64{0, float64(len(s.labels) - 1)}
}

func (s *OrdinalScale) Ticks() []Tick {
	ticks := make([]Tick, len(s.labels))
	for i, label := range s.labels {
		ticks[i] = Tick{
			Value:    float64(i),
			Position: s.Map(float64(i)),
			Label:    label,
		}
	}
	return ticks
}

// niceStep calculates a "nice" step size for tick marks.
// It rounds to multiples of 1, 2, or 5 × 10^n.
func niceStep(dataRange float64, targetTicks int) float64 {
	if dataRange <= 0 || targetTicks <= 0 {
		return 1
	}
	rough := dataRange / float64(targetTicks)
	exp := math.Floor(math.Log10(rough))
	pow := math.Pow(10, exp)
	frac := rough / pow

	var nice float64
	switch {
	case frac <= 1.5:
		nice = 1
	case frac <= 3.5:
		nice = 2
	case frac <= 7.5:
		nice = 5
	default:
		nice = 10
	}
	return nice * pow
}
