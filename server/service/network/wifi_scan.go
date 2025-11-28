package network

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	singleHexEscapeRegex = regexp.MustCompile(`"ssid":"([^"]*\\x[0-9a-fA-F]{2}[^"]*)"`)
)

func (s *Service) ScanWifi(c *gin.Context) {
	var rsp proto.Response

	if isAPMode() {
		rsp.ErrRsp(c, -1, "AP mode is enabled")
		return
	}

	out, err := exec.Command(WiFiScript, "try_scan").CombinedOutput()
	if err != nil {
		rsp.ErrRsp(c, -2, "scan wifi failed")
		log.Errorf("scan wifi failed, %s", err)
		return
	}
	data := fixInvalidJSONEscapes(out)

	var wifiList []*proto.WiFi
	if err := json.Unmarshal(data, &wifiList); err != nil {
		rsp.ErrRsp(c, -3, "scan wifi failed")
		log.Errorf("parse wifi failed: %s", err)
		return
	}

	wifiList = removeConnected(wifiList)

	rsp.OkRspWithData(c, &proto.ScanWifiRsp{
		WifiList: wifiList,
	})
	log.Debugf("scan wifi successfully, total: %d", len(wifiList))
}

func removeConnected(wifiList []*proto.WiFi) []*proto.WiFi {
	if len(wifiList) == 0 {
		return wifiList
	}

	details, err := utils.GetWifiDetails()
	if err != nil || details.WPAState != "COMPLETED" || details.IPAddress == "" {
		return wifiList
	}

	for i, wifi := range wifiList {
		if wifi.BSSID == details.BSSID {
			return append(wifiList[:i], wifiList[i+1:]...)
		}
	}

	return wifiList
}

func fixInvalidJSONEscapes(data []byte) []byte {
	return singleHexEscapeRegex.ReplaceAllFunc(data, func(match []byte) []byte {
		parts := bytes.SplitN(match, []byte(`":"`), 2)
		if len(parts) != 2 {
			log.Warnf("failed to parse ssid field from: %s", string(match))
			return match
		}

		fieldName := string(bytes.Trim(parts[0], `"`))
		value := bytes.TrimSuffix(parts[1], []byte(`"`))

		var result []byte
		for i := 0; i < len(value); {
			if value[i] == '\\' && i+1 < len(value) {
				switch value[i+1] {
				case 'x':
					if i+3 < len(value) {
						if b, err := strconv.ParseUint(string(value[i+2:i+4]), 16, 8); err == nil {
							result = append(result, byte(b))
							i += 4
							continue
						}
					}
				case '\\':
					result = append(result, '\\')
					i += 2
					continue
				case '"':
					result = append(result, '"')
					i += 2
					continue
				}
			}
			result = append(result, value[i])
			i++
		}

		cleanedStr := sanitizeString(result)

		return []byte(fmt.Sprintf(`"%s":"%s"`, fieldName, cleanedStr))
	})
}

func sanitizeString(data []byte) string {
	if utf8.Valid(data) {
		str := string(data)
		str = removeControlCharacters(str)
		str = escapeJSONString(str)
		return str
	}

	str := strings.ToValidUTF8(string(data), "�")
	str = removeControlCharacters(str)
	str = escapeJSONString(str)

	log.Debugf("converted invalid UTF-8 sequence to: %s", str)
	return str
}

func removeControlCharacters(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		if r >= 32 && r != 127 || r == '\t' || r == '\n' || r == '\r' {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

func escapeJSONString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
