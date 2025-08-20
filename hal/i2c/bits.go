// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c

type CON uint32

const (
	MASTER_MODE               CON = 0x01 << 0  //+ This bit controls whether the DW_apb_i2c master is enabled. NOTE: Software should ensure that if this bit is written with '1' then bit 6 should also be written with a '1'.
	SPEED                     CON = 0x03 << 1  //+ These bits control at which speed the DW_apb_i2c operates; its setting is relevant only if one is operating the DW_apb_i2c in master mode. Hardware protects against illegal values being programmed by software. These bits must be programmed appropriately for slave mode also, as it is used to capture correct value of spike filter as per the speed mode. This register should be programmed only with a value in the range of 1 to MAX_SPEED_MODE; otherwise, hardware updates this register with the value of MAX_SPEED_MODE. 1: standard mode (100 kbit/s) 2: fast mode (<=400 kbit/s) or fast mode plus (<=1000Kbit/s) 3: high speed mode (3.4 Mbit/s) Note: This field is not applicable when ULTRA_FAST_MODE=1
	STANDARD                  CON = 0x01 << 1  //  Standard Speed mode of operation
	FAST                      CON = 0x02 << 1  //  Fast or Fast Plus mode of operation
	HIGH                      CON = 0x03 << 1  //  High Speed mode of operation
	SLAVE_10BITADDR           CON = 0x01 << 3  //+ When acting as a slave, this bit controls whether the DW_apb_i2c responds to 7- or 10-bit addresses. - 0: 7-bit addressing. The DW_apb_i2c ignores transactions that involve 10-bit addressing; for 7-bit addressing, only the lower 7 bits of the SAR register are compared. - 1: 10-bit addressing. The DW_apb_i2c responds to only 10-bit addressing transfers that match the full 10 bits of the SAR register.
	MASTER_10BITADDR          CON = 0x01 << 4  //+ Controls whether the DW_apb_i2c starts its transfers in 7- or 10-bit addressing mode when acting as a master. - 0: 7-bit addressing - 1: 10-bit addressing
	RESTART_EN                CON = 0x01 << 5  //+ Determines whether RESTART conditions may be sent when acting as a master. Some older slaves do not support handling RESTART conditions; however, RESTART conditions are used in several DW_apb_i2c operations. When RESTART is disabled, the master is prohibited from performing the following functions: - Sending a START BYTE - Performing any high-speed mode operation - High-speed mode operation - Performing direction changes in combined format mode - Performing a read operation with a 10-bit address By replacing RESTART condition followed by a STOP and a subsequent START condition, split operations are broken down into multiple DW_apb_i2c transfers. If the above operations are performed, it will result in setting bit 6 (TX_ABRT) of the RAW_INTR_STAT register. Reset value: ENABLED
	SLAVE_DISABLE             CON = 0x01 << 6  //+ This bit controls whether I2C has its slave disabled, which means once the presetn signal is applied, then this bit is set and the slave is disabled. If this bit is set (slave is disabled), DW_apb_i2c functions only as a master and does not perform any action that requires a slave. NOTE: Software should ensure that if this bit is written with 0, then bit 0 should also be written with a 0.
	STOP_DET_IFADDRESSED      CON = 0x01 << 7  //+ In slave mode: - 1'b1: issues the STOP_DET interrupt only when it is addressed. - 1'b0: issues the STOP_DET irrespective of whether it's addressed or not. Reset value: 0x0 NOTE: During a general call address, this slave does not issue the STOP_DET interrupt if STOP_DET_IF_ADDRESSED = 1'b1, even if the slave responds to the general call address by generating ACK. The STOP_DET interrupt is generated only when the transmitted address matches the slave address (SAR).
	TX_EMPTY_CTRL             CON = 0x01 << 8  //+ This bit controls the generation of the TX_EMPTY interrupt, as described in the RAW_INTR_STAT register. Reset value: 0x0.
	RX_FIFO_FULL_HLD_CTRL     CON = 0x01 << 9  //+ This bit controls whether DW_apb_i2c should hold the bus when the Rx FIFO is physically full to its RX_BUFFER_DEPTH, as described in the RX_FULL_HLD_BUS_EN parameter. Reset value: 0x0.
	STOP_DET_IF_MASTER_ACTIVE CON = 0x01 << 10 //+ Master issues the STOP_DET interrupt irrespective of whether master is active or not
)

