package proto

type GetHidModeRsp struct {
	Mode string `json:"mode"` // normal or hid-only
}

type SetHidModeReq struct {
	Mode string `validate:"required"` // normal or hid-only
}

type ShortcutKey struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type Shortcut struct {
	ID   string        `json:"id"`
	Keys []ShortcutKey `json:"keys"`
}

type GetShortcutsRsp struct {
	Shortcuts []Shortcut `json:"shortcuts"`
}

type AddShortcutReq struct {
	Keys []ShortcutKey `validate:"required"`
}

type DeleteShortcutReq struct {
	ID string `validate:"required"`
}

type SetLeaderKeyReq struct {
	Key string `validate:"omitempty"`
}

type GetLeaderKeyRsp struct {
	Key string `json:"key"`
}

type Macro struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Script string `json:"script"`
}

type GetMacrosRsp struct {
	Macros []Macro `json:"macros"`
}

type AddMacroReq struct {
	Name   string `json:"name" validate:"required,max=64"`
	Script string `json:"script" validate:"required,min=1,max=16384"`
}

type UpdateMacroReq struct {
	ID     string `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required,max=64"`
	Script string `json:"script" validate:"required,min=1,max=16384"`
}

type DeleteMacroReq struct {
	ID string `json:"id" validate:"required"`
}
