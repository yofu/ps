package ps

import (
	"fmt"
	"strings"
)

type Paper struct {
	name string
	width int
	height int
	Portrait bool
}

var (
	A4Portrait  = Paper{"A4", 595, 842, true}
	A4Landscape = Paper{"A4", 595, 842, false}
	A3Portrait  = Paper{"A3", 842, 1190, true}
	A3Landscape = Paper{"A3", 842, 1190, false}
)

func (p Paper) Size() (int, int) {
	return p.width, p.height
}

func (p Paper) Orientation() string {
	if p.Portrait {
		return "Orientation: Portrait\n"
	} else {
		return "Orientation: Landscape\n"
	}
}

func (p Paper) DocumentMedia() string {
	return fmt.Sprintf("DocumentMedia: %s %d %d 80 () ()\n", strings.ToLower(p.name), p.width, p.height)
}

func (p Paper) SetPageDevice() string {
	if p.Portrait {
		return fmt.Sprintf("<< /PageSize [%d %d] /Orientation 0 >> setpagedevice\n", p.width, p.height)
	} else {
		return fmt.Sprintf("<< /PageSize [%d %d] /Orientation 3 >> setpagedevice\n", p.width, p.height)
	}
}

func (p Paper) PageSetup() string {
	if p.Portrait {
		return ""
	} else {
		return fmt.Sprintf("90 rotate 0 -%d translate\n", p.width)
	}
}
