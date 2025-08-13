package vfd

import "errors"

// SetCursor moves the cursor to a 1-based position (column, row).
// This uses the command sequence US $ n m.
func (d *Display) SetCursor(column, row int) error {
	if d.cursorColumn == column && d.cursorRow == row {
		return nil
	}
	if column < 1 || row < 1 {
		return errors.New("column/row must be >= 1")
	}
	if d.columns > 0 && column > d.columns {
		return errors.New("column exceeds configured width")
	}
	if d.rows > 0 && row > d.rows {
		return errors.New("row exceeds configured height")
	}
	if column > 255 || row > 255 {
		return errors.New("column/row out of device range")
	}
	d.cursorColumn = column
	d.cursorRow = row
	if d.commandProfile == nil || d.commandProfile.Protocol == nil {
		return errors.New("no command profile set")
	}
	cmd := d.commandProfile.Protocol.MoveCursor(column, row)
	if cmd == nil {
		return errors.New("invalid cursor position for this protocol")
	}
	return d.writeBytes(cmd)
}

// GetCursor returns the current cursor position.
func (d *Display) GetCursor() (int, int) {
	return d.cursorColumn, d.cursorRow
}

// advanceCursorBy updates internal cursor position after writing a number of characters
func (d *Display) advanceCursorBy(chars int) {
	if chars <= 0 {
		return
	}
	if d.columns <= 0 || d.rows <= 0 {
		return
	}
	// Normalize to 1-based starting point if not yet set
	if d.cursorColumn < 1 {
		d.cursorColumn = 1
	}
	if d.cursorRow < 1 {
		d.cursorRow = 1
	}

	zeroBasedCol := d.cursorColumn - 1
	total := zeroBasedCol + chars
	rowsAdded := total / d.columns
	newColZero := total % d.columns

	newRow := d.cursorRow + rowsAdded
	// Wrap rows cyclically within [1..rows]
	if d.rows > 0 {
		newRow = ((newRow - 1) % d.rows) + 1
	}

	d.cursorColumn = newColZero + 1
	d.cursorRow = newRow
}
