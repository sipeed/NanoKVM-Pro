# Video IPC transport probe

This separate executable runs the real Server stream handlers and controls on
`127.0.0.1:18080`, without changing the installed Server or its authentication.
Build `./test/videoipc` for the board with the normal ARM64 CGO toolchain and
Opus dependency. At runtime, supply the matching real libkvm, not the link stub:

```sh
KVM_VIDEO_SOCKET_DIR=/run/kvm-stage4 LD_LIBRARY_PATH=/tmp/kvm-stage4-validation ./videoipc
```

Stop old video hardware consumers before starting a temporary
`kvm_vin --video-ipc --socket-dir /run/kvm-stage4`, and restore those services
after testing. Do not run two hardware owners at once.

The probe exposes the existing `/api/stream/*` paths, `/probe/state`, and
`/probe/hold?ms=2200` (bounded to 0–3000 ms). The hold endpoint subscribes but
pauses consumption to exercise backpressure through IPC into Axera.
`/probe.html` optionally serves a local browser test page from the working directory.

From the host, tunnel the loopback port and receive H.264 with PLI/FIR feedback:

```sh
go run ./test/webrtcreceiver -output received.h264
ffmpeg -v error -f h264 -i received.h264 -f null -
```

The receiver defaults to 90 frames and emits JSON statistics. It sends 20 PLI
requests after frame 15 and one FIR after frame 30. Use a long GOP to distinguish
requested keyframes from periodic ones. See the sibling nanokvmpro repository's
`docs/video-pipeline-stage4/` for validation results and scope.
