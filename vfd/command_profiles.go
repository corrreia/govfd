package vfd

import "vfd/vfd/commands/escpos"

// CommandProtocolName represents a specific command protocol identifier.
type CommandProtocolName string

// Supported command protocols
const (
	// ESC/POS protocol - most common for VFD displays
	ProtocolESCPOS CommandProtocolName = "ESC_POS"

	// Add more protocols here as support is added
	// ProtocolCustom CommandProtocolName = "CUSTOM"
	// ProtocolAPI CommandProtocolName = "API"
	// ProtocolWebSocket CommandProtocolName = "WEBSOCKET"
)

// CommandProtocol interface defines high-level operations that any VFD command protocol must implement.
// This abstracts away the specific implementation (bytes, API calls, etc.) used by different protocols.
type CommandProtocol interface {
	// Protocol identification
	GetName() string
	GetDescription() string

	// Display control operations
	Clear() []byte                     // Initialize/clear display
	FormFeed() []byte                  // Clear screen content
	MoveCursor(column, row int) []byte // Move cursor to position (1-based)
	WriteText(text string) []byte      // Write text at current position

	// Display settings
	SetBrightness(level int) []byte // Set brightness (1-4 typically)
	SetBlink(intervalMs int) []byte // Set cursor blink (0=off)
	SetCharset(page int) []byte     // Set character encoding table

	// Utility operations
	SelfTest() []byte // Execute self-test
}

// CommandProfile contains the implementation for a specific protocol.
type CommandProfile struct {
	// Protocol identification
	Name        string
	Description string

	// Protocol implementation
	Protocol CommandProtocol
}

// commandRegistry contains profiles for all supported command protocols.
var commandRegistry = map[CommandProtocolName]*CommandProfile{
	ProtocolESCPOS: {
		Name:        "ESC/POS",
		Description: "Standard ESC/POS command set for VFD displays",
		Protocol:    &escpos.ESCPOSProtocol{},
	},
}

// GetCommandProfile returns the command profile for the specified protocol.
func GetCommandProfile(protocolName CommandProtocolName) (*CommandProfile, bool) {
	profile, exists := commandRegistry[protocolName]
	return profile, exists
}

// GetSupportedProtocols returns a list of all supported command protocols.
func GetSupportedProtocols() []CommandProtocolName {
	protocols := make([]CommandProtocolName, 0, len(commandRegistry))
	for protocol := range commandRegistry {
		protocols = append(protocols, protocol)
	}
	return protocols
}
