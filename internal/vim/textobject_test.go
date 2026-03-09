package vim

import (
	"testing"
)

func TestByteOffsetFromLineCol(t *testing.T) {
	text := "hello\nworld\nfoo"
	tests := []struct {
		line, col    int
		wantOffset   int
	}{
		{0, 0, 0},
		{0, 5, 5},
		{1, 0, 6},
		{1, 5, 11},
		{2, 0, 12},
		{2, 3, 15},
		// Clamping
		{10, 0, 12}, // line clamped to last
		{0, 100, 5}, // col clamped to line length
	}
	for _, tt := range tests {
		got := ByteOffsetFromLineCol(text, tt.line, tt.col)
		if got != tt.wantOffset {
			t.Errorf("ByteOffsetFromLineCol(%q, %d, %d) = %d, want %d", text, tt.line, tt.col, got, tt.wantOffset)
		}
	}
}

func TestLineColFromByteOffset(t *testing.T) {
	text := "hello\nworld\nfoo"
	tests := []struct {
		offset     int
		wantLine   int
		wantCol    int
	}{
		{0, 0, 0},
		{5, 0, 5},
		{6, 1, 0},
		{11, 1, 5},
		{12, 2, 0},
		{15, 2, 3},
	}
	for _, tt := range tests {
		line, col := LineColFromByteOffset(text, tt.offset)
		if line != tt.wantLine || col != tt.wantCol {
			t.Errorf("LineColFromByteOffset(%q, %d) = (%d,%d), want (%d,%d)",
				text, tt.offset, line, col, tt.wantLine, tt.wantCol)
		}
	}
}

