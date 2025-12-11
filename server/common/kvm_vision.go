package common

/*
	#cgo CFLAGS: -I../include
	#cgo LDFLAGS: -L../dl_lib -lkvm
	#include "kvm_vision.h"
*/
import "C"
import (
	"strings"
	"sync"
	"unsafe"

	"NanoKVM-Server/config"

	log "github.com/sirupsen/logrus"
)

type KvmVision struct {
	StreamType uint8
}

const (
	RATE_CONTROL_CBR uint8 = iota
	RATE_CONTROL_VBR
)

const (
	IMG_MJPEG_TYPE uint8 = iota
	IMG_H264_TYPE_SPS
	IMG_H264_TYPE_PPS
	IMG_H264_TYPE_IF
	IMG_H264_TYPE_PF
	IMG_H265_TYPE_SPS
	IMG_H265_TYPE_PPS
	IMG_H265_TYPE_IF
	IMG_H265_TYPE_PF
)

const (
	STREAM_TYPE_MJPEG = iota
	STREAM_TYPE_H264_WEBRTC
	STREAM_TYPE_H264_DIRECT
	STREAM_TYPE_H265_WEBRTC
	STREAM_TYPE_H265_DIRECT
)

var (
	kvmVision     *KvmVision
	kvmVisionOnce sync.Once
)

func GetKvmVision() *KvmVision {
	kvmVisionOnce.Do(func() {
		kvmVision = &KvmVision{
			StreamType: STREAM_TYPE_H264_WEBRTC,
		}

		conf := config.GetInstance()
		logLevel := strings.ToLower(conf.Logger.Level)

		logEnable := C.uint8_t(0)
		if logLevel == "debug" {
			logEnable = C.uint8_t(1)
		}

		C.kvmv_init(logEnable)
		log.Debugf("kvm vision initialized")
	})

	return kvmVision
}

func (k *KvmVision) SetStreamType(streamType uint8) {
	k.StreamType = streamType
}

func (k *KvmVision) SetRateControl(mode uint8) int {
	if mode != RATE_CONTROL_CBR && mode != RATE_CONTROL_VBR {
		log.Debugf("invalid rate control mode: %d", mode)
		return -1
	}

	result := int(C.kvmv_set_rate_control(
		C.uint8_t(mode),
	))
	if result < 0 {
		log.Debugf("failed to set rate control mode: %d", result)
		return result
	}

	return result
}

func (k *KvmVision) ReadMjpeg(width uint16, height uint16, quality uint16) (data []byte, result int) {
	var (
		kvmData  *C.uint8_t
		dataSize C.uint32_t
	)

	result = int(C.kvmv_read_img(
		C.uint16_t(width),
		C.uint16_t(height),
		C.uint8_t(IMG_MJPEG_TYPE),
		C.uint16_t(quality),
		&kvmData,
		&dataSize,
	))
	if result < 0 {
		log.Errorf("failed to read MJPEG: %v", result)
		return
	}
	if kvmData == nil || dataSize == 0 {
		return
	}

	data = C.GoBytes(unsafe.Pointer(kvmData), C.int(dataSize))
	return
}

func (k *KvmVision) ReadH264(width uint16, height uint16, bitRate uint16) (data []byte, result int) {
	var (
		kvmData  *C.uint8_t
		dataSize C.uint32_t
	)

	result = int(C.kvmv_read_img(
		C.uint16_t(width),
		C.uint16_t(height),
		C.uint8_t(IMG_H264_TYPE_SPS),
		C.uint16_t(bitRate),
		&kvmData,
		&dataSize,
	))
	if result < 0 {
		log.Errorf("failed to read H.264: %d", result)
		return
	}
	if kvmData == nil || dataSize == 0 {
		return
	}

	data = C.GoBytes(unsafe.Pointer(kvmData), C.int(dataSize))
	return
}

func (k *KvmVision) ReadH265(width uint16, height uint16, bitRate uint16) (data []byte, result int) {
	var (
		kvmData  *C.uint8_t
		dataSize C.uint32_t
	)

	result = int(C.kvmv_read_img(
		C.uint16_t(width),
		C.uint16_t(height),
		C.uint8_t(IMG_H265_TYPE_SPS),
		C.uint16_t(bitRate),
		&kvmData,
		&dataSize,
	))
	if result < 0 {
		log.Errorf("failed to read H.265: %d", result)
		return
	}
	if kvmData == nil || dataSize == 0 {
		return
	}

	data = C.GoBytes(unsafe.Pointer(kvmData), C.int(dataSize))
	return
}

func (k *KvmVision) ReadAudio() (data []byte, result int) {
	var (
		kvmData  *C.uint8_t
		dataSize C.uint32_t
	)

	result = int(C.kvmv_read_audio(&kvmData, &dataSize))
	if result < 0 {
		log.Errorf("failed to read audio: %d", result)
		return
	}
	if kvmData == nil || dataSize == 0 {
		return
	}

	data = C.GoBytes(unsafe.Pointer(kvmData), C.int(dataSize))
	return
}

func (k *KvmVision) GetFps() int {
	return int(C.kvmv_get_fps())
}

func (k *KvmVision) SetHDMI(enable bool) int {
	hdmiEnable := C.uint8_t(0)
	if enable {
		hdmiEnable = C.uint8_t(1)
	}

	result := int(C.kvmv_hdmi_control(hdmiEnable))
	if result < 0 {
		log.Errorf("failed to set hdmi to %t", enable)
		return result
	}

	return result
}

func (k *KvmVision) SetGop(gop uint8) {
	_gop := C.uint8_t(gop)
	C.kvmv_set_gop(_gop)
}

func (k *KvmVision) Close() {
	C.kvmv_deinit()
	log.Debugf("stop kvm vision...")
}