const (
	MASTER_MODEn               = 0
	SPEEDn                     = 1
	SLAVE_10BITADDRn           = 3
	MASTER_10BITADDRn          = 4
	RESTART_ENn                = 5
	SLAVE_DISABLEn             = 6
	STOP_DET_IFADDRESSEDn      = 7
	TX_EMPTY_CTRLn             = 8
	RX_FIFO_FULL_HLD_CTRLn     = 9
	STOP_DET_IF_MASTER_ACTIVEn = 10
)

type TAR uint32

const (
	ADDR        TAR = 0x3FF << 0 //+ This is the target address for any master transaction. When transmitting a General Call, these bits are ignored. To generate a START BYTE, the CPU needs to write only once into these bits. If the TAR and SAR are the same, loopback exists but the FIFOs are shared between master and slave, so full loopback is not feasible. Only one direction loopback mode is supported (simplex), not duplex. A master cannot transmit to itself; it can transmit to only a slave.
	GC_OR_START TAR = 0x01 << 10 //+ If bit 11 (SPECIAL) is set to 1 and bit 13(Device-ID) is set to 0, then this bit indicates whether a General Call or START byte command is to be performed by the DW_apb_i2c. - 0: General Call Address - after issuing a General Call, only writes may be performed. Attempting to issue a read command results in setting bit 6 (TX_ABRT) of the RAW_INTR_STAT register. The DW_apb_i2c remains in General Call mode until the SPECIAL bit value (bit 11) is cleared. - 1: START BYTE Reset value: 0x0
	SPECIAL     TAR = 0x01 << 11 //+ This bit indicates whether software performs a Device-ID or General Call or START BYTE command. - 0: ignore bit 10 GC_OR_START and use TAR normally - 1: perform special I2C command as specified in Device_ID or GC_OR_START bit Reset value: 0x0
)

const (
	ADDRn        = 0
	GC_OR_STARTn = 10
	SPECIALn     = 11
)

// DATA_CMD bits
const (
	DAT             uint32 = 0xFF << 0  //+ This register contains the data to be transmitted or received on the I2C bus. If you are writing to this register and want to perform a read, bits 7:0 (DAT) are ignored by the DW_apb_i2c. However, when you read this register, these bits return the value of data received on the DW_apb_i2c interface. Reset value: 0x0
	CMD             uint32 = 0x01 << 8  //+ This bit controls whether a read or a write is performed. This bit does not control the direction when the DW_apb_i2con acts as a slave. It controls only the direction when it acts as a master. When a command is entered in the TX FIFO, this bit distinguishes the write and read commands. In slave-receiver mode, this bit is a 'don't care' because writes to this register are not required. In slave-transmitter mode, a '0' indicates that the data in DATA_CMD is to be transmitted. When programming this bit, you should remember the following: attempting to perform a read operation after a General Call command has been sent results in a TX_ABRT interrupt (bit 6 of the RAW_INTR_STAT register), unless bit 11 (SPECIAL) in the TAR register has been cleared. If a '1' is written to this bit after receiving a RD_REQ interrupt, then a TX_ABRT interrupt occurs. Reset value: 0x0
	STOP            uint32 = 0x01 << 9  //+ This bit controls whether a STOP is issued after the byte is sent or received. - 1 - STOP is issued after this byte, regardless of whether or not the Tx FIFO is empty. If the Tx FIFO is not empty, the master immediately tries to start a new transfer by issuing a START and arbitrating for the bus. - 0 - STOP is not issued after this byte, regardless of whether or not the Tx FIFO is empty. If the Tx FIFO is not empty, the master continues the current transfer by sending/receiving data bytes according to the value of the CMD bit. If the Tx FIFO is empty, the master holds the SCL line low and stalls the bus until a new command is available in the Tx FIFO. Reset value: 0x0
	RESTART         uint32 = 0x01 << 10 //+ This bit controls whether a RESTART is issued before the byte is sent or received. 1 - If RESTART_EN is 1, a RESTART is issued before the data is sent/received (according to the value of CMD), regardless of whether or not the transfer direction is changing from the previous command; if RESTART_EN is 0, a STOP followed by a START is issued instead. 0 - If RESTART_EN is 1, a RESTART is issued only if the transfer direction is changing from the previous command; if RESTART_EN is 0, a STOP followed by a START is issued instead. Reset value: 0x0
	FIRST_DATA_BYTE uint32 = 0x01 << 11 //+ Indicates the first data byte received after the address phase for receive transfer in Master receiver or Slave receiver mode. Reset value : 0x0 NOTE: In case of APB_DATA_WIDTH=8, 1. The user has to perform two APB Reads to DATA_CMD in order to get status on 11 bit. 2. In order to read the 11 bit, the user has to perform the first data byte read [7:0] (offset 0x10) and then perform the second read [15:8] (offset 0x11) in order to know the status of 11 bit (whether the data received in previous read is a first data byte or not). 3. The 11th bit is an optional read field, user can ignore 2nd byte read [15:8] (offset 0x11) if not interested in FIRST_DATA_BYTE status.
)

