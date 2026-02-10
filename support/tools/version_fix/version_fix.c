#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>

#include <linux/i2c-dev.h>
#include <linux/i2c.h>
#include <sys/ioctl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/time.h>
#include <time.h>

#include "version_fix.h"

// =======================================================================================================

int client;
// old offset set to 0xff first
static uint8_t old_offset = 0xff;

// I2C write function
int i2c_write_byte(uint8_t offset, uint8_t reg, uint8_t data) 
{
    // reg buffer
    uint8_t reg_buf[2] = {0};

    // if offset is changed, write it first
    if (offset != old_offset) {
        old_offset = offset;
        reg_buf[0] = LT6911_REG_OFFSET; // Set the offset register
        reg_buf[1] = offset;    // Set the register to read
        if (write(client, reg_buf, 2) != 2) {
            perror("Failed to write offset to the i2c bus");
            return -1;
        }
    }

    // write the data to the i2c bus
    reg_buf[0] = reg; // Set the register to write
    reg_buf[1] = data; // Set the data to write
    if (write(client, reg_buf, 2) != 2) {
        perror("Failed to write to the i2c bus");
        return -1;
    }
    return 0;
}

int i2c_write_bytes(uint8_t offset, uint8_t reg, const uint8_t *data, size_t len) 
{
    // check len
    if (len == 0) {
        fprintf(stderr, "Data length must be greater than 0.\n");
        return -1;
    }

    // reg buffer
    uint8_t reg_buf[1 + len];

    // if offset is changed, write it first
    if (offset != old_offset) {
        old_offset = offset;
        reg_buf[0] = LT6911_REG_OFFSET; // Set the offset register
        reg_buf[1] = offset;    // Set the register to read
        if (write(client, reg_buf, 2) != 2) {
            perror("Failed to write offset to the i2c bus");
            return -1;
        }
    }

    // write the data to the i2c bus
    reg_buf[0] = reg;
    for (size_t i = 0; i < len; i++) {
        reg_buf[i + 1] = data[i];
    }

    // write to the I2C bus
    if (write(client, reg_buf, 1 + len) != 1 + len) {
        perror("Failed to write to the i2c bus");
        return -1;
    }

    return 0;
}

// I2C read function
int i2c_read_byte(uint8_t offset, uint8_t reg, uint8_t *data) 
{
    // reg buffer
    uint8_t reg_buf[2] = {0};

    // if offset is changed, write it first
    if (offset != old_offset) {
        old_offset = offset;
        reg_buf[0] = LT6911_REG_OFFSET; // Set the offset register
        reg_buf[1] = offset;    // Set the register to read
        if (write(client, reg_buf, 2) != 2) {
            perror("Failed to write offset to the i2c bus");
            return -1;
        }
    }

    // write the register address to read
    reg_buf[0] = reg; // Set the register to read
    if (write(client, reg_buf, 1) != 1) {
        perror("Failed to write register address to the i2c bus");
        return -1;
    }

    // read the data from the i2c bus
    if (read(client, data, 1) != 1 ) {
        perror("Failed to read from the i2c bus");
        return -1;
    }

    return 0;
}

int i2c_read_bytes(uint8_t offset, uint8_t reg, uint8_t *data, size_t len)
{
    // Check data length
    if (len == 0) {
        fprintf(stderr, "Data length must be greater than 0.\n");
        return -1;
    }

    uint8_t reg_buf[2]; // Create register buffer

    // If offset has changed, write the offset first
    if (offset != old_offset) {
        old_offset = offset;
        reg_buf[0] = LT6911_REG_OFFSET; // Set offset register
        reg_buf[1] = offset; // Set offset value
        if (write(client, reg_buf, 2) != 2) {
            perror("Failed to write offset to the i2c bus");
            return -1;
        }
    }

    // Write the register address to read
    reg_buf[0] = reg; // Set the register to read
    if (write(client, reg_buf, 1) != 1) {
        perror("Failed to write register address to the i2c bus");
        return -1;
    }

    // Read data from I2C bus
    if (read(client, data, len) != len) {
        perror("Failed to read from the i2c bus");
        return -1;
    }

    return 0;
}

