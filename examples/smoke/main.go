package main

import (
	"fmt"
	"os"
	"time"

	"github.com/corrreia/govfd"
	"github.com/corrreia/govfd/types"
)

func main() {
	ports := []string{"/dev/ttyUSB0", "/dev/ttyUSB1"}

	for i, port := range ports {
		fmt.Printf("Trying display %d on %s...\n", i+1, port)

		display, err := govfd.OpenModel(port, types.ModelEpsonDMD110)
		if err != nil {
			fmt.Printf("  SKIP: %v\n", err)
			continue
		}

		// Clear and identify
		display.Clear()
		time.Sleep(200 * time.Millisecond)

		display.WriteText(fmt.Sprintf("Display %d", i+1))
		display.SetCursor(1, 2)
		display.WriteText(port)

		time.Sleep(2 * time.Second)

		// Test Latin encoding
		display.Clear()
		time.Sleep(200 * time.Millisecond)
		display.WriteText("Latin: cafe")
		display.SetCursor(1, 2)
		display.WriteText("PT: acao")

		time.Sleep(2 * time.Second)

		// Test Portuguese (requires charset switch)
		display.Clear()
		time.Sleep(200 * time.Millisecond)
		display.WriteText("Portugues:")
		display.SetCursor(1, 2)

		err = display.WriteText("ação coração")
		if err != nil {
			fmt.Printf("  WriteText error: %v\n", err)
		}

		time.Sleep(2 * time.Second)

		// Final message
		display.Clear()
		time.Sleep(200 * time.Millisecond)
		display.WriteText("govfd OK!")
		display.SetCursor(1, 2)
		display.WriteText("Display " + fmt.Sprintf("%d", i+1) + " works!")

		display.Close()
		fmt.Printf("  Done.\n")
	}

	if len(os.Args) > 1 && os.Args[1] == "--wait" {
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
	}
}
