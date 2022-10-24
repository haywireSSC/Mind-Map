package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/slices"
)

var idealCornerSize float64 = 50
var colourIndex int

var PALETTE = map[string]rl.Color{
	"Bakerloo":     rl.Color{135, 78, 36, 255},   //baker loo line (brown)rl.Color{178, 99, 0, 255}
	"Hammersmith":  rl.Color{244, 169, 190, 255}, //hammersmith and city (pink)
	"Piccadilly":   rl.Color{0, 25, 168, 255},    //picadilly (dark blue)
	"Central":      rl.Color{220, 36, 31, 255},   //central line(red)
	"Jubilee":      rl.Color{161, 165, 167, 255}, //jubilee line (gray)
	"Victoria":     rl.Color{0, 152, 216, 255},   //victoria line(light blue)
	"Circle":       rl.Color{255, 211, 41, 255},  //circle line (yellow)
	"Metropolitan": rl.Color{155, 0, 88, 255},    //metropolian(dark purple)
	"Waterloo":     rl.Color{147, 206, 186, 255}, //waterloo and city line(aqua)
	"District":     rl.Color{0, 125, 50, 255},    //district line (green)
	"Northern":     rl.Color{0, 0, 0, 255},       //northern line (black)

	"Overground": rl.Color{232, 106, 16, 255},  //(orange)rl.Color{239, 123, 16, 255}
	"Elizabeth":  rl.Color{147, 100, 204, 255}, //(purple)
	"Saftey":     rl.Color{0, 96, 168, 255},    //(blue)
	"Tramlink":   rl.Color{0, 189, 25, 255},    //(lime)
	"Docklands":  rl.Color{0, 175, 173, 255},   //(aqua)
	"White":      rl.Color{255, 255, 255, 255}}

type Outline struct {
	Inner string
	Outer string
}

var OutlineColours = [6]Outline{
	Outline{"Jubilee", "District"},

	Outline{"Piccadilly", "Circle"},
	Outline{"Piccadilly", "Central"},

	Outline{"Victoria", "Overground"},
	Outline{"Circle", "Northern"},
	Outline{"Waterloo", "Metropolitan"}}

var LineColours = [11]string{
	"Bakerloo",
	"Hammersmith",
	"Piccadilly",
	"Central",
	"Jubilee",
	"Victoria",
	"Circle",
	"Metropolitan",
	"Waterloo",
	"District",
	"Northern"}

func GrowAround(point, origin rl.Vector2, length float32) rl.Vector2 {
	newPoint := rl.Vector2Subtract(point, origin)
	newPoint = rl.Vector2Normalize(newPoint)
	newPoint = rl.Vector2Multiply(newPoint, rl.Vector2{length, length})
	newPoint = rl.Vector2Add(origin, newPoint)
	return newPoint
}

type AngleVec struct {
	Pos rl.Vector2
	Dir rl.Vector2
}

type BezierCurve struct {
	start     rl.Vector2
	end       rl.Vector2
	startCtrl rl.Vector2
	endCtrl   rl.Vector2
}

func NormalizdDiff(start, end rl.Vector2) rl.Vector2 {
	return rl.Vector2Normalize(rl.Vector2Subtract(end, start))
}

func InterpolateBezier(c BezierCurve, amount float32) rl.Vector2 {
	//caluclate inner lines
	a1 := rl.Vector2Lerp(c.start, c.startCtrl, amount)
	a2 := rl.Vector2Lerp(c.startCtrl, c.endCtrl, amount)
	a3 := rl.Vector2Lerp(c.endCtrl, c.end, amount)

	//calculate inner inners
	b1 := rl.Vector2Lerp(a1, a2, amount)
	b2 := rl.Vector2Lerp(a2, a3, amount)

	// interpolate that
	return rl.Vector2Lerp(b1, b2, amount)
}

func GenerateBezier(curve BezierCurve, pointAmount int) (points []rl.Vector2) {
	var amount float32
	step := 1 / float32(pointAmount)
	for i := 0; i < pointAmount; i++ {
		amount += step
		points = append(points, InterpolateBezier(curve, amount))
	}
	return
}

