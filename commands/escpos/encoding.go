package escpos

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// CharsetEncoder handles character encoding conversion from UTF-8 to legacy charsets.
// Supports auto-detection for Latin characters (Portuguese, Spanish, French, German, Italian).
type CharsetEncoder struct {
	currentCharset int
	encoder        *encoding.Encoder
}

// NewCharsetEncoder creates a new character encoder with default charset (PC437).
func NewCharsetEncoder() *CharsetEncoder {
	e := &CharsetEncoder{
		currentCharset: chartablePC437,
	}
	e.updateEncoder()
	return e
}

// SetCharset sets the current character encoding table.
func (e *CharsetEncoder) SetCharset(charset int) {
	e.currentCharset = charset
	e.updateEncoder()
}

// updateEncoder updates the internal encoder based on current charset.
func (e *CharsetEncoder) updateEncoder() {
	e.encoder = encoderForCharset(e.currentCharset)
}

// encoderForCharset creates a fresh encoder for the given charset page
// without mutating any state.
func encoderForCharset(charset int) *encoding.Encoder {
	switch charset {
	case chartablePC850:
		return charmap.CodePage850.NewEncoder()
	case chartablePC860:
		return charmap.CodePage860.NewEncoder()
	case chartablePC858:
		return charmap.CodePage858.NewEncoder()
	default:
		return charmap.CodePage437.NewEncoder()
	}
}

// CharsetSwitcher defines the interface for charset switching on the display.
type CharsetSwitcher interface {
	SetCharacterCodeTableInternal(page int) error
}

// EncodeTextWithAutoCharsetSwitching encodes UTF-8 text for a VFD display,
// automatically detecting the best charset and switching hardware if needed.
//
// The function tries the current charset first, then auto-detects a better one.
// Unrepresentable characters are replaced with '?' rather than sending raw UTF-8
// bytes that the display firmware cannot interpret.
func (e *CharsetEncoder) EncodeTextWithAutoCharsetSwitching(text string, display CharsetSwitcher) ([]byte, error) {
	if !utf8.ValidString(text) {
		return []byte(text), nil
	}

	// Try current charset first.
	if e.encoder != nil {
		if encoded, err := e.encoder.String(text); err == nil {
			return []byte(encoded), nil
		}
	}

	// Auto-detect best charset for this text.
	bestCharset := e.detectBestCharset(text)
	if bestCharset == e.currentCharset {
		// Already using the best candidate but it failed above.
		return SanitizeForDisplay(text), nil
	}

	// Try encoding with the candidate charset without mutating encoder state.
	testEnc := encoderForCharset(bestCharset)
	encoded, err := testEnc.String(text)
	if err != nil {
		return SanitizeForDisplay(text), nil
	}

	// Encoding succeeded — switch hardware and encoder state atomically.
	if err := display.SetCharacterCodeTableInternal(bestCharset); err != nil {
		return nil, fmt.Errorf("charset switch failed: %w", err)
	}

	return []byte(encoded), nil
}

// SanitizeForDisplay replaces any non-ASCII rune with '?' so that only
// characters guaranteed to be representable in any single-byte codepage
// are sent to the hardware.
func SanitizeForDisplay(text string) []byte {
	result := make([]byte, 0, len(text))
	for _, r := range text {
		if r <= 127 {
			result = append(result, byte(r))
		} else {
			result = append(result, '?')
		}
	}
	return result
}

// detectBestCharset analyzes Latin text and returns the best charset to use.
func (e *CharsetEncoder) detectBestCharset(text string) int {
	hasPortuguese := false
	hasEuro := false
	hasLatin := false

	for _, r := range text {
		if r > 127 {
			switch {
			case r == '€':
				hasEuro = true
			case r >= 'À' && r <= 'ÿ':
				hasLatin = true
				if r == 'ã' || r == 'õ' || r == 'ç' || r == 'Ã' || r == 'Õ' || r == 'Ç' {
					hasPortuguese = true
				}
			}
		}
	}

	if hasPortuguese {
		return chartablePC860
	}
	if hasEuro {
		return chartablePC858
	}
	if hasLatin {
		return chartablePC850
	}

	return chartablePC437
}
