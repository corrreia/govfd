package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/corrreia/govfd"
	"github.com/corrreia/govfd/types"
)

func main() {
	// Change COM port as needed
	portName := "COM3" // <-- set to your actual COM port

	// Open using model-specific defaults for Epson DM-D110
	// This automatically sets: 20x2 display, 9600 baud, 8N1
	display, err := govfd.OpenModel(portName, types.ModelEpsonDMD110)

	// Alternative approaches:
	// 1. Use model defaults but override specific settings:
	//    opts := &govfd.Options{BaudRate: 19200}  // override baud rate
	//    display, err := govfd.OpenModelWithOptions(portName, types.ModelEpsonDMD110, opts)
	//
	// 2. Manual configuration (legacy approach):
	//    opts := &govfd.Options{BaudRate: 9600, DataBits: 8, Columns: 20, Rows: 2}
	//    display, err := govfd.Open(portName, opts)
	if err != nil {
		fmt.Println("Error opening port:", err)
		return
	}
	defer display.Close()

	fmt.Println("Connected to", portName, "using Epson DM-D110 profile")
	fmt.Println("Type 'help' for commands.")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		args := strings.Split(line, " ")
		cmd := strings.ToLower(args[0])

		switch cmd {
		case "quit", "exit":
			fmt.Println("Bye.")
			return

		case "help":
			fmt.Println("Commands:")
			fmt.Println("  clear                         - ESC @ (initialize)")
			fmt.Println("  clr                           - Form feed (0x0C)")
			fmt.Println("  pos <col> <row>               - Move cursor (1-based)")
			fmt.Println("  cur                           - Show current cursor (1-based)")
			fmt.Println("  text <message>                - Send text with AUTOMATIC encoding! 🚀")
			fmt.Println("  bright <1-4>                  - Set brightness")
			fmt.Println("  getbright                     - Show last set brightness (0 if unknown)")
			fmt.Println("  blink <ms>                    - Set cursor blink period in ms (0=off)")
			fmt.Println("  selftest, test                - Execute display self-test (US @)")
			fmt.Println("  models                        - List all supported VFD models")
			fmt.Println("  protocols                     - List all supported command protocols")
			fmt.Println("  info                          - Show current display model information")
			fmt.Println("  exit, quit                    - Exit program")
			fmt.Println("  ")
			fmt.Println("💡 TIP: Just type 'text café' or 'text ação' - encoding is automatic!")

		case "clear":
			if err := display.Clear(); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Display cleared.")

		case "clr":
			if err := display.FormFeed(); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Screen cleared.")

		case "pos":
			if len(args) != 3 {
				fmt.Println("Usage: pos <col> <row>")
				continue
			}
			col, _ := strconv.Atoi(args[1])
			row, _ := strconv.Atoi(args[2])
			if err := display.SetCursor(col, row); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("Moved to col=%d row=%d\n", col, row)

		case "cur":
			col, row := display.GetCursor()
			fmt.Printf("Cursor at col=%d row=%d\n", col, row)

		case "text":
			if len(args) < 2 {
				fmt.Println("Usage: text <message>")
				continue
			}
			msg := strings.Join(args[1:], " ")
			if err := display.WriteText(msg); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Text sent.")

		case "bright":
			if len(args) != 2 {
				fmt.Println("Usage: bright <1-4>")
				continue
			}
			level, _ := strconv.Atoi(args[1])
			if err := display.SetBrightness(level); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("Brightness set to %d\n", level)

		case "getbright":
			fmt.Printf("Current brightness (last set): %d\n", display.GetBrightness())

		case "getdims":
			cols, rows := display.Dimensions()
			if cols == 0 || rows == 0 {
				fmt.Println("Dimensions not set.")
			} else {
				fmt.Printf("Dimensions: %dx%d (cols x rows)\n", cols, rows)
			}

		case "blink":
			if len(args) != 2 {
				fmt.Println("Usage: blink <ms> (0=off)")
				continue
			}
			ms, err := strconv.Atoi(args[1])
			if err != nil || ms < 0 {
				fmt.Println("blink <ms> (0=off)")
				continue
			}
			if err := display.SetBlink(ms); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			steps := ms / 50
			fmt.Printf("Blink set to %d ms (steps=%d)\n", steps*50, steps)

		case "selftest", "test":
			if err := display.SelfTest(); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Self-test executed.")

		case "models":
			fmt.Println("Supported VFD models:")
			models := govfd.GetSupportedModels()
			for _, model := range models {
				if profile, exists := govfd.GetModelProfile(model); exists {
					fmt.Printf("  %s - %s (%dx%d, %d baud, %s)\n",
						model, profile.Name, profile.Columns, profile.Rows,
						profile.DefaultBaudRate, profile.CommandProtocol)
				}
			}

		case "protocols":
			fmt.Println("Supported command protocols:")
			protocols := govfd.GetSupportedProtocols()
			for _, protocolName := range protocols {
				if protocol, exists := govfd.GetProtocol(protocolName); exists {
					fmt.Printf("  %s - %s\n", protocolName, protocol.GetDescription())
				}
			}

		case "info":
			// Show info for the current model (Epson DM-D110)
			if profile, exists := govfd.GetModelProfile(types.ModelEpsonDMD110); exists {
				fmt.Printf("Current Display Model: %s\n", profile.Name)
				fmt.Printf("  Manufacturer: %s\n", profile.Manufacturer)
				fmt.Printf("  Model: %s\n", profile.Model)
				fmt.Printf("  Dimensions: %d columns x %d rows\n", profile.Columns, profile.Rows)
				fmt.Printf("  Default Baud Rate: %d\n", profile.DefaultBaudRate)
				fmt.Printf("  Command Protocol: %s\n", profile.CommandProtocol)
				if protocol, exists := govfd.GetProtocol(profile.CommandProtocol); exists {
					fmt.Printf("  Protocol Description: %s\n", protocol.GetDescription())
				}
				fmt.Printf("  Brightness Levels: %d (1-%d)\n", profile.BrightnessLevels, profile.BrightnessLevels)
				fmt.Printf("  Supports: ")
				features := []string{}
				if profile.SupportsBrightness {
					features = append(features, "Brightness")
				}
				if profile.SupportsCursorBlink {
					features = append(features, "Cursor Blink")
				}
				if profile.SupportsCharsetTable {
					features = append(features, "Character Sets")
				}
				if profile.SupportsSelfTest {
					features = append(features, "Self Test")
				}
				fmt.Println(strings.Join(features, ", "))
				if profile.DocumentationURL != "" {
					fmt.Printf("  Documentation: %s\n", profile.DocumentationURL)
				}
			}

		default:
			fmt.Println("Unknown command. Type 'help'.")
		}
	}
}
