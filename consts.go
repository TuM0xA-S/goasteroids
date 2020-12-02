package main

import (
	"image/color"
	"math"

	. "github.com/TuM0xA-S/goasteroids/vobj"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

//game modes
const (
	ModeStart = iota
	ModePlay
	ModeSmash
	ModeGameOver
)

const (
	ScreenWidth           = 1200
	ScreenHeight          = 800
	AsteroidSpawnTime     = 1
	AsteroidCountMax      = 6
	AsteroidSpeedMin      = 30
	AsteroidSpeedMax      = 180
	AsteroidRadiusMin     = 50
	AsteroidRadiusMax     = 120
	AsteroidRadiusGap     = 0.15
	AsteroidVertexCount   = 12
	AsteroidDestroyPoints = 10

	CooldownTime                   = 1.1
	BulletLifeTime                 = 1
	BulletSpeed                    = 1000
	TitleTextY                     = 250
	TitleTextValue                 = "ASTEROIDS"
	PressSpaceTextValue            = "PRESS SPACE TO CONTINUE"
	PressSpaceTextY                = 700
	ControlsTextValue              = "CONTROLS: A, D - ROTATE; W - ACCELERATE; SPACEBAR - FIRE"
	ControlsTextY                  = ScreenHeight - 20
	RecordTextY                    = ScreenHeight/2 + 150
	ScoreInGameY                   = 20
	GameOverTextY                  = ScreenHeight/2 - 150
	GameOverTextValue              = "GAME OVER"
	ScoreGameOverY                 = GameOverTextY + 80
	NewRecordTextValue             = "NEW RECORD"
	NewRecordTextY                 = ScoreGameOverY + 40
	TraceLifeTime                  = 0.4
	TraceCooldownTime              = 0.07
	TraceDisappearanceSpeed        = 500
	TraceGrowthSpeed               = 100
	ExplosionDisappearanceSpeed    = 300
	ExplosionGrowthSpeed           = 225
	ExplosionLifeTime              = 0.7
	EffectsDetalization            = 20
	MaxSpeed                       = 400
	TraceSpeed                     = MaxSpeed
	TraceDeltaY                    = 40
	SafeRange                      = 300
	SmashTime                      = 1.2
	EnergyBlockSpeedMin            = 20
	EnergyBlockSpeedMax            = 50
	EnergyBlockLifetime            = 10
	InitialEnergy                  = 100
	EnergyConsumptionSpeed         = 2
	EnergyPerBlock                 = 20
	EnergyTextY                    = ScreenHeight - 20
	EnergyGainedLifeTime           = 0.3
	EnergyGainedDisappearanceSpeed = 400
	EnergyGainedGrowthSpeed        = 200
	EnergyBlockSpawnRate           = 0.3
	PointsPerEnergyBlock           = 50
)

const (
	EffectTrace = iota
	EffectExplosion
	EffectEnergyGained
	EffectEnergyLost
)

var DefaultExplosion = ExplosionData{
	LifeTime:           ExplosionLifeTime,
	GrowthSpeed:        ExplosionGrowthSpeed,
	DisappearanceSpeed: ExplosionDisappearanceSpeed,
	Color:              color.RGBA{255, 123, 0, 255},
}

var EnergyGainedExplosion = ExplosionData{
	LifeTime:           EnergyGainedLifeTime,
	GrowthSpeed:        EnergyGainedGrowthSpeed,
	DisappearanceSpeed: EnergyGainedDisappearanceSpeed,
	Color:              color.RGBA{0, 0, 255, 255},
}

var TraceExplosion = ExplosionData{
	LifeTime:           TraceLifeTime,
	GrowthSpeed:        TraceGrowthSpeed,
	DisappearanceSpeed: TraceDisappearanceSpeed,
	Color:              ColorYellow,
}

var EnergyLostExplosion = ExplosionData{
	LifeTime:           EnergyGainedLifeTime,
	GrowthSpeed:        EnergyGainedGrowthSpeed,
	DisappearanceSpeed: EnergyGainedDisappearanceSpeed,
	Color:              color.RGBA{255, 123, 0, 255},
}

var RocketObj = VectorObject{
	Position:          Vec2{ScreenWidth / 2, ScreenHeight / 2},
	AccelerationValue: 140,
	Geometry: []Vec2{
		{0, -2}, {1, 2}, {-1, 2},
	},
	Scale:       23,
	RotateSpeed: 1.2 * math.Pi,
	Color:       ColorGreen,
	MaxSpeed:    MaxSpeed,
}

var BulletObj = VectorObject{
	Geometry: []Vec2{
		{0, -2}, {1, 2}, {-1, 2},
	},

	Scale: 8,
	Color: ColorRed,
}

//record consts
const (
	RecordFileName = ".asteroids_record"
)

var (
	PressSpaceTextX int
	TitleTextX      int
	ControlsTextX   int
	GameOverTextX   int
	NewRecordTextX  int
)

type ExplosionData struct {
	GrowthSpeed, DisappearanceSpeed, LifeTime float64
	Color                                     color.RGBA
}

var getExplosionData = map[int]ExplosionData{
	EffectEnergyGained: EnergyGainedExplosion,
	EffectExplosion:    DefaultExplosion,
	EffectTrace:        TraceExplosion,
	EffectEnergyLost:   EnergyLostExplosion,
}

//key to cipher record file
var (
	SecretKey = []byte("A4A4OCHENTUPO5A5")
)

func getCentredPosForText(txt string, f font.Face) int {
	W := text.BoundString(f, txt).Max.X
	X := ScreenWidth/2 - W/2
	return X
}
