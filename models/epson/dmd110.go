package epson

import (
	"github.com/corrreia/govfd/types"

	"go.bug.st/serial"
)

// DMD110Spec contains the specification for the Epson DM-D110.
// Note: Charset selection is now fully automatic - no need to specify!
var DMD110Spec = struct {
	Name                 string
	Manufacturer         string
	Model                string
	Columns              int
	Rows                 int
	DefaultBaudRate      int
	DefaultDataBits      int
	DefaultParity        serial.Parity
	DefaultStopBits      serial.StopBits
	CommandProtocol      string
	SupportsBrightness   bool
	BrightnessLevels     int
	SupportsCursorBlink  bool
	SupportsCharsetTable bool
	SupportsSelfTest     bool
	DocumentationURL     string
}{
	Name:                 "Epson DM-D110",
	Manufacturer:         "Epson",
	Model:                "DM-D110",
	Columns:              20,
	Rows:                 2,
	DefaultBaudRate:      9600,
	DefaultDataBits:      8,
	DefaultParity:        serial.NoParity,
	DefaultStopBits:      serial.OneStopBit,
	CommandProtocol:      types.ProtocolESCPOS,
	SupportsBrightness:   true,
	BrightnessLevels:     4,
	SupportsCursorBlink:  true,
	SupportsCharsetTable: true, // But handled automatically!
	SupportsSelfTest:     true,
	DocumentationURL:     "https://download4.epson.biz/sec_pubs/pos/reference_en/escpos_dm/commands.html",
}
