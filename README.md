# hidraw macro pad

A simple macro pad using the hidraw interface on Linux. This repo will contain all code for both the firmware and the host-side software.
Currently a few tests using evdev have been conducted; after getting familiar with evdev, the plan is to move to using hidraw for more complex and reliable functionality.

The final HID device will contain some standard keys (next track, previous track, etc.). hidraw will be used to power the little display + rotary encoder. Potentially other features as well.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
