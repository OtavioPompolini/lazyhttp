package vim

import (
	"strings"
	"unicode"
)

// Parser accumulates key sequences and resolves them against a KeyMap.
type Parser struct {
	buf       string // accumulated key sequence (post-leader substitution)
	count     int    // accumulated numeric count (0 = no count typed)
	leaderKey string
}

// Feed processes one raw key string against km.
// Returns:
//
//	action  — resolved Action or ActionNone
//	count   — resolved count (≥1 after resolution; 0 if no count was typed → caller treats as 1)
//	pending — buf is a valid prefix; wait for more keys
//	reset   — sequence discarded (no match and no prefix)
func (p *Parser) Feed(key string, km KeyMap) (action Action, count int, pending bool, reset bool) {
	// Count accumulation: digits accumulate into count.
	// Special case: "0" when count==0 and buf=="" is NOT a count digit — it's the line-start motion.
	if len(key) == 1 && unicode.IsDigit([]rune(key)[0]) {
		digit := int([]rune(key)[0] - '0')
		if digit == 0 && p.count == 0 && p.buf == "" {
			// "0" as motion — fall through to normal processing
		} else {
			p.count = p.count*10 + digit
			return ActionNone, 0, true, false
		}
	}

	p.buf += key

	// Check for exact match
	if act, ok := km[p.buf]; ok {
		resolvedCount := p.count
		p.Reset()
		return act, resolvedCount, false, false
	}

	// Check if buf is a valid prefix of any binding
	for binding := range km {
		if strings.HasPrefix(binding, p.buf) && binding != p.buf {
			return ActionNone, 0, true, false
		}
	}

	// No match and no prefix — reset
	p.Reset()
	return ActionNone, 0, false, true
}

// Reset clears the parser state.
func (p *Parser) Reset() {
	p.buf = ""
	p.count = 0
}

// HasPending returns true if there is an in-progress key sequence.
func (p *Parser) HasPending() bool {
	return p.buf != ""
}

// PendingBuf returns the current accumulated key sequence.
func (p *Parser) PendingBuf() string {
	return p.buf
}

// SetLeader configures the leader key string.
func (p *Parser) SetLeader(key string) {
	p.leaderKey = key
}
