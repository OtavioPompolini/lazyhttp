package vim

import (
	"strings"
	"unicode"
)

// isWordChar returns true for letters, digits, and underscore (vim's \w).
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isWhitespace returns true for space/tab.
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

// runesOf converts text to a rune slice.
func runesOf(text string) []rune {
	return []rune(text)
}

// byteOffsetToRuneIndex converts a byte offset to a rune index.
func byteOffsetToRuneIndex(text string, byteOffset int) int {
	if byteOffset <= 0 {
		return 0
	}
	idx := 0
	for i := range text {
		if i >= byteOffset {
			break
		}
		idx++
	}
	return idx
}

// runeIndexToByteOffset converts a rune index to a byte offset.
func runeIndexToByteOffset(runes []rune, runeIdx int) int {
	if runeIdx <= 0 {
		return 0
	}
	return len(string(runes[:runeIdx]))
}

// WordBounds returns the byte range [start, end) of the word/whitespace run
// containing cursorByte.
//
// inner=true matches vim's "iw": includes the word plus trailing whitespace on
// the same line (or leading whitespace if no trailing), matching vim exactly.
// When cursor is on whitespace, returns the contiguous whitespace run.
//
// inner=false matches vim's "aw": includes word plus all adjacent whitespace.
func WordBounds(text string, cursorByte int, inner bool) (start, end int) {
	runes := runesOf(text)
	if len(runes) == 0 {
		return 0, 0
	}
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes) {
		ri = len(runes) - 1
	}

	cur := runes[ri]
	if isWhitespace(cur) {
		// On whitespace: return the contiguous whitespace run (same for iw/aw).
		s := ri
		for s > 0 && isWhitespace(runes[s-1]) {
			s--
		}
		e := ri
		for e < len(runes)-1 && isWhitespace(runes[e+1]) {
			e++
		}
		return runeIndexToByteOffset(runes, s), runeIndexToByteOffset(runes, e+1)
	}

	if isWordChar(cur) {
		// Find word boundaries.
		s := ri
		for s > 0 && isWordChar(runes[s-1]) {
			s--
		}
		e := ri
		for e < len(runes)-1 && isWordChar(runes[e+1]) {
			e++
		}
		if inner {
			// iw: prefer trailing whitespace on the same line; fall back to leading.
			if e+1 < len(runes) && isWhitespace(runes[e+1]) {
				for e < len(runes)-1 && isWhitespace(runes[e+1]) {
					e++
				}
			} else if s > 0 && isWhitespace(runes[s-1]) {
				for s > 0 && isWhitespace(runes[s-1]) {
					s--
				}
			}
		} else {
			// aw: include all surrounding whitespace.
			for e < len(runes)-1 && isWhitespace(runes[e+1]) {
				e++
			}
			for s > 0 && isWhitespace(runes[s-1]) {
				s--
			}
		}
		return runeIndexToByteOffset(runes, s), runeIndexToByteOffset(runes, e+1)
	}

	// Punctuation/symbol run.
	s := ri
	for s > 0 && !isWordChar(runes[s-1]) && !isWhitespace(runes[s-1]) {
		s--
	}
	e := ri
	for e < len(runes)-1 && !isWordChar(runes[e+1]) && !isWhitespace(runes[e+1]) {
		e++
	}
	if inner {
		if e+1 < len(runes) && isWhitespace(runes[e+1]) {
			for e < len(runes)-1 && isWhitespace(runes[e+1]) {
				e++
			}
		} else if s > 0 && isWhitespace(runes[s-1]) {
			for s > 0 && isWhitespace(runes[s-1]) {
				s--
			}
		}
	} else {
		for e < len(runes)-1 && isWhitespace(runes[e+1]) {
			e++
		}
		for s > 0 && isWhitespace(runes[s-1]) {
			s--
		}
	}
	return runeIndexToByteOffset(runes, s), runeIndexToByteOffset(runes, e+1)
}

