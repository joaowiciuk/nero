package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func onDialogDelete(dialog *gtk.Dialog) bool {
	dialog.Hide()
	return true
}

// onMainWindowDestory is the callback that is linked to the
// on_main_window_destroy handler. It is not required to map this,
// and is here to simply demo how to hook-up custom callbacks.
func onMainWindowDestroy() {
	log.Println("Leaving application. Good bye!")
}
