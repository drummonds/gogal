package gogal

import (
	"fmt"
	"io"
)

// svgWriter wraps an io.Writer for generating SVG markup.
type svgWriter struct {
	w   io.Writer
	err error
}

func newSVGWriter(w io.Writer) *svgWriter {
	return &svgWriter{w: w}
}

func (sw *svgWriter) printf(format string, args ...any) {
	if sw.err != nil {
		return
	}
	_, sw.err = fmt.Fprintf(sw.w, format, args...)
}

func (sw *svgWriter) openSVG(width, height float64, accessible bool, title string) {
	sw.printf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f"`,
		width, height)
	if accessible {
		sw.printf(` role="img"`)
		if title != "" {
			sw.printf(` aria-labelledby="chart-title"`)
		}
	}
	sw.printf(">\n")
}

func (sw *svgWriter) closeSVG() {
	sw.printf("</svg>\n")
}

func (sw *svgWriter) title(text string) {
	sw.printf(`  <title id="chart-title">%s</title>`+"\n", xmlEscape(text))
}

func (sw *svgWriter) desc(text string) {
	sw.printf(`  <desc>%s</desc>`+"\n", xmlEscape(text))
}

func (sw *svgWriter) style(css string) {
	sw.printf("  <style>\n%s  </style>\n", css)
}

func (sw *svgWriter) rect(x, y, w, h float64, attrs string) {
	sw.printf(`  <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f"%s/>`+"\n",
		x, y, w, h, attrs)
}

func (sw *svgWriter) line(x1, y1, x2, y2 float64, attrs string) {
	sw.printf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"%s/>`+"\n",
		x1, y1, x2, y2, attrs)
}

func (sw *svgWriter) path(d string, attrs string) {
	sw.printf(`  <path d="%s"%s/>`+"\n", d, attrs)
}

func (sw *svgWriter) circle(cx, cy, r float64, attrs string) {
	sw.printf(`  <circle cx="%.2f" cy="%.2f" r="%.1f"%s/>`+"\n", cx, cy, r, attrs)
}

func (sw *svgWriter) text(x, y float64, text string, attrs string) {
	sw.printf(`  <text x="%.2f" y="%.2f"%s>%s</text>`+"\n",
		x, y, attrs, xmlEscape(text))
}

func (sw *svgWriter) openGroup(attrs string) {
	sw.printf(`  <g%s>`+"\n", attrs)
}

func (sw *svgWriter) closeGroup() {
	sw.printf("  </g>\n")
}

// xmlEscape escapes special XML characters.
func xmlEscape(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			result = append(result, []byte("&amp;")...)
		case '<':
			result = append(result, []byte("&lt;")...)
		case '>':
			result = append(result, []byte("&gt;")...)
		case '"':
			result = append(result, []byte("&quot;")...)
		case '\'':
			result = append(result, []byte("&#39;")...)
		default:
			result = append(result, s[i])
		}
	}
	return string(result)
}
