package ai

import "fmt"

type typeChar struct {
	modifiers int
	code      int
}

const maxKeyboardKeys = 6

var keyUpReport = []byte{0, 0, 0, 0, 0, 0, 0, 0}

var modifierMap = map[string]byte{
	"ControlLeft":  1 << 0,
	"ShiftLeft":    1 << 1,
	"AltLeft":      1 << 2,
	"MetaLeft":     1 << 3,
	"ControlRight": 1 << 4,
	"ShiftRight":   1 << 5,
	"AltRight":     1 << 6,
	"MetaRight":    1 << 7,
}

var keyCodeMap = map[string]byte{
	"KeyA": 0x04, "KeyB": 0x05, "KeyC": 0x06, "KeyD": 0x07, "KeyE": 0x08,
	"KeyF": 0x09, "KeyG": 0x0a, "KeyH": 0x0b, "KeyI": 0x0c, "KeyJ": 0x0d,
	"KeyK": 0x0e, "KeyL": 0x0f, "KeyM": 0x10, "KeyN": 0x11, "KeyO": 0x12,
	"KeyP": 0x13, "KeyQ": 0x14, "KeyR": 0x15, "KeyS": 0x16, "KeyT": 0x17,
	"KeyU": 0x18, "KeyV": 0x19, "KeyW": 0x1a, "KeyX": 0x1b, "KeyY": 0x1c,
	"KeyZ": 0x1d,

	"Digit1": 0x1e, "Digit2": 0x1f, "Digit3": 0x20, "Digit4": 0x21, "Digit5": 0x22,
	"Digit6": 0x23, "Digit7": 0x24, "Digit8": 0x25, "Digit9": 0x26, "Digit0": 0x27,

	"Enter": 0x28, "Escape": 0x29, "Backspace": 0x2a, "Tab": 0x2b, "Space": 0x2c,
	"Minus": 0x2d, "Equal": 0x2e, "BracketLeft": 0x2f, "BracketRight": 0x30,
	"Backslash": 0x31, "IntlHash": 0x32, "Semicolon": 0x33, "Quote": 0x34,
	"Backquote": 0x35, "Comma": 0x36, "Period": 0x37, "Slash": 0x38,
	"CapsLock": 0x39,

	"F1": 0x3a, "F2": 0x3b, "F3": 0x3c, "F4": 0x3d, "F5": 0x3e, "F6": 0x3f,
	"F7": 0x40, "F8": 0x41, "F9": 0x42, "F10": 0x43, "F11": 0x44,
	"F12": 0x45,

	"PrintScreen": 0x46, "ScrollLock": 0x47, "Pause": 0x48, "Insert": 0x49,
	"Home": 0x4a, "PageUp": 0x4b, "Delete": 0x4c, "End": 0x4d,
	"PageDown": 0x4e, "ArrowRight": 0x4f, "ArrowLeft": 0x50,
	"ArrowDown": 0x51, "ArrowUp": 0x52,

	"NumLock": 0x53, "NumpadDivide": 0x54, "NumpadMultiply": 0x55,
	"NumpadSubtract": 0x56, "NumpadAdd": 0x57, "NumpadEnter": 0x58,
	"Numpad1": 0x59, "Numpad2": 0x5a, "Numpad3": 0x5b, "Numpad4": 0x5c,
	"Numpad5": 0x5d, "Numpad6": 0x5e, "Numpad7": 0x5f, "Numpad8": 0x60,
	"Numpad9": 0x61, "Numpad0": 0x62, "NumpadDecimal": 0x63,
	"IntlBackslash": 0x64, "ContextMenu": 0x65, "Power": 0x66,
	"NumpadEqual": 0x67,

	"F13": 0x68, "F14": 0x69, "F15": 0x6a, "F16": 0x6b, "F17": 0x6c,
	"F18": 0x6d, "F19": 0x6e, "F20": 0x6f, "F21": 0x70, "F22": 0x71,
	"F23": 0x72, "F24": 0x73,

	"Execute": 0x74, "Help": 0x75, "Props": 0x76, "Select": 0x77,
	"Stop": 0x78, "Again": 0x79, "Undo": 0x7a, "Cut": 0x7b,
	"Copy": 0x7c, "Paste": 0x7d, "Find": 0x7e,

	"AudioVolumeMute": 0x7f, "AudioVolumeUp": 0x80, "AudioVolumeDown": 0x81,
	"VolumeMute": 0x7f, "VolumeUp": 0x80, "VolumeDown": 0x81,
	"LockingCapsLock": 0x82, "LockingNumLock": 0x83, "LockingScrollLock": 0x84,
	"NumpadComma": 0x85, "NumpadEqual2": 0x86, "IntlRo": 0x87,
	"KanaMode": 0x88, "IntlYen": 0x89, "Convert": 0x8a, "NonConvert": 0x8b,
	"International6": 0x8c, "International7": 0x8d, "International8": 0x8e,
	"International9": 0x8f, "Lang1": 0x90, "Lang2": 0x91, "Lang3": 0x92,
	"Lang4": 0x93, "Lang5": 0x94, "Lang6": 0x95, "Lang7": 0x96,
	"Lang8": 0x97, "Lang9": 0x98,

	"NumpadParenLeft": 0xb6, "NumpadParenRight": 0xb7,
	"NumpadBackspace": 0xbb, "NumpadMemoryStore": 0xd0,
	"NumpadMemoryRecall": 0xd1, "NumpadMemoryClear": 0xd2,
	"NumpadMemoryAdd": 0xd3, "NumpadMemorySubtract": 0xd4,
	"NumpadClear": 0xd8, "NumpadClearEntry": 0xd9,

	"MediaPlayPause": 0xe8, "MediaStop": 0xe9, "MediaTrackPrevious": 0xea,
	"MediaTrackNext": 0xeb, "Eject": 0xec, "MediaSelect": 0xed,
	"LaunchMail": 0xee, "LaunchApp1": 0xef, "BrowserSearch": 0xf0,
	"LaunchApp2": 0xf0, "BrowserHome": 0xf1, "BrowserBack": 0xf2,
	"BrowserForward": 0xf3, "BrowserStop": 0xf4, "BrowserRefresh": 0xf5,
	"BrowserFavorites": 0xf6, "Sleep": 0xf8, "Wake": 0xf9,
	"MediaRewind": 0xfa, "MediaFastForward": 0xfb,
}

