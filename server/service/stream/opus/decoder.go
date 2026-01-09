package opus

/*
#cgo LDFLAGS: -lopus -lm
#include <opus/opus.h>
#include <stdlib.h>
#include <string.h>

static OpusDecoder* g_opus_decoder = NULL;

// Initialize Opus decoder
// Returns 0 on success, negative error code on failure
int opus_init_decoder(int sample_rate, int channels) {
    int error;
    g_opus_decoder = opus_decoder_create(sample_rate, channels, &error);
    if (error != OPUS_OK) {
        g_opus_decoder = NULL;
        return error;
    }
    return 0;
}

// Decode Opus frame to PCM
// Returns number of samples decoded per channel, or negative error code
int opus_do_decode(const unsigned char* data, int len, short* pcm, int max_samples) {
    if (g_opus_decoder == NULL) {
        return -1;
    }
    return opus_decode(g_opus_decoder, data, len, pcm, max_samples, 0);
}

// Cleanup decoder
void opus_cleanup_decoder() {
    if (g_opus_decoder != NULL) {
        opus_decoder_destroy(g_opus_decoder);
        g_opus_decoder = NULL;
    }
}

// Get error string
const char* opus_err_string(int error) {
    return opus_strerror(error);
}
*/
import "C"

import (
	"unsafe"
)

const (
	SampleRate   = 48000
	Channels     = 2
	MaxFrameSize = 5760
)

type Decoder struct {
	inited    bool
	pcmBuffer []int16
}

func NewDecoder() *Decoder {
	return &Decoder{
		pcmBuffer: make([]int16, MaxFrameSize*Channels),
	}
}

func (d *Decoder) Init() error {
	if d.inited {
		return nil
	}

	ret := C.opus_init_decoder(C.int(SampleRate), C.int(Channels))
	if ret < 0 {
		return &OpusError{Code: int(ret), Message: C.GoString(C.opus_err_string(ret))}
	}
	d.inited = true
	return nil
}

// Decode decodes Opus data to PCM
func (d *Decoder) Decode(opusData []byte) ([]byte, error) {
	if !d.inited || len(opusData) == 0 {
		return nil, nil
	}

	samples := C.opus_do_decode(
		(*C.uchar)(unsafe.Pointer(&opusData[0])),
		C.int(len(opusData)),
		(*C.short)(unsafe.Pointer(&d.pcmBuffer[0])),
		C.int(MaxFrameSize),
	)

	if samples < 0 {
		return nil, &OpusError{Code: int(samples), Message: C.GoString(C.opus_err_string(samples))}
	}

	bytesToReturn := int(samples) * Channels * 2
	pcmBytes := (*[MaxFrameSize * Channels * 2]byte)(unsafe.Pointer(&d.pcmBuffer[0]))[:bytesToReturn]

	return pcmBytes, nil
}

func (d *Decoder) Cleanup() {
	if d.inited {
		C.opus_cleanup_decoder()
		d.inited = false
	}
}

func (d *Decoder) IsInited() bool {
	return d.inited
}

type OpusError struct {
	Code    int
	Message string
}

func (e *OpusError) Error() string {
	return e.Message
}
