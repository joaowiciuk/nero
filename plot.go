package main

import (
	"image/color"
	"io"
	"log"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func plotTemperatures(readings []reading) (io.WriterTo, error) {
	xticks := plot.TimeTicks{
		Format: "15:04 05s",
		Time:   plot.UnixTimeIn(time.Local),
	}

	data := make(plotter.XYs, len(readings))
	for i, reading := range readings {
		data[i].X = float64(reading.when.Unix())
		data[i].Y = reading.value
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.X.Tick.Marker = xticks
	p.Y.Min = 0
	p.Y.Max = 120
	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = color.RGBA{B: 255, A: 255}
	points.Shape = draw.CrossGlyph{}
	points.Color = color.RGBA{B: 255, A: 255}

	p.Add(line, points)

	return p.WriterTo(10*vg.Centimeter, 5*vg.Centimeter, "png")
}
