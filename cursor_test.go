package govfd

import (
	"errors"
	"testing"
	"time"

	"go.bug.st/serial"
)

// mockPort implements serial.Port for testing without real hardware.
type mockPort struct {
	written    []byte
	failOnNext bool
}

func (m *mockPort) Read(p []byte) (int, error) { return 0, nil }
func (m *mockPort) Write(p []byte) (int, error) {
	if m.failOnNext {
		m.failOnNext = false
		return 0, errors.New("write failed")
	}
	m.written = append(m.written, p...)
	return len(p), nil
}
func (m *mockPort) ResetInputBuffer() error                          { return nil }
func (m *mockPort) ResetOutputBuffer() error                         { return nil }
func (m *mockPort) SetDTR(dtr bool) error                            { return nil }
func (m *mockPort) SetRTS(rts bool) error                            { return nil }
func (m *mockPort) SetReadTimeout(t time.Duration) error             { return nil }
func (m *mockPort) GetModemStatusBits() (*serial.ModemStatusBits, error) { return nil, nil }
func (m *mockPort) Close() error                                     { return nil }
func (m *mockPort) Break(t time.Duration) error                      { return nil }
func (m *mockPort) Drain() error                                     { return nil }
func (m *mockPort) SetMode(mode *serial.Mode) error                  { return nil }

// newTestDisplay creates a Display with a mock serial port and protocol for testing.
func newTestDisplay(cols, rows int) (*Display, *mockPort) {
	port := &mockPort{}
	protocol, _ := GetProtocol("ESC/POS")
	d := &Display{
		port:     port,
		portName: "test",
		columns:  cols,
		rows:     rows,
		protocol: protocol,
	}
	return d, port
}

func TestSetCursorBasic(t *testing.T) {
	d, _ := newTestDisplay(20, 2)

	if err := d.SetCursor(5, 1); err != nil {
		t.Fatalf("SetCursor(5,1) error: %v", err)
	}
	col, row := d.GetCursor()
	if col != 5 || row != 1 {
		t.Errorf("cursor = (%d,%d), want (5,1)", col, row)
	}
}

func TestSetCursorNilProtocol(t *testing.T) {
	d := &Display{columns: 20, rows: 2}

	err := d.SetCursor(1, 1)
	if err == nil {
		t.Fatal("expected error for nil protocol, got nil")
	}
}

func TestSetCursorInvalidPosition(t *testing.T) {
	d, _ := newTestDisplay(20, 2)

	tests := []struct {
		col, row int
		desc     string
	}{
		{0, 1, "column zero"},
		{1, 0, "row zero"},
		{-1, 1, "negative column"},
		{1, -1, "negative row"},
		{21, 1, "column exceeds width"},
		{1, 3, "row exceeds height"},
	}

	for _, tt := range tests {
		if err := d.SetCursor(tt.col, tt.row); err == nil {
			t.Errorf("SetCursor(%d,%d) [%s]: expected error, got nil", tt.col, tt.row, tt.desc)
		}
	}
}

func TestSetCursorNoOpWhenAlreadyThere(t *testing.T) {
	d, port := newTestDisplay(20, 2)

	d.SetCursor(3, 1)
	port.written = nil // clear

	// Same position — should be a no-op, no bytes written
	if err := d.SetCursor(3, 1); err != nil {
		t.Fatalf("SetCursor same position error: %v", err)
	}
	if len(port.written) != 0 {
		t.Errorf("wrote %d bytes for no-op cursor move, want 0", len(port.written))
	}
}

func TestSetCursorStateNotUpdatedOnWriteFailure(t *testing.T) {
	d, port := newTestDisplay(20, 2)

	// Move to (1,1) first
	d.SetCursor(1, 1)

	// Now fail the next write
	port.failOnNext = true
	err := d.SetCursor(5, 2)
	if err == nil {
		t.Fatal("expected write error, got nil")
	}

	// Cursor should still be at (1,1)
	col, row := d.GetCursor()
	if col != 1 || row != 1 {
		t.Errorf("cursor = (%d,%d) after failed write, want (1,1)", col, row)
	}

	// Retry should now work (not no-op'd)
	port.failOnNext = false
	if err := d.SetCursor(5, 2); err != nil {
		t.Fatalf("retry SetCursor error: %v", err)
	}
	col, row = d.GetCursor()
	if col != 5 || row != 2 {
		t.Errorf("cursor = (%d,%d) after retry, want (5,2)", col, row)
	}
}

func TestAdvanceCursorBy(t *testing.T) {
	tests := []struct {
		startCol, startRow int
		advance            int
		wantCol, wantRow   int
		desc               string
	}{
		{1, 1, 5, 6, 1, "simple advance"},
		{1, 1, 20, 1, 2, "wrap to next row"},
		{1, 1, 40, 1, 1, "wrap around both rows"},
		{10, 1, 15, 5, 2, "mid-line wrap"},
		{1, 1, 0, 1, 1, "zero advance"},
		{20, 2, 1, 1, 1, "wrap from end of display"},
	}

	for _, tt := range tests {
		d, _ := newTestDisplay(20, 2)
		d.cursorColumn = tt.startCol
		d.cursorRow = tt.startRow

		d.advanceCursorBy(tt.advance)

		col, row := d.GetCursor()
		if col != tt.wantCol || row != tt.wantRow {
			t.Errorf("[%s] advance(%d) from (%d,%d): got (%d,%d), want (%d,%d)",
				tt.desc, tt.advance, tt.startCol, tt.startRow, col, row, tt.wantCol, tt.wantRow)
		}
	}
}

func TestAdvanceCursorNoDimensionsIsNoOp(t *testing.T) {
	d := &Display{cursorColumn: 1, cursorRow: 1}

	d.advanceCursorBy(10)

	col, row := d.GetCursor()
	if col != 1 || row != 1 {
		t.Errorf("cursor = (%d,%d) without dimensions set, want (1,1) unchanged", col, row)
	}
}
