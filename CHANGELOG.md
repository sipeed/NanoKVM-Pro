## 1.2.11 (2026-01-17)

### Features

- Added USB network adapter support for sharing KVM network to the remote host
- Added menu bar display mode options (auto-hide/always visible)

### Bug Fixes

- Fixed a keyboard issue where IME composition events were not handled correctly
- Resolved an issue where the `AltGr` key was not recognized on the Windows system
- Fixed a bug where `Command` key combinations were not fully released on the macOS system

### Chores

- Bump react-router-dom from 6.30.1 to 6.30.3

## 1.2.10 (2026-01-09)

> This release includes a firmware update (v1.0.13). The device will automatically reboot once the update is complete.

### Features

- added support for mounting a virtual microphone on the remote host
- enabled bidirectional audio transmission (Speaker and Microphone) via WebRTC
- supported persistent storage for menu bar item configurations

### Bug Fixes

- fixed functionality issues with the AIC8800D80 Wi-Fi module
- resolved an issue where USB UAC2 failed to work on Windows systems
- fixed unresponsiveness when using two-finger scrolling on touchscreen

### Improvements

- optimized CMM memory usage
- enhanced menu bar interaction and performance

## 1.2.9 (2026-01-05)

### Bug Fixes

- fixed an issue preventing network configuration in WiFi AP mode

### Improvements

- expanded keyboard HID mapping to support more international characters
- added mouse HID support for "Forward" and "Back" buttons
- the menu bar now automatically hides when it loses focus

## 1.2.8 (2025-12-31)

> This release includes a firmware update (v1.0.12). The device will automatically reboot once the update is complete.

### Features

- added support for customizing FPS settings
- added support for user-defined keyboard shortcuts
- added support for SDIO Bluetooth

### Bug Fixes

- fixed the issue where a black screen appears when the remote host video frame rate is too low
- fixed an incorrect path for `axemac.sh` in `rc.local`
- fixed an issue where the resolution displayed on the LCD screen was not updated

### Improvements

