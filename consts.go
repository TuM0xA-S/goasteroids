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
	AsteroidSpawnTime     = 2
	AsteroidCountMax      = 6
	AsteroidSpeedMin      = 20
	AsteroidSpeedMax      = 50
	AsteroidRadiusMin     = 50
	AsteroidRadiusMax     = 120
	AsteroidRadiusGap     = 0.15
	AsteroidVertexCount   = 12
	AsteroidDestroyPoints = 5

	CooldownTime        = 0.8
	BulletLifeTime      = 1
	BulletSpeed         = 1000
	TitleTextY          = 250
	TitleTextValue      = "ASTEROIDS"
	PressSpaceTextValue = "PRESS SPACE TO CONTINUE"
	PressSpaceTextY     = 700
	ControlsTextValue   = "CONTROLS: A, D - ROTATE; W - ACCELERATE; SPACEBAR - FIRE"
	ControlsTextY       = ScreenHeight - 20
	RecordTextY         = ScreenHeight/2 + 150
	ScoreInGameY        = 20
	GameOverTextY       = ScreenHeight/2 - 150
	GameOverTextValue   = "GAME OVER"
	ScoreGameOverY      = GameOverTextY + 80
	NewRecordTextValue  = "NEW RECORD"
	NewRecordTextY      = ScoreGameOverY + 40
)

var RocketObj = VectorObject{
	Position:          Vec2{ScreenWidth / 2, ScreenHeight / 2},
	AccelerationValue: 80,
	Geometry: []Vec2{
		{0, -2}, {1, 2}, {-1, 2},
	},
	Scale:       20,
	RotateSpeed: 1.2 * math.Pi,
	Color:       ColorGreen,
}

var BulletObj = VectorObject{
	Geometry: []Vec2{
		{0, -2}, {1, 2}, {-1, 2},
	},

	Scale: 10,
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
