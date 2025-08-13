package vfd

import (
	"vfd/vfd/models/epson"

	"go.bug.st/serial"
)

// Model represents a specific VFD display model.
type Model string

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
	CommandProtocol CommandProtocolName // Which command set this display uses

	// Display capabilities
	SupportsBrightness   bool
	BrightnessLevels     int // Number of brightness levels (e.g., 1-4)
	SupportsCursorBlink  bool
	SupportsCharsetTable bool
	SupportsSelfTest     bool

	// Character encoding
	DefaultCharset int // Default character code table page

	// Documentation reference
	DocumentationURL string
}

// Supported VFD display models
const (
	// Epson DM-D110 customer display
	ModelEpsonDMD110 Model = "EPSON_DM_D110"

	// Add more models here as support is added
	// ModelEpsonDMD210 Model = "EPSON_DM_D210"
	// ModelBixolonBCD1100 Model = "BIXOLON_BCD_1100"
)

// modelRegistry contains profiles for all supported VFD models.
var modelRegistry = map[Model]*ModelProfile{
	ModelEpsonDMD110: {
		Name:                 epson.DMD110Spec.Name,
		Manufacturer:         epson.DMD110Spec.Manufacturer,
		Model:                epson.DMD110Spec.Model,
		Columns:              epson.DMD110Spec.Columns,
		Rows:                 epson.DMD110Spec.Rows,
		DefaultBaudRate:      epson.DMD110Spec.DefaultBaudRate,
		DefaultDataBits:      epson.DMD110Spec.DefaultDataBits,
		DefaultParity:        epson.DMD110Spec.DefaultParity,
		DefaultStopBits:      epson.DMD110Spec.DefaultStopBits,
		CommandProtocol:      CommandProtocolName(epson.DMD110Spec.CommandProtocol),
		SupportsBrightness:   epson.DMD110Spec.SupportsBrightness,
		BrightnessLevels:     epson.DMD110Spec.BrightnessLevels,
		SupportsCursorBlink:  epson.DMD110Spec.SupportsCursorBlink,
		SupportsCharsetTable: epson.DMD110Spec.SupportsCharsetTable,
		SupportsSelfTest:     epson.DMD110Spec.SupportsSelfTest,
		DefaultCharset:       epson.DMD110Spec.DefaultCharset,
		DocumentationURL:     epson.DMD110Spec.DocumentationURL,
	},
}

// GetModelProfile returns the profile for the specified model.
func GetModelProfile(model Model) (*ModelProfile, bool) {
	profile, exists := modelRegistry[model]
	return profile, exists
}

// GetSupportedModels returns a list of all supported VFD models.
func GetSupportedModels() []Model {
	models := make([]Model, 0, len(modelRegistry))
	for model := range modelRegistry {
		models = append(models, model)
	}
	return models
}

// IsModelSupported checks if a model is supported by this library.
func IsModelSupported(model Model) bool {
	_, exists := modelRegistry[model]
	return exists
}

// GetModelSpecs returns the basic specifications for a model.
func GetModelSpecs(model Model) (columns, rows int, found bool) {
	if profile, exists := modelRegistry[model]; exists {
		return profile.Columns, profile.Rows, true
	}
	return 0, 0, false
}

// GetModelDefaults returns default serial communication settings for a model.
func GetModelDefaults(model Model) (*Options, bool) {
	profile, exists := modelRegistry[model]
	if !exists {
		return nil, false
	}

	return &Options{
		BaudRate: profile.DefaultBaudRate,
		DataBits: profile.DefaultDataBits,
		Parity:   profile.DefaultParity,
		StopBits: profile.DefaultStopBits,
		Columns:  profile.Columns,
		Rows:     profile.Rows,
	}, true
}