const (
	DATn             = 0
	CMDn             = 8
	STOPn            = 9
	RESTARTn         = 10
	FIRST_DATA_BYTEn = 11
)

type INTR uint32

const (
	RX_UNDER    INTR = 0x01 << 0  //+ Set if the processor attempts to read the receive buffer when it is empty by reading from the DATA_CMD register. If the module is disabled (ENABLE[0]=0), this bit keeps its level until the master or slave state machines go into idle, and when ic_en goes to 0, this interrupt is cleared. Reset value: 0x0
	RX_OVER     INTR = 0x01 << 1  //+ Set if the receive buffer is completely filled to RX_BUFFER_DEPTH and an additional byte is received from an external I2C device. The DW_apb_i2c acknowledges this, but any data bytes received after the FIFO is full are lost. If the module is disabled (ENABLE[0]=0), this bit keeps its level until the master or slave state machines go into idle, and when ic_en goes to 0, this interrupt is cleared. Note: If bit 9 of the CON register (RX_FIFO_FULL_HLD_CTRL) is programmed to HIGH, then the RX_OVER interrupt never occurs, because the Rx FIFO never overflows. Reset value: 0x0
	RX_FULL     INTR = 0x01 << 2  //+ Set when the receive buffer reaches or goes above the RX_TL threshold in the RX_TL register. It is automatically cleared by hardware when buffer level goes below the threshold. If the module is disabled (ENABLE[0]=0), the RX FIFO is flushed and held in reset; therefore the RX FIFO is not full. So this bit is cleared once the ENABLE bit 0 is programmed with a 0, regardless of the activity that continues. Reset value: 0x0
	TX_OVER     INTR = 0x01 << 3  //+ Set during transmit if the transmit buffer is filled to TX_BUFFER_DEPTH and the processor attempts to issue another I2C command by writing to the DATA_CMD register. When the module is disabled, this bit keeps its level until the master or slave state machines go into idle, and when ic_en goes to 0, this interrupt is cleared. Reset value: 0x0
	TX_EMPTY    INTR = 0x01 << 4  //+ The behavior of the TX_EMPTY interrupt status differs based on the TX_EMPTY_CTRL selection in the CON register. - When TX_EMPTY_CTRL = 0: This bit is set to 1 when the transmit buffer is at or below the threshold value set in the TX_TL register. - When TX_EMPTY_CTRL = 1: This bit is set to 1 when the transmit buffer is at or below the threshold value set in the TX_TL register and the transmission of the address/data from the internal shift register for the most recently popped command is completed. It is automatically cleared by hardware when the buffer level goes above the threshold. When ENABLE[0] is set to 0, the TX FIFO is flushed and held in reset. There the TX FIFO looks like it has no data within it, so this bit is set to 1, provided there is activity in the master or slave state machines. When there is no longer any activity, then with ic_en=0, this bit is set to 0. Reset value: 0x0.
	RD_REQ      INTR = 0x01 << 5  //+ This bit is set to 1 when DW_apb_i2c is acting as a slave and another I2C master is attempting to read data from DW_apb_i2c. The DW_apb_i2c holds the I2C bus in a wait state (SCL=0) until this interrupt is serviced, which means that the slave has been addressed by a remote master that is asking for data to be transferred. The processor must respond to this interrupt and then write the requested data to the DATA_CMD register. This bit is set to 0 just after the processor reads the CLR_RD_REQ register. Reset value: 0x0
	TX_ABRT     INTR = 0x01 << 6  //+ This bit indicates if DW_apb_i2c, as an I2C transmitter, is unable to complete the intended actions on the contents of the transmit FIFO. This situation can occur both as an I2C master or an I2C slave, and is referred to as a 'transmit abort'. When this bit is set to 1, the TX_ABRT_SOURCE register indicates the reason why the transmit abort takes places. Note: The DW_apb_i2c flushes/resets/empties the TX_FIFO and RX_FIFO whenever there is a transmit abort caused by any of the events tracked by the TX_ABRT_SOURCE register. The FIFOs remains in this flushed state until the register CLR_TX_ABRT is read. Once this read is performed, the Tx FIFO is then ready to accept more data bytes from the APB interface. Reset value: 0x0
	RX_DONE     INTR = 0x01 << 7  //+ When the DW_apb_i2c is acting as a slave-transmitter, this bit is set to 1 if the master does not acknowledge a transmitted byte. This occurs on the last byte of the transmission, indicating that the transmission is done. Reset value: 0x0
	ACTIVITY    INTR = 0x01 << 8  //+ This bit captures DW_apb_i2c activity and stays set until it is cleared. There are four ways to clear it: - Disabling the DW_apb_i2c - Reading the CLR_ACTIVITY register - Reading the CLR_INTR register - System reset Once this bit is set, it stays set unless one of the four methods is used to clear it. Even if the DW_apb_i2c module is idle, this bit remains set until cleared, indicating that there was activity on the bus. Reset value: 0x0
	STOP_DET    INTR = 0x01 << 9  //+ Indicates whether a STOP condition has occurred on the I2C interface regardless of whether DW_apb_i2c is operating in slave or master mode. In Slave Mode: - If CON[7]=1'b1 (STOP_DET_IFADDRESSED), the STOP_DET interrupt will be issued only if slave is addressed. Note: During a general call address, this slave does not issue a STOP_DET interrupt if STOP_DET_IF_ADDRESSED=1'b1, even if the slave responds to the general call address by generating ACK. The STOP_DET interrupt is generated only when the transmitted address matches the slave address (SAR). - If CON[7]=1'b0 (STOP_DET_IFADDRESSED), the STOP_DET interrupt is issued irrespective of whether it is being addressed. In Master Mode: - If CON[10]=1'b1 (STOP_DET_IF_MASTER_ACTIVE),the STOP_DET interrupt will be issued only if Master is active. - If CON[10]=1'b0 (STOP_DET_IFADDRESSED),the STOP_DET interrupt will be issued irrespective of whether master is active or not. Reset value: 0x0
	START_DET   INTR = 0x01 << 10 //+ Indicates whether a START or RESTART condition has occurred on the I2C interface regardless of whether DW_apb_i2c is operating in slave or master mode. Reset value: 0x0
	GEN_CALL    INTR = 0x01 << 11 //+ Set only when a General Call address is received and it is acknowledged. It stays set until it is cleared either by disabling DW_apb_i2c or when the CPU reads bit 0 of the CLR_GEN_CALL register. DW_apb_i2c stores the received data in the Rx buffer. Reset value: 0x0
	RESTART_DET INTR = 0x01 << 12 //+ Indicates whether a RESTART condition has occurred on the I2C interface when DW_apb_i2c is operating in Slave mode and the slave is being addressed. Enabled only when SLV_RESTART_DET_EN=1. Note: However, in high-speed mode or during a START BYTE transfer, the RESTART comes before the address field as per the I2C protocol. In this case, the slave is not the addressed slave when the RESTART is issued, therefore DW_apb_i2c does not generate the RESTART_DET interrupt. Reset value: 0x0
)

