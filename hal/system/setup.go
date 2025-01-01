// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package system

/*
  pico-sdk initialization sequence

  pico_crt0/crt0.S   _entry_point
  newlib_interface.c   runtime_init
  runtime.c              runtime_run_initializers
  runtime.c                runtime_run_initializers_from
  runtime_init.c             runtime_init_bootrom_reset
  runtime_init.c             runtime_init_early_resets
  runtime_init.c             runtime_init_usb_power_down
  runtime_init_clocks.c      runtime_init_clocks
  runtime_init.c             runtime_init_post_clock_resets
  boot_lock.c                runtime_init_boot_locks_reset
  runtime_init.c             runtime_init_spin_locks_reset
  bootrom_lock.c             runtime_init_bootrom_locking_enable
  mutex.c                    runtime_init_mutex
  runtime_init.c             runtime_init_install_ram_vector_table
  time.c                     runtime_init_default_alarm_pool
  runtime.c                  first_per_core_initializer
  runtime_init.c             runtime_init_per_core_bootrom_reset
  runtime_init.c             runtime_init_per_core_enable_coprocessors
  sync_spin_lock.c           spinlock_set_extexclall
  irq.c                      runtime_init_per_core_irq_priorities
*/

// SetupPico2 initializes the whole system assuming it is RPI Pico 2
// compatible: 12 MHz XOSC, 133 MHz QSPI Flash (66 MHz Read Data).
func SetupPico2() {
}
