package main

import (
	"image/color"
	"io"
	"log"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func plotTemperatures(temperatures []float64) (io.WriterTo, error) {
	rnd := rand.New(rand.NewSource(1))

	// xticks defines how we convert and display time.Time values.
	xticks := plot.TimeTicks{Format: "15:04:05"}

	// randomPoints returns some random x, y points
	// with some interesting kind of trend.
	randomPoints := func(n int) plotter.XYs {
		const (
			month = 1
			day   = 1
			hour  = 1
			min   = 1
			sec   = 1
			nsec  = 1
		)
		pts := make(plotter.XYs, n)
		for i := range pts {
			date := time.Date(2007+i, month, day, hour, min, sec, nsec, time.UTC).Unix()
			pts[i].X = float64(date)
			pts[i].Y = float64(pts[i].X+10*rnd.Float64()) * 1e-9
		}
		return pts
	}

	n := 10
	data := randomPoints(n)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Time Series"
	p.X.Tick.Marker = xticks
	p.Y.Label.Text = "Number of Gophers\n(Billions)"
	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = color.RGBA{G: 255, A: 255}
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, points)

	return p.WriterTo(10*vg.Centimeter, 5*vg.Centimeter, "png")
}