const (
	RX_UNDERn    = 0
	RX_OVERn     = 1
	RX_FULLn     = 2
	TX_OVERn     = 3
	TX_EMPTYn    = 4
	RD_REQn      = 5
	TX_ABRTn     = 6
	RX_DONEn     = 7
	ACTIVITYn    = 8
	STOP_DETn    = 9
	START_DETn   = 10
	GEN_CALLn    = 11
	RESTART_DETn = 12
)

type ENABLE uint32

const (
	EN           ENABLE = 0x01 << 0 //+ Controls whether the DW_apb_i2c is enabled. - 0: Disables DW_apb_i2c (TX and RX FIFOs are held in an erased state) - 1: Enables DW_apb_i2c Software can disable DW_apb_i2c while it is active. However, it is important that care be taken to ensure that DW_apb_i2c is disabled properly. A recommended procedure is described in 'Disabling DW_apb_i2c'. When DW_apb_i2c is disabled, the following occurs: - The TX FIFO and RX FIFO get flushed. - Status bits in the INTR_STAT register are still active until DW_apb_i2c goes into IDLE state. If the module is transmitting, it stops as well as deletes the contents of the transmit buffer after the current transfer is complete. If the module is receiving, the DW_apb_i2c stops the current transfer at the end of the current byte and does not acknowledge the transfer. In systems with asynchronous pclk and ic_clk when CLK_TYPE parameter set to asynchronous (1), there is a two ic_clk delay when enabling or disabling the DW_apb_i2c. For a detailed description on how to disable DW_apb_i2c, refer to 'Disabling DW_apb_i2c' Reset value: 0x0
	ABORT        ENABLE = 0x01 << 1 //+ When set, the controller initiates the transfer abort. - 0: ABORT not initiated or ABORT done - 1: ABORT operation in progress The software can abort the I2C transfer in master mode by setting this bit. The software can set this bit only when ENABLE is already set; otherwise, the controller ignores any write to ABORT bit. The software cannot clear the ABORT bit once set. In response to an ABORT, the controller issues a STOP and flushes the Tx FIFO after completing the current transfer, then sets the TX_ABORT interrupt after the abort operation. The ABORT bit is cleared automatically after the abort operation. For a detailed description on how to abort I2C transfers, refer to 'Aborting I2C Transfers'. Reset value: 0x0
	TX_CMD_BLOCK ENABLE = 0x01 << 2 //+ In Master mode: - 1'b1: Blocks the transmission of data on I2C bus even if Tx FIFO has data to transmit. - 1'b0: The transmission of data starts on I2C bus automatically, as soon as the first data is available in the Tx FIFO. Note: To block the execution of Master commands, set the TX_CMD_BLOCK bit only when Tx FIFO is empty (STATUS[2]==1) and Master is in Idle state (STATUS[5] == 0). Any further commands put in the Tx FIFO are not executed until TX_CMD_BLOCK bit is unset. Reset value: TX_CMD_BLOCK_DEFAULT
)

