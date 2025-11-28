package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"
)

type WiFiDetails struct {
	BSSID          string `json:"bssid"`
	SSID           string `json:"ssid"`
	ID             string `json:"id"`
	Mode           string `json:"mode"`
	WiFiGeneration string `json:"wifi_generation"`
	PairwiseCipher string `json:"pairwise_cipher"`
	GroupCipher    string `json:"group_cipher"`
	KeyMgmt        string `json:"key_mgmt"`
	WPAState       string `json:"wpa_state"`
	IPAddress      string `json:"ip_address"`
	P2PDeviceAddr  string `json:"p2p_device_address"`
	Address        string `json:"address"`
	UUID           string `json:"uuid"`
	Freq           string `json:"freq"`
}

func GetWifiDetails() (*WiFiDetails, error) {
	cmd := exec.Command("wpa_cli", "-i", "wlan0", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute wpa_cli status: %v", err)
	}

	details := &WiFiDetails{}
	lines := bytes.Split(output, []byte("\n"))

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) == 2 {
			key := string(bytes.TrimSpace(parts[0]))
			value := sanitizeString(bytes.TrimSpace(parts[1]))

			switch key {
			case "bssid":
				details.BSSID = value
			case "ssid":
				details.SSID = value
			case "id":
				details.ID = value
			case "mode":
				details.Mode = value
			case "wifi_generation":
				details.WiFiGeneration = value
			case "pairwise_cipher":
				details.PairwiseCipher = value
			case "group_cipher":
				details.GroupCipher = value
			case "key_mgmt":
				details.KeyMgmt = value
			case "wpa_state":
				details.WPAState = value
			case "ip_address":
				details.IPAddress = value
			case "p2p_device_address":
				details.P2PDeviceAddr = value
			case "address":
				details.Address = value
			case "uuid":
				details.UUID = value
			case "freq":
				details.Freq = value
			}
		}
	}

	return details, nil
}

func GetConfiguredNetworks() ([]string, error) {
	cmd := exec.Command("wpa_cli", "-i", "wlan0", "list_networks")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute wpa_cli: %v, output: %s", err, string(output))
	}

	lines := bytes.Split(output, []byte("\n"))
	var networks []string

	for i, line := range lines {
		if i == 0 || len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		fields := bytes.Fields(line)
		if len(fields) >= 2 {
			ssid := string(fields[1])
			networks = append(networks, ssid)
		}
	}

	return networks, nil
}

func sanitizeString(data []byte) string {
	decoded := decodeHexEscapes(data)

	if utf8.Valid(decoded) {
		return string(decoded)
	}

	var buf strings.Builder
	buf.Grow(len(decoded))

	for len(decoded) > 0 {
		r, size := utf8.DecodeRune(decoded)
		if r == utf8.RuneError && size == 1 {
			decoded = decoded[1:]
			continue
		}
		buf.WriteRune(r)
		decoded = decoded[size:]
	}

	return buf.String()
}

func decodeHexEscapes(data []byte) []byte {
	str := string(data)
	if !strings.Contains(str, "\\") {
		return data
	}

	var result []byte
	i := 0
	for i < len(str) {
		if str[i] == '\\' {
			if i+4 < len(str) && str[i+1] == '\\' && str[i+2] == 'x' {
				hex := str[i+3 : i+5]
				var b byte
				if n, err := fmt.Sscanf(hex, "%02x", &b); err == nil && n == 1 {
					result = append(result, b)
					i += 5
					continue
				}
			} else if i+3 < len(str) && str[i+1] == 'x' {
				hex := str[i+2 : i+4]
				var b byte
				if n, err := fmt.Sscanf(hex, "%02x", &b); err == nil && n == 1 {
					result = append(result, b)
					i += 4
					continue
				}
			} else if i+1 < len(str) && str[i+1] == '\\' {
				result = append(result, '\\')
				i += 2
				continue
			}
		}
		result = append(result, str[i])
		i++
	}

	return result
}
