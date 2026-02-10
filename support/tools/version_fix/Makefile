CROSS_COMPILE = aarch64-none-linux-gnu-

CFLAGS += -I/usr/local/ARM-toolchain/gcc-arm-9.2-2019.12-x86_64-aarch64-none-linux-gnu/aarch64-none-linux-gnu/libc/usr/include

CC = $(CROSS_COMPILE)gcc

SRC = version_fix.c

TARGET = version_fix

all: $(TARGET)

$(TARGET): $(SRC)
	# $(CC) -march=rv64gc -o $@ $^
	$(CC) -o $@ $^

clean:
	rm -f $(TARGET)

.PHONY: all clean