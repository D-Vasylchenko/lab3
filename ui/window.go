package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title: pw.Title,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) drawDefaultUI() {
	// Заповнюємо фон зеленим кольором
	pw.w.Fill(pw.sz.Bounds(), color.RGBA{G: 0xff, A: 0xff}, draw.Src)

	// Визначаємо розміри фігури "Т" (максимальний розмір не більше половини вікна)
	shapeWidth := pw.sz.Bounds().Dx() / 2
	shapeHeight := pw.sz.Bounds().Dy() / 2
	verticalRectWidth := shapeWidth / 3
	horizontalRectHeight := shapeHeight / 3

	// Використовуємо оновлені координати для малювання "Т"
	centerX := pw.pos.Min.X
	centerY := pw.pos.Min.Y

	// Малюємо вертикальну частину "Т"
	verticalRect := image.Rect(centerX-verticalRectWidth/2, centerY-(shapeHeight/2), centerX+verticalRectWidth/2, centerY+(shapeHeight/2))
	pw.w.Fill(verticalRect, color.RGBA{R: 255, G: 255, B: 0, A: 255}, draw.Src)

	// Малюємо горизонтальну частину "Т"
	horizontalRect := image.Rect(centerX-(shapeWidth/2), centerY-(horizontalRectHeight/2)-shapeHeight/2, centerX+(shapeWidth/2), centerY+(horizontalRectHeight/2)-shapeHeight/2)
	pw.w.Fill(horizontalRect, color.RGBA{R: 255, G: 255, B: 0, A: 255}, draw.Src)

	// Малюємо білу рамку навколо вікна
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {
	case size.Event: // Window size change
		pw.sz = e
	case error:
		log.Printf("ERROR: %s", e)
	case mouse.Event:
		// Left mouse button click event
		if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
			// Оновлення позиції "Т" на основі координат миші
			pw.pos.Min.X = int(e.X) - pw.pos.Dx()/2 // Центруємо фігуру
			pw.pos.Min.Y = int(e.Y) - pw.pos.Dy()/2

			// Перемалюємо вікно, щоб відобразити нову позицію
			pw.w.Send(paint.Event{})
		}
	case paint.Event:
		// Малюємо оновлений контент у вікні
		if t == nil {
			// Якщо t == nil, це перша подія малювання, тоді малюємо фігуру в центрі
			pw.drawDefaultUI()
		} else {
			// Якщо є текстура, використовуємо її
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish() // Оновлюємо вікно
	}
}