func buildKeyComboReport(keys []string) ([]byte, error) {
	report := make([]byte, 8)
	seenKeys := make(map[byte]struct{})
	keyIndex := 2

	for _, key := range keys {
		if modifier, ok := modifierMap[key]; ok {
			report[0] |= modifier
			continue
		}

		code, ok := keyCodeMap[key]
		if !ok {
			return nil, fmt.Errorf("unknown key: %s", key)
		}
		if _, ok := seenKeys[code]; ok {
			continue
		}
		if keyIndex >= 2+maxKeyboardKeys {
			return nil, fmt.Errorf("too many non-modifier keys")
		}
		seenKeys[code] = struct{}{}
		report[keyIndex] = code
		keyIndex++
	}

	return report, nil
}

func buildTypeReports(text string) ([][]byte, []rune) {
	reports := make([][]byte, 0, len(text)*2)
	skipped := make([]rune, 0)

	for _, char := range text {
		key, ok := typeCharMap[char]
		if !ok {
			skipped = append(skipped, char)
			continue
		}
		reports = append(reports,
			[]byte{byte(key.modifiers), 0x00, byte(key.code), 0x00, 0x00, 0x00, 0x00, 0x00},
			keyUpReport,
		)
	}

	return reports, skipped
}

var typeCharMap = map[rune]typeChar{
	'a': {0, 4}, 'b': {0, 5}, 'c': {0, 6}, 'd': {0, 7}, 'e': {0, 8},
	'f': {0, 9}, 'g': {0, 10}, 'h': {0, 11}, 'i': {0, 12}, 'j': {0, 13},
	'k': {0, 14}, 'l': {0, 15}, 'm': {0, 16}, 'n': {0, 17}, 'o': {0, 18},
	'p': {0, 19}, 'q': {0, 20}, 'r': {0, 21}, 's': {0, 22}, 't': {0, 23},
	'u': {0, 24}, 'v': {0, 25}, 'w': {0, 26}, 'x': {0, 27}, 'y': {0, 28},
	'z': {0, 29},

	'A': {2, 4}, 'B': {2, 5}, 'C': {2, 6}, 'D': {2, 7}, 'E': {2, 8},
	'F': {2, 9}, 'G': {2, 10}, 'H': {2, 11}, 'I': {2, 12}, 'J': {2, 13},
	'K': {2, 14}, 'L': {2, 15}, 'M': {2, 16}, 'N': {2, 17}, 'O': {2, 18},
	'P': {2, 19}, 'Q': {2, 20}, 'R': {2, 21}, 'S': {2, 22}, 'T': {2, 23},
	'U': {2, 24}, 'V': {2, 25}, 'W': {2, 26}, 'X': {2, 27}, 'Y': {2, 28},
	'Z': {2, 29},

	'1': {0, 30}, '2': {0, 31}, '3': {0, 32}, '4': {0, 33}, '5': {0, 34},
	'6': {0, 35}, '7': {0, 36}, '8': {0, 37}, '9': {0, 38}, '0': {0, 39},

	'!': {2, 30}, '@': {2, 31}, '#': {2, 32}, '$': {2, 33}, '%': {2, 34},
	'^': {2, 35}, '&': {2, 36}, '*': {2, 37}, '(': {2, 38}, ')': {2, 39},

	'\n': {0, 40}, '\t': {0, 43}, ' ': {0, 44}, '-': {0, 45},
	'=': {0, 46}, '[': {0, 47}, ']': {0, 48}, '\\': {0, 49},
	';': {0, 51}, '\'': {0, 52}, '`': {0, 53}, ',': {0, 54},
	'.': {0, 55}, '/': {0, 56},

	'_': {2, 45}, '+': {2, 46}, '{': {2, 47}, '}': {2, 48},
	'|': {2, 49}, ':': {2, 51}, '"': {2, 52}, '~': {2, 53},
	'<': {2, 54}, '>': {2, 55}, '?': {2, 56},
}
