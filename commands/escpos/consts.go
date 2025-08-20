package escpos

// ESC/POS Command byte constants for VFD display operations.
// These represent the hexadecimal command sequences used to control ESC/POS compatible VFD displays.
// Reference: https://download4.epson.biz/sec_pubs/pos/reference_en/escpos_dm/commands.html

// ASCII Control Characters
const (
	// Form Feed - clears the screen
	CmdFormFeed = 0x0C
)

// Escape Sequence Prefixes
const (
	// Escape character - used as prefix for escape sequences
	CmdEscape = 0x1B // ESC

	// Unit Separator - used as prefix for device-specific commands
	CmdUnitSeparator = 0x1F // US
)

// Escape Sequence Commands (ESC + command)
const (
	// ESC @ - Initialize/clear display state
	CmdEscInitialize = 0x40 // @ - used with ESC (0x1B)

	// ESC t - Set character code table page
	CmdEscCharsetTable = 0x74 // t - used with ESC (0x1B)
)

// Unit Separator Commands (US + command)
const (
	// US $ - Set cursor position (followed by column, row bytes)
	CmdUSSetCursor = 0x24 // $ - used with US (0x1F)

	// US @ - Execute display self-test
	CmdUSSelfTest = 0x40 // @ - used with US (0x1F)

	// US E - Set cursor blink period (followed by step value)
	CmdUSSetBlink = 0x45 // E - used with US (0x1F)

	// US X - Set brightness level (followed by level 1-4)
	CmdUSSetBrightness = 0x58 // X - used with US (0x1F)
)

// Complete Command Sequences as byte arrays for convenience
var (
	// Clear/Initialize display: ESC @
	SeqClear = []byte{CmdEscape, CmdEscInitialize}

	// Form feed to clear screen
	SeqFormFeed = []byte{CmdFormFeed}

	// Self-test: US @
	SeqSelfTest = []byte{CmdUnitSeparator, CmdUSSelfTest}
)

// Command sequence builders (return byte arrays for specific operations)

// BuildSetCursorSeq creates the command sequence to set cursor position.
// Returns: US $ column row
func BuildSetCursorSeq(column, row byte) []byte {
	return []byte{CmdUnitSeparator, CmdUSSetCursor, column, row}
}

// BuildSetBrightnessSeq creates the command sequence to set brightness.
// Returns: US X level
func BuildSetBrightnessSeq(level byte) []byte {
	return []byte{CmdUnitSeparator, CmdUSSetBrightness, level}
}

// BuildSetBlinkSeq creates the command sequence to set cursor blink.
// Returns: US E steps
func BuildSetBlinkSeq(steps byte) []byte {
	return []byte{CmdUnitSeparator, CmdUSSetBlink, steps}
}

// BuildSetCharsetSeq creates the command sequence to set character code table.
// Returns: ESC t page
func BuildSetCharsetSeq(page byte) []byte {
	return []byte{CmdEscape, CmdEscCharsetTable, page}
}
