package main

import (
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

	CooldownTime                = 1
	BulletLifeTime              = 1
	BulletSpeed                 = 1000
	TitleTextY                  = 250
	TitleTextValue              = "ASTEROIDS"
	PressSpaceTextValue         = "PRESS SPACE TO CONTINUE"
	PressSpaceTextY             = 700
	ControlsTextValue           = "CONTROLS: A, D - ROTATE; W - ACCELERATE; SPACEBAR - FIRE"
	ControlsTextY               = ScreenHeight - 20
	RecordTextY                 = ScreenHeight/2 + 150
	ScoreInGameY                = 20
	GameOverTextY               = ScreenHeight/2 - 150
	GameOverTextValue           = "GAME OVER"
	ScoreGameOverY              = GameOverTextY + 80
	NewRecordTextValue          = "NEW RECORD"
	NewRecordTextY              = ScoreGameOverY + 40
	TraceLifeTime               = 0.3
	TraceCooldownTime           = 0.08
	TraceDisappearanceSpeed     = 500
	TraceGrowthSpeed            = 100
	ExplosionDisappearanceSpeed = 300
	ExplosionGrowthSpeed        = 200
	ExplosionLifeTime           = 0.5
	EffectsDetalization         = 20
	MaxSpeed                    = 400
	TraceSpeed                  = MaxSpeed
	TraceDeltaY                 = 40
	SafeRange                   = 300
	SmashTime                   = 1
)

var RocketObj = VectorObject{
	Position:          Vec2{ScreenWidth / 2, ScreenHeight / 2},
	AccelerationValue: 100,
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

//key to cipher record file
var (
	SecretKey = []byte("A5A5TUPOOCHEN4A4")
)

func getCentredPosForText(txt string, f font.Face) int {
	W := text.BoundString(f, txt).Max.X
	X := ScreenWidth/2 - W/2
	return X
}
