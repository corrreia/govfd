package escpos

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// CharsetEncoder handles character encoding conversion from UTF-8 to legacy charsets
// Features smart auto-detection for Latin characters (Portuguese, Spanish, French, German, Italian)
type CharsetEncoder struct {
	currentCharset int
	encoder        *encoding.Encoder
}

// NewCharsetEncoder creates a new character encoder with default charset (PC437)
func NewCharsetEncoder() *CharsetEncoder {
	e := &CharsetEncoder{
		currentCharset: chartablePC437, // Default to PC437
	}
	e.updateEncoder()
	return e
}

// SetCharset sets the current character encoding table
func (e *CharsetEncoder) SetCharset(charset int) {
	e.currentCharset = charset
	e.updateEncoder()
}

// updateEncoder updates the internal encoder based on current charset (Latin only)
func (e *CharsetEncoder) updateEncoder() {
	switch e.currentCharset {
	case chartablePC850: // CP850 Multilingual Latin
		e.encoder = charmap.CodePage850.NewEncoder()
	case chartablePC860: // CP860 Portuguese
		e.encoder = charmap.CodePage860.NewEncoder()
	case chartablePC858: // PC858 Euro
		e.encoder = charmap.CodePage858.NewEncoder()
	default: // PC437 default
		e.encoder = charmap.CodePage437.NewEncoder()
	}
}

// EncodeText converts UTF-8 text to the current character encoding
// 🎯 SMART MODE: Automatically handles everything for you!
func (e *CharsetEncoder) EncodeText(text string) []byte {
	if !utf8.ValidString(text) {
		// If not valid UTF-8, assume it's already encoded and return as-is
		return []byte(text)
	}

	// 🚀 SMART AUTO-DETECTION: Find the best encoding strategy
	return e.smartEncode(text)
}

// CharsetSwitcher defines the interface for charset switching on the display
type CharsetSwitcher interface {
	SetCharacterCodeTableInternal(page int) error
}

// EncodeTextWithAutoCharsetSwitching is the ultimate automatic encoding function
// It can actually change the display's charset automatically for perfect results!
func (e *CharsetEncoder) EncodeTextWithAutoCharsetSwitching(text string, display CharsetSwitcher) ([]byte, error) {
	if !utf8.ValidString(text) {
		// If not valid UTF-8, assume it's already encoded and return as-is
		return []byte(text), nil
	}

	// Step 1: Try current charset first
	if e.encoder != nil {
		if encoded, err := (*e.encoder).String(text); err == nil {
			// Success! Current charset works perfectly
			return []byte(encoded), nil
		}
	}

	// Step 2: Auto-detect best charset for this text
	bestCharset := e.detectBestCharset(text)

	if bestCharset != e.currentCharset {
		// Try the better charset
		oldCharset := e.currentCharset
		e.SetCharset(bestCharset)

		if e.encoder != nil {
			if encoded, err := (*e.encoder).String(text); err == nil {
				// Great! This charset works - switch the hardware too!
				if err := display.SetCharacterCodeTableInternal(bestCharset); err != nil {
					// If hardware switch fails, restore old charset and continue
					e.SetCharset(oldCharset)
				} else {
					// Success! Hardware and software are now in sync
					return []byte(encoded), nil
				}
			}
		}

		// Restore original charset if switching didn't work
		e.SetCharset(oldCharset)
	}

	// Step 3: Fallback - return text as-is (ASCII-safe)
	return []byte(text), nil
}

// smartEncode is the magic function that handles everything automatically
func (e *CharsetEncoder) smartEncode(text string) []byte {
	// Step 1: Try current charset first
	if e.encoder != nil {
		if encoded, err := (*e.encoder).String(text); err == nil {
			// Success! Current charset works perfectly
			return []byte(encoded)
		}
	}

	// Step 2: Auto-detect best charset for this text
	bestCharset := e.detectBestCharset(text)
	if bestCharset != e.currentCharset {
		// Temporarily switch to better charset
		oldCharset := e.currentCharset
		e.SetCharset(bestCharset)

		if e.encoder != nil {
			if encoded, err := (*e.encoder).String(text); err == nil {
				// Great! Found a charset that works
				return []byte(encoded)
			}
		}

		// Restore original charset
		e.SetCharset(oldCharset)
	}

	// Step 3: Simple fallback - return as ASCII
	return []byte(text)
}

// detectBestCharset analyzes Latin text and returns the best charset to use
func (e *CharsetEncoder) detectBestCharset(text string) int {
	hasPortuguese := false
	hasEuro := false
	hasLatin := false

	for _, r := range text {
		if r > 127 { // Non-ASCII character
			switch {
			case r == '€':
				hasEuro = true
			case r >= 'À' && r <= 'ÿ': // Latin extended only
				hasLatin = true
				// Portuguese-specific characters get priority
				if r == 'ã' || r == 'õ' || r == 'ç' || r == 'Ã' || r == 'Õ' || r == 'Ç' {
					hasPortuguese = true
				}
			}
		}
	}

	// Simple charset selection for Latin scripts only
	if hasPortuguese {
		return chartablePC860 // Portuguese
	}
	if hasEuro {
		return chartablePC858 // Euro support
	}
	if hasLatin {
		return chartablePC850 // Multilingual Latin
	}

	return chartablePC437 // Default
}
