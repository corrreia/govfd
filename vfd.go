// Package govfd provides a smart VFD (Vacuum Fluorescent Display) library for Go
// with automatic Latin character encoding. Just send UTF-8 text and it works perfectly!
//
// Features:
//   - Smart Latin character encoding (Portuguese, Spanish, French, German, Italian)
//   - Automatic charset detection and hardware switching
//   - Model-based architecture with automatic configuration
//   - ESC/POS protocol support
//   - Easy-to-use API with zero configuration needed
//
// Example:
//
//	display, err := govfd.OpenModel("COM3", types.ModelEpsonDMD110)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	display.WriteText("Café")  // Works perfectly!
package govfd

import (
	"errors"

	"github.com/corrreia/govfd/commands/escpos"
	"github.com/corrreia/govfd/types"

	"go.bug.st/serial"
)

// Display represents an open connection to a VFD display over a serial port.
type Display struct {
	port         serial.Port
	portName     string
	columns      int
	rows         int
	cursorColumn int
	cursorRow    int
	brightness   int
	blinkMs      int
	protocol     Protocol               // Command protocol for this display
	encoder      *escpos.CharsetEncoder // Character encoding handler
}

// DefaultOptions returns commonly used defaults (9600 8N1).
func DefaultOptions() *Options {
	return &Options{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
}

// OpenModel establishes a connection to a VFD using model-specific defaults.
// This is the recommended way to open a VFD display as it automatically
// configures the correct serial settings and dimensions for the specified model.
func OpenModel(portName string, model types.Model) (*Display, error) {
	if portName == "" {
		return nil, errors.New("portName is required")
	}

	modelProfile, exists := GetModelProfile(model)
	if !exists {
		return nil, errors.New("unsupported VFD model: " + string(model))
	}

	opts, _ := GetModelDefaults(model)
	display, err := Open(portName, opts)
	if err != nil {
		return nil, err
	}

	// Set the command protocol for this model
	protocol, exists := GetProtocol(modelProfile.CommandProtocol)
	if !exists {
		return nil, errors.New("unsupported command protocol: " + modelProfile.CommandProtocol)
	}
	display.protocol = protocol

	// Initialize character encoding
	display.encoder = escpos.NewCharsetEncoder()

	return display, nil
}

// OpenModelWithOptions establishes a connection to a VFD using model defaults
// but allows overriding specific options. Model defaults are used for any
// options that are not explicitly set (zero values) in the provided opts.
func OpenModelWithOptions(portName string, model types.Model, opts *Options) (*Display, error) {
	if portName == "" {
		return nil, errors.New("portName is required")
	}

	modelProfile, exists := GetModelProfile(model)
	if !exists {
		return nil, errors.New("unsupported VFD model: " + string(model))
	}

	defaults, _ := GetModelDefaults(model)

	// Merge user options with model defaults
	if opts == nil {
		opts = defaults
	} else {
		// Use model defaults for unspecified options
		if opts.BaudRate == 0 {
			opts.BaudRate = defaults.BaudRate
		}
		if opts.DataBits == 0 {
			opts.DataBits = defaults.DataBits
		}
		if opts.Parity == 0 {
			opts.Parity = defaults.Parity
		}
		if opts.StopBits == 0 {
			opts.StopBits = defaults.StopBits
		}
		if opts.Columns == 0 {
			opts.Columns = defaults.Columns
		}
		if opts.Rows == 0 {
			opts.Rows = defaults.Rows
		}
	}

	display, err := Open(portName, opts)
	if err != nil {
		return nil, err
	}

	// Set the command protocol for this model
	protocol, exists := GetProtocol(modelProfile.CommandProtocol)
	if !exists {
		return nil, errors.New("unsupported command protocol: " + modelProfile.CommandProtocol)
	}
	display.protocol = protocol

	// Initialize character encoding
	display.encoder = escpos.NewCharsetEncoder()

	return display, nil
}

// Open establishes a connection to the VFD at the given serial port.
// If opts is nil, DefaultOptions() are used.
//
// Note: For better compatibility, consider using OpenModel() or OpenModelWithOptions()
// which automatically configure settings for specific VFD models.
func Open(portName string, opts *Options) (*Display, error) {
	if portName == "" {
		return nil, errors.New("portName is required")
	}
	if opts == nil {
		opts = DefaultOptions()
	}

	mode := &serial.Mode{
		BaudRate: opts.BaudRate,
		DataBits: opts.DataBits,
		Parity:   opts.Parity,
		StopBits: opts.StopBits,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, errors.New("open serial port " + portName + ": " + err.Error())
	}

	d := &Display{port: port, portName: portName}
	if opts.Columns > 0 {
		d.columns = opts.Columns
	}
	if opts.Rows > 0 {
		d.rows = opts.Rows
	}

	// Set default command protocol (ESC/POS) for manual connections
	if protocol, exists := GetProtocol(types.ProtocolESCPOS); exists {
		d.protocol = protocol
	}

	// Initialize character encoding
	d.encoder = escpos.NewCharsetEncoder()

	return d, nil
}

// Close closes the underlying serial port.
func (d *Display) Close() error {
	if d == nil || d.port == nil {
		return nil
	}
	return d.port.Close()
}

// writeBytes centralizes writes to the serial port with simple nil checks.
func (d *Display) writeBytes(payload []byte) error {
	if d == nil || d.port == nil {
		return errors.New("display is not open")
	}
	_, err := d.port.Write(payload)
	return err
}
