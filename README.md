# VFD Serial Display Library (Go)

A Go library and CLI tool for controlling VFD (Vacuum Fluorescent Display) devices over serial connections.

## Features

- **Model-based configuration**: Automatic setup for supported VFD models
- **Command protocol flexibility**: Support for different command sets (ESC/POS by default)
- **Multiple display support**: Easily extensible for different VFD models and protocols
- **Comprehensive command set**: Clear, cursor positioning, brightness, character sets, and more

## Supported Models

- **Epson DM-D110**: 20x2 customer display (9600 baud, 8N1, ESC/POS protocol)

## Command Protocols

- **ESC/POS**: Standard ESC/POS command set for VFD displays (default)

## Library Usage

### Quick Start (Recommended)

```go
package main

import (
    "log"
    "vfd/vfd"
)

func main() {
    // Open with model-specific defaults (recommended)
    d, err := vfd.OpenModel("COM3", vfd.ModelEpsonDMD110)
    if err != nil {
        log.Fatal(err)
    }
    defer d.Close()

    if err := d.Clear(); err != nil { log.Fatal(err) }
    if err := d.SetCursor(1, 1); err != nil { log.Fatal(err) }
    if err := d.WriteText("Hello VFD"); err != nil { log.Fatal(err) }
    if err := d.SetBrightness(3); err != nil { log.Fatal(err) }
}
```

### Advanced Usage

```go
// Override model defaults
opts := &vfd.Options{BaudRate: 19200} // Use 19200 instead of model default
d, err := vfd.OpenModelWithOptions("COM3", vfd.ModelEpsonDMD110, opts)

// Manual configuration (legacy approach)
opts := &vfd.Options{BaudRate: 9600, DataBits: 8, Columns: 20, Rows: 2}
d, err := vfd.Open("COM3", opts)
```

## Build

```bash
# You compile it; the assistant will not build automatically.
go mod tidy

go build -o vfd.exe ./
```

## Notes

- If unsure which COM port your PL2303HXA uses, check Windows Device Manager under "Ports (COM & LPT)".
- Typical VFD modules accept 9600 8N1, but confirm your device's requirements. Adjust `-baud`, `-dataBits`, `-stopBits`, and `-parity` as needed.
- If nothing shows on the display, try adding `-crlf` or `-cr`/`-lf` depending on the module.
