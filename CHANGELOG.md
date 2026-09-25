# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

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
