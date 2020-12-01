package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	. "github.com/TuM0xA-S/goasteroids/vobj"
	"github.com/hajimehoshi/ebiten"
)

const (
	ScreenWidth       = 1000
	ScreenHeight      = 800
	AsteroidSpawnTime = 2
	AsteroidCountMax  = 5
	AsteroidSpeedMin  = 20
	AsteroidSpeedMax  = 50
	AsteroidRadiusMin = 50
	AsteroidRadiusMax = 120
	AsteroidRadiusGap = 0.15
	AsteroidVertexCount = 12
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Game asteroids
type Game struct {
	rocket             VectorObject
	beginLast          time.Time
	asteroids          map[*VectorObject]bool
	asteroidSpawnTimer float64
}


func generateCircle(deviation float64, cnt int) (circle []Vec2) {
	dangle := 2 * math.Pi / float64(cnt)
	v := Vec2{0, -1}
	for i := 0; i < cnt; i++ {
		r := 1 - deviation + rand.Float64() * 2 * deviation
		circle = append(circle, v.Mult(r))
		v.Rotate(dangle)
	}
	return
}

func (g *Game) generateAsteroid() {
	isHorizontal := rand.Intn(2)
	var pos Vec2
	if isHorizontal == 1 {
		pos.X = float64(rand.Intn(ScreenWidth))
	} else {
		pos.Y = float64(rand.Intn(ScreenHeight))
	}
	angle := rand.Float64() * 2 * math.Pi
	speed := Vec2{0, -1}
	speed.Rotate(angle)
	speed = speed.Mult(rand.Float64()*(AsteroidSpeedMax-AsteroidSpeedMin) + AsteroidSpeedMin)

	circle := generateCircle(AsteroidRadiusGap, AsteroidVertexCount)
	angle = rand.Float64() * 2 * math.Pi
	for _, v := range circle {
		v.Rotate(angle)
	}
	asteroid := &VectorObject{
		Speed: speed,
		Geometry: circle,
		Scale:    rand.Float64() * (AsteroidRadiusMax - AsteroidRadiusMin) + AsteroidRadiusMin,
		Color:    ColorGray,
		Position: pos,
	}
	g.asteroids[asteroid] = true
}

//Update ...
func (g *Game) Update() error {
	dt := time.Since(g.beginLast).Seconds()
	g.asteroidSpawnTimer += dt
	if g.asteroidSpawnTimer >= AsteroidSpawnTime {
		if len(g.asteroids) < AsteroidCountMax {
			g.generateAsteroid()
		}
		g.asteroidSpawnTimer = 0
	}
	g.beginLast = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.rocket.Rotate(dt, DirLeft)
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.rocket.Rotate(dt, DirRight)
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.rocket.Accelerate(dt)
	}

	g.rocket.Color = ColorGreen
	for k := range g.asteroids {
		if g.rocket.Collides(k) {
			g.rocket.Color = ColorRed
			break
		}
	}

	g.rocket.Move(dt)
	for k := range g.asteroids {
		k.Move(dt)
	}

	return nil
}

//Draw ....
func (g *Game) Draw(dst *ebiten.Image) {
	for k := range g.asteroids {
		k.Draw(dst)
	}
	g.rocket.Draw(dst)
}

//Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Asteroids")
	MaxX = ScreenWidth
	MaxY = ScreenHeight
	MaxSpeed = 400
	g := &Game{
		rocket: VectorObject{
			Position:          Vec2{ScreenWidth / 2, ScreenHeight / 2},
			AccelerationValue: 80,
			Geometry: []Vec2{
				{0, -2}, {1, 2}, {-1, 2},
			},
			Scale:       20,
			RotateSpeed: 1.2 * math.Pi,
			Color:       ColorGreen,
		},
		asteroids: make(map[*VectorObject]bool),
		beginLast: time.Now(),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
