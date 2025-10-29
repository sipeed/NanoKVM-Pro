# NanoKVM-Pro

> **Note:** This project is currently in development and will be officially open-sourced next month (November).
>
> We are excited to share our work with the community soon! Stay tuned for the official release.

## Introduction

NanoKVM-Pro is the continuation of NanoKVM, inheriting the extreme compactness and powerful expandability of the NanoKVM series as an IP-KVM product.
It has made a significant leap in performance, making it more suitable for remote working scenarios.

To meet different user needs, NanoKVM-Pro offers two forms: NanoKVM-Desk and NanoKVM-ATX:

![NanoKVM-Pro Desktop and ATX versions side by side](https://wiki.sipeed.com/hardware/assets/NanoKVM/pro/introduce/combine.png)

- **NanoKVM-Desk** is the desktop version of NanoKVM-Pro, featuring an anodized matte metal shell. The front panel has a 1.47-inch touchscreen that displays core KVM information and allows for easy hardware function settings or can be used as a mini secondary screen, providing a more tactile user experience with the left-side infinite knob.

- **NanoKVM-ATX** is the internal version of NanoKVM-Pro, equipped with half-height/full-height brackets for installation inside a case. It allows for easier installation for host users with built-in USB cables and power control interfaces. Remote control can be achieved via external HDMI, network, and USB connections.

NanoKVM-Pro uses the AX630 as its main control core, featuring an ARM 1.2G dual-core A53 CPU. The external 1GB LPDDR4 memory provides strong computing support for remote desktop connections. It has built-in HDMI loop-out and capture chips, offering up to 4K60FPS HDMI loop-out and 4K45FPS video capture. Thanks to AX630's efficient and powerful image processing architecture, NanoKVM-Pro can transmit high-resolution images with very low latency, with typical delays as low as 60ms at 2K resolution.

## Specifications

| Product       | NanoKVM-Pro | NanoKVM      | GxxKVM      | JxxKVM      |
|---------------|----------|--------------|-------------|-------------|
| Main Control  | AX630C   | SG2002       | RV1126      | RV1106      |
| Core          | 2xA53@1.2G | 1xC906@1.0G  | 4xA7@1.5G   | 1xA7@1.2G   |
| Memory        | 1G LPDDR4X | 256M DDR3    | 1G DDR3     | 256M DDR3   |
| Storage       | 32G eMMC | 32G microSD  | 8G eMMC     | 16G eMMC    |
| System        | NanoKVM+PIKVM | NanoKVM      | GxxKVM      | JxxKVM      |
| Resolution    | 4K@45fps | 1080P@60fps | 4K@30fps, 2K@60fps | 1080P@60fps |
| HDMI Loop-Out | 4K Loop-Out | ×            | ×           | ×           |
| Video Encoding | MJPG/H264 | MJPG/H264    | MJPG/H264   | MJPG/H264   |
| Audio Transmission | ✓        | ×            | ✓           | ×           |
| UEFI/BIOS Support | ✓        | ✓            | ✓           | ✓           |
| Simulated USB Keyboard/Mouse | ✓ | ✓          | ✓           | ✓           |
| Simulated USB ISO | ✓        | ✓            | ✓           | ✓           |
| IPMI          | ✓        | ✓            | ✓           | ×           |
| Wake-on-LAN (WOL) | ✓        | ✓            | ✓           | ✓           |
| WebSSH        | ✓        | ✓            | ✓           | ✓           |
| Custom Scripts | ✓        | ✓            | ×           | ×           |
| Serial Terminal | 2 Channels | 2 Channels   | None        | 1 Channel   |
| Storage Performance | 32G eMMC 300MB/s | 32G MicroSD 12MB/s | 8G eMMC 120MB/s | 8G eMMC 60MB/s |
| Ethernet      | 1000M    | 100M         | 1000M       | 100M        |
| Internal Form Factor | Optional ATX version | Optional PCIe version | ×           | ×           |
| WiFi          | Optional WiFi6 | Optional WiFi6 | ×           | ×           |
| MicroSD Expansion | ✓        | ×            | ×           | ×           |
| ATX Power Control | ✓        | ✓            | +15$        | +10$        |
| Display       | 1.47-inch 320x172 LCD<br>0.96-inch 128x64 OLED | 0.96-inch 128x64 OLED | None | 1.66-inch 280x240 |
| Additional Features | Synchronized LED effects, Smart Assistant | –        | –           | –           |
| Power Consumption | 0.6A@5V  | 0.2A@5V      | 0.4A@5V     | 0.2A@5V     |
| Power Input   | USB-C/PoE | USB-C/PoE/PCIe | USB-C       | USB-C       |
| Dimensions     | 65x65x28mm | 40x36x36mm   | 80x60x7.5mm | 60x6x24-30mm |
