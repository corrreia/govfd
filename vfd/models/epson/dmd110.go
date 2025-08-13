package epson

import (
	"vfd/vfd/commands/escpos"

	"go.bug.st/serial"
)

// DMD110Spec contains the specification for the Epson DM-D110.
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
	DefaultCharset       int
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
	CommandProtocol:      "ESC_POS", // This will be converted to CommandProtocolName
	SupportsBrightness:   true,
	BrightnessLevels:     4,
	SupportsCursorBlink:  true,
	SupportsCharsetTable: true,
	SupportsSelfTest:     true,
	DefaultCharset:       escpos.CharsetPC437,
	DocumentationURL:     "https://download4.epson.biz/sec_pubs/pos/reference_en/escpos_dm/commands.html",
}
