package escpos

// ESCPOSProtocol implements the CommandProtocol interface for ESC/POS displays.
type ESCPOSProtocol struct{}

// GetName returns the protocol name.
func (p *ESCPOSProtocol) GetName() string {
	return "ESC/POS"
}

// GetDescription returns the protocol description.
func (p *ESCPOSProtocol) GetDescription() string {
	return "Standard ESC/POS command set for VFD displays"
}

// Clear returns the command sequence to initialize/clear the display.
func (p *ESCPOSProtocol) Clear() []byte {
	return SeqClear
}

// FormFeed returns the command sequence to clear screen content.
func (p *ESCPOSProtocol) FormFeed() []byte {
	return SeqFormFeed
}

// MoveCursor returns the command sequence to move cursor to position (1-based).
func (p *ESCPOSProtocol) MoveCursor(column, row int) []byte {
	if column < 1 || column > 255 || row < 1 || row > 255 {
		return nil // Invalid position
	}
	return BuildSetCursorSeq(byte(column), byte(row))
}

// WriteText returns the text as bytes (ESC/POS displays accept raw text).
func (p *ESCPOSProtocol) WriteText(text string) []byte {
	return []byte(text)
}

// SetBrightness returns the command sequence to set brightness level.
func (p *ESCPOSProtocol) SetBrightness(level int) []byte {
	if level < 1 || level > 4 {
		return nil // Invalid brightness level
	}
	return BuildSetBrightnessSeq(byte(level))
}

// SetBlink returns the command sequence to set cursor blink period.
func (p *ESCPOSProtocol) SetBlink(intervalMs int) []byte {
	if intervalMs < 0 {
		return nil // Invalid interval
	}
	steps := byte(intervalMs / 50) // ESC/POS uses 50ms steps
	return BuildSetBlinkSeq(steps)
}

// SetCharset returns the command sequence to set character encoding table.
func (p *ESCPOSProtocol) SetCharset(page int) []byte {
	if page < 0 || page > 255 {
		return nil // Invalid charset page
	}
	return BuildSetCharsetSeq(byte(page))
}

// SelfTest returns the command sequence to execute self-test.
func (p *ESCPOSProtocol) SelfTest() []byte {
	return SeqSelfTest
}