const (
	ENn           = 0
	ABORTn        = 1
	TX_CMD_BLOCKn = 2
)

type STATUS uint32

const (
	ACT     STATUS = 0x01 << 0 //+ I2C Activity Status. Reset value: 0x0
	TFNF    STATUS = 0x01 << 1 //+ Transmit FIFO Not Full. Set when the transmit FIFO contains one or more empty locations, and is cleared when the FIFO is full. - 0: Transmit FIFO is full - 1: Transmit FIFO is not full Reset value: 0x1
	TFE     STATUS = 0x01 << 2 //+ Transmit FIFO Completely Empty. When the transmit FIFO is completely empty, this bit is set. When it contains one or more valid entries, this bit is cleared. This bit field does not request an interrupt. - 0: Transmit FIFO is not empty - 1: Transmit FIFO is empty Reset value: 0x1
	RFNE    STATUS = 0x01 << 3 //+ Receive FIFO Not Empty. This bit is set when the receive FIFO contains one or more entries; it is cleared when the receive FIFO is empty. - 0: Receive FIFO is empty - 1: Receive FIFO is not empty Reset value: 0x0
	RFF     STATUS = 0x01 << 4 //+ Receive FIFO Completely Full. When the receive FIFO is completely full, this bit is set. When the receive FIFO contains one or more empty location, this bit is cleared. - 0: Receive FIFO is not full - 1: Receive FIFO is full Reset value: 0x0
	MST_ACT STATUS = 0x01 << 5 //+ Master FSM Activity Status. When the Master Finite State Machine (FSM) is not in the IDLE state, this bit is set. - 0: Master FSM is in IDLE state so the Master part of DW_apb_i2c is not Active - 1: Master FSM is not in IDLE state so the Master part of DW_apb_i2c is Active Note: STATUS[0]-that is, ACTIVITY bit-is the OR of SLV_ACTIVITY and MST_ACTIVITY bits. Reset value: 0x0
	SLV_ACT STATUS = 0x01 << 6 //+ Slave FSM Activity Status. When the Slave Finite State Machine (FSM) is not in the IDLE state, this bit is set. - 0: Slave FSM is in IDLE state so the Slave part of DW_apb_i2c is not Active - 1: Slave FSM is not in IDLE state so the Slave part of DW_apb_i2c is Active Reset value: 0x0
)

const (
	ACTn     = 0
	TFNFn    = 1
	TFEn     = 2
	RFNEn    = 3
	RFFn     = 4
	MST_ACTn = 5
	SLV_ACTn = 6
)

