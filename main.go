package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WIDTH  = 480
	HEIGHT = 600
)

const (
	MENU = iota
	PLAY
	GAME_OVER
)

type Player struct {
	img         [3]*ebiten.Image
	current_img int
	x, y        float64
	vx          float64
	ax          float64
}

func (p *Player) update() {

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && p.x < WIDTH-64 {
		p.vx += 5
		p.current_img = 2
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && p.x > 64 {
		p.vx -= 5
		p.current_img = 0
	} else {
		p.current_img = 1
	}

	p.ax = -p.vx * 0.3
	p.vx += p.ax
	p.x += p.vx
}

type Tree struct {
	img  *ebiten.Image
	x, y float64
}

type Flag struct {
	img  *ebiten.Image
	x, y float64
}

type Game struct {
	Player
	Tree
	Flag
	f     font.Face
	ft    font.Face
	score int
	bg    *ebiten.Image
	scene int
	trees []Tree
	flags []Flag
}

func (g *Game) Update() error {
	g.Player.update()

	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.scene == 0 {
		g.scene = 1
	}

	// get player img width and height for collision detection
	pw, ph := g.Player.img[g.Player.current_img].Size()

	// push tree jika kurang dari 5
	for len(g.trees) <= 6 {
		rand.Seed(time.Now().UnixNano())
		g.trees = append(g.trees, Tree{
			g.Tree.img,
			float64(rand.Intn(WIDTH-42) + 42),
			float64(rand.Intn(HEIGHT+100-HEIGHT+50) + HEIGHT + 50),
		})
	}

	if g.scene == 1 {
		// update posisi tree
		for i := 0; i < len(g.trees); i++ {
			g.trees[i].y -= 2
			if g.trees[i].y < -50 {
				g.trees = append(g.trees[:i], g.trees[i+1:]...)
			}
		}
	}

	// check collision antara player dan pohon
	for i := 0; i < len(g.trees); i++ {
		tw, th := g.trees[i].img.Size()
		if g.Player.x+float64(pw) >= g.trees[i].x &&
			g.Player.x <= g.trees[i].x+float64(tw) &&
			g.Player.y+float64(ph) >= g.trees[i].y &&
			g.Player.y <= g.trees[i].y+float64(th) {
			g.scene = 2
			g.trees = append(g.trees[:i], g.trees[i+1:]...)
		}
	}

	// push flag jika kurang dari 3
	for len(g.flags) <= 3 {
		rand.Seed(time.Now().UnixNano())
		g.flags = append(g.flags, Flag{
			g.Flag.img,
			float64(rand.Intn(WIDTH-42) + 42),
			float64(rand.Intn(HEIGHT+100-HEIGHT+50) + HEIGHT + 50),
		})
	}

	if g.scene == 1 {
		// update posisi flag
		for i := 0; i < len(g.flags); i++ {
			g.flags[i].y -= 2
			if g.flags[i].y < -50 {
				g.flags = append(g.flags[:i], g.flags[i+1:]...)
			}
		}
	}

	// check collision antara player dan flag
	for i := 0; i < len(g.flags); i++ {
		fw, fh := g.flags[i].img.Size()
		if g.Player.x+float64(pw) >= g.flags[i].x &&
			g.Player.x <= g.flags[i].x+float64(fw) &&
			g.Player.y+float64(ph) >= g.flags[i].y &&
			g.Player.y <= g.flags[i].y+float64(fh) {
			g.score += 10
			g.flags = append(g.flags[:i], g.flags[i+1:]...)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) && g.scene == 2 {
		g.scene, g.score = 1, 0
		g.trees, g.flags = nil, nil
		g.Player.x = WIDTH / 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.bg, &ebiten.DrawImageOptions{})
	switch g.scene {
	case MENU:
		text.Draw(screen, "In the Ski", g.ft, WIDTH/2-56, HEIGHT/4, color.Black)
		text.Draw(screen, "Tekan \"SPACE\" untuk play!", g.f, WIDTH/2-120, HEIGHT/4+80, color.Black)
		text.Draw(screen, "created by aji mustofa @pepega90", g.f, 80, HEIGHT-25, color.Black)
	case PLAY:
		// draw player
		pp := &ebiten.DrawImageOptions{}
		pp.GeoM.Translate(g.Player.x, g.Player.y)
		screen.DrawImage(g.Player.img[g.Player.current_img], pp)

		// draw pohon
		for i := 0; i < len(g.trees); i++ {
			tp := &ebiten.DrawImageOptions{}
			tp.GeoM.Translate(g.trees[i].x, g.trees[i].y)
			screen.DrawImage(g.trees[i].img, tp)
		}

		// draw flag
		for i := 0; i < len(g.flags); i++ {
			tp := &ebiten.DrawImageOptions{}
			tp.GeoM.Translate(g.flags[i].x, g.flags[i].y)
			screen.DrawImage(g.flags[i].img, tp)
		}

		// draw score
		score_text := "Score: " + strconv.Itoa(g.score)
		text.Draw(screen, score_text, g.f, 10, 30, color.Black)
	case GAME_OVER:
		text.Draw(screen, "Game Over", g.ft, WIDTH/2-56, HEIGHT/4, color.Black)
		text.Draw(screen, "Score Kamu: "+strconv.Itoa(g.score), g.f, WIDTH/2-56, HEIGHT/4+80, color.Black)
		text.Draw(screen, "Tekan \"R\" untuk restart", g.f, WIDTH/2-100, HEIGHT/2, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("In the Ski")

	// load font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	// load images player
	img, _, _ := ebitenutil.NewImageFromFile("./images/skier_forward.png")
	right_img, _, _ := ebitenutil.NewImageFromFile("./images/skier_right1.png")
	left_img, _, _ := ebitenutil.NewImageFromFile("./images/skier_left1.png")
	bg, _, _ := ebitenutil.NewImageFromFile("./images/bg.png")
	tree_img, _, _ := ebitenutil.NewImageFromFile("./images/tree.png")
	flag_img, _, _ := ebitenutil.NewImageFromFile("./images/flag.png")

	game := &Game{}
	// player
	game.Player.img = [3]*ebiten.Image{left_img, img, right_img}
	game.Player.current_img = 1
	game.Player.x = WIDTH / 2
	game.Player.y = 40

	// tree
	game.Tree.img = tree_img
	game.Flag.img = flag_img

	// etc
	game.bg = bg
	game.f, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	game.ft, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
