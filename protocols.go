package govfd

import (
	"github.com/corrreia/govfd/commands/escpos"
)

// Protocol interface defines high-level operations that any VFD command protocol must implement.
// This abstracts away the specific implementation (bytes, API calls, etc.) used by different protocols.
type Protocol interface {
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

// commandRegistry contains implementations for all supported command protocols.
var commandRegistry = []Protocol{
	&escpos.ESCPOSProtocol{},
}

// GetProtocol returns the command protocol implementation for the specified protocol name.
func GetProtocol(protocolName string) (Protocol, bool) {
	for _, protocol := range commandRegistry {
		if protocol.GetName() == protocolName {
			return protocol, true
		}
	}
	return nil, false
}

// GetSupportedProtocols returns a list of all supported command protocol names.
func GetSupportedProtocols() []string {
	protocolList := make([]string, 0, len(commandRegistry))
	for _, protocol := range commandRegistry {
		protocolList = append(protocolList, protocol.GetName())
	}
	return protocolList
}
