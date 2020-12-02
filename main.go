package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/text"

	. "github.com/TuM0xA-S/goasteroids/vobj"
	"github.com/hajimehoshi/ebiten"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	EffectTrace = iota
	EffectExplosion
	EffectEnergyGained
)

//Game asteroids
type Game struct {
	rocket             VectorObject
	beginLast          time.Time
	asteroids          map[*VectorObject]bool
	bullets            map[*VectorObject]bool
	energyBlocks       map[*VectorObject]bool
	lifetime           map[*VectorObject]float64
	asteroidSpawnTimer float64
	cooldownTimer      float64
	mode               int
	score              float64
	record             int
	recordUpdated      bool
	effects            map[*VectorObject]int
	traceTimer         float64
	idFor              map[*VectorObject]int
	nextID             int
	smashTimer         float64
	energy             float64
}

func (g *Game) generateEnergyBlock(pos Vec2) {
	angle := rand.Float64() * 2 * math.Pi
	speed := Vec2{0, -1}
	speed.Rotate(angle)
	speed = speed.Mult(rand.Float64()*(EnergyBlockSpeedMax-EnergyBlockSpeedMin) + EnergyBlockSpeedMin)

	block := &VectorObject{
		Position: pos,
		Speed:    speed,
		Geometry: []Vec2{
			{-1, -1}, {1, -1}, {1, 1}, {-1, 1},
		},
		Scale: 20,
		Color: color.RGBA{0, 0, 255, 255},
	}
	g.lifetime[block] = EnergyBlockLifetime
	g.energyBlocks[block] = true
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

func (g *Game) generateRocketTrace() {
	v := MakeVector(g.rocket.Angle+math.Pi, TraceSpeed)
	v.SetLen(TraceDeltaY)
	t := &VectorObject{
		Geometry: generateCircle(0, EffectsDetalization),
		Speed:    MakeVector(g.rocket.Angle+math.Pi, TraceSpeed),
		Position: g.rocket.Position.Add(v),
		Color:    ColorYellow,
		Scale:    1,
	}
	g.lifetime[t] = TraceLifeTime
	g.effects[t] = EffectTrace
}

func (g *Game) generateEnergyBlockCathed(pos Vec2) {
	t := &VectorObject{
		Geometry: generateCircle(0, EffectsDetalization),
		Position: pos,
		Color:    color.RGBA{0, 0, 255, 255},
		Scale:    1,
	}
	g.lifetime[t] = EnergyGainedLifeTime
	g.effects[t] = EffectEnergyGained

}

func (g *Game) generateExplosion(pos Vec2) {
	t := &VectorObject{
		Geometry: generateCircle(0, EffectsDetalization),
		Position: pos,
		Color:    color.RGBA{255, 123, 0, 255},
		Scale:    1,
	}
	g.lifetime[t] = ExplosionLifeTime
	g.effects[t] = EffectExplosion
}

func (g *Game) processEffects(dt float64) {
	for k, v := range g.effects {
		k.Move(dt)
		var gs float64
		var ds float64
		if v == EffectTrace {
			gs = TraceGrowthSpeed
			ds = TraceDisappearanceSpeed
		} else if v == EffectExplosion {
			gs = ExplosionGrowthSpeed
			ds = ExplosionDisappearanceSpeed
		} else if v == EffectEnergyGained {
			gs = EnergyGainedGrowthSpeed
			ds = EnergyGainedDisappearanceSpeed
		}
		k.Scale += dt * gs
		clr := k.Color.(color.RGBA)
		untransp := float64(clr.A)
		untransp -= dt * ds
		if untransp < 0 {
			untransp = 0
		}
		clr.A = uint8(untransp)
		k.Color = clr
	}
}

func (g *Game) generateAsteroid() {
	ok := false
	var pos Vec2
	for !ok {
		pos = Vec2{}
		isHorizontal := rand.Intn(2)
		if isHorizontal == 1 {
			pos.X = float64(rand.Intn(ScreenWidth))
		} else {
			pos.Y = float64(rand.Intn(ScreenHeight))
		}
		ok = true
		for _, delta := range []Vec2{{0, 0}, {-MaxX, 0}, {MaxX, 0}, {0, MaxY}, {0, -MaxY}} {
			p := pos
			p.X += delta.X
			p.Y += delta.Y
			if math.Hypot(p.X-g.rocket.Position.X, p.Y-g.rocket.Position.Y) < SafeRange {
				ok = false
				break
			}
		}
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

	b := uint8(rand.Intn(101) + 100)
	asteroid := &VectorObject{
		Speed:    speed,
		Geometry: circle,
		Scale:    rand.Float64()*(AsteroidRadiusMax-AsteroidRadiusMin) + AsteroidRadiusMin,
		Color:    color.RGBA{b, b, b, 255},
		Position: pos,
	}
	g.asteroids[asteroid] = true
	g.idFor[asteroid] = g.nextID
	g.nextID++
}

func (g *Game) spawnBullet() {
	bullet := BulletObj
	bullet.Angle = g.rocket.Angle
	bullet.Speed = MakeVector(bullet.Angle, BulletSpeed)
	bullet.Position = g.rocket.Position
	g.bullets[&bullet] = true
	g.lifetime[&bullet] = BulletLifeTime
}

func (g *Game) startPlay() {
	g.asteroids = make(map[*VectorObject]bool)
	g.bullets = make(map[*VectorObject]bool)
	g.lifetime = make(map[*VectorObject]float64)
	g.effects = make(map[*VectorObject]int)
	g.energyBlocks = make(map[*VectorObject]bool)

	g.cooldownTimer = CooldownTime
	g.asteroidSpawnTimer = AsteroidSpawnTime / 2
	g.rocket = RocketObj
	g.score = 0
	g.recordUpdated = false
	g.nextID = 0
	g.smashTimer = SmashTime
	g.energy = 100
}

//Draw ....
func (g *Game) Draw(dst *ebiten.Image) {
	for k := range g.energyBlocks {
		k.Draw(dst)
	}
	asteroids := []*VectorObject{}
	for k := range g.asteroids {
		asteroids = append(asteroids, k)
	}
	sort.Slice(asteroids, func(i, j int) bool {
		return g.idFor[asteroids[i]] < g.idFor[asteroids[j]]
	})
	for _, a := range asteroids {
		a.Draw(dst)
	}
	for k := range g.bullets {
		k.Draw(dst)
	}
	for k := range g.effects {
		k.Draw(dst)
	}

	switch g.mode {
	case ModePlay, ModeSmash:
		if g.mode == ModePlay {
			g.rocket.Draw(dst)
		}
		scoreText := fmt.Sprint("SCORE: ", int(g.score))
		text.Draw(dst, scoreText, FontSmall, getCentredPosForText(scoreText, FontSmall), ScoreInGameY, ColorYellow)
		energyText := fmt.Sprint("ENERGY: ", int(g.energy))
		text.Draw(dst, energyText, FontSmall, getCentredPosForText(energyText, FontSmall), EnergyTextY, ColorYellow)
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
		effects:            map[*VectorObject]int{},
		idFor:              map[*VectorObject]int{},
		bullets:            map[*VectorObject]bool{},
		energyBlocks:       map[*VectorObject]bool{},
	}
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Asteroids beta v0.8")
	MaxX = ScreenWidth
	MaxY = ScreenHeight
	g := getGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	saveRecord(g.record)

}
