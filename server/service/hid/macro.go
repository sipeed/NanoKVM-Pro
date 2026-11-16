package hid

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	macroFile  = "/etc/kvm/macros.json"
	macroMutex = sync.RWMutex{}
)

type MacroStore struct {
	Macros []proto.Macro `json:"macros"`
}

func (s *Service) GetMacros(c *gin.Context) {
	var rsp proto.Response

	macros, err := listMacros()
	if err != nil {
		log.Errorf("failed to get macros: %v", err)
		rsp.ErrRsp(c, -1, "get macros failed")
		return
	}

	rsp.OkRspWithData(c, &proto.GetMacrosRsp{
		Macros: macros,
	})
	log.Debugf("get macros success, total: %d", len(macros))
}

func (s *Service) AddMacro(c *gin.Context) {
	var req proto.AddMacroReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if err := validateMacro(req.Name, req.Script); err != nil {
		rsp.ErrRsp(c, -1, err.Error())
		return
	}

	macro, err := addMacro(req.Name, req.Script)
	if err != nil {
		log.Errorf("failed to add macro: %v", err)
		rsp.ErrRsp(c, -2, "add macro failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("add macro %s", macro.ID)
}

func (s *Service) UpdateMacro(c *gin.Context) {
	var req proto.UpdateMacroReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if err := validateMacro(req.Name, req.Script); err != nil {
		rsp.ErrRsp(c, -1, err.Error())
		return
	}

	if err := updateMacro(req.ID, req.Name, req.Script); err != nil {
		log.Errorf("failed to update macro: %v", err)
		if errors.Is(err, errMacroNotFound) {
			rsp.ErrRsp(c, -2, "macro not found")
			return
		}
		rsp.ErrRsp(c, -2, "update macro failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("update macro %s", req.ID)
}

func (s *Service) DeleteMacro(c *gin.Context) {
	var req proto.DeleteMacroReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if err := deleteMacro(req.ID); err != nil {
		log.Errorf("failed to delete macro: %v", err)
		if errors.Is(err, errMacroNotFound) {
			rsp.ErrRsp(c, -2, "macro not found")
			return
		}
		rsp.ErrRsp(c, -2, "delete macro failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("delete macro %s", req.ID)
}

var errMacroNotFound = errors.New("macro not found")

func validateMacro(name, script string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	trimmed := strings.TrimSpace(script)
	if trimmed == "" {
		return errors.New("script is required")
	}

	if utf8.RuneCountInString(script) > 16384 {
		return errors.New("script too long")
	}

	return nil
}

func loadMacros() (*MacroStore, error) {
	if _, err := os.Stat(macroFile); os.IsNotExist(err) {
		return &MacroStore{Macros: []proto.Macro{}}, nil
	}

	data, err := os.ReadFile(macroFile)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return &MacroStore{Macros: []proto.Macro{}}, nil
	}

	var store MacroStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	return &store, nil
}

func saveMacros(store *MacroStore) error {
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}

	return os.WriteFile(macroFile, data, 0644)
}

func listMacros() ([]proto.Macro, error) {
	macroMutex.RLock()
	defer macroMutex.RUnlock()

	store, err := loadMacros()
	if err != nil {
		return nil, err
	}

	return store.Macros, nil
}

func addMacro(name, script string) (*proto.Macro, error) {
	macroMutex.Lock()
	defer macroMutex.Unlock()

	store, err := loadMacros()
	if err != nil {
		return nil, err
	}

	macro := proto.Macro{
		ID:     uuid.New().String(),
		Name:   strings.TrimSpace(name),
		Script: script,
	}

	store.Macros = append(store.Macros, macro)

	if err := saveMacros(store); err != nil {
		return nil, err
	}

	return &macro, nil
}

func updateMacro(id, name, script string) error {
	macroMutex.Lock()
	defer macroMutex.Unlock()

	store, err := loadMacros()
	if err != nil {
		return err
	}

	found := false
	for i, macro := range store.Macros {
		if macro.ID == id {
			store.Macros[i].Name = strings.TrimSpace(name)
			store.Macros[i].Script = script
			found = true
			break
		}
	}

	if !found {
		return errMacroNotFound
	}

	return saveMacros(store)
}

func deleteMacro(id string) error {
	macroMutex.Lock()
	defer macroMutex.Unlock()

	store, err := loadMacros()
	if err != nil {
		return err
	}

	found := false
	newMacros := make([]proto.Macro, 0, len(store.Macros))
	for _, macro := range store.Macros {
		if macro.ID == id {
			found = true
			continue
		}
		newMacros = append(newMacros, macro)
	}

	if !found {
		return errMacroNotFound
	}

	store.Macros = newMacros
	return saveMacros(store)
}
