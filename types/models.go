package types

// Model represents a specific VFD display model.
type Model string

// Supported VFD display models
const (
	// Epson DM-D110 customer display
	ModelEpsonDMD110 Model = "EPSON_DM_D110"

	// Add more models here as support is added
	// ModelEpsonDMD210 Model = "EPSON_DM_D210"
	// ModelBixolonBCD1100 Model = "BIXOLON_BCD_1100"
)
