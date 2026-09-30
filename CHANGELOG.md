# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

## [0.2.0] - 2026-09-30

 - Y tick cap (WithMaxYTicks) and non-finite domain guard in LinearScale.Ticks

### Added
- `WithPoints(false)` draws no per-point markers (nor their CSS), cutting a 1100-point line from ~150 KB to ~20 KB
- `WithIntegerY(true)` steps Y ticks in whole numbers, for counts
- `WithEndLabels(true)` pins the first and last dates to the ends of a temporal X axis, dropping calendar ticks that would collide with them
- `WithMaxYTicks(n)` caps the number of labelled Y ticks (at least 3), using the finest 1-2-5 step that fits
- `LinearScale.Nice`, `LinearScale.SetIntegerTicks`, `LinearScale.SetMaxTicks`, `TemporalScale.TicksWithEnds`

### Changed
- Line and step charts widen the Y domain to whole ticks, so both ends of the Y axis are labelled

### Fixed
- Temporal ticks are never finer than the time format can show: a date-only format no longer repeats the same date at sub-day ticks
- Plain XY data (no times) gets a numeric X axis instead of 1970-01-01 labels
- `LinearScale.Ticks` returns no ticks for an infinite or NaN domain instead of appending NaN ticks until memory runs out

## [0.1.8] - 2026-09-25

 - Deterministic example 06 screenshot

### Fixed
- Example 06 screenshot generator uses a fixed random seed, so docs builds no longer dirty the tree

## [0.1.7] - 2026-09-25

 - Add step chart for balances, WASM demo for example 05, fix 01a screenshot

## [0.1.6] - 2026-09-25

 - Updating htmx demo

### Added
- Step chart: `NewStepChart` holds each value until the next point then jumps, for balances and other event-driven state; example 11

## [0.1.5] - 2026-09-25

 - Add bar and pie charts, fix temporal tick labels

### Added
- Pie chart: `NewPieChart`, `AddSlice`; proportional sectors, legend with percentages, hover titles (#2)
- `WithColors` palette override for any chart type
- Bar chart: `NewBarChart`, `AddCategories`, `WithValueLabels`; vertical bars with a zero baseline, value labels, grouped multi-series and a sparkline variant (#1)

### Fixed
- Temporal axis ticks are formatted in the data's time zone (or `WithLocation`) and stepped in calendar units (5 min, 1 h, 1 day, quarters, years) instead of raw-second multiples (#3)

## [0.1.4] - 2026-09-21

 - Migrate module path from codeberg.org to git.bytestone.uk

## [0.1.3] - 2026-04-07

 - Improve documentation: branded HTML, fixed nav, lofigui examples, url2svg captures

## [0.1.2] - 2026-04-06

 - Updating 06 dmeo

## [0.1.1] - 2026-04-06

 - Adding setttings

## [0.1.0] - 2026-04-05

 - Adding favicon

### Added
- Project scaffolding and research
- Core architecture: Scale interface, Layout/Render separation
