package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"io"
	"log"
	"math/rand"
	"os"
)

type Game struct {
	background *ebiten.Image
	bgX        float64

	firetruck      *ebiten.Image
	truckX, truckY float64

	bullets   []*Bullet
	waterdrop *ebiten.Image

	flame   *ebiten.Image
	enemies []*Enemy

	score int
}

type Bullet struct {
	x, y float64
}

func (b *Bullet) Update() {
	b.x += 4
}

func (b *Bullet) Draw(screen *ebiten.Image, waterdropImg *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(waterdropImg, op)
}

type Enemy struct {
	x, y float64
}

func (e *Enemy) Update() {
	e.x -= 2
}

func (e *Enemy) Draw(screen *ebiten.Image, flameImg *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.x, e.y)
	screen.DrawImage(flameImg, op)
}

var audioContext *audio.Context
var audioPlayer *audio.Player

func loadSoundEffect() {
	audioContext = audio.NewContext(44100)

	f, err := os.Open("collision.mp3")
	if err != nil {
		log.Fatal(err)
	}

	// Read the entire file into memory
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		return
	} // Close the file after reading its content

	// Decode the MP3 data from memory
	soundStream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	audioPlayer, err = audioContext.NewPlayer(soundStream)
	if err != nil {
		log.Fatal(err)
	}
}

func playSoundEffect() {
	err := audioPlayer.Rewind()
	if err != nil {
		return
	}
	audioPlayer.Play()
}

func (g *Game) Update() error {
	g.bgX -= 2
	if g.bgX <= -800 {
		g.bgX = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.truckY -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.truckY += 2
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.bullets = append(g.bullets, &Bullet{x: g.truckX + 50, y: g.truckY + 25})
	}

	for i := 0; i < len(g.bullets); i++ {
		bullet := g.bullets[i]
		bullet.Update()

		for j := 0; j < len(g.enemies); j++ {
			enemy := g.enemies[j]
			if bullet.x < enemy.x+float64(g.flame.Bounds().Dx()) &&
				bullet.x+float64(g.waterdrop.Bounds().Dx()) > enemy.x &&
				bullet.y < enemy.y+float64(g.flame.Bounds().Dy()) &&
				bullet.y+float64(g.waterdrop.Bounds().Dy()) > enemy.y {
				// Collision detected
				playSoundEffect()
				g.score++
				g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
				g.enemies = append(g.enemies[:j], g.enemies[j+1:]...)
				i-- // Adjust index after removal
				break
			}
		}
	}

	if rand.Float64() < 0.02 {
		g.enemies = append(g.enemies, &Enemy{x: 800, y: rand.Float64() * 600})
	}
	for _, enemy := range g.enemies {
		enemy.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.bgX, 0)
	screen.DrawImage(g.background, op)
	op.GeoM.Translate(float64(g.background.Bounds().Dx()), 0)
	screen.DrawImage(g.background, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.truckX, g.truckY)
	screen.DrawImage(g.firetruck, op)

	for _, bullet := range g.bullets {
		bullet.Draw(screen, g.waterdrop)
	}

	for _, enemy := range g.enemies {
		enemy.Draw(screen, g.flame)
	}

	// Draw score
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func main() {
	backgroundImg, _, err := ebitenutil.NewImageFromFile("street.png")
	if err != nil {
		log.Fatalf("Failed to load the background image: %v", err)
	}

	truckImg, _, err := ebitenutil.NewImageFromFile("firetruck.png") // Corrected this line
	if err != nil {
		log.Fatalf("Failed to load the firetruck image: %v", err)
	}

	waterdropImg, _, err := ebitenutil.NewImageFromFile("waterdrop.png") // Corrected this line
	if err != nil {
		log.Fatalf("Failed to load the waterdrop image: %v", err)
	}

	flameImg, _, err := ebitenutil.NewImageFromFile("fire.png") // Corrected this line
	if err != nil {
		log.Fatalf("Failed to load the flame image: %v", err)
	}

	loadSoundEffect()

	game := &Game{
		background: backgroundImg,
		firetruck:  truckImg,
		truckX:     50,
		truckY:     250,
		waterdrop:  waterdropImg,
		flame:      flameImg,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Firefighter Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Failed to run the game: %v", err)
	}
}
