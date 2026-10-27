//go:build windows

package app

import (
	"sort"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

const (
	defaultCharset = 1
	lfFaceName     = 32
)

type logfont struct {
	Height         int32
	Width          int32
	Escapement     int32
	Orientation    int32
	Weight         int32
	Italic         byte
	Underline      byte
	StrikeOut      byte
	CharSet        byte
	OutPrecision   byte
	ClipPrecision  byte
	Quality        byte
	PitchAndFamily byte
	FaceName       [lfFaceName]uint16
}

var (
	gdi32              = syscall.NewLazyDLL("gdi32.dll")
	user32             = syscall.NewLazyDLL("user32.dll")
	enumFontFamiliesEx = gdi32.NewProc("EnumFontFamiliesExW")
	getDC              = user32.NewProc("GetDC")
	releaseDC          = user32.NewProc("ReleaseDC")
	fontsOnce          sync.Once
	fonts              []string
)

func systemFonts() []string {
	fontsOnce.Do(func() {
		fonts = installedFonts()
	})
	return fonts
}

func installedFonts() []string {
	hdc, _, _ := getDC.Call(0)
	if hdc == 0 {
		return nil
	}
	defer releaseDC.Call(0, hdc)

	fonts := map[string]struct{}{}
	callback := syscall.NewCallback(func(font *logfont, _ uintptr, _ uint32, _ uintptr) uintptr {
		if font != nil {
			name := syscall.UTF16ToString(font.FaceName[:])
			if name != "" && !strings.HasPrefix(name, "@") {
				fonts[name] = struct{}{}
			}
		}
		return 1
	})
	var filter logfont
	filter.CharSet = defaultCharset
	enumFontFamiliesEx.Call(hdc, uintptr(unsafe.Pointer(&filter)), callback, 0, 0)

	names := make([]string, 0, len(fonts))
	for name := range fonts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
