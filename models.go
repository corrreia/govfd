package govfd

import (
	"github.com/corrreia/govfd/models/epson"
	"github.com/corrreia/govfd/types"

	"go.bug.st/serial"
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
	CommandProtocol string // Which command set this display uses

	// Display capabilities
	SupportsBrightness   bool
	BrightnessLevels     int // Number of brightness levels (e.g., 1-4)
	SupportsCursorBlink  bool
	SupportsCharsetTable bool // Charset is now handled automatically!
	SupportsSelfTest     bool

	// Documentation reference
	DocumentationURL string
}

// Options configures serial connection parameters for the VFD display.
type Options struct {
	BaudRate int
	DataBits int
	Parity   serial.Parity
	StopBits serial.StopBits
	// Optional logical screen size for bounds validation.
	// If zero, bounds are not enforced beyond device limits (1..255).
	Columns int
	Rows    int
}

// modelRegistry contains profiles for all supported VFD models.
var modelRegistry = map[types.Model]*ModelProfile{
	types.ModelEpsonDMD110: {
		Name:                 epson.DMD110Spec.Name,
		Manufacturer:         epson.DMD110Spec.Manufacturer,
		Model:                epson.DMD110Spec.Model,
		Columns:              epson.DMD110Spec.Columns,
		Rows:                 epson.DMD110Spec.Rows,
		DefaultBaudRate:      epson.DMD110Spec.DefaultBaudRate,
		DefaultDataBits:      epson.DMD110Spec.DefaultDataBits,
		DefaultParity:        epson.DMD110Spec.DefaultParity,
		DefaultStopBits:      epson.DMD110Spec.DefaultStopBits,
		CommandProtocol:      epson.DMD110Spec.CommandProtocol,
		SupportsBrightness:   epson.DMD110Spec.SupportsBrightness,
		BrightnessLevels:     epson.DMD110Spec.BrightnessLevels,
		SupportsCursorBlink:  epson.DMD110Spec.SupportsCursorBlink,
		SupportsCharsetTable: epson.DMD110Spec.SupportsCharsetTable,
		SupportsSelfTest:     epson.DMD110Spec.SupportsSelfTest,
		DocumentationURL:     epson.DMD110Spec.DocumentationURL,
	},
}

// GetModelProfile returns the profile for the specified model.
func GetModelProfile(model types.Model) (*ModelProfile, bool) {
	profile, exists := modelRegistry[model]
	return profile, exists
}

// GetSupportedModels returns a list of all supported VFD models.
func GetSupportedModels() []types.Model {
	models := make([]types.Model, 0, len(modelRegistry))
	for model := range modelRegistry {
		models = append(models, model)
	}
	return models
}

// IsModelSupported checks if a model is supported by this library.
func IsModelSupported(model types.Model) bool {
	_, exists := modelRegistry[model]
	return exists
}

// GetModelSpecs returns the basic specifications for a model.
func GetModelSpecs(model types.Model) (columns, rows int, found bool) {
	if profile, exists := modelRegistry[model]; exists {
		return profile.Columns, profile.Rows, true
	}
	return 0, 0, false
}

// GetModelDefaults returns default serial communication settings for a model.
func GetModelDefaults(model types.Model) (*Options, bool) {
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