type SDA_HOLD uint32

const (
	SDA_TX_HOLD SDA_HOLD = 0xFFFF << 0 //+ Sets the required SDA hold time in units of ic_clk period, when DW_apb_i2c acts as a transmitter. Reset value: DEFAULT_SDA_HOLD[15:0].
	SDA_RX_HOLD SDA_HOLD = 0xFF << 16  //+ Sets the required SDA hold time in units of ic_clk period, when DW_apb_i2c acts as a receiver. Reset value: DEFAULT_SDA_HOLD[23:16].
)

const (
	SDA_TX_HOLDn = 0
	SDA_RX_HOLDn = 16
)

type TX_ABRT_SOURCE uint32

const (
	ABRT_7B_ADDR_NOACK   TX_ABRT_SOURCE = 0x01 << 0   //+ This field indicates that the Master is in 7-bit addressing mode and the address sent was not acknowledged by any slave. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter or Master-Receiver
	ABRT_10ADDR1_NOACK   TX_ABRT_SOURCE = 0x01 << 1   //+ This field indicates that the Master is in 10-bit address mode and the first 10-bit address byte was not acknowledged by any slave. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter or Master-Receiver
	ABRT_10ADDR2_NOACK   TX_ABRT_SOURCE = 0x01 << 2   //+ This field indicates that the Master is in 10-bit address mode and that the second address byte of the 10-bit address was not acknowledged by any slave. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter or Master-Receiver
	ABRT_TXDATA_NOACK    TX_ABRT_SOURCE = 0x01 << 3   //+ This field indicates the master-mode only bit. When the master receives an acknowledgement for the address, but when it sends data byte(s) following the address, it did not receive an acknowledge from the remote slave(s). Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter
	ABRT_GCALL_NOACK     TX_ABRT_SOURCE = 0x01 << 4   //+ This field indicates that DW_apb_i2c in master mode has sent a General Call and no slave on the bus acknowledged the General Call. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter
	ABRT_GCALL_READ      TX_ABRT_SOURCE = 0x01 << 5   //+ This field indicates that DW_apb_i2c in the master mode has sent a General Call but the user programmed the byte following the General Call to be a read from the bus (DATA_CMD[9] is set to 1). Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter
	ABRT_HS_ACKDET       TX_ABRT_SOURCE = 0x01 << 6   //+ This field indicates that the Master is in High Speed mode and the High Speed Master code was acknowledged (wrong behavior). Reset value: 0x0 Role of DW_apb_i2c: Master
	ABRT_SBYTE_ACKDET    TX_ABRT_SOURCE = 0x01 << 7   //+ This field indicates that the Master has sent a START Byte and the START Byte was acknowledged (wrong behavior). Reset value: 0x0 Role of DW_apb_i2c: Master
	ABRT_HS_NORSTRT      TX_ABRT_SOURCE = 0x01 << 8   //+ This field indicates that the restart is disabled (RESTART_EN bit (CON[5]) =0) and the user is trying to use the master to transfer data in High Speed mode. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter or Master-Receiver
	ABRT_SBYTE_NORSTRT   TX_ABRT_SOURCE = 0x01 << 9   //+ To clear Bit 9, the source of the ABRT_SBYTE_NORSTRT must be fixed first; restart must be enabled (CON[5]=1), the SPECIAL bit must be cleared (TAR[11]), or the GC_OR_START bit must be cleared (TAR[10]). Once the source of the ABRT_SBYTE_NORSTRT is fixed, then this bit can be cleared in the same manner as other bits in this register. If the source of the ABRT_SBYTE_NORSTRT is not fixed before attempting to clear this bit, bit 9 clears for one cycle and then gets reasserted. When this field is set to 1, the restart is disabled (RESTART_EN bit (CON[5]) =0) and the user is trying to send a START Byte. Reset value: 0x0 Role of DW_apb_i2c: Master
	ABRT_10B_RD_NORSTRT  TX_ABRT_SOURCE = 0x01 << 10  //+ This field indicates that the restart is disabled (RESTART_EN bit (CON[5]) =0) and the master sends a read command in 10-bit addressing mode. Reset value: 0x0 Role of DW_apb_i2c: Master-Receiver
	ABRT_MASTER_DIS      TX_ABRT_SOURCE = 0x01 << 11  //+ This field indicates that the User tries to initiate a Master operation with the Master mode disabled. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter or Master-Receiver
	ARB_LOST             TX_ABRT_SOURCE = 0x01 << 12  //+ This field specifies that the Master has lost arbitration, or if TX_ABRT_SOURCE[14] is also set, then the slave transmitter has lost arbitration. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter or Slave-Transmitter
	ABRT_SLVFLUSH_TXFIFO TX_ABRT_SOURCE = 0x01 << 13  //+ This field specifies that the Slave has received a read command and some data exists in the TX FIFO, so the slave issues a TX_ABRT interrupt to flush old data in TX FIFO. Reset value: 0x0 Role of DW_apb_i2c: Slave-Transmitter
	ABRT_SLV_ARBLOST     TX_ABRT_SOURCE = 0x01 << 14  //+ This field indicates that a Slave has lost the bus while transmitting data to a remote master. TX_ABRT_SOURCE[12] is set at the same time. Note: Even though the slave never 'owns' the bus, something could go wrong on the bus. This is a fail safe check. For instance, during a data transmission at the low-to-high transition of SCL, if what is on the data bus is not what is supposed to be transmitted, then DW_apb_i2c no longer own the bus. Reset value: 0x0 Role of DW_apb_i2c: Slave-Transmitter
	ABRT_SLVRD_INTX      TX_ABRT_SOURCE = 0x01 << 15  //+ 1: When the processor side responds to a slave mode request for data to be transmitted to a remote master and user writes a 1 in CMD (bit 8) of DATA_CMD register. Reset value: 0x0 Role of DW_apb_i2c: Slave-Transmitter
	ABRT_USER_ABRT       TX_ABRT_SOURCE = 0x01 << 16  //+ This is a master-mode-only bit. Master has detected the transfer abort (ENABLE[1]) Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter
	TX_FLUSH_CNT         TX_ABRT_SOURCE = 0x1FF << 23 //+ This field indicates the number of Tx FIFO Data Commands which are flushed due to TX_ABRT interrupt. It is cleared whenever I2C is disabled. Reset value: 0x0 Role of DW_apb_i2c: Master-Transmitter or Slave-Transmitter
)

