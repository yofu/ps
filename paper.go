package ps

import (
	"fmt"
	"strings"
)

type Paper struct {
	name string
	width int
	height int
	portrait bool
}

var (
	A4Portrait  = Paper{"A4", 595, 842, true}
	A4Landscape = Paper{"A4", 595, 842, false}
)

func (p Paper) Orientation() string {
	if p.portrait {
		return "Orientation: Portrait\n"
	} else {
		return "Orientation: Landscape\n"
	}
}

func (p Paper) DocumentMedia() string {
	return fmt.Sprintf("DocumentMedia: %s %d %d 80 () ()\n", strings.ToLower(p.name), p.width, p.height)
}

func (p Paper) SetPageDevice() string {
	if p.portrait {
		return fmt.Sprintf("<< /PageSize [%d %d] /Orientation 0 >> setpagedevice\n", p.width, p.height)
	} else {
		return fmt.Sprintf("<< /PageSize [%d %d] /Orientation 3 >> setpagedevice\n", p.width, p.height)
	}
}

func (p Paper) PageSetup() string {
	if p.portrait {
		return ""
	} else {
		return fmt.Sprintf("90 rotate 0 -%d translate\n", p.width)
	}
}
