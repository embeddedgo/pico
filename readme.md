## Support for Raspberry PI microcontrollers

Embedded Go supports the RP2350 family, aka Pico 2 (there is no support for RP2040).

Your program will run on both ARM cores (the RISCV mode is not supported).

### Prerequisites

1. Go complier.

   You can download it from [go.dev/dl](https://go.dev/dl/).

2. Git command.

   To instll git on Linux use the package manager provided by your Linux distribution (apt, pacman, rpm, ...).

   Windows users may check the [git for Windows](https://gitforwindows.org/) website.

   The Mac users may use the git command provided by the [Xcode](https://developer.apple.com/xcode/) commandline tools. Another way is to use the [Homebrew](https://brew.sh/) package manager.

### Getting started

1. Install the Embedded Go toolchain.

   Make sure the `$GOPATH/bin` directory is in your `PATH`, as tools installed with the `go install` command will be placed here. If you didn't set the `GOPATH` environment variable manually you can find its default value using the `go env GOPATH` command.

   Then install the Embedded Go toolchain using the following two commands:

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

4. Initialize your project.

   ```sh
   go mod init firstprog
   go mod tidy
   ```

5. Copy the `go.env` file suitable for your board (here is one for [Pico 2](devboard/pico2/examples/go.env) and another one for a [board with 16 MB flash](devboard/weacta10/examples/go.env)).

6. Compile your first program.

   ```sh
   export GOENV=go.env
   go build
   ```

   or

   ```sh
   GOENV=go.env go build
   ```

   or

   ```sh
   egtool build
   ```

   The last one is like `GOENV=go.env go build` but looks for the `go.env` file up the current module directory tree.

7. Connect your Pico 2 to the computer in the BOOT mode (press the onboard button while connecting it to the USB).

8. Load and run.

   ```sh
   egtool load
   ```

9. See the [Embedded Go](https://embeddedgo.github.io/) website for more information.

### Examples

See more example code for [supported develompent boards](devboard).

