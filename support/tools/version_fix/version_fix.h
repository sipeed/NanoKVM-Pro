#ifndef VERSION_FIX_H
#define VERSION_FIX_H

#define I2C_DEVICE              "/dev/i2c-0"    // I2C Device path (adjust based on your system)

#define I2C_ADDRESS             0x2b            // I2C device address (adjust based on actual device)
#define LT6911_REG_OFFSET       0xFF            // LT6911UXC register offset address
#define LT6911_SYS_OFFSET       0x80            // LT6911UXC register offset address
#define LT6911_SYS2_OFFSET      0x90            // LT6911UXC register offset address
#define LT6911_SYS3_OFFSET      0x81            // LT6911UXC register offset address
#define LT6911_SYS4_OFFSET      0xA0            // LT6911UXC register offset address
#define LT6911_CSI_INFO_OFFSET  0x85            // LT6911UXC CSI interface information register offset address
#define LT6911_HDMI_INFO_OFFSET 0x86            // LT6911UXC HDMI information register offset address
#define LT6911_CSI_TOTAL_OFFSET 0xD4            // LT6911UXC CSI bus statistics information
#define LT6911_AUDIO_INFO_OFFSET 0xB0            // LT6911UXC audio information register offset address
#define LT6911C_HDMI_INFO_OFFSET 0xD2            // LT6911C HDMI information register offset address
#define LT6911C_AUDIO_INFO_OFFSET 0xD1            // LT6911C audio information register offset address
#define LT6911C_CSI_INFO_OFFSET 0xC2            // LT6911C CSI information register offset address
#define LT6911D_MANAGE_OFFSET   0xE0            // LT6911D management register offset address
#define LT6911D_DATA_OFFSET     0xE1            // LT6911D data register offset address

#define VERSION_BUFFER_SIZE     288             // Maximum supported bytes
#define LT6911UXC_WR_SIZE       32              // LT6911UXC maximum bytes for single read/write
#define LT6911D_WR_SIZE         32              // LT6911D maximum bytes for single read/write
#define LT6911C_WR_SIZE         16              // LT6911C maximum bytes for single read/write

#endif  // VERSION_FIX_H