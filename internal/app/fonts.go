package app

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

func systemFonts() []string {
	families := ui.FontFamilies()
	names := make([]string, 0, len(families))
	for _, family := range families {
		if family.Name != "" && !strings.HasPrefix(family.Name, "@") {
			names = append(names, family.Name)
		}
	}
	return names
}
