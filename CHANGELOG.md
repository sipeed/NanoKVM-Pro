## 1.2.4 (2025-11-28)

### Features

- Add support for hardware revision 30126F-S

### Bug Fixes

- Fix resource leak issue caused by improper server shutdown

### Performance Improvements

- Reduce power consumption by disabling the encoder when idle
- Optimize Wi-Fi connection
- Implement delay interval after failed login attempts

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
