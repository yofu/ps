package ps

import (
	"bytes"
	"fmt"
	"io"
)

type Doc struct {
	dsc    *DSC
	Canvas *Canvas
}

func NewDoc(title string) *Doc {
	d := new(Doc)
	d.dsc = NewDSC()
	d.Canvas = NewCanvas()
	d.dsc.Stuck(fmt.Sprintf("%%!PS-Adobe-3.0\n%%%%Title: %s\n", title))
	return d
}

func (d *Doc) WriteTo(otp io.Writer) (int64, error) {
	var tmp, rtn int64
	var err error
	if d.Canvas.page > 0 {
		val, err := d.Dsc(fmt.Sprintf("Pages: %d\n", d.Canvas.page))
		if err != nil {
			return rtn, err
		}
		rtn += int64(val)
	}
	tmp, err = d.dsc.WriteTo(otp)
	if err != nil {
		return rtn, err
	}
	rtn += tmp
	tmp, err = d.Canvas.WriteTo(otp)
	if err != nil {
		return rtn, err
	}
	rtn += tmp
	val, err := otp.Write([]byte("showpage\n"))
	if err != nil {
		return rtn, err
	}
	rtn += int64(val)
	return rtn, nil
}

func (d *Doc) Dsc(str string) (int, error) {
	return d.dsc.Stuck(fmt.Sprintf("%%%%%s", str))
}

func (d *Doc) Setup(str string) (int, error) {
	return d.dsc.Setup(str)
}

func (d *Doc) SetPaperSize(paper Paper) error {
	var err error
	_, err = d.Dsc(paper.Orientation())
	if err != nil {
		return err
	}
	_, err = d.Dsc(paper.DocumentMedia())
	if err != nil {
		return err
	}
	_, err = d.Setup(paper.SetPageDevice())
	if err != nil {
		return err
	}
	return nil
}

type DSC struct {
	stuck bytes.Buffer
	setup bytes.Buffer
}

func NewDSC() *DSC {
	return new(DSC)
}

func (d *DSC) WriteTo(otp io.Writer) (int64, error) {
	var tmp, rtn int64
	var err error
	tmp, err = d.stuck.WriteTo(otp)
	if err != nil {
		return rtn, err
	}
	rtn += tmp
	val, err := otp.Write([]byte("%%BeginSetup\n"))
	if err != nil {
		return rtn, err
	}
	rtn += int64(val)
	tmp, err = d.setup.WriteTo(otp)
	if err != nil {
		return rtn, err
	}
	rtn += tmp
	val, err = otp.Write([]byte("%%EndSetup\n"))
	if err != nil {
		return rtn, err
	}
	rtn += int64(val)
	return rtn, nil
}

func (d *DSC) Stuck(str string) (int, error) {
	return d.stuck.WriteString(str)
}

func (d *DSC) Setup(str string) (int, error) {
	return d.setup.WriteString(str)
}

type Canvas struct {
	stuck bytes.Buffer
	page  int
}

func NewCanvas() *Canvas {
	c := new(Canvas)
	c.page = 0
	return c
}

func MoveTo(x, y int) string {
	return fmt.Sprintf("%d %d moveto\n", x, y)
}
func FMoveTo(x, y float64) string {
	return fmt.Sprintf("%f %f moveto\n", x, y)
}

func RMoveTo(x, y int) string {
	return fmt.Sprintf("%d %d rmoveto\n", x, y)
}
func FRMoveTo(x, y float64) string {
	return fmt.Sprintf("%f %f rmoveto\n", x, y)
}

func LineTo(x, y int) string {
	return fmt.Sprintf("%d %d lineto\n", x, y)
}
func FLineTo(x, y float64) string {
	return fmt.Sprintf("%f %f lineto\n", x, y)
}

func RLineTo(x, y int) string {
	return fmt.Sprintf("%d %d rlineto\n", x, y)
}
func FRLineTo(x, y float64) string {
	return fmt.Sprintf("%f %f rlineto\n", x, y)
}

func Arc(x, y, r, start, end int) string {
	return fmt.Sprintf("%d %d %d %d %d arc\n", x, y, r, start, end)
}

