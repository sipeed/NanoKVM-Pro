
#include <fcntl.h> /* low-level i/o */
#include <pthread.h>
#include <stdexcept>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#ifdef __cplusplus
extern "C" {
#endif

// #ifdef KVMIMPL_THROW
// #pragma message("Throw Impl")
// #else
// #pragma message("Empty Impl")
// #endif

static int __throw_error() {
#ifdef KVMIMPL_THROW
  throw std::runtime_error(
      "[Error] This shared library is exclusively for linking with GO "
      "applications. "
      "Please compile the implementation from support/libkvm directory.");
  return -1;
#else
  printf("[Error] This shared library is exclusively for linking with GO "
         "applications. "
         "Please compile the implementation from support/libkvm directory.");
  exit(-1);
  return -1;
#endif
}

void kvmv_init(uint8_t _debug_info_en) {
  __throw_error();
}

void kvmv_free_all_data() {
  __throw_error();
}

void kvmv_deinit() {
  __throw_error();
}

int kvmv_read_img(uint16_t _width,
                  uint16_t _height,
                  uint8_t _type,
                  uint16_t _qlty,
                  uint8_t** _pp_kvm_data,
                  uint32_t* _p_kvmv_data_size) {
  return __throw_error();
}

int kvmv_get_sps_frame(uint8_t** _pp_kvm_data, uint32_t* _p_kvmv_data_size) {
  return __throw_error();
}

int kvmv_get_pps_frame(uint8_t** _pp_kvm_data, uint32_t* _p_kvmv_data_size) {
  return __throw_error();
}

int kvmv_free_data(uint8_t** _pp_kvm_data) {
  return __throw_error();
}

int kvmv_set_fps(uint8_t _fps) {
  return __throw_error();
}

int kvmv_get_fps(void) {
  return __throw_error();
}

int kvmv_set_gop(uint8_t _gop) {
  return __throw_error();
}

int kvmv_hdmi_control(uint8_t _en) {
  return __throw_error();
}

int kvmv_read_audio(uint8_t** _pp_kvm_data, uint32_t* _p_kvmv_data_size) {
  return __throw_error();
}

int kvmv_set_rate_control(uint8_t mode) {
  return __throw_error();
}

#ifdef __cplusplus
}
#endif