func getDistances(points []rl.Vector2) (dists []float32) {
	dists = make([]float32, len(points))
	var totalDist float32

	for i, v := range points {
		if i > 0 {
			totalDist += rl.Vector2Distance(v, points[i-1]) //distance between current and prev
		}
		dists[i] = totalDist
	}

	return
}

func LerpAlongPoints(points []rl.Vector2, amount float32) (angleVec AngleVec) {
	//array of point and total dist
	//get total length(last item)
	//do amount * length
	//get that line and interpolate along

	dists := getDistances(points)
	var totalDist float32

	for i, v := range points {
		if i > 0 {
			totalDist += rl.Vector2Distance(v, points[i-1]) //distance between current and prev
		}
		dists[i] = totalDist
	}

	targetDist := amount * dists[len(dists)-1] //times amount by total length
	for i := len(dists) - 1; i >= 0; i-- {
		if dists[i] <= targetDist { //lerp along inner and get angle
			if i == len(dists)-1 {
				return
			}
			innerAmount := (targetDist - dists[i]) / (dists[i+1] - dists[i])
			angleVec.Pos = rl.Vector2Lerp(points[i], points[i+1], innerAmount)
			angleVec.Dir = NormalizdDiff(points[i], points[i+1])
			return
		}
	}
	return
}

func DrawAngleTriangle(pos, dir rl.Vector2, width, height float32, colour rl.Color) {

	var v1, v2, v3, end rl.Vector2
	v2 = pos
	end = pos

	end.X -= dir.X * height
	end.Y -= dir.Y * height

	v1 = end
	v3 = end

	v1.X += dir.Y * width
	v1.Y += dir.X * -width

	v3.X -= dir.Y * width
	v3.Y -= dir.X * -width
	rl.DrawTriangle(v1, v2, v3, colour) //sometimes no work due to direction
	//rl.DrawCircleV(v1, 5, rl.Red)
	//rl.DrawCircleV(v2, 5, rl.Green)
	//rl.DrawCircleV(v3, 5, rl.Blue)
}

func FindEdge(from, to *Node) (pos rl.Vector2) {
	if from.Rect.Width != 0 && from.Rect.Height != 0 {

		xdiff := float64(from.Center.X - to.Center.X)
		ydiff := float64(from.Center.Y - to.Center.Y)
		xlen := math.Abs(xdiff)
		ylen := math.Abs(ydiff)

		rect := from.Rect

		pos = from.Center

		if xlen > ylen {
			if xdiff > 0 {
				pos.X -= rect.Width/2 - float32(from.Theme.Margin)
			} else {
				pos.X += rect.Width/2 - float32(from.Theme.Margin)
			}
		} else {
			if ydiff > 0 {
				pos.Y -= rect.Height/2 - float32(from.Theme.Margin)
			} else {
				pos.Y += rect.Height/2 - float32(from.Theme.Margin)
			}
		}
	} else {
		pos = from.Center
	}

	return
}

func FindEdgeToPoint(from *Node, to rl.Vector2) (pos rl.Vector2) {
	if from.Rect.Width != 0 && from.Rect.Height != 0 {

		xdiff := float64(from.Center.X - to.X)
		ydiff := float64(from.Center.Y - to.Y)
		xlen := math.Abs(xdiff)
		ylen := math.Abs(ydiff)

		rect := from.Rect

		pos = from.Center

		if xlen > ylen {
			if xdiff > 0 {
				pos.X -= rect.Width/2 - float32(from.Theme.Margin)
			} else {
				pos.X += rect.Width/2 - float32(from.Theme.Margin)
			}
		} else {
			if ydiff > 0 {
				pos.Y -= rect.Height/2 - float32(from.Theme.Margin)
			} else {
				pos.Y += rect.Height/2 - float32(from.Theme.Margin)
			}
		}
	} else {
		pos = from.Center
	}

	return
}