// WORDBounds returns the byte range [start, end) of the WORD (non-whitespace
// run) containing cursorByte. inner flag works same as WordBounds (iW/aW).
func WORDBounds(text string, cursorByte int, inner bool) (start, end int) {
	runes := runesOf(text)
	if len(runes) == 0 {
		return 0, 0
	}
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes) {
		ri = len(runes) - 1
	}

	cur := runes[ri]
	if isWhitespace(cur) {
		// On whitespace: return whitespace run.
		s := ri
		for s > 0 && isWhitespace(runes[s-1]) {
			s--
		}
		e := ri
		for e < len(runes)-1 && isWhitespace(runes[e+1]) {
			e++
		}
		return runeIndexToByteOffset(runes, s), runeIndexToByteOffset(runes, e+1)
	}

	// Non-whitespace (WORD) run.
	s := ri
	for s > 0 && !isWhitespace(runes[s-1]) {
		s--
	}
	e := ri
	for e < len(runes)-1 && !isWhitespace(runes[e+1]) {
		e++
	}
	if inner {
		// iW: prefer trailing whitespace; fall back to leading.
		if e+1 < len(runes) && isWhitespace(runes[e+1]) {
			for e < len(runes)-1 && isWhitespace(runes[e+1]) {
				e++
			}
		} else if s > 0 && isWhitespace(runes[s-1]) {
			for s > 0 && isWhitespace(runes[s-1]) {
				s--
			}
		}
	} else {
		// aW: include all surrounding whitespace.
		for e < len(runes)-1 && isWhitespace(runes[e+1]) {
			e++
		}
		for s > 0 && isWhitespace(runes[s-1]) {
			s--
		}
	}
	return runeIndexToByteOffset(runes, s), runeIndexToByteOffset(runes, e+1)
}

// LineBounds returns the byte range [start, end) of the line containing
// cursorByte. inner=true excludes the newline character.
func LineBounds(text string, cursorByte int, inner bool) (start, end int) {
	if cursorByte > len(text) {
		cursorByte = len(text)
	}
	// Find start of line
	s := cursorByte
	for s > 0 && text[s-1] != '\n' {
		s--
	}
	// Find end of line
	e := cursorByte
	for e < len(text) && text[e] != '\n' {
		e++
	}
	if !inner && e < len(text) {
		e++ // include newline
	}
	return s, e
}

// QuoteBounds returns the byte range [start, end) of the quoted string
// surrounding cursorByte. Returns ok=false if no enclosing quotes found.
// inner=true excludes the quote characters.
func QuoteBounds(text string, cursorByte int, quote rune, inner bool) (start, end int, ok bool) {
	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes) {
		ri = len(runes) - 1
	}

	// Search left for opening quote
	left := -1
	for i := ri; i >= 0; i-- {
		if runes[i] == quote {
			left = i
			break
		}
		if runes[i] == '\n' {
			break
		}
	}
	if left < 0 {
		return 0, 0, false
	}

	// Search right for closing quote
	right := -1
	for i := left + 1; i < len(runes); i++ {
		if runes[i] == quote {
			right = i
			break
		}
		if runes[i] == '\n' {
			break
		}
	}
	if right < 0 {
		return 0, 0, false
	}

	if inner {
		return runeIndexToByteOffset(runes, left+1), runeIndexToByteOffset(runes, right), true
	}
	return runeIndexToByteOffset(runes, left), runeIndexToByteOffset(runes, right+1), true
}

// BracketBounds returns the byte range [start, end) of the bracket pair
// surrounding cursorByte. open is one of '(', '{', '[', '<'.
// inner=true excludes the bracket characters.
func BracketBounds(text string, cursorByte int, open rune, inner bool) (start, end int, ok bool) {
	var close rune
	switch open {
	case '(':
		close = ')'
	case '{':
		close = '}'
	case '[':
		close = ']'
	case '<':
		close = '>'
	default:
		return 0, 0, false
	}

	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes) {
		ri = len(runes) - 1
	}

	// Search left for opening bracket (track depth)
	depth := 0
	left := -1
	for i := ri; i >= 0; i-- {
		if runes[i] == close {
			depth++
		} else if runes[i] == open {
			if depth == 0 {
				left = i
				break
			}
			depth--
		}
	}
	if left < 0 {
		return 0, 0, false
	}

	// Search right for matching close bracket
	depth = 0
	right := -1
	for i := left + 1; i < len(runes); i++ {
		if runes[i] == open {
			depth++
		} else if runes[i] == close {
			if depth == 0 {
				right = i
				break
			}
			depth--
		}
	}
	if right < 0 {
		return 0, 0, false
	}

	if inner {
		return runeIndexToByteOffset(runes, left+1), runeIndexToByteOffset(runes, right), true
	}
	return runeIndexToByteOffset(runes, left), runeIndexToByteOffset(runes, right+1), true
}

// WordForwardByte returns the byte offset of the start of the next word.
func WordForwardByte(text string, cursorByte int) int {
	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes)-1 {
		return len(text)
	}

	cur := runes[ri]
	// Skip current word/punct chars
	if isWordChar(cur) {
		for ri < len(runes) && isWordChar(runes[ri]) {
			ri++
		}
	} else if !isWhitespace(cur) {
		for ri < len(runes) && !isWordChar(runes[ri]) && !isWhitespace(runes[ri]) {
			ri++
		}
	}
	// Skip whitespace
	for ri < len(runes) && isWhitespace(runes[ri]) {
		ri++
	}
	return runeIndexToByteOffset(runes, ri)
}