func FArc(x, y, r, start, end float64) string {
	return fmt.Sprintf("%f %f %f %f %f arc\n", x, y, r, start, end)
}

func (cvs *Canvas) WriteTo(otp io.Writer) (int64, error) {
	return cvs.stuck.WriteTo(otp)
}

func (cvs *Canvas) Stuck(str string) (int, error) {
	return cvs.stuck.WriteString(str)
}

func (cvs *Canvas) NewPage(label string, paper Paper, setup ...string) (int, error) {
	var val, rtn int
	var err error
	var tmp bytes.Buffer
	if cvs.page > 0 {
		val, err = tmp.WriteString("showpage\n")
		if err != nil {
			return rtn, err
		}
		rtn += val
	}
	cvs.page++
	val, err = tmp.WriteString(fmt.Sprintf("%%%%Page: (%s) %d\n", label, cvs.page))
	if err != nil {
		return rtn, err
	}
	rtn += val
	if len(setup) > 0 || !paper.Portrait {
		val, err = tmp.WriteString("%%BeginPageSetup\n")
		if err != nil {
			return rtn, err
		}
		rtn += val
		val, err = tmp.WriteString(paper.PageSetup())
		if err != nil {
			return rtn, err
		}
		for _, s := range setup {
			val, err = tmp.WriteString(s)
			if err != nil {
				return rtn, err
			}
			rtn += val
		}
		val, err = tmp.WriteString("%%EndPageSetup\n")
		if err != nil {
			return rtn, err
		}
	}
	return cvs.Stuck(tmp.String())
}

func (cvs *Canvas) Page(label string, lines ...string) (int, error) {
	var val, rtn int
	var err error
	var tmp bytes.Buffer
	cvs.page++
	val, err = tmp.WriteString(fmt.Sprintf("%%%%Page: (%s) %d\n", label, cvs.page))
	if err != nil {
		return rtn, err
	}
	rtn += val
	for _, l := range lines {
		val, err = tmp.WriteString(l)
		if err != nil {
			return rtn, err
		}
		rtn += val
	}
	val, err = tmp.WriteString("showpage\n")
	if err != nil {
		return rtn, err
	}
	rtn += val
	rtn, err = cvs.Stuck(tmp.String())
	if err != nil {
		return rtn, err
	}
	return rtn, nil
}

func (cvs *Canvas) Def(key string, value ...string) (int, error) {
	switch len(value) {
	case 0:
		return 0, nil
	case 1:
		return cvs.Stuck(fmt.Sprintf("/%s %s def\n", key, value))
	default:
		var val, rtn int
		var err error
		var tmp bytes.Buffer
		val, err = tmp.WriteString(fmt.Sprintf("/%s {\n", key))
		if err != nil {
			return rtn, err
		}
		rtn += val
		for _, v := range value {
			val, err = tmp.WriteString(v)
			if err != nil {
				return rtn, err
			}
			rtn += val
		}
		val, err = tmp.WriteString(fmt.Sprintf("}def\n"))
		if err != nil {
			return rtn, err
		}
		return cvs.Stuck(tmp.String())
	}
}

func (cvs *Canvas) ForAll(list string, value ...string) (int, error) {
	var val, rtn int
	var err error
	var tmp bytes.Buffer
	val, err = tmp.WriteString(fmt.Sprintf("%s\n{\n", list))
	if err != nil {
		return rtn, err
	}
	rtn += val
	for _, v := range value {
		val, err = tmp.WriteString(v)
		if err != nil {
			return rtn, err
		}
		rtn += val
	}
	val, err = tmp.WriteString(fmt.Sprintf("}forall\n"))
	if err != nil {
		return rtn, err
	}
	rtn += val
	return cvs.Stuck(tmp.String())
}

func (cvs *Canvas) LineWidth(width int) (int, error) {
	return cvs.Stuck(fmt.Sprintf("%d setlinewidth\n", width))
}

func (cvs *Canvas) SetRGBColor(r, g, b float64) (int, error) {
	return cvs.Stuck(fmt.Sprintf("%f %f %f setrgbcolor\n", r, g, b))
}

