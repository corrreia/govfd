package govfd

import "errors"

// Clear sends ESC @ to initialize/clear the display state.
func (d *Display) Clear() error {
	if d.protocol == nil {
		return errors.New("no command protocol set")
	}
	return d.writeBytes(d.protocol.Clear())
}

// FormFeed sends a form feed (0x0C) to clear the screen.
func (d *Display) FormFeed() error {
	if d.protocol == nil {
		return errors.New("no command protocol set")
	}
	return d.writeBytes(d.protocol.FormFeed())
}

// WriteText writes a string to the display at the current cursor position.
// 🎯 FULLY AUTOMATIC: The system handles ALL character encoding automatically!
// Just send any UTF-8 text and it works perfectly - zero configuration needed!
func (d *Display) WriteText(message string) error {
	if d.protocol == nil {
		return errors.New("no command protocol set")
	}

	// 🚀 SMART AUTO-ENCODING: Automatically handles charset selection and encoding
	encodedBytes, err := d.smartEncodeText(message)
	if err != nil {
		return err
	}

	if err := d.writeBytes(encodedBytes); err != nil {
		return err
	}
	d.advanceCursorBy(len(encodedBytes))
	return nil
}

// smartEncodeText is the magic function that handles everything automatically
func (d *Display) smartEncodeText(text string) ([]byte, error) {
	if d.encoder == nil {
		// Fallback if encoder not available
		return []byte(text), nil
	}

	// Let the encoder handle everything - it will automatically:
	// 1. Try current charset
	// 2. Auto-detect best charset for the text
	// 3. Switch hardware charset if needed
	// 4. Fall back to ASCII transliteration if needed
	return d.encoder.EncodeTextWithAutoCharsetSwitching(text, d)
}

// WriteRawBytes writes raw bytes directly to the display at the current cursor position.
// This allows sending specific byte sequences without protocol interpretation.
func (d *Display) WriteRawBytes(data []byte) error {
	if err := d.writeBytes(data); err != nil {
		return err
	}
	d.advanceCursorBy(len(data))
	return nil
}

// SetBrightness sets the display brightness (expected 1..4 on many VFDs).
// This uses the command sequence US X n.
func (d *Display) SetBrightness(level int) error {
	if d.protocol == nil {
		return errors.New("no command protocol set")
	}
	if level < 1 || level > 4 {
		return errors.New("brightness level must be between 1 and 4")
	}
	cmd := d.protocol.SetBrightness(level)
	if cmd == nil {
		return errors.New("invalid brightness level for this protocol")
	}
	if err := d.writeBytes(cmd); err != nil {
		return err
	}
	d.brightness = level
	return nil
}

// GetBrightness returns the current brightness level.
func (d *Display) GetBrightness() int {
	return d.brightness
}

// SetBlink sets the cursor blink period in milliseconds (0 to disable).
// Device expects 50ms steps (n = ms / 50), sent as US E n (0x1F 0x45 n).
func (d *Display) SetBlink(ms int) error {
	if d == nil || d.port == nil {
		return errors.New("display is not open")
	}
	if d.protocol == nil {
		return errors.New("no command protocol set")
	}
	if ms < 0 {
		return errors.New("blink ms must be >= 0")
	}
	cmd := d.protocol.SetBlink(ms)
	if cmd == nil {
		return errors.New("invalid blink interval for this protocol")
	}
	if err := d.writeBytes(cmd); err != nil {
		return err
	}
	// Protocol handles the actual conversion, so we store the requested value
	d.blinkMs = ms
	return nil
}

// GetBlinkMs returns the last set blink period in milliseconds (0 if unknown).
func (d *Display) GetBlinkMs() int {
	return d.blinkMs
}

// Dimensions returns the current logical dimensions. Zero means unspecified.
func (d *Display) Dimensions() (int, int) {
	return d.columns, d.rows
}

// SelfTest executes the display's built-in self-test.
// This uses the command sequence US @ (0x1F 0x40).
func (d *Display) SelfTest() error {
	if d.protocol == nil {
		return errors.New("no command protocol set")
	}
	return d.writeBytes(d.protocol.SelfTest())
}

// SetCharacterCodeTableInternal selects the character code table page (INTERNAL USE ONLY).
// This is now handled automatically by the smart encoding system.
// This method implements the CharsetSwitcher interface.
func (d *Display) SetCharacterCodeTableInternal(page int) error {
	if d.protocol == nil {
		return errors.New("no command protocol set")
	}
	if page < 0 || page > 255 {
		return errors.New("page must be between 0 and 255")
	}
	cmd := d.protocol.SetCharset(page)
	if cmd == nil {
		return errors.New("invalid charset page for this protocol")
	}

	// Update the character encoder's charset
	if d.encoder != nil {
		d.encoder.SetCharset(page)
	}

	return d.writeBytes(cmd)
}
