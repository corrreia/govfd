package vfd

import "errors"

// Clear sends ESC @ to initialize/clear the display state.
func (d *Display) Clear() error {
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	return d.writeBytes(d.commandProfile.Protocol.Clear())
}

// FormFeed sends a form feed (0x0C) to clear the screen.
func (d *Display) FormFeed() error {
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	return d.writeBytes(d.commandProfile.Protocol.FormFeed())
}

// WriteText writes a string to the display at the current cursor position.
func (d *Display) WriteText(message string) error {
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	if err := d.writeBytes(d.commandProfile.Protocol.WriteText(message)); err != nil {
		return err
	}
	d.advanceCursorBy(len(message))
	return nil
}

// SetBrightness sets the display brightness (expected 1..4 on many VFDs).
// This uses the command sequence US X n.
func (d *Display) SetBrightness(level int) error {
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	if level < 1 || level > 4 {
		return errors.New("brightness level must be between 1 and 4")
	}
	cmd := d.commandProfile.Protocol.SetBrightness(level)
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
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	if ms < 0 {
		return errors.New("blink ms must be >= 0")
	}
	cmd := d.commandProfile.Protocol.SetBlink(ms)
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
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	return d.writeBytes(d.commandProfile.Protocol.SelfTest())
}

// SetCharacterCodeTable selects the character code table page.
// This uses the command sequence ESC t n (0x1B 0x74 n).
// Parameter n specifies the page number (0-255).
// See character code table constants for common pages.
func (d *Display) SetCharacterCodeTable(page int) error {
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	if page < 0 || page > 255 {
		return errors.New("page must be between 0 and 255")
	}
	cmd := d.commandProfile.Protocol.SetCharset(page)
	if cmd == nil {
		return errors.New("invalid charset page for this protocol")
	}
	return d.writeBytes(cmd)
}