// =======================================================================================================

int lt6911_enable(void) {
    // Enable the LT6911UXC by writing to the appropriate register
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0xEE, 0x01) != 0) {
        fprintf(stderr, "Failed to enable LT6911UXC\n");
        return -1;
    }
    return 0;
}

int lt6911_disable(void) {
    // Disable the LT6911UXC by writing to the appropriate register
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0xEE, 0x00) != 0) {
        fprintf(stderr, "Failed to disable LT6911UXC\n");
        return -1;
    }
    return 0;
}

int lt6911uxc_version_write(uint8_t *version_data, uint16_t version_size)
{
    uint8_t i;
    int ret;
    uint8_t chip_data[16] = {0};
    uint8_t wr_count = version_size / LT6911UXC_WR_SIZE;
    uint8_t version_str[32] = {0};

    fprintf(stdout, "Writing version data....\n");
    
    // Start
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0xFF, 0x80) != 0) return -1;
    lt6911_enable();
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5E, 0xDF) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x58, 0x00) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x59, 0x51) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x10) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x00) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x58, 0x21) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0xFF, 0x80) != 0) return -1;
    lt6911_enable();
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x84) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5B, 0x01) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5C, 0x80) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5D, 0x00) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x81) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
    // Waiting for erasure
    usleep(500000);

    if (i2c_read_byte(LT6911_SYS3_OFFSET, 0x08, chip_data) != 0) return -1;
    if (chip_data[0] != 0xEE) {
        fprintf(stderr, "Unsupported chip version:%x\n", chip_data[0]);
        return -1;
    } 
    if (i2c_write_byte(LT6911_SYS3_OFFSET, 0x08, 0xAE) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS3_OFFSET, 0x08, 0xEE) != 0) return -1;
    // Write
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0xFF, 0x80) != 0) return -1;
    lt6911_enable();
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x84) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x84) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
    for (i = 0; i < wr_count; i++) {
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5E, 0xDF) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x20) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x00) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x58, 0x21) != 0) return -1;
        // Write data block
        if (i2c_write_bytes(LT6911_SYS_OFFSET, 0x59, version_data+(LT6911UXC_WR_SIZE*i), LT6911UXC_WR_SIZE) != 0) return -1;
        // Write Address
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5B, 0x01) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5C, 0x80+(i/(0xff/LT6911UXC_WR_SIZE))) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5D, 0x00+(i*LT6911UXC_WR_SIZE)) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5E, 0xC0) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x90) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
        // Finalize write
        if (i != wr_count-1) {
            if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x84) != 0) return -1;
        } else {
            if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x88) != 0) return -1;
        }
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
    }
    if (i2c_read_byte(LT6911_SYS3_OFFSET, 0x08, chip_data) != 0) return -1;
    if (chip_data[0] != 0xEE) return -1;
    if (i2c_write_byte(LT6911_SYS3_OFFSET, 0x08, 0xAE) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS3_OFFSET, 0x08, 0xEE) != 0) return -1;

    fprintf(stdout, "Version write completed\n");

    return 0;
}

int lt6911uxc_version_read(uint8_t *version_data, uint16_t version_size)
{
    uint8_t i;
    int ret;
    uint8_t chip_data[16] = {0};
    uint8_t wr_count = version_size / LT6911UXC_WR_SIZE;

    // Read version data from LT6911UXC
    fprintf(stdout, "Reading version data...\n");
    // Read
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0xFF, 0x80) != 0) return -1;
    lt6911_enable();
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x84) != 0) return -1;
    if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
    for (i = 0; i < wr_count; i++) {
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5E, 0x5F) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0xA0) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5B, 0x01) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5C, 0x80+(i/(0xff/LT6911UXC_WR_SIZE))) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5D, 0x00+(i*LT6911UXC_WR_SIZE)) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x90) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x5A, 0x80) != 0) return -1;
        if (i2c_write_byte(LT6911_SYS_OFFSET, 0x58, 0x21) != 0) return -1;
        if (i2c_read_bytes(LT6911_SYS_OFFSET, 0x5F, version_data+(i*LT6911UXC_WR_SIZE), LT6911UXC_WR_SIZE) != 0) return -1;
    }

    return 0;
}