func TestWordBounds_InnerWord(t *testing.T) {
	// inner=true matches vim's "iw": word + trailing whitespace (or leading if no trailing).
	tests := []struct {
		name      string
		text      string
		cursor    int // byte offset
		wantStart int
		wantEnd   int
	}{
		{
			// cursor on "hello", trailing space exists → include trailing space
			name:      "first word includes trailing space",
			text:      "hello world",
			cursor:    2,
			wantStart: 0,
			wantEnd:   6, // "hello "
		},
		{
			// cursor on "world", no trailing space → include leading space
			name:      "last word includes leading space",
			text:      "hello world",
			cursor:    7,
			wantStart: 5,
			wantEnd:   11, // " world"
		},
		{
			// cursor on space between words → returns whitespace run
			name:      "on space returns whitespace run",
			text:      "hello   world",
			cursor:    6,
			wantStart: 5,
			wantEnd:   8,
		},
		{
			// single word, no surrounding whitespace → just the word
			name:      "single word no surrounding space",
			text:      "hello",
			cursor:    0,
			wantStart: 0,
			wantEnd:   5,
		},
		{
			// last word at end of file → no trailing → include leading space
			name:      "word at end no trailing",
			text:      "foo bar",
			cursor:    5,
			wantStart: 3,
			wantEnd:   7, // " bar"
		},
		{
			// unicode word with trailing space
			name:      "unicode word trailing space",
			text:      "héllo wörld",
			cursor:    0,
			wantStart: 0,
			wantEnd:   len("héllo "), // "héllo "
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := WordBounds(tt.text, tt.cursor, true)
			if s != tt.wantStart || e != tt.wantEnd {
				t.Errorf("WordBounds(%q, cursor=%d, inner=true) = [%d,%d), want [%d,%d)",
					tt.text, tt.cursor, s, e, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestWordBounds_OuterWord(t *testing.T) {
	// inner=false matches vim's "aw": word + all surrounding whitespace.
	tests := []struct {
		name      string
		text      string
		cursor    int
		wantStart int
		wantEnd   int
	}{
		{
			name:      "first word with surrounding space",
			text:      "hello world",
			cursor:    2,
			wantStart: 0,
			wantEnd:   6, // "hello "
		},
		{
			name:      "last word with leading space",
			text:      "hello world",
			cursor:    7,
			wantStart: 5,
			wantEnd:   11, // " world"
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := WordBounds(tt.text, tt.cursor, false)
			if s != tt.wantStart || e != tt.wantEnd {
				t.Errorf("WordBounds outer(%q, cursor=%d) = [%d,%d), want [%d,%d)",
					tt.text, tt.cursor, s, e, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestQuoteBounds(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		cursor    int
		quote     rune
		inner     bool
		wantStart int
		wantEnd   int
		wantOK    bool
	}{
		{
			name:      "double quotes inner",
			text:      `hello "world" foo`,
			cursor:    9, // inside "world"
			quote:     '"',
			inner:     true,
			wantStart: 7,
			wantEnd:   12,
			wantOK:    true,
		},
		{
			name:      "double quotes outer",
			text:      `hello "world" foo`,
			cursor:    9,
			quote:     '"',
			inner:     false,
			wantStart: 6,
			wantEnd:   13,
			wantOK:    true,
		},
		{
			name:      "no quotes",
			text:      "hello world",
			cursor:    5,
			quote:     '"',
			inner:     true,
			wantOK:    false,
		},
		{
			name:      "single quotes inner",
			text:      "it's 'fine' okay",
			cursor:    7, // inside 'fine'
			quote:     '\'',
			inner:     true,
			wantStart: 6,
			wantEnd:   10,
			wantOK:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e, ok := QuoteBounds(tt.text, tt.cursor, tt.quote, tt.inner)
			if ok != tt.wantOK {
				t.Errorf("QuoteBounds(%q, %d, %c, %v) ok=%v, want %v", tt.text, tt.cursor, tt.quote, tt.inner, ok, tt.wantOK)
				return
			}
			if ok && (s != tt.wantStart || e != tt.wantEnd) {
				t.Errorf("QuoteBounds(%q, %d, %c, %v) = [%d,%d), want [%d,%d)",
					tt.text, tt.cursor, tt.quote, tt.inner, s, e, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestBracketBounds(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		cursor    int
		open      rune
		inner     bool
		wantStart int
		wantEnd   int
		wantOK    bool
	}{
		{
			name:      "simple parens inner",
			text:      "foo(bar)baz",
			cursor:    5, // inside parens
			open:      '(',
			inner:     true,
			wantStart: 4,
			wantEnd:   7,
			wantOK:    true,
		},
		{
			name:      "simple parens outer",
			text:      "foo(bar)baz",
			cursor:    5,
			open:      '(',
			inner:     false,
			wantStart: 3,
			wantEnd:   8,
			wantOK:    true,
		},
		{
			name:      "nested parens",
			text:      "foo((bar))baz",
			cursor:    6, // inside inner parens
			open:      '(',
			inner:     true,
			wantStart: 5,
			wantEnd:   8,
			wantOK:    true,
		},
		{
			name:      "curly braces",
			text:      `{"key": "val"}`,
			cursor:    5,
			open:      '{',
			inner:     true,
			wantStart: 1,
			wantEnd:   13,
			wantOK:    true,
		},
		{
			name:   "no brackets",
			text:   "hello world",
			cursor: 5,
			open:   '(',
			inner:  true,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e, ok := BracketBounds(tt.text, tt.cursor, tt.open, tt.inner)
			if ok != tt.wantOK {
				t.Errorf("BracketBounds(%q, %d, %c, %v) ok=%v, want %v", tt.text, tt.cursor, tt.open, tt.inner, ok, tt.wantOK)
				return
			}
			if ok && (s != tt.wantStart || e != tt.wantEnd) {
				t.Errorf("BracketBounds(%q, %d, %c, %v) = [%d,%d), want [%d,%d)",
					tt.text, tt.cursor, tt.open, tt.inner, s, e, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestWordForwardByte(t *testing.T) {
	tests := []struct {
		text       string
		cursor     int
		wantOffset int
	}{
		{"hello world", 0, 6},    // "hello" → "world"
		{"hello world", 6, 11},   // "world" → end
		{"hello  world", 0, 7},   // skip double space
		{"  hello", 0, 2},        // from leading space
		{"hello", 2, 5},          // end of single word
	}
	for _, tt := range tests {
		got := WordForwardByte(tt.text, tt.cursor)
		if got != tt.wantOffset {
			t.Errorf("WordForwardByte(%q, %d) = %d, want %d", tt.text, tt.cursor, got, tt.wantOffset)
		}
	}
}

func TestWordBackByte(t *testing.T) {
	tests := []struct {
		text       string
		cursor     int
		wantOffset int
	}{
		{"hello world", 6, 0},   // from "world" → start of "hello"
		{"hello world", 11, 6},  // from end → start of "world"
		{"hello world", 5, 0},   // from space → start of "hello"
		{"hello world", 0, 0},   // already at start
	}
	for _, tt := range tests {
		got := WordBackByte(tt.text, tt.cursor)
		if got != tt.wantOffset {
			t.Errorf("WordBackByte(%q, %d) = %d, want %d", tt.text, tt.cursor, got, tt.wantOffset)
		}
	}
}

func TestWordEndByte(t *testing.T) {
	tests := []struct {
		text       string
		cursor     int
		wantOffset int
	}{
		{"hello world", 0, 5},   // end of "hello"
		{"hello world", 5, 11},  // from space → end of "world"
		{"hello world", 6, 11},  // from "w" → end of "world"
	}
	for _, tt := range tests {
		got := WordEndByte(tt.text, tt.cursor)
		if got != tt.wantOffset {
			t.Errorf("WordEndByte(%q, %d) = %d, want %d", tt.text, tt.cursor, got, tt.wantOffset)
		}
	}
}

func TestLineBounds(t *testing.T) {
	text := "hello\nworld\nfoo"
	tests := []struct {
		cursor     int
		inner      bool
		wantStart  int
		wantEnd    int
	}{
		{2, true, 0, 5},    // middle of "hello", inner (no newline)
		{2, false, 0, 6},   // middle of "hello", outer (includes \n)
		{6, true, 6, 11},   // start of "world"
		{14, true, 12, 15}, // "foo" (last line)
		{14, false, 12, 15}, // last line, outer (no trailing \n)
	}
	for _, tt := range tests {
		s, e := LineBounds(text, tt.cursor, tt.inner)
		if s != tt.wantStart || e != tt.wantEnd {
			t.Errorf("LineBounds(%q, cursor=%d, inner=%v) = [%d,%d), want [%d,%d)",
				text, tt.cursor, tt.inner, s, e, tt.wantStart, tt.wantEnd)
		}
	}
}

func TestEmptyText(t *testing.T) {
	// None of these should panic
	WordBounds("", 0, true)
	WORDBounds("", 0, true)
	WordForwardByte("", 0)
	WordBackByte("", 0)
	WordEndByte("", 0)
	LineBounds("", 0, true)
	LineColFromByteOffset("", 0)
	ByteOffsetFromLineCol("", 0, 0)
}
