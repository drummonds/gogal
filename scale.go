package gogal

import (
	"fmt"
	"math"
	"time"
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

// TemporalScale wraps LinearScale and formats tick labels as human-readable times.
// Ticks are stepped in calendar units (seconds, minutes, hours, days, months,
// years) and aligned to those units in loc, so a 14-minute span gives
// 15:00, 15:05, 15:10 rather than 200-second multiples.
type TemporalScale struct {
	*LinearScale
	timeFormat string
	loc        *time.Location
}

// NewTemporalScale creates a temporal scale with the given time format.
// Tick labels are formatted in loc; nil means time.Local.
func NewTemporalScale(ls *LinearScale, timeFormat string, loc *time.Location) *TemporalScale {
	if loc == nil {
		loc = time.Local
	}
	return &TemporalScale{LinearScale: ls, timeFormat: timeFormat, loc: loc}
}

func (s *TemporalScale) Ticks() []Tick {
	step := temporalStep(s.domainMax-s.domainMin, s.tickCount)
	first := time.Unix(int64(math.Ceil(s.domainMin)), 0).In(s.loc)
	last := time.Unix(int64(math.Floor(s.domainMax)), 0).In(s.loc)

	var ticks []Tick
	for t := step.align(first); !t.After(last); t = step.next(t) {
		v := float64(t.Unix())
		ticks = append(ticks, Tick{
			Value:    v,
			Position: s.Map(v),
			Label:    t.Format(s.timeFormat),
		})
	}
	return ticks
}

// timeStep is one rung of the temporal tick ladder: either a fixed duration
// (sub-day steps) or a whole number of calendar months or years.
type timeStep struct {
	duration time.Duration // for steps below one month
	months   int
	years    int
}

// seconds is the nominal length of the step, used to pick a rung.
func (st timeStep) seconds() float64 {
	switch {
	case st.years > 0:
		return float64(st.years) * 365 * 86400
	case st.months > 0:
		return float64(st.months) * 30 * 86400
	default:
		return st.duration.Seconds()
	}
}

// temporalLadder lists the tick steps to choose from, smallest first.
var temporalLadder = func() []timeStep {
	var ladder []timeStep
	for _, d := range []time.Duration{
		time.Second, 2 * time.Second, 5 * time.Second, 10 * time.Second, 15 * time.Second, 30 * time.Second,
		time.Minute, 2 * time.Minute, 5 * time.Minute, 10 * time.Minute, 15 * time.Minute, 30 * time.Minute,
		time.Hour, 2 * time.Hour, 3 * time.Hour, 6 * time.Hour, 12 * time.Hour,
		24 * time.Hour, 2 * 24 * time.Hour, 7 * 24 * time.Hour, 14 * 24 * time.Hour,
	} {
		ladder = append(ladder, timeStep{duration: d})
	}
	for _, m := range []int{1, 2, 3, 6} {
		ladder = append(ladder, timeStep{months: m})
	}
	for _, y := range []int{1, 2, 5, 10, 20, 50, 100} {
		ladder = append(ladder, timeStep{years: y})
	}
	return ladder
}()

// temporalStep picks the smallest ladder step that yields at most
// targetTicks intervals across a span of the given seconds.
func temporalStep(spanSeconds float64, targetTicks int) timeStep {
	if targetTicks <= 0 {
		targetTicks = 1
	}
	for _, st := range temporalLadder {
		if spanSeconds/st.seconds() <= float64(targetTicks) {
			return st
		}
	}
	return temporalLadder[len(temporalLadder)-1]
}

// align returns the first tick at or after t, on a step boundary in t's location.
func (st timeStep) align(t time.Time) time.Time {
	loc := t.Location()
	switch {
	case st.years > 0:
		y := t.Year() / st.years * st.years
		aligned := time.Date(y, 1, 1, 0, 0, 0, 0, loc)
		if aligned.Before(t) {
			aligned = aligned.AddDate(st.years, 0, 0)
		}
		return aligned
	case st.months > 0:
		m := (int(t.Month())-1)/st.months*st.months + 1
		aligned := time.Date(t.Year(), time.Month(m), 1, 0, 0, 0, 0, loc)
		if aligned.Before(t) {
			aligned = aligned.AddDate(0, st.months, 0)
		}
		return aligned
	case st.duration >= 24*time.Hour:
		aligned := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		if aligned.Before(t) {
			aligned = aligned.AddDate(0, 0, 1)
		}
		return aligned
	default:
		// Step from local midnight so ticks land on wall-clock boundaries.
		midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		n := math.Ceil(float64(t.Sub(midnight)) / float64(st.duration))
		return midnight.Add(time.Duration(n) * st.duration)
	}
}

// next advances t by one step.
func (st timeStep) next(t time.Time) time.Time {
	switch {
	case st.years > 0:
		return t.AddDate(st.years, 0, 0)
	case st.months > 0:
		return t.AddDate(0, st.months, 0)
	case st.duration >= 24*time.Hour:
		return t.AddDate(0, 0, int(st.duration/(24*time.Hour)))
	default:
		return t.Add(st.duration)
	}
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
