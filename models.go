package govfd

import (
	"github.com/corrreia/govfd/models/epson"
	"github.com/corrreia/govfd/types"

	"go.bug.st/serial"
)

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
var modelRegistry = map[types.Model]*types.ModelProfile{
	types.ModelEpsonDMD110: &epson.DMD110Profile,
}

// GetModelProfile returns the profile for the specified model.
func GetModelProfile(model types.Model) (*types.ModelProfile, bool) {
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