func (cvs *Canvas) SetCMYKColor(c, m, y, k float64) (int, error) {
	return cvs.Stuck(fmt.Sprintf("%f %f %f %f setrgbcolor\n", c, m, y, k))
}

func (cvs *Canvas) Path(closed bool, name string, lines ...string) (int, error) {
	var val, rtn int
	var err error
	var tmp bytes.Buffer
	val, err = tmp.WriteString("newpath\n")
	if err != nil {
		return rtn, err
	}
	rtn += val
	for _, l := range lines {
		val, err = tmp.WriteString(l)
		if err != nil {
			return rtn, err
		}
		rtn += val
	}
	if closed {
		val, err = tmp.WriteString("closepath\n")
		if err != nil {
			return rtn, err
		}
		rtn += val
	}
	val, err = tmp.WriteString(fmt.Sprintf("%s\n", name))
	if err != nil {
		return rtn, err
	}
	return cvs.Stuck(tmp.String())
}
func (cvs *Canvas) Stroke(closed bool, lines ...string) (int, error) {
	return cvs.Path(closed, "stroke", lines...)
}
func (cvs *Canvas) Fill(closed bool, lines ...string) (int, error) {
	return cvs.Path(closed, "fill", lines...)
}
func (cvs *Canvas) EOFill(closed bool, lines ...string) (int, error) {
	return cvs.Path(closed, "eofill", lines...)
}

func (cvs *Canvas) Line(x0, y0, x1, y1 int) (int, error) {
	return cvs.Stroke(false,
		MoveTo(x0, y0),
		LineTo(x1, y1))
}
func (cvs *Canvas) FLine(x0, y0, x1, y1 float64) (int, error) {
	return cvs.Stroke(false,
		FMoveTo(x0, y0),
		FLineTo(x1, y1))
}

func (cvs *Canvas) Polyline(coord [][]int) (int, error) {
	if len(coord) < 2 {
		return 0, fmt.Errorf("Polyline: not enough length")
	}
	lines := make([]string, len(coord))
	lines[0] = MoveTo(coord[0][0], coord[0][1])
	for i, c := range coord[1:] {
		lines[i+1] = LineTo(c[0], c[1])
	}
	return cvs.Stroke(false, lines...)
}
func (cvs *Canvas) FPolyline(coord [][]float64) (int, error) {
	if len(coord) < 2 {
		return 0, fmt.Errorf("FPolyline: not enough length")
	}
	lines := make([]string, len(coord))
	lines[0] = FMoveTo(coord[0][0], coord[0][1])
	for i, c := range coord[1:] {
		lines[i+1] = FLineTo(c[0], c[1])
	}
	return cvs.Stroke(false, lines...)
}

func (cvs *Canvas) Polygon(coord [][]int) (int, error) {
	if len(coord) < 2 {
		return 0, fmt.Errorf("Polyline: not enough length")
	}
	lines := make([]string, len(coord))
	lines[0] = MoveTo(coord[0][0], coord[0][1])
	for i, c := range coord[1:] {
		lines[i+1] = LineTo(c[0], c[1])
	}
	return cvs.Fill(false, lines...)
}
func (cvs *Canvas) FPolygon(coord [][]float64) (int, error) {
	if len(coord) < 2 {
		return 0, fmt.Errorf("FPolyline: not enough length")
	}
	lines := make([]string, len(coord))
	lines[0] = FMoveTo(coord[0][0], coord[0][1])
	for i, c := range coord[1:] {
		lines[i+1] = FLineTo(c[0], c[1])
	}
	return cvs.Fill(false, lines...)
}

func (cvs *Canvas) Circle(x, y, r int) (int, error) {
	return cvs.Stroke(true, Arc(x, y, r, 0, 360))
}
func (cvs *Canvas) FCircle(x, y, r float64) (int, error) {
	return cvs.Stroke(true, FArc(x, y, r, 0, 360))
}

func (cvs *Canvas) FilledCircle(x, y, r int) (int, error) {
	return cvs.Fill(true, Arc(x, y, r, 0, 360))
}
func (cvs *Canvas) FFilledCircle(x, y, r float64) (int, error) {
	return cvs.Fill(true, FArc(x, y, r, 0, 360))
}
