package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/text"

	. "github.com/TuM0xA-S/goasteroids/vobj"
	"github.com/hajimehoshi/ebiten"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Game asteroids
type Game struct {
	rocket             VectorObject
	beginLast          time.Time
	asteroids          map[*VectorObject]bool
	bullets            map[*VectorObject]bool
	bulletLifetime     map[*VectorObject]float64
	asteroidSpawnTimer float64
	cooldownTimer      float64
	mode               int
	score              float64
	record             int
	recordUpdated      bool
}

func generateCircle(deviation float64, cnt int) (circle []Vec2) {
	dangle := 2 * math.Pi / float64(cnt)
	v := Vec2{0, -1}
	for i := 0; i < cnt; i++ {
		r := 1 - deviation + rand.Float64()*2*deviation
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
		Speed:    speed,
		Geometry: circle,
		Scale:    rand.Float64()*(AsteroidRadiusMax-AsteroidRadiusMin) + AsteroidRadiusMin,
		Color:    ColorGray,
		Position: pos,
	}
	g.asteroids[asteroid] = true
}

func (g *Game) spawnBullet() {
	bullet := BulletObj
	bullet.Angle = g.rocket.Angle
	bullet.Speed = MakeVector(bullet.Angle, BulletSpeed)
	bullet.Position = g.rocket.Position
	g.bullets[&bullet] = true
	g.bulletLifetime[&bullet] = BulletLifeTime
}

func (g *Game) startPlay() {
	g.asteroids = make(map[*VectorObject]bool)
	g.bullets = make(map[*VectorObject]bool)
	g.bulletLifetime = make(map[*VectorObject]float64)
	g.cooldownTimer = CooldownTime
	g.asteroidSpawnTimer = AsteroidSpawnTime / 2
	g.rocket = RocketObj
	g.score = 0
	g.recordUpdated = false
}

//Draw ....
func (g *Game) Draw(dst *ebiten.Image) {
	for k := range g.asteroids {
		k.Draw(dst)
	}
	switch g.mode {
	case ModePlay:
		for k := range g.bullets {
			k.Draw(dst)
		}
		g.rocket.Draw(dst)
		scoreText := fmt.Sprint("SCORE: ", int(g.score))
		text.Draw(dst, scoreText, FontSmall, getCentredPosForText(scoreText, FontSmall), ScoreInGameY, ColorYellow)
	case ModeStart:
		text.Draw(dst, TitleTextValue, FontBig, TitleTextX, TitleTextY, ColorGold)
		text.Draw(dst, PressSpaceTextValue, FontSmall, PressSpaceTextX, PressSpaceTextY, ColorYellow)
		text.Draw(dst, ControlsTextValue, FontSmall, ControlsTextX, ControlsTextY, ColorOrange)
		recordText := fmt.Sprint("RECORD: ", g.record)
		text.Draw(dst, recordText, FontMedium, getCentredPosForText(recordText, FontMedium), RecordTextY, ColorGreen)
	case ModeGameOver:
		text.Draw(dst, GameOverTextValue, FontMedium, GameOverTextX, GameOverTextY, ColorRed)
		text.Draw(dst, PressSpaceTextValue, FontSmall, PressSpaceTextX, PressSpaceTextY, ColorYellow)
		scoreText := fmt.Sprint("SCORE: ", int(g.score))
		text.Draw(dst, scoreText, FontMedium, getCentredPosForText(scoreText, FontMedium), ScoreGameOverY, ColorYellow)

		if g.recordUpdated {
			text.Draw(dst, NewRecordTextValue, FontMedium, NewRecordTextX, NewRecordTextY, ColorGreen)
		}
	}
}

//Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func getGame() *Game {
	return &Game{
		asteroids:          make(map[*VectorObject]bool),
		beginLast:          time.Now(),
		mode:               ModeStart,
		record:             loadRecord(),
		asteroidSpawnTimer: AsteroidSpawnTime / 2,
	}
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Asteroids")
	MaxX = ScreenWidth
	MaxY = ScreenHeight
	MaxSpeed = 400
	g := getGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	saveRecord(g.record)

}
