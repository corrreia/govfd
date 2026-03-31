package epson

import (
	"github.com/corrreia/govfd/types"

	"go.bug.st/serial"
)

// DMD110Profile contains the specification for the Epson DM-D110.
var DMD110Profile = types.ModelProfile{
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
	SupportsCharsetTable: true,
	SupportsSelfTest:     true,
	DocumentationURL:     "https://download4.epson.biz/sec_pubs/pos/reference_en/escpos_dm/commands.html",
}
