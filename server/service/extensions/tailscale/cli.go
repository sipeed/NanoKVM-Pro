package tailscale

import (
	"NanoKVM-Server/utils"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Cli struct{}

const (
	serveTargetLocalhost = "localhost:443"
	serveTargetLoopback  = "127.0.0.1:443"
	serveBaseCommand     = "tailscale serve --bg --https=443 https+insecure://localhost:443"
)

type TsStatus struct {
	BackendState string `json:"BackendState"`

	Self struct {
		HostName     string   `json:"HostName"`
		TailscaleIPs []string `json:"TailscaleIPs"`
	} `json:"Self"`

	CurrentTailnet struct {
		Name string `json:"Name"`
	} `json:"CurrentTailnet"`
}

func NewCli() *Cli {
	return &Cli{}
}

func (c *Cli) Start() error {
	return utils.StartService(Tailscaled, true)
}

func (c *Cli) Restart() error {
	return utils.RestartService(Tailscaled)
}

func (c *Cli) Stop() error {
	return utils.StopService(Tailscaled, true)
}

func (c *Cli) Up() error {
	command := "tailscale up --accept-dns=false"
	return exec.Command("sh", "-c", command).Run()
}

func (c *Cli) Down() error {
	command := "tailscale down"
	return exec.Command("sh", "-c", command).Run()
}

func (c *Cli) Status() (*TsStatus, error) {
	command := "tailscale status --json"
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// output is not in standard json format
	if outputStr := string(output); !strings.HasPrefix(outputStr, "{") {
		index := strings.Index(outputStr, "{")
		if index == -1 {
			return nil, errors.New("unknown output")
		}

		output = []byte(outputStr[index:])
	}

	var status TsStatus
	err = json.Unmarshal(output, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Cli) Login() (string, error) {
	command := "tailscale login --accept-dns=false --timeout=10m"
	cmd := exec.Command("sh", "-c", command)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = stderr.Close()
	}()

	go func() {
		_ = cmd.Run()
	}()

	reader := bufio.NewReader(stderr)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		if strings.Contains(line, "https") {
			reg := regexp.MustCompile(`\s+`)
			url := reg.ReplaceAllString(line, "")
			return url, nil
		}
	}
}

func (c *Cli) Logout() error {
	command := "tailscale logout"
	return exec.Command("sh", "-c", command).Run()
}

func (c *Cli) ServeStatus() (bool, string, error) {
	command := "tailscale serve status"
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, "", err
	}

	statusOutput := string(output)
	isForwardingToLocal := strings.Contains(statusOutput, serveTargetLocalhost) || strings.Contains(statusOutput, serveTargetLoopback)
	serveUrl := ""

	if isForwardingToLocal {
		// Matches the HTTPS serve URL (e.g. https://hostname.tail12345.ts.net)
		serveUrlPattern := regexp.MustCompile(`https://[^\s]+\.ts\.net`)
		matches := serveUrlPattern.FindStringSubmatch(statusOutput)
		if len(matches) > 0 {
			serveUrl = matches[0]
		}
	}

	return isForwardingToLocal, serveUrl, nil
}

func (c *Cli) Serve(enable bool) (string, error) {
	command := serveBaseCommand
	if !enable {
		command += " off"
	}

	cmd := exec.Command("sh", "-c", command)

	// Merge stdout into stderr so we capture everything in one stream.
	// stderr is where interactive prompts/info usually go, but just in case.
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	cmd.Stdout = cmd.Stderr

	defer func() {
		_ = stderr.Close()
	}()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	// if disable, wait for command to finish
	if !enable {
		return "", cmd.Wait()
	}

	// if enable, check if need auth
	reader := bufio.NewReader(stderr)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// EOF reached — check whether the process actually succeeded.
			if waitErr := cmd.Wait(); waitErr != nil {
				log.Errorf("tailscale serve command failed: %v", waitErr)
				return "", fmt.Errorf("tailscale serve failed: %w", waitErr)
			}
			// Exit code 0 — command succeeded without needing auth.
			return "", nil
		}

		if strings.Contains(line, "https://login.tailscale.com/f/serve") {
			authUrlPattern := regexp.MustCompile(`(https://login\.tailscale\.com/f/serve[^\s]+)`)
			matches := authUrlPattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				authUrl := matches[1]
				// Fire-and-forget: we've consumed the data we need; let cmd clean up asynchronously
				go cmd.Wait()
				return authUrl, nil
			}
		}
	}
}
