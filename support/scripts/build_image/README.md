# NanoKVM-Pro Image Build Tool

[中文文档](README_zh.md)

## Overview

This directory contains script tools for converting the base AXP image compiled from [`maix_ax620e_sdk`](https://github.com/sipeed/maix_ax620e_sdk) into a complete image that can be directly flashed to NanoKVM-Pro devices.

## Features

This script provides a complete image packaging process that can add NanoKVM-Pro APP and configure required system components starting from the maix_ax620e_sdk base image, generating a final image file ready for flashing.

## Prerequisites

- Successfully compiled maix_ax620e_sdk and obtained the base AXP image
- Linux build environment (Ubuntu 22.04 recommended)
- Sufficient disk space for image processing

### Install Dependencies

```bash
sudo apt update

sudo apt install -y android-sdk-libsparse-utils qemu-user-static

pip3 install axp-tools
```

## Quick Start

1. Build the maix_ax620e_sdk base AXP image. For detailed steps, please refer to the [official documentation](https://github.com/sipeed/maix_ax620e_sdk/blob/main/README.md) of maix_ax620e_sdk.

2. Copy the generated base AXP image to this directory, for example with the filename `image.axp`.

3. Download the latest software package from the Release page and extract it:

    ```bash
    wget https://github.com/sipeed/NanoKVM-Pro/releases/download/x.x.x/nanokvm_pro_x.x.x.tar.gz

    tar -xvf nanokvm_pro_x.x.x.tar.gz

    # After extraction, you will get the nanokvm_pro_x.x.x directory containing .deb files
    ```

4. Execute the command:

   ```bash
   ./build_image.py image.axp --app nanokvm_pro_x.x.x -o nano_kvm_pro_image.axp
   ```

   This script will process the base image and generate the final NanoKVM-Pro image file, outputting as `nano_kvm_pro_image.axp`.

   Convert the axp image to img format:

   ```bash
   axp2img -i nano_kvm_pro_image.axp
   ```

## Output

The final generated image can be directly written to the NanoKVM-Pro device using flashing tools.

For specific burning methods, please refer to: [https://wiki.sipeed.com/hardware/en/kvm/NanoKVM_Pro/faq.html#Image-Burning-Methods](https://wiki.sipeed.com/hardware/en/kvm/NanoKVM_Pro/faq.html#Image-Burning-Methods)
