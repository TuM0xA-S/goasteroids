package vobj

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/vector"

	"github.com/hajimehoshi/ebiten"
)

const (
	DirLeft  = -1
	DirRight = +1
)

//Consts
var (
	MaxX     float64
	MaxY     float64
)

//Vec2 represents 2d vector
type Vec2 struct {
	X, Y float64
}

//Rotate ...
func (v *Vec2) Rotate(angle float64) {
	g := ebiten.GeoM{}
	g.Rotate(angle)
	v.X, v.Y = g.Apply(v.X, v.Y)
}

//Move ...
func (v *Vec2) Move(dx, dy float64) {
	g := ebiten.GeoM{}
	g.Translate(dx, dy)
	v.X, v.Y = g.Apply(v.X, v.Y)
}

//Scale ...
func (v *Vec2) Scale(factor float64) {
	g := ebiten.GeoM{}
	g.Scale(factor, factor)
	v.X, v.Y = g.Apply(v.X, v.Y)
}

//Add ...
func (v Vec2) Add(other Vec2) (res Vec2) {
	res.X = v.X + other.X
	res.Y = v.Y + other.Y
	return
}

//Mult ...
func (v Vec2) Mult(scalar float64) (res Vec2) {
	res.X = v.X * scalar
	res.Y = v.Y * scalar
	return
}

//Len ...
func (v *Vec2) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

//SetLen ...
func (v *Vec2) SetLen(newLen float64) {
	len := v.Len()
	v.X /= len
	v.Y /= len
	*v = v.Mult(newLen)
}

//VectorObject is a universal object representing any vector object in game
type VectorObject struct {
	Position          Vec2
	Speed             Vec2
	AccelerationValue float64
	Geometry          []Vec2
	Scale             float64
	Color             color.RGBA
	RotateSpeed       float64
	Angle             float64
	MaxSpeed          float64
}

//Rotate acceleration vector
func (v *VectorObject) Rotate(dt float64, direction float64) {
	v.Angle += v.RotateSpeed * dt * direction
	if v.Angle >= 2*math.Pi {
		v.Angle -= 2 * math.Pi
	}
	if v.Angle < 0 {
		v.Angle += 2 * math.Pi
	}
}

//Acceleration
func (v *VectorObject) Acceleration() (res Vec2) {
	res.X = 0
	res.Y = -1
	res.Rotate(v.Angle)
	res = res.Mult(v.AccelerationValue)
	return
}

//Accelerate object
func (v *VectorObject) Accelerate(dt float64) {
	v.Speed = v.Speed.Add(v.Acceleration().Mult(dt))
	if v.Speed.Len() > v.MaxSpeed {
		v.Speed.SetLen(v.MaxSpeed)
	}
}

//Move in current direction
func (v *VectorObject) Move(dt float64) {
	v.Position = v.Position.Add(v.Speed.Mult(dt))

	if v.Position.X >= MaxX {
		v.Position.X -= MaxX
	}
	if v.Position.X < 0 {
		v.Position.X += MaxX
	}

	if v.Position.Y >= MaxY {
		v.Position.Y -= MaxY
	}
	if v.Position.Y < 0 {
		v.Position.Y += MaxY
	}

}

//Draw on dst...
func (v *VectorObject) Draw(dst *ebiten.Image) {
	vecs := v.GetTransformed()
	for _, delta := range []Vec2{{0, 0}, {-MaxX, 0}, {MaxX, 0}, {0, MaxY}, {0, -MaxY}} {
		figure := vector.Path{}
		figure.MoveTo(float32(vecs[0].X+delta.X), float32(vecs[0].Y+delta.Y))
		for _, p := range vecs[1:] {
			figure.LineTo(float32(p.X+delta.X), float32(p.Y+delta.Y))
		}
		figure.Fill(dst, &vector.FillOptions{Color: v.Color})
	}

}

//GetTransformed gives transformed object
func (v *VectorObject) GetTransformed() []Vec2 {
	vecs := append([]Vec2{}, v.Geometry...)
	for i := range vecs {
		vecs[i].Rotate(v.Angle)
		vecs[i].Scale(v.Scale)
		vecs[i].Move(v.Position.X, v.Position.Y)
	}
	return vecs
}

func MakeVector(angle, lenght float64) Vec2 {
	v := Vec2{0, -lenght}
	v.Rotate(angle)

	return v
}

func getAllSegs(points []Vec2) (segs [][2]Vec2) {
	for _, delta := range []Vec2{{0, 0}, {-MaxX, 0}, {MaxX, 0}, {0, MaxY}, {0, -MaxY}} {
		for i := 0; i < len(points); i++ {
			p0 := points[i]
			p0.X += delta.X
			p0.Y += delta.Y
			p1 := points[(i+1)%len(points)]
			p1.X += delta.X
			p1.Y += delta.Y

			segs = append(segs, [2]Vec2{p0, p1})
		}
	}

	return
}

//Collides cheks circuit intersection
func (v *VectorObject) Collides(other *VectorObject) bool {
	a := getAllSegs(v.GetTransformed())
	b := getAllSegs(other.GetTransformed())
	for i := range a {
		for j := range b {
			if SectionIntersects(a[i][0], a[i][1], b[j][0], b[j][1]) {
				return true
			}
		}
	}

	return false
}

func minmax(a, b float64) (float64, float64) {
	return math.Min(a, b), math.Max(a, b)
}

func between(a, b, c float64) bool {
	a, b = minmax(a, b)
	return c >= a*(1-2e-15) && c <= b*(1+2e-15)
}

func SectionIntersects(v0, v1, u0, u1 Vec2) bool {
	a1 := v0.Y - v1.Y
	b1 := v1.X - v0.X
	c1 := a1*v1.X + b1*v1.Y
	a2 := u0.Y - u1.Y
	b2 := u1.X - u0.X
	c2 := a2*u1.X + b2*u1.Y

	d := a1*b2 - b1*a2
	if math.Abs(d) < 1e-9 {
		return false
	}
	dx := c1*b2 - b1*c2
	x := dx / d
	dy := a1*c2 - c1*a2
	y := dy / d

	return between(v0.X, v1.X, x) && between(v0.Y, v1.Y, y) &&
		between(u0.X, u1.X, x) && between(u0.Y, u1.Y, y)
}
