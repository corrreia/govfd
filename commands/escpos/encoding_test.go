package escpos

import (
	"errors"
	"testing"
)

// mockDisplay implements CharsetSwitcher for testing.
// It mirrors what Display.SetCharacterCodeTableInternal does:
// on success, it updates both hardware page and encoder state.
type mockDisplay struct {
	encoder     *CharsetEncoder // linked encoder to update on success
	currentPage int
	failOnPage  int // return error when switching to this page (-1 = never fail)
	switchCount int
}

func newMockDisplay() *mockDisplay {
	return &mockDisplay{failOnPage: -1}
}

func (m *mockDisplay) SetCharacterCodeTableInternal(page int) error {
	m.switchCount++
	if page == m.failOnPage {
		return errors.New("hardware switch failed")
	}
	m.currentPage = page
	if m.encoder != nil {
		m.encoder.SetCharset(page)
	}
	return nil
}

func TestNewCharsetEncoder(t *testing.T) {
	enc := NewCharsetEncoder()
	if enc == nil {
		t.Fatal("NewCharsetEncoder returned nil")
	}
	if enc.currentCharset != chartablePC437 {
		t.Errorf("default charset = %d, want %d (PC437)", enc.currentCharset, chartablePC437)
	}
	if enc.encoder == nil {
		t.Fatal("encoder is nil after init")
	}
}

func TestEncodeASCII(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()

	result, err := enc.EncodeTextWithAutoCharsetSwitching("Hello World", display)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "Hello World" {
		t.Errorf("got %q, want %q", result, "Hello World")
	}
	if display.switchCount != 0 {
		t.Errorf("hardware switch count = %d, want 0 for ASCII text", display.switchCount)
	}
}

func TestEncodeLatinNoSwitchNeeded(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()

	// café, Müller, etc. encode fine in PC437 — no switch needed
	result, err := enc.EncodeTextWithAutoCharsetSwitching("café", display)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 4 {
		t.Errorf("len = %d, want 4", len(result))
	}
	if display.switchCount != 0 {
		t.Errorf("hardware switch count = %d, want 0 (café works in PC437)", display.switchCount)
	}
}

func TestEncodePortuguese(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()

	// ação contains ã which is NOT in PC437 — requires switch to PC860
	_, err := enc.EncodeTextWithAutoCharsetSwitching("ação", display)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if display.currentPage != chartablePC860 {
		t.Errorf("hardware page = %d, want %d (PC860 Portuguese)", display.currentPage, chartablePC860)
	}
	if display.switchCount != 1 {
		t.Errorf("hardware switch count = %d, want 1", display.switchCount)
	}
}

func TestEncodeEuro(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()

	// € is NOT in PC437 — requires switch to PC858
	_, err := enc.EncodeTextWithAutoCharsetSwitching("€19.99", display)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if display.currentPage != chartablePC858 {
		t.Errorf("hardware page = %d, want %d (PC858 Euro)", display.currentPage, chartablePC858)
	}
}

func TestNoHardwareSwitchWhenCurrentCharsetWorks(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()
	display.encoder = enc // link so SetCharacterCodeTableInternal updates encoder

	// First call switches to PC860 for Portuguese
	enc.EncodeTextWithAutoCharsetSwitching("ação", display)
	count := display.switchCount

	// Second call with same charset should not switch again
	enc.EncodeTextWithAutoCharsetSwitching("coração", display)
	if display.switchCount != count {
		t.Errorf("hardware switched again when current charset already worked")
	}
}

func TestHardwareSwitchFailureReturnsError(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()
	display.failOnPage = chartablePC860 // fail when trying to switch to PC860

	// ação requires PC860 which will fail
	_, err := enc.EncodeTextWithAutoCharsetSwitching("ação", display)
	if err == nil {
		t.Fatal("expected error when hardware switch fails, got nil")
	}
}

func TestEncoderStateNotCorruptedOnSwitchFailure(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()
	display.failOnPage = chartablePC860

	// This should fail — but encoder should remain on PC437
	enc.EncodeTextWithAutoCharsetSwitching("ação", display)

	if enc.currentCharset != chartablePC437 {
		t.Errorf("encoder charset = %d after failed switch, want %d (PC437 unchanged)",
			enc.currentCharset, chartablePC437)
	}

	// ASCII should still work fine
	display.failOnPage = -1
	result, err := enc.EncodeTextWithAutoCharsetSwitching("Hello", display)
	if err != nil {
		t.Fatalf("ASCII encoding failed after switch failure: %v", err)
	}
	if string(result) != "Hello" {
		t.Errorf("got %q, want %q", result, "Hello")
	}
}

func TestSanitizeForDisplay(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello", "Hello"},
		{"café", "caf?"},
		{"日本語", "???"},
		{"abc日def", "abc?def"},
		{"€19.99", "?19.99"},
		{"", ""},
	}

	for _, tt := range tests {
		got := string(SanitizeForDisplay(tt.input))
		if got != tt.want {
			t.Errorf("SanitizeForDisplay(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSanitizeFallbackForUnrepresentableChars(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()

	// CJK characters can't be encoded in any Latin codepage.
	// Should sanitize to '?' instead of sending raw UTF-8.
	result, err := enc.EncodeTextWithAutoCharsetSwitching("日本語", display)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "???" {
		t.Errorf("got %q, want %q (sanitized)", result, "???")
	}
	if display.switchCount != 0 {
		t.Errorf("hardware switch count = %d, want 0 for unrepresentable chars", display.switchCount)
	}
}

func TestDetectBestCharset(t *testing.T) {
	enc := NewCharsetEncoder()

	tests := []struct {
		text string
		want int
	}{
		{"Hello", chartablePC437},
		{"café", chartablePC850},  // detected as Latin (even though PC437 handles it)
		{"ação", chartablePC860},  // Portuguese-specific
		{"€100", chartablePC858},  // Euro
		{"café ação", chartablePC860}, // Portuguese takes priority
	}

	for _, tt := range tests {
		got := enc.detectBestCharset(tt.text)
		if got != tt.want {
			t.Errorf("detectBestCharset(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}

func TestEncodedOutputIsSingleBytePerChar(t *testing.T) {
	enc := NewCharsetEncoder()
	display := newMockDisplay()

	tests := []struct {
		text    string
		wantLen int
	}{
		{"Hello", 5},
		{"café", 4},   // 4 display chars (é is single byte in CP437)
		{"ação", 4},   // 4 display chars after switch to PC860
		{"€19.99", 6}, // 6 display chars after switch to PC858
		{"日本語", 3},   // 3 sanitized '?' chars
	}

	for _, tt := range tests {
		enc.SetCharset(chartablePC437)
		display.failOnPage = -1

		result, err := enc.EncodeTextWithAutoCharsetSwitching(tt.text, display)
		if err != nil {
			t.Fatalf("EncodeText(%q): unexpected error: %v", tt.text, err)
		}
		if len(result) != tt.wantLen {
			t.Errorf("EncodeText(%q): len = %d, want %d (single-byte-per-char)",
				tt.text, len(result), tt.wantLen)
		}
	}
}