func DrawPath(start, end rl.Vector2, isInverted bool, colour rl.Color, outlined bool, outlineColour rl.Color, dotted bool) {
	xdiff := float64(end.X - start.X)
	ydiff := float64(end.Y - start.Y)

	//quick fix to straight and diag not work
	var points []rl.Vector2
	if !(math.Abs(xdiff) == math.Abs(ydiff) || xdiff == 0 || ydiff == 0) {
		var mid rl.Vector2
		if (xdiff > 0) == (ydiff > 0) {
			if isInverted {

				if math.Abs(xdiff) < math.Abs(ydiff) {
					mid.X = end.X
					mid.Y = end.Y - float32(ydiff-xdiff)
				} else {
					mid.Y = end.Y
					mid.X = end.X - float32(xdiff-ydiff)
				}
			} else {

				if math.Abs(xdiff) < math.Abs(ydiff) {
					mid.X = start.X
					mid.Y = start.Y + float32(ydiff-xdiff)
				} else {
					mid.Y = start.Y
					mid.X = start.X + float32(xdiff-ydiff)
				}
			}
		} else {
			if isInverted {

				if math.Abs(xdiff) < math.Abs(ydiff) {
					mid.X = end.X
					mid.Y = end.Y - float32(ydiff+xdiff)
				} else {
					mid.Y = end.Y
					mid.X = end.X - float32(xdiff+ydiff)
				}
			} else {

				if math.Abs(xdiff) < math.Abs(ydiff) {
					mid.X = start.X
					mid.Y = start.Y + float32(ydiff+xdiff)
				} else {
					mid.Y = start.Y
					mid.X = start.X + float32(xdiff+ydiff)
				}
			}
		}

		midToStart := rl.Vector2Subtract(mid, start)
		midToEnd := rl.Vector2Subtract(mid, end)

		cornerSize := float32(math.Floor(math.Min(math.Min(float64(rl.Vector2Length(midToStart)), float64(rl.Vector2Length(midToEnd))), idealCornerSize))) - 1

		dir := rl.Vector2Normalize(midToStart)
		startMid := mid
		startMid.X -= dir.X * cornerSize
		startMid.Y -= dir.Y * cornerSize

		dir = rl.Vector2Normalize(midToEnd)
		endMid := mid
		endMid.X -= dir.X * cornerSize
		endMid.Y -= dir.Y * cornerSize

		//making curve and points list
		curve := BezierCurve{startMid, endMid, GrowAround(start, startMid, -cornerSize), GrowAround(end, endMid, -cornerSize)}
		points = GenerateBezier(curve, 10)
	}
	points = slices.Insert(points, 0, start)
	points = append(points, end)

	//drawing triangle
	triangle := LerpAlongPoints(points, 0.5)

	if outlined {
		DrawLinesFromPoints(points, 15, outlineColour)
		DrawAngleTriangle(triangle.Pos, triangle.Dir, 20, -30, colour)
	} else {
		DrawAngleTriangle(triangle.Pos, triangle.Dir, 10, -20, colour)
	}

	if dotted {
		DrawLinesDotted(points, 5, colour)
	} else {
		DrawLinesFromPoints(points, 5, colour)
	}
}

func DrawLinesFromPoints(points []rl.Vector2, thickness float32, colour rl.Color) {
	halfThick := thickness / 2
	for i := 0; i < len(points)-1; i++ {
		rl.DrawCircleV(points[i], halfThick, colour)
		rl.DrawLineEx(points[i], points[i+1], thickness, colour)
	}
}

func LerpAlongPointsWithDist(points []rl.Vector2, targetDist float32, dists []float32) (angleVec AngleVec) {
	//array of point and total dist
	//get total length(last item)
	//do amount * length
	//get that line and interpolate along
	for i := len(dists) - 1; i >= 0; i-- {
		if dists[i] <= targetDist { //lerp along inner and get angle
			if i == len(dists)-1 {
				return
			}
			innerAmount := (targetDist - dists[i]) / (dists[i+1] - dists[i])
			angleVec.Pos = rl.Vector2Lerp(points[i], points[i+1], innerAmount)
			angleVec.Dir = NormalizdDiff(points[i], points[i+1])
			return
		}
	}
	return
}

func DrawLinesDotted(points []rl.Vector2, thickness float32, colour rl.Color) {

	dists := getDistances(points)

	stepSize := float32(10)

	totalDist := dists[len(dists)-1]
	prevPos := points[0]
	i := 0
	for d := float32(0); d <= totalDist; d += stepSize {
		i += 1
		pos := LerpAlongPointsWithDist(points, d, dists).Pos
		if i%2 == 0 {
			rl.DrawLineEx(prevPos, pos, thickness, colour)
		}
		prevPos = pos
	}
}
