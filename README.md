# Circus

## Sparkfun Bluetooth board based on NRF52840.
Overall design is to use a skull (fake) to imitate a ouija board where the user must select the correct order of items to unlock the door and escape the room.

GPIO's are used to indicate when items are in place, RFID's allow the user to select each item in the proper order to open the door.

Requires a micro controller such as the Raspberry PI to monitor the GPIO's and enable the Skull.  The sparkfun board sits in the skull and is powered by Lithium ION battery
so it can operate wirelessly.

A number of items are left to be done - with this replacing the original C++ code and LAMP setup for interacting with the door device.
