package api

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kraken-hpc/imageapi/models"
	"github.com/kraken-hpc/uinit"
)

// Script wraps uinit.Script with running logic for container hooks
type Script struct {
	Script *uinit.Script
	Must   bool
}

// NewScript builds a new script based on script, mountpoint, log, and keyvalue store
// Note: we need mountpoint so mountpoint-relative paths can be resolved
func NewScript(script *models.ContainerScript, mountpoint string, ctx *uinit.ModuleContext) (*Script, error) {
	s := &Script{}
	if err := s.Build(script, mountpoint, ctx); err != nil {
		return nil, err
	}
	return s, nil
}

// Run the script; only return an error if Must
func (s *Script) Run(ctx *uinit.ModuleContext) error {
	s.Script.Context = ctx
	err := s.Script.Run()
	if err != nil && s.Must {
		return err
	}
	return nil
}

// Build initializes a Script based on a ContainerScript specification
func (s *Script) Build(script *models.ContainerScript, mountpoint string, ctx *uinit.ModuleContext) error {
	var err error
	var data []byte
	var ns *uinit.Script
	s.Must = script.Must != nil && *script.Must
	switch script.Encoding {
	case models.ContainerScriptEncodingFile:
		if data, err = ioutil.ReadFile(script.Script); err != nil {
			err = fmt.Errorf("failed to read script file %s: %v", script.Script, err)
			goto error
		}
	case models.ContainerScriptEncodingContainerFile:
		if data, err = ioutil.ReadFile(filepath.Join(mountpoint, script.Script)); err != nil {
			err = fmt.Errorf("failed to read script file %s: %v", filepath.Join(mountpoint, script.Script), err)
			goto error
		}
	case models.ContainerScriptEncodingPlain:
		data = ([]byte)(script.Script)
	case models.ContainerScriptEncodingBase64:
		if data, err = base64.StdEncoding.DecodeString(script.Script); err != nil {
			err = fmt.Errorf("failed to decode base64 script: %v", err)
			goto error
		}
	case models.ContainerScriptEncodingGzip:
		if data, err = base64.StdEncoding.DecodeString(script.Script); err != nil {
			err = fmt.Errorf("failed to decode base64/gzip script: %v", err)
			goto error
		}
		var r *gzip.Reader
		if r, err = gzip.NewReader(bytes.NewReader(data)); err != nil {
			err = fmt.Errorf("could not decompress gzip script: %v", err)
			goto error
		}
		if data, err = ioutil.ReadAll(r); err != nil {
			err = fmt.Errorf("could not decompress gzip script: %v", err)
			goto error
		}
	case models.ContainerScriptEncodingBzip2:
		if data, err = base64.StdEncoding.DecodeString(script.Script); err != nil {
			err = fmt.Errorf("failed to decode base64/gzip script: %v", err)
			goto error
		}
		r := bzip2.NewReader(bytes.NewReader(data))
		if data, err = ioutil.ReadAll(r); err != nil {
			err = fmt.Errorf("failed to decode base64/bzip2 script: %v", err)
			goto error
		}
	default:
		return fmt.Errorf("unrecognized encoding type: %v", script.Encoding)
	}
	if ns, err = uinit.NewScript(data, nil); err != nil {
		err = fmt.Errorf("failed to parse script: %v", err)
		goto error
	}
	s.Script = ns
	s.Script.Context = ctx
	return nil

error:
	// we set these even if it's not fatal
	script.Success = new(bool)
	script.LastError = err.Error()
	if s.Must {
		return err
	}
	ctx.Log.Printf("failed to parse non-mandatory script: %v", err)
	return nil
}

// A Hook is a sequence of scripts, includes error handling logic for running
type Hook struct {
	Scripts       []*Script
	Context       *uinit.ModuleContext
	Mountpoint    string
	DefaultScript string
}

// NewHook creates an initialized, built hook
func NewHook(csh *models.ContainerScriptHook, defaultScript, mountpoint string, l *log.Logger, vars map[string]string) (hook *Hook, err error) {
	hook = &Hook{
		Scripts: []*Script{},
		Context: &uinit.ModuleContext{
			Vars: uinit.NewSimpleKV(),
			Log:  l,
		},
		Mountpoint:    mountpoint,
		DefaultScript: defaultScript,
	}
	if hook.Context.Log == nil {
		hook.Context.Log = log.New(os.Stdout, "", 0)
	}
	if hook.Context.Vars == nil {
		hook.Context.Vars = uinit.NewSimpleKV()
	}
	for k, v := range vars {
		hook.Context.Vars.Set(k, v)
	}
	if csh != nil {
		if err = hook.Build(csh); err != nil {
			return nil, err
		}
	}
	return
}

func (h *Hook) Run() error {
	for _, script := range h.Scripts {
		if err := script.Run(h.Context); err != nil {
			return fmt.Errorf("mandatory script failed: %v", err)
		}
	}
	return nil
}

func (h *Hook) Build(hook *models.ContainerScriptHook) (err error) {
	h.Scripts = []*Script{}
	if hook == nil {
		hook = &models.ContainerScriptHook{
			Scripts:         []*models.ContainerScript{},
			DisableDefaults: new(bool),
		}
	}
	if (hook.DisableDefaults == nil || !*hook.DisableDefaults) && h.DefaultScript != "" {
		var s *Script
		if s, _ = NewScript(&models.ContainerScript{
			Encoding: models.ContainerScriptEncodingFile,
			Script:   h.DefaultScript,
			Must:     new(bool), // should this be hard-coded?
		}, h.Mountpoint, h.Context); s.Script != nil {
			h.Scripts = append(h.Scripts, s)
		} else {
			h.Context.Log.Printf("not adding default script, non-fatal error")
		}
	}
	for i, script := range hook.Scripts {
		var s *Script
		if s, err = NewScript(script, h.Mountpoint, h.Context); err != nil {
			return err
		}
		if s.Script != nil { // this happens if a non-fatal error occurred
			h.Scripts = append(h.Scripts, s)
		} else {
			h.Context.Log.Printf("not adding script %d, non-fatal error", i)
		}
	}
	return
}
