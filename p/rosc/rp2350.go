// DO NOT EDIT THIS FILE. GENERATED BY svdxgen.

//go:build rp2350

// Package rosc provides access to the registers of the ROSC peripheral.
//
// Instances:
//
//	ROSC  ROSC_BASE  -  -
//
// Registers:
//
//	0x000 32  CTRL       Ring Oscillator control
//	0x004 32  FREQA      The FREQA & FREQB registers control the frequency by controlling the drive strength of each stage The drive strength has 4 levels determined by the number of bits set Increasing the number of bits set increases the drive strength and increases the oscillation frequency 0 bits set is the default drive strength 1 bit set doubles the drive strength 2 bits set triples drive strength 3 bits set quadruples drive strength For frequency randomisation set both DS0_RANDOM=1 & DS1_RANDOM=1
//	0x008 32  FREQB      For a detailed description see freqa register
//	0x00C 32  RANDOM     Loads a value to the LFSR randomiser
//	0x010 32  DORMANT    Ring Oscillator pause control
//	0x014 32  DIV        Controls the output divider
//	0x018 32  PHASE      Controls the phase shifted output
//	0x01C 32  STATUS     Ring Oscillator Status
//	0x020 32  RANDOMBIT  This just reads the state of the oscillator output so randomness is compromised if the ring oscillator is stopped or run at a harmonic of the bus frequency
//	0x024 32  COUNT      A down counter running at the ROSC frequency which counts to zero and stops. To start the counter write a non-zero value. Can be used for short software pauses when setting up time sensitive hardware.
//
// Import:
//
//	github.com/embeddedgo/pico/p/mmap
package rosc

const (
	FREQ_RANGE CTRL = 0xFFF << 0 //+ Controls the number of delay stages in the ROSC ring LOW uses stages 0 to 7 MEDIUM uses stages 2 to 7 HIGH uses stages 4 to 7 TOOHIGH uses stages 6 to 7 and should not be used because its frequency exceeds design specifications The clock output will not glitch when changing the range up one step at a time The clock output will glitch when changing the range down Note: the values here are gray coded which is why HIGH comes before TOOHIGH
	LOW        CTRL = 0xFA4 << 0
	MEDIUM     CTRL = 0xFA5 << 0
	TOOHIGH    CTRL = 0xFA6 << 0
	HIGH       CTRL = 0xFA7 << 0
	ENABLE     CTRL = 0xFFF << 12 //+ On power-up this field is initialised to ENABLE The system clock must be switched to another source before setting this field to DISABLE otherwise the chip will lock up The 12-bit code is intended to give some protection against accidental writes. An invalid setting will enable the oscillator.
	DISABLE    CTRL = 0xD1E << 12
	ENABLE     CTRL = 0xFAB << 12
)

const (
	FREQ_RANGEn = 0
	ENABLEn     = 12
)

const (
	DS0        FREQA = 0x07 << 0    //+ Stage 0 drive strength
	DS0_RANDOM FREQA = 0x01 << 3    //+ Randomises the stage 0 drive strength
	DS1        FREQA = 0x07 << 4    //+ Stage 1 drive strength
	DS1_RANDOM FREQA = 0x01 << 7    //+ Randomises the stage 1 drive strength
	DS2        FREQA = 0x07 << 8    //+ Stage 2 drive strength
	DS3        FREQA = 0x07 << 12   //+ Stage 3 drive strength
	PASSWD     FREQA = 0xFFFF << 16 //+ Set to 0x9696 to apply the settings Any other value in this field will set all drive strengths to 0
	PASS       FREQA = 0x9696 << 16
)

const (
	DS0n        = 0
	DS0_RANDOMn = 3
	DS1n        = 4
	DS1_RANDOMn = 7
	DS2n        = 8
	DS3n        = 12
	PASSWDn     = 16
)

const (
	DS4    FREQB = 0x07 << 0    //+ Stage 4 drive strength
	DS5    FREQB = 0x07 << 4    //+ Stage 5 drive strength
	DS6    FREQB = 0x07 << 8    //+ Stage 6 drive strength
	DS7    FREQB = 0x07 << 12   //+ Stage 7 drive strength
	PASSWD FREQB = 0xFFFF << 16 //+ Set to 0x9696 to apply the settings Any other value in this field will set all drive strengths to 0
	PASS   FREQB = 0x9696 << 16
)

const (
	DS4n    = 0
	DS5n    = 4
	DS6n    = 8
	DS7n    = 12
	PASSWDn = 16
)

const (
	DORMANT DORMANT = 0xFFFFFFFF << 0 //+ This is used to save power by pausing the ROSC On power-up this field is initialised to WAKE An invalid write will also select WAKE Warning: setup the irq before selecting dormant mode
	DORMANT DORMANT = 0x636F6D61 << 0
	WAKE    DORMANT = 0x77616B65 << 0
)

const (
	DORMANTn = 0
)

const (
	DIV  DIV = 0xFFFF << 0 //+ set to 0xaa00 + div where div = 0 divides by 128 div = 1-127 divides by div any other value sets div=128 this register resets to div=32
	PASS DIV = 0xAA00 << 0
)

const (
	DIVn = 0
)

const (
	SHIFT  PHASE = 0x03 << 0 //+ phase shift the phase-shifted output by SHIFT input clocks this can be changed on-the-fly must be set to 0 before setting div=1
	FLIP   PHASE = 0x01 << 2 //+ invert the phase-shifted output this is ignored when div=1
	ENABLE PHASE = 0x01 << 3 //+ enable the phase-shifted output this can be changed on-the-fly
	PASSWD PHASE = 0xFF << 4 //+ set to 0xaa any other value enables the output with shift=0
)

const (
	SHIFTn  = 0
	FLIPn   = 2
	ENABLEn = 3
	PASSWDn = 4
)

const (
	ENABLED     STATUS = 0x01 << 12 //+ Oscillator is enabled but not necessarily running and stable this resets to 0 but transitions to 1 during chip startup
	DIV_RUNNING STATUS = 0x01 << 16 //+ post-divider is running this resets to 0 but transitions to 1 during chip startup
	BADWRITE    STATUS = 0x01 << 24 //+ An invalid value has been written to CTRL_ENABLE or CTRL_FREQ_RANGE or FREQA or FREQB or DIV or PHASE or DORMANT
	STABLE      STATUS = 0x01 << 31 //+ Oscillator is running and stable
)

const (
	ENABLEDn     = 12
	DIV_RUNNINGn = 16
	BADWRITEn    = 24
	STABLEn      = 31
)
