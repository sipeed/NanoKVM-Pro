export function isWebrtcSupported() {
  return 'RTCPeerConnection' in window;
}

export function isDecoderSupported() {
  return 'VideoDecoder' in window;
}

export function isH265Supported() {
  if ('getCapabilities' in RTCRtpReceiver) {
    const capabilities = RTCRtpReceiver.getCapabilities('video');
    if (capabilities && capabilities.codecs) {
      const h265Codec = capabilities.codecs.find((codec) => codec.mimeType === 'video/H265');
      return !!h265Codec;
    }
  }

  return false;
}

export function getSupportedVideoModes() {
  const webrtcSupported = isWebrtcSupported();
  const decoderSupported = isDecoderSupported();
  const h265Supported = isH265Supported();

  const videoModes = ['mjpeg'];

  if (webrtcSupported) {
    videoModes.push('h264-webrtc');

    if (h265Supported) {
      videoModes.push('h265-webrtc');
    }
  }

  if (decoderSupported) {
    videoModes.push('h264-direct');

    if (h265Supported) {
      videoModes.push('h265-direct');
    }
  }

  return videoModes;
}