// WordEndByte returns the byte offset of the end of the current/next word.
func WordEndByte(text string, cursorByte int) int {
	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes)-1 {
		return len(text)
	}

	// Move one step forward first
	ri++
	// Skip whitespace
	for ri < len(runes) && isWhitespace(runes[ri]) {
		ri++
	}
	if ri >= len(runes) {
		return len(text)
	}

	cur := runes[ri]
	if isWordChar(cur) {
		for ri < len(runes)-1 && isWordChar(runes[ri+1]) {
			ri++
		}
	} else {
		for ri < len(runes)-1 && !isWordChar(runes[ri+1]) && !isWhitespace(runes[ri+1]) {
			ri++
		}
	}
	return runeIndexToByteOffset(runes, ri+1)
}

// WordBackByte returns the byte offset of the start of the previous word.
func WordBackByte(text string, cursorByte int) int {
	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri <= 0 {
		return 0
	}

	ri-- // step back one
	// Skip whitespace
	for ri > 0 && isWhitespace(runes[ri]) {
		ri--
	}
	if isWhitespace(runes[ri]) {
		return runeIndexToByteOffset(runes, ri)
	}

	cur := runes[ri]
	if isWordChar(cur) {
		for ri > 0 && isWordChar(runes[ri-1]) {
			ri--
		}
	} else {
		for ri > 0 && !isWordChar(runes[ri-1]) && !isWhitespace(runes[ri-1]) {
			ri--
		}
	}
	return runeIndexToByteOffset(runes, ri)
}

// WORDForwardByte returns the byte offset of the start of the next WORD (non-whitespace run).
func WORDForwardByte(text string, cursorByte int) int {
	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes)-1 {
		return len(text)
	}
	// Skip current WORD
	for ri < len(runes) && !isWhitespace(runes[ri]) {
		ri++
	}
	// Skip whitespace
	for ri < len(runes) && isWhitespace(runes[ri]) {
		ri++
	}
	return runeIndexToByteOffset(runes, ri)
}

// WORDEndByte returns the byte offset of the end of the current/next WORD.
func WORDEndByte(text string, cursorByte int) int {
	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri >= len(runes)-1 {
		return len(text)
	}
	ri++
	// Skip whitespace
	for ri < len(runes) && isWhitespace(runes[ri]) {
		ri++
	}
	if ri >= len(runes) {
		return len(text)
	}
	// Advance to end of WORD
	for ri < len(runes)-1 && !isWhitespace(runes[ri+1]) {
		ri++
	}
	return runeIndexToByteOffset(runes, ri+1)
}

// WORDBackByte returns the byte offset of the start of the previous WORD.
func WORDBackByte(text string, cursorByte int) int {
	runes := runesOf(text)
	ri := byteOffsetToRuneIndex(text, cursorByte)
	if ri <= 0 {
		return 0
	}
	ri--
	// Skip whitespace
	for ri > 0 && isWhitespace(runes[ri]) {
		ri--
	}
	// Move to start of WORD
	for ri > 0 && !isWhitespace(runes[ri-1]) {
		ri--
	}
	return runeIndexToByteOffset(runes, ri)
}

// ByteOffsetFromLineCol converts a 0-based (line, col) pair to a byte offset.
func ByteOffsetFromLineCol(text string, line, col int) int {
	lines := strings.Split(text, "\n")
	if line < 0 {
		line = 0
	}
	if line >= len(lines) {
		line = len(lines) - 1
	}
	offset := 0
	for i := 0; i < line; i++ {
		offset += len(lines[i]) + 1 // +1 for \n
	}
	lineRunes := []rune(lines[line])
	if col < 0 {
		col = 0
	}
	if col > len(lineRunes) {
		col = len(lineRunes)
	}
	offset += len(string(lineRunes[:col]))
	return offset
}

// LineColFromByteOffset converts a byte offset to 0-based (line, col).
func LineColFromByteOffset(text string, offset int) (line, col int) {
	if offset > len(text) {
		offset = len(text)
	}
	if offset < 0 {
		offset = 0
	}
	line = strings.Count(text[:offset], "\n")
	lastNL := strings.LastIndex(text[:offset], "\n")
	if lastNL < 0 {
		col = len([]rune(text[:offset]))
	} else {
		col = len([]rune(text[lastNL+1 : offset]))
	}
	return
}
