package opus

import (
	"io"
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
)

type AudioInputPlayer struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	mutex   sync.Mutex
	active  bool
	decoder *Decoder
}

var (
	audioInputPlayer     *AudioInputPlayer
	audioInputPlayerOnce sync.Once
)

func GetAudioInputPlayer() *AudioInputPlayer {
	audioInputPlayerOnce.Do(func() {
		audioInputPlayer = &AudioInputPlayer{
			decoder: NewDecoder(),
		}
	})
	return audioInputPlayer
}

// Start initializes the Opus decoder and aplay process
func (p *AudioInputPlayer) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.active {
		return nil
	}

	if err := p.decoder.Init(); err != nil {
		log.Errorf("Opus decoder init failed: %s", err)
		return err
	}

	// Start aplay process - receives raw PCM from decoded Opus
	p.cmd = exec.Command("aplay",
		"-D", "plughw:UAC2Gadget,0",
		"-r", "48000",
		"-f", "S16_LE",
		"-c", "2",
		"-t", "raw",
		"--buffer-size=4096",
	)

	var err error
	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		log.Errorf("failed to get stdin pipe: %s", err)
		p.decoder.Cleanup()
		return err
	}

	if err := p.cmd.Start(); err != nil {
		log.Errorf("failed to start aplay: %s", err)
		p.decoder.Cleanup()
		return err
	}

	p.active = true
	log.Debug("audio input player started")

	go func() {
		if err := p.cmd.Wait(); err != nil {
			log.Debugf("aplay process exited: %s", err)
		}
		p.mutex.Lock()
		p.active = false
		p.mutex.Unlock()
	}()

	return nil
}

// DecodeAndWrite decodes Opus data and writes PCM to aplay
func (p *AudioInputPlayer) DecodeAndWrite(opusData []byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.active || !p.decoder.IsInited() || p.stdin == nil {
		return nil
	}

	pcmBytes, err := p.decoder.Decode(opusData)
	if err != nil {
		log.Debugf("Opus decode failed: %s", err)
		return nil
	}

	if pcmBytes == nil {
		return nil
	}

	_, err = p.stdin.Write(pcmBytes)
	return err
}

// Stop terminates the audio pipeline
func (p *AudioInputPlayer) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.stdin != nil {
		_ = p.stdin.Close()
		p.stdin = nil
	}

	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		p.cmd = nil
	}

	p.decoder.Cleanup()
	p.active = false
	log.Debug("audio input player stopped")
}

// IsActive returns whether the player is currently active
func (p *AudioInputPlayer) IsActive() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.active
}
