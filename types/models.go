package types

import "go.bug.st/serial"

// Model represents a specific VFD display model.
type Model string

// Supported VFD display models
const (
	// Epson DM-D110 customer display
	ModelEpsonDMD110 Model = "EPSON_DM_D110"
)

// ModelProfile contains all specifications and default settings for a VFD model.
type ModelProfile struct {
	// Display identification
	Name         string
	Manufacturer string
	Model        string

	// Physical specifications
	Columns int // Number of columns (characters per line)
	Rows    int // Number of rows (lines)

	// Serial communication defaults
	DefaultBaudRate int
	DefaultDataBits int
	DefaultParity   serial.Parity
	DefaultStopBits serial.StopBits

	// Command protocol
	CommandProtocol string

	// Display capabilities
	SupportsBrightness   bool
	BrightnessLevels     int
	SupportsCursorBlink  bool
	SupportsCharsetTable bool
	SupportsSelfTest     bool

	// Documentation reference
	DocumentationURL string
}
