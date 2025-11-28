
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

void kvmv_init([[__maybe_unused__]] uint8_t _debug_info_en) { __throw_error(); }
int kvmv_read_img([[__maybe_unused__]] uint16_t _width,
                  [[__maybe_unused__]] uint16_t _height,
                  [[__maybe_unused__]] uint8_t _type,
                  [[__maybe_unused__]] uint16_t _qlty,
                  [[__maybe_unused__]] uint8_t **_pp_kvm_data,
                  [[__maybe_unused__]] uint32_t *_p_kvmv_data_size) {
  return __throw_error();
}
int kvmv_get_sps_frame([[__maybe_unused__]] uint8_t **_pp_kvm_data,
                       [[__maybe_unused__]] uint32_t *_p_kvmv_data_size) {
  return __throw_error();
}
int kvmv_get_pps_frame([[__maybe_unused__]] uint8_t **_pp_kvm_data,
                       [[__maybe_unused__]] uint32_t *_p_kvmv_data_size) {
  return __throw_error();
}
int free_kvmv_data([[__maybe_unused__]] uint8_t **_pp_kvm_data) {
  return __throw_error();
}
void free_all_kvmv_data() { __throw_error(); }
void set_gop([[__maybe_unused__]] uint8_t _gop) { __throw_error(); }
int get_fps() {
  return __throw_error();
}
void kvmv_deinit() { __throw_error(); }
uint8_t kvmv_hdmi_control([[__maybe_unused__]] uint8_t _en) {
  return __throw_error();
}

int kvma_read_audio(uint8_t **_pp_kvm_data, uint32_t *_p_kvma_data_size) {
  return __throw_error();
}

#ifdef __cplusplus
}
#endif
