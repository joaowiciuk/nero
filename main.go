package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	util "github.com/joaowiciuk/gotk-util"
)

const appID = "com.github.joaowiciuk.nero"

var application *gtk.Application
var builder *gtk.Builder
var mainWindow *gtk.ApplicationWindow

func main() {

	// Create a new application.
	var err error
	application, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	errorCheck(err)

	// Connect function to application startup event, this is not required.
	application.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	application.Connect("activate", func() {
		log.Println("application activate")

		// Get the GtkBuilder UI definition in the glade file.
		builder, err = gtk.BuilderNewFromFile("resources/ui/main.glade")
		errorCheck(err)

		// Get the object with the id of "main_window".
		obj, err := builder.GetObject("main_window")
		errorCheck(err)
		// Verify that the object is a pointer to a gtk.ApplicationWindow.
		mainWindow, err = util.IsWindow(obj)
		errorCheck(err)

		// Map the handlers to callback functions, and connect the signals
		// to the Builder.
		signals := map[string]interface{}{
			"on_main_window_destroy": onMainWindowDestroy,
		}
		builder.ConnectSignals(signals)

		obj, err = builder.GetObject("label_temperature")
		errorCheck(err)
		labelTemperature, err := util.IsLabel(obj)
		errorCheck(err)

		/* obj, err = builder.GetObject("drawing_area")
		errorCheck(err)
		drawingArea, err := util.IsDrawingArea(obj)
		errorCheck(err) */

		temperatures := make([]float64, 0)
		temperatureChan := watchTemperature("sensors", `Core 0:\ +(\+.*?)°C`)
		errorCheck(err)
		go func() {
			for {
				select {
				case temperature := <-temperatureChan:
					temperatures = append(temperatures, temperature)
					glib.IdleAdd(labelTemperature.SetText, fmt.Sprintf("%.2f° C", temperature))
					w, err := plotTemperatures(temperatures)
					errorCheck(err)

					var buffer bytes.Buffer
					n, err := w.WriteTo(&buffer)
					errorCheck(err)
					log.Printf("%d bytes writen to buffer\n", n)
				}
			}
		}()

		// Show the Window and all of its components.
		//mainWindow.Maximize()
		mainWindow.Show()
		application.AddWindow(mainWindow)
	})

	// Connect function to application shutdown event, this is not required.
	application.Connect("shutdown", func() {
		log.Println("Application finished.")
	})

	// Launch the application
	os.Exit(application.Run(os.Args))
}
