// DO NOT EDIT THIS FILE. GENERATED BY svdxgen.

//go:build rp2350

// Package bootram provides access to the registers of the BOOTRAM peripheral.
//
// Instances:
//
//	BOOTRAM  BOOTRAM_BASE  -  -  Additional registers mapped adjacent to the bootram, for use by the bootrom.
//
// Registers:
//
//	0x800 32  WRITE_ONCE0    This registers always ORs writes into its current contents. Once a bit is set, it can only be cleared by a reset.
//	0x804 32  WRITE_ONCE1    This registers always ORs writes into its current contents. Once a bit is set, it can only be cleared by a reset.
//	0x808 32  BOOTLOCK_STAT  Bootlock status register. 1=unclaimed, 0=claimed. These locks function identically to the SIO spinlocks, but are reserved for bootrom use.
//	0x80C 32  BOOTLOCK0      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//	0x810 32  BOOTLOCK1      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//	0x814 32  BOOTLOCK2      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//	0x818 32  BOOTLOCK3      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//	0x81C 32  BOOTLOCK4      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//	0x820 32  BOOTLOCK5      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//	0x824 32  BOOTLOCK6      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//	0x828 32  BOOTLOCK7      Read to claim and check. Write to unclaim. The value returned on successful claim is 1 << n, and on failed claim is zero.
//
// Import:
//
//	github.com/embeddedgo/pico/p/mmap
package bootram