// =======================================================================================================

uint8_t version_check(uint8_t *version_data)
{
    uint8_t version_char[32] = {0};     // from version_data[256] to version_data[288]
    memcpy(version_char, version_data + 256, 32);
    printf("Current version code: %s\n", version_char);
    if (strstr((char*)version_char, "(Desk-") != NULL || strstr((char*)version_char, "(ATX-") != NULL) {
        printf("The current version information is valid.\n");
        return 0;
    } else {
        fprintf(stderr, "Version string does not contain valid prefix.\n");
        return -1;
    }
}

void generate_kvm_string(int type, char* buffer) {
    // Initialize random seed (if not already set)
    static int seeded = 0;
    if (!seeded) {
        srand((unsigned int)time(NULL));
        seeded = 1;
    }
    
    // Generate a random number between 0-65535
    uint16_t random_value = (uint16_t)(rand() % 65536);
    
    // Select device type based on type parameter
    const char* device_type = (type == 1) ? "ATX" : "Desk";
    
    // Format the string
    snprintf(buffer, 32, "NanoKVM_Pro (%s-X) NxaL0%04X\n", device_type, random_value);
}

void fix_data_with_kvm_string(uint8_t *version_data, const char* kvm_string) {
    // Write kvm_string to the specified position in version_data
    memcpy(version_data + 256, kvm_string, 32);
}

void print_help(const char* program_name) {
    printf("============================================================================\n");
    printf("                      DEVICE FIRMWARE UPDATE UTILITY                        \n");
    printf("============================================================================\n\n");
    
    printf("IMPORTANT PREREQUISITES:\n");
    printf("──────────────────────────\n");
    printf("• Ensure the kvmcomm.service is stopped\n");
    printf("• Ensure the lt6911_manage module is unloaded\n");
    printf("\n");
    printf("Execute the following commands BEFORE running this program:\n");
    printf("    systemctl stop kvmcomm.service && rmmod lt6911_manage\n");
    printf("\n");
    printf("After completion, restart the service with:\n");
    printf("    systemctl start kvmcomm.service\n");
    printf("\n");
    
    printf("USAGE:\n");
    printf("──────\n");
    printf("    %s [OPTIONS]\n\n", program_name);
    
    printf("OPTIONS:\n");
    printf("────────\n");
    printf("    -f, --force          Force modification mode\n");
    printf("                         Ignores existing valid information in current version\n");
    printf("\n");
    printf("    -t, --type <value>   Specify device type\n");
    printf("                         Valid values: 'Desk' or 'ATX'\n");
    printf("\n");
    printf("    -h, --help           Display this help message\n");
    printf("\n");
    
    printf("NOTES:\n");
    printf("──────\n");
    printf("• Running without any arguments will perform a non-forced write of ATX version\n");
    printf("• Use caution when using force mode as it may overwrite existing data\n");
    printf("\n");
    
    printf("EXAMPLES:\n");
    printf("─────────\n");
    printf("    %s -t Desk          # Write Desk version (non-forced)\n", program_name);
    printf("    %s --type ATX       # Write ATX version (non-forced)\n", program_name);
    printf("    %s -f -t Desk       # Force write Desk version\n", program_name);
    printf("    %s --force --type ATX  # Force write ATX version\n", program_name);
    printf("\n");
    
    printf("============================================================================\n");
    printf("                          WARNING: USE WITH CAUTION                         \n");
    printf("============================================================================\n");
}

// =======================================================================================================

