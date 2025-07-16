## Support for Raspberry PI microcontrollers

Embedded Go supports the RP2350 family, aka Pico 2 (there is no support for RP2040).

Your program will run on both ARM cores (the RISCV mode is not supported).

### Getting started

1. Install the Embedded Go toolchain.

   ```sh
   go install github.com/embeddedgo/dl/go1.24.5-embedded@latest
   go1.24.5-embedded download
   ```

2. Install egtool.

   ```sh
   go install github.com/embeddedgo/tools/egtool@latest
   ```

3. Create a project directory containing the `main.go` file with your first Go program for RPi Pico 2.

   ```go
   package main

   import (
   	"time"

   	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
   )

   func main() {
   	for {
   		leds.User.Toggle()
   		time.Sleep(time.Second/2)
   	}
   }
   ```

4. Initialize your project

   ```sh
   go mod init firstprog
   go mod tidy
   ```

5. Copy the `go.env` file suitable for your board (here is one for [Pico 2](https://github.com/embeddedgo/pico/tree/master/devboard/pico2/examples/go.env) and another one for a [board with 16 MB flash](https://github.com/embeddedgo/pico/tree/master/devboard/wiacta10/examples/go.env)).

6. Compile your first program

   ```sh
   export GOENV=go.env
   go build
   ```

   or

   ```sh
   GOENV=go.env go build
   ```

7. Connect your Pico 2 to your computer in the BOOT mode (press the onboard button while connecting it to the USB).

8. Load and run.

   ```sh
   egtool load
   ```

### Examples

See more example code for [supported develompent boards](devboard).