const (
	ABRT_7B_ADDR_NOACKn   = 0
	ABRT_10ADDR1_NOACKn   = 1
	ABRT_10ADDR2_NOACKn   = 2
	ABRT_TXDATA_NOACKn    = 3
	ABRT_GCALL_NOACKn     = 4
	ABRT_GCALL_READn      = 5
	ABRT_HS_ACKDETn       = 6
	ABRT_SBYTE_ACKDETn    = 7
	ABRT_HS_NORSTRTn      = 8
	ABRT_SBYTE_NORSTRTn   = 9
	ABRT_10B_RD_NORSTRTn  = 10
	ABRT_MASTER_DISn      = 11
	ARB_LOSTn             = 12
	ABRT_SLVFLUSH_TXFIFOn = 13
	ABRT_SLV_ARBLOSTn     = 14
	ABRT_SLVRD_INTXn      = 15
	ABRT_USER_ABRTn       = 16
	TX_FLUSH_CNTn         = 23
)

type SLV_DATA_NACK_ONLY uint32

const (
	NACK SLV_DATA_NACK_ONLY = 0x01 << 0 //+ Generate NACK. This NACK generation only occurs when DW_apb_i2c is a slave-receiver. If this register is set to a value of 1, it can only generate a NACK after a data byte is received; hence, the data transfer is aborted and the data received is not pushed to the receive buffer. When the register is set to a value of 0, it generates NACK/ACK, depending on normal criteria. - 1: generate NACK after data byte received - 0: generate NACK/ACK normally Reset value: 0x0
)

const (
	NACKn = 0
)

type DMA_CR uint32

const (
	RDMAE DMA_CR = 0x01 << 0 //+ Receive DMA Enable. This bit enables/disables the receive FIFO DMA channel. Reset value: 0x0
	TDMAE DMA_CR = 0x01 << 1 //+ Transmit DMA Enable. This bit enables/disables the transmit FIFO DMA channel. Reset value: 0x0
)

const (
	RDMAEn = 0
	TDMAEn = 1
)

type ACK_GENERAL_CALL uint32

const (
	ACK_GEN_CALL ACK_GENERAL_CALL = 0x01 << 0 //+ ACK General Call. When set to 1, DW_apb_i2c responds with a ACK (by asserting ic_data_oe) when it receives a General Call. Otherwise, DW_apb_i2c responds with a NACK (by negating ic_data_oe).
)

