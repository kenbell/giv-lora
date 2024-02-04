# Hardware

## OEM

The Module appears to be a GivEnergy branded variant of the [Aurtron LoRa Serial Converter CC32LR](http://www.aurtron.com/en/h-col-262.html).

The manual reveals that:
1. Firmware upgrades may be supported by `LoRa OTA` and `RS-485`
2. 'Maintenance' is possible using 'AT commands'

### Owens Brothers Variant

In addition to the GivEnergy variant, another variant appears to be the Owens Brothers [OBM-LoRa-223](https://owen-brothers.com/obm-lora-223-long-range-wireless-transceiver.html).  The GivEnergy variant appears to be less configurable than the Owens Brothers.  Screen captures in the manual for the Owens Brothers module does reveal some interesting info for helping to understand the radio protocol:

1. The default radio address is likely `0x1001`
2. There is a 2-byte 'Network ID'
3. The default LoRa Spreading Factor is 7
4. The default data rate is CR4/5
5. The default bandwidth is 125kHz

This corresponds closely with the observed behaviour of the GivEnergy modules.

### BMS Parts

A variant also appears to be available from [BMS Parts](https://bmsparts.co.uk/shop-maxking-limited/lora/lora-wireless-rs485-transmitter-and-receiver-2-units-lcd-screen/).  The documentation and videos reveal that some variants expose the module over modbus itself, with and the default address set to 0.  This presumably disables modbus rather than listening on address 0 which is reserved for broadcast only.

## PCBs

Internally the module consists of two PCBs.  One PCB contains a powersupply and RS485 interfacing logic.  The other 'logic' PCB seems to work exclusively at 3.3V and contains the combined wireless and MCU module from Lanvee.

The logic PCB has a debug port, with this pinout:

  1. GND
  2. +3.3V
  3. SWCLK
  4. SWDIO
  5. RST (? not verified)

## Chip
Module: Lanvee s34s
Link: http://www.lanvee.com/h-col-131.html

    $ st-info --probe --hot-plug
    Found 1 stlink programmers
    version:    V2J29S7
    serial:     XXXXXXXXXXXXXXXXXXXXXXXX
    flash:      131072 (pagesize: 128)
    sram:       20480
    chipid:     0x0447
    descr:      L0xx Category 5