int main(int argc, char* argv[]) 
{
    // root@kvm-72d6:~# systemctl stop kvmcomm.service 
    // root@kvm-72d6:~# rmmod lt6911_manage

    int force_flag = 0;      // Whether to enable force mode
    char* type_value = NULL; // Type value
    int show_help = 0;       // Whether to display help
    
    // Parse command-line arguments
    for (int i = 1; i < argc; i++) {
        // --force or -f option
        if (strcmp(argv[i], "-f") == 0 || strcmp(argv[i], "--force") == 0) {
            force_flag = 1;
        }
        // --type or -t option
        else if (strcmp(argv[i], "-t") == 0 || strcmp(argv[i], "--type") == 0) {
            // Check if there is a next argument
            if (i + 1 < argc) {
                type_value = argv[i + 1];
                i++; // Skip the next argument
            } else {
                printf("Error: --type option requires an argument\n");
                print_help(argv[0]);
                return 1;
            }
        }
        // --help or -h option
        else if (strcmp(argv[i], "-h") == 0 || strcmp(argv[i], "--help") == 0) {
            show_help = 1;
        }
    }
    
    // Display help information
    if (show_help) {
        print_help(argv[0]);
        return 0;
    }
    
    // Process -f/--force option
    if (force_flag) {
        printf("Force mode enabled\n");
    }
    
    // Process -t/--type option
    if (type_value != NULL) {
        printf("Type value is: %s\n", type_value);
    }


    uint8_t version_data[VERSION_BUFFER_SIZE] = {0};
    uint8_t version_read[VERSION_BUFFER_SIZE] = {0};

    // open i2c device
    if ((client = open(I2C_DEVICE, O_RDWR)) < 0) {
        perror("Failed to open the i2c bus");
        return -1;
    }

    // set the lt6911uxc address
    if (ioctl(client, I2C_SLAVE, I2C_ADDRESS) < 0) {
        perror("Failed to acquire bus access and/or talk to slave");
        close(client);
        return -1;
    }

    // read version data from chip
    if (lt6911uxc_version_read(version_data, VERSION_BUFFER_SIZE) != 0) {
        fprintf(stderr, "Failed to read version data from LT6911UXC\n");
        close(client);
        return -1;
    }

    if (version_check(version_data) == 0) {
        // Version is valid, no need to update
        if (force_flag) {
            fprintf(stdout, "Version data is valid, forced write.\n");
        } else {
            fprintf(stdout, "Version data is valid, no update needed.\n");
            close(client);
            return EXIT_SUCCESS;
        }
    }

    // generate version string
    char kvm_string[32] = {0};
    uint8_t type = 1; // 1 for ATX (default), 0 for Desk
    if ( type_value != NULL ) {
        if ( strcmp(type_value, "Desk") == 0 ) {
            type = 0;
        } else if ( strcmp(type_value, "ATX") == 0 ) {
            type = 1;
        } else {
            fprintf(stderr, "Invalid argument. Use 'desk' or 'atx'. Defaulting to 'atx'.\n");
            type = 1;
        }
    }
    generate_kvm_string(type, kvm_string);  // type 1 for ATX, 0 for Desk
    printf("Generated new KVM string: %s\n", kvm_string);
    fix_data_with_kvm_string(version_data, kvm_string);

    // version write
    if (lt6911uxc_version_write(version_data, VERSION_BUFFER_SIZE) != 0) {
        fprintf(stderr, "Failed to write version data to LT6911UXC\n");
        close(client);
        return -1;
    }
    // sleep 1s
    sleep(1);

    // // read version data back
    if (lt6911uxc_version_read(version_read, VERSION_BUFFER_SIZE) != 0) {
        fprintf(stderr, "Failed to read version data from LT6911UXC\n");
        close(client);
        return -1;
    }

    // check if the read data matches the written data
    if (memcmp(version_data, version_read, VERSION_BUFFER_SIZE) != 0) {
        fprintf(stderr, "Version data mismatch after write/read cycle\n");
        // print version data with 16*16 hex format
        fprintf(stderr, "Written version data:\n");
        for (int i = 0; i < VERSION_BUFFER_SIZE; i++) {
            fprintf(stderr, "%02X ", version_data[i]);
            if ((i + 1) % 16 == 0) {
                fprintf(stderr, "\n");
            }
        }
        fprintf(stderr, "\nRead version data:\n");
        for (int i = 0; i < VERSION_BUFFER_SIZE; i++) {
            fprintf(stderr, "%02X ", version_read[i]);
            if ((i + 1) % 16 == 0) {
                fprintf(stderr, "\n");
            }
        }
        close(client);
        return -1;
    } else {
        fprintf(stdout, "Version data verified successfully\n");
    }

    // close client
    close(client);

    return EXIT_SUCCESS;
}