- improved the logic for HDMI status detection
- optimized resolution detection for HDMI video
- optimized mouse wheel scrolling speed for better control
- optimized the menubar and settings UI
- updated the Russian translation (thanks to [@Tail870](https://github.com/Tail870))

### Chores

- added WireGuard support to the kernel
- upgraded antd to v5.29.3 and refactored code to replace deprecated APIs

## 1.2.7 (2025-12-18)

### Features

- added Low Power Mode under `Device > Advanced Settings`
- added support for configuring the mouse wheel scroll direction
- added a refresh button to manually update the virtual disk and image list in the `Mount Image` page

### Bug Fixes

- fixed an issue where the mouse wheel scroll direction was inverted
- resolved a bug where the serial terminal connection did not disconnect properly when closing the page

### Improvements

- optimized video scaling options
- improved compatibility with HTTP mode (thanks to [@0x656E](https://github.com/0x656E))
- updated the Russian translation (thanks to [@Tail870](https://github.com/Tail870))

## 1.2.6 (2025-12-11)

### Features

- add configuration support for video rate control modes and change the default from CBR (Constant Bit Rate) to VBR (Variable Bit Rate)
- add support for custom screen scaling ratios
- simplify the screen operations menu and introduce a new tab for advanced screen settings
- add support for mounting eMMC as a virtual disk

### Bug Fixes

- resolve an issue where closing the web terminal might spawn a zombie process (thanks to [@li20034](https://github.com/li20034))
- fix simultaneous key press conflicts and correct recognition errors for specific special characters

### Code Refactoring

- refactor the mouse module to improve response latency
- refactor the keyboard module to address various key mapping stability issues
- refactor the WebSocket module to transmit keyboard and mouse data using the standard **HID** format

### Improvements

- optimize animations for the LED screen interface
- reduce the overall installation package size
- remove the redundant password encryption step
- standardize function naming conventions within `kvmvision`

## 1.2.5 (2025-12-02)

### Bug Fixed

- Patched a potential security vulnerability
- Resolved an issue causing Wi-Fi connection failures

## 1.2.4 (2025-11-28)

### Features

- Add support for hardware revision 30126F-S

### Bug Fixes

- Fix resource leak issue caused by improper server shutdown

### Performance Improvements

- reduce power consumption by disabling the encoder when idle
- optimize Wi-Fi connection
- implement delay interval after failed login attempts

## 1.2.3 (2025-11-20)

### Features

- add support for uploading custom EDID
- enable static IP configuration via the Web UI
- add boot logos for PiKVM and flashing-mode

### Bug Fixes

- fix potential memory leaks
- ensure hostname and EDID are correctly reset during factory reset
- fix unresponsive right-side modifier keys
- fix issue where OLED sleep timeout failed to update in certain scenarios

### Performance Improvements

- optimize network throughput in Wi-Fi AP mode
- optimize reading and writing of LED UI configuration file

## 1.2.2 (2025-11-14)

### Bug Fixes

- added a mechanism to automatically reboot the device when a CMM resource leak is detected
- fixed a potential service crash that could occur when using MJPEG video mode
- resolved an issue where PiKVM UAC failed to enumerate on Windows system

### Performance Improvements

- optimized the boot sequence for a faster startup time

## 1.2.1 (2025-11-12)

### Features

- UserApp now includes an APP Hub for easy discovery and management of applications
- pre-installed Samba application in UserApp
- pre-installed ATX power control application in UserApp
- support for read-only image mounting

### Bug Fixes

- fixed a Wi-Fi disconnection bug
- resolved an issue where H.264 Direct video would stop when the web page was in the background
- fixed a bug where setting the web page title had no effect
- corrected wrong links on the settings page

## 1.2.0 (2025-11-07)

### Features

- display status and black screen on HDMI signal loss (H.264 WebRTC only)
- support resetting static IP settings on the LCD screen
- allow customizing USB information
- add browser compatibility warning for non-Chrome users
- support time format modification in the Web UI

### Bug Fixes

- prevent service status errors when `/bin/sh` is modified
- correct vertical mouse drift when remote desktop is scaled
- resolve issue where keyboard modifiers occasionally fail
- optimize UserApp UI and fix logic bugs
- ensure non-ASCII WiFi automatically reconnects after reboot
- remove duplicate scroll bar in the web terminal

### Performance Improvements

- optimize H.264 encoding performance
- rename the EDID on the web UI
- remove MJPEG frame difference detection

## 1.1.9 (2025-10-31)

### Features

- add support for 3840x2400@30FPS (16:10) and 3840x1600@50FPS (24:10) resolutions
- the UserApp now supports touch knob events and included a corresponding demo
- the touchscreen UI can now remember Wi-Fi passwords for easier connectivity

### Bug Fixes

- resolved an issue where kvm_vin would sometimes fail to retrieve HDMI data
- resolved an issue preventing 3440x1440 resolution from working on certain systems
- ensured that custom port for the reverse proxy are now correctly applied
- fixed the incorrect HDMI passthrough state, and the signal disappears after clicking
- fixed a bug that prevented exiting the UserApp with the knob unless the app was selected
- fixed a bug causing the new UserApp demo to be copied repeatedly
- fixed a bug where image download failures did not have a corresponding notification

## 1.1.8 (2025-10-29)

### Features

- add support for long screen resolutions, including 3440x1440@60Hz and 2560x1080@75Hz
- the device now securely saves the password for connected Wi-Fi networks

### Bug Fixes

- Tailscale would fail to start automatically after the device was rebooted
- H.264 WebRTC incompatible with the Safari browser
- certain screen resolutions would cause the display to repeatedly reset
- fix compatibility with special characters in Wi-Fi password
- correct the link to the CHANGELOG file

### Performance Improvements

- disable Wi-Fi scanning while in AP mode
- remove touchscreen operations on phone
- implement a retry mechanism for installation failure

## 1.1.7 (2025-10-24)

### Features

- add support for touchscreen
- add support for Wi-Fi scanning and connect to open SSID
- add a Tomato Clock user application demo to showcase

### Bug Fixes

- HDMI Passthrough state was being displayed incorrectly
- ensured compatibility with special characters in Wi-Fi SSID
- fix the PiKVM password restore issue

## 1.1.6 (2025-10-22)

### Features

- add a switch to toggle between 12-hour and 24-hour time formats
- add support for HID reset
- add script for configuration of a static IP address
- display the device's time and allows for manual synchronization in web UI

### Bug Fixes

- incorrect backlight display on the LCD screen
- host cannot be woken up after entering sleep mode

### Performance Improvements

- add support for deleting images in web UI
- the default GOP has been set to 50
- the default audio bitrate has been increased from 32kbps to 128kbps

## 1.1.5 (2025-10-17)

### Features

- users can now create and apply their own custom UI applications
- add NanoKVM Batch Manage for easier management of multiple devices
- the LCD screen now has an automatic sleep function
- add support for serial port terminal

### Bug Fixes

- the mDNS button would occasionally disappear
- the currently selected EDID is now displayed correctly

### Performance Improvements

- optimize HID Only mode for better compatibility with BIOS
- application upgrade experience has been streamlined and improved
- to prevent video stream loss, only one video mode can be enabled at a time

## 1.1.4 (2025-10-10)

- feat: update new experimental AI Agent! (Computer Use Agent)
- feat: add customize auxiliary screen functions
- perf: optimize WebRTC latency, now 1080P~4K have the same latency
