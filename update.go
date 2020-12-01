package main

import (
	. "github.com/TuM0xA-S/goasteroids/vobj"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"

)

//Update ...
func (g *Game) Update() error {
	dt := time.Since(g.beginLast).Seconds()
	g.beginLast = time.Now()
	g.asteroidSpawnTimer += dt
	if g.asteroidSpawnTimer >= AsteroidSpawnTime {
		if len(g.asteroids) < AsteroidCountMax {
			g.generateAsteroid()
		}
		g.asteroidSpawnTimer = 0
	}
	for k := range g.asteroids {
		k.Move(dt)
	}
	g.processEffects(dt)
	switch g.mode {
	case ModePlay:
		g.score += dt
		if g.cooldownTimer > 0 {
			g.cooldownTimer -= dt
		}
		if g.traceTimer > 0 {
			g.traceTimer -= dt
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.rocket.Rotate(dt, DirLeft)
		}

		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.rocket.Rotate(dt, DirRight)
		}

		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.rocket.Accelerate(dt)
			if g.traceTimer <= 0 {
				g.generateRocketTrace()
				g.traceTimer = TraceCooldownTime
			}
		}

		for k := range g.lifetime {
			g.lifetime[k] -= dt
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if g.cooldownTimer <= 0 {
				g.spawnBullet()
				g.cooldownTimer = CooldownTime
			}
		}

		toRemoveAsteroids := map[*VectorObject]bool{}
		toRemoveBullets := map[*VectorObject]bool{}
		for asteroid := range g.asteroids {
			if g.rocket.Collides(asteroid) {
				g.mode = ModeSmash
				g.generateExplosion(g.rocket.Position)
				return nil
			}
		}

		for k, v := range g.lifetime {
			if v <= 0 {
				toRemoveBullets[k] = true
			}
		}

		for b := range g.bullets {
			for a := range g.asteroids {
				if b.Collides(a) {
					g.score += AsteroidDestroyPoints
					toRemoveAsteroids[a] = true
					toRemoveBullets[b] = true
					g.generateExplosion(a.Position)
					break
				}
			}
		}

		for k := range toRemoveBullets {
			delete(g.bullets, k)
			delete(g.lifetime, k)
		}

		for k := range toRemoveAsteroids {
			delete(g.asteroids, k)
		}

		for k := range g.bullets {
			k.Move(dt)
		}
		g.rocket.Move(dt)

	case ModeStart:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.mode = ModePlay
			g.startPlay()
		}
	case ModeGameOver:
		if int(g.score) > g.record {
			g.record = int(g.score)
			g.recordUpdated = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.mode = ModeStart
		}

	case ModeSmash:
		g.smashTimer -= dt
		if g.smashTimer <= 0 {
			g.mode = ModeGameOver
		}

	}
	return nil
}