const (
	ACK_GEN_CALLn = 0
)

type ENABLE_STATUS uint32

const (
	ENABLED                 ENABLE_STATUS = 0x01 << 0 //+ ic_en Status. This bit always reflects the value driven on the output port ic_en. - When read as 1, DW_apb_i2c is deemed to be in an enabled state. - When read as 0, DW_apb_i2c is deemed completely inactive. Note: The CPU can safely read this bit anytime. When this bit is read as 0, the CPU can safely read SLV_RX_DATA_LOST (bit 2) and SLV_DISABLED_WHILE_BUSY (bit 1). Reset value: 0x0
	SLV_DISABLED_WHILE_BUSY ENABLE_STATUS = 0x01 << 1 //+ Slave Disabled While Busy (Transmit, Receive). This bit indicates if a potential or active Slave operation has been aborted due to the setting bit 0 of the ENABLE register from 1 to 0. This bit is set when the CPU writes a 0 to the ENABLE register while: (a) DW_apb_i2c is receiving the address byte of the Slave-Transmitter operation from a remote master; OR, (b) address and data bytes of the Slave-Receiver operation from a remote master. When read as 1, DW_apb_i2c is deemed to have forced a NACK during any part of an I2C transfer, irrespective of whether the I2C address matches the slave address set in DW_apb_i2c (SAR register) OR if the transfer is completed before ENABLE is set to 0 but has not taken effect. Note: If the remote I2C master terminates the transfer with a STOP condition before the DW_apb_i2c has a chance to NACK a transfer, and ENABLE[0] has been set to 0, then this bit will also be set to 1. When read as 0, DW_apb_i2c is deemed to have been disabled when there is master activity, or when the I2C bus is idle. Note: The CPU can safely read this bit when EN (bit 0) is read as 0. Reset value: 0x0
	SLV_RX_DATA_LOST        ENABLE_STATUS = 0x01 << 2 //+ Slave Received Data Lost. This bit indicates if a Slave-Receiver operation has been aborted with at least one data byte received from an I2C transfer due to the setting bit 0 of ENABLE from 1 to 0. When read as 1, DW_apb_i2c is deemed to have been actively engaged in an aborted I2C transfer (with matching address) and the data phase of the I2C transfer has been entered, even though a data byte has been responded with a NACK. Note: If the remote I2C master terminates the transfer with a STOP condition before the DW_apb_i2c has a chance to NACK a transfer, and ENABLE[0] has been set to 0, then this bit is also set to 1. When read as 0, DW_apb_i2c is deemed to have been disabled without being actively involved in the data phase of a Slave-Receiver transfer. Note: The CPU can safely read this bit when EN (bit 0) is read as 0. Reset value: 0x0
)

const (
	ENABLEDn                 = 0
	SLV_DISABLED_WHILE_BUSYn = 1
	SLV_RX_DATA_LOSTn        = 2
)

type COMP_PARAM_1 uint32

const (
	APB_DATA_WIDTH     COMP_PARAM_1 = 0x03 << 0  //+ APB data bus width is 32 bits
	MAX_SPEED_MODE     COMP_PARAM_1 = 0x03 << 2  //+ MAX SPEED MODE = FAST MODE
	HC_COUNT_VALUES    COMP_PARAM_1 = 0x01 << 4  //+ Programmable count values for each mode.
	INTR_IO            COMP_PARAM_1 = 0x01 << 5  //+ COMBINED Interrupt outputs
	HAS_DMA            COMP_PARAM_1 = 0x01 << 6  //+ DMA handshaking signals are enabled
	ADD_ENCODED_PARAMS COMP_PARAM_1 = 0x01 << 7  //+ Encoded parameters not visible
	RX_BUFFER_DEPTH    COMP_PARAM_1 = 0xFF << 8  //+ RX Buffer Depth = 16
	TX_BUFFER_DEPTH    COMP_PARAM_1 = 0xFF << 16 //+ TX Buffer Depth = 16
)

const (
	APB_DATA_WIDTHn     = 0
	MAX_SPEED_MODEn     = 2
	HC_COUNT_VALUESn    = 4
	INTR_IOn            = 5
	HAS_DMAn            = 6
	ADD_ENCODED_PARAMSn = 7
	RX_BUFFER_DEPTHn    = 8
	TX_BUFFER_DEPTHn    = 16
)
