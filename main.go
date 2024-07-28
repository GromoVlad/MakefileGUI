package main

import (
	"go_gui/internal/service"
)

func main() {
	service.Bootstrap()
	window, myApp := service.CreateWindowAndApp()
	label := service.CreateLabel()
	window = service.AddContent(label, window, myApp)
	service.StartingListeners(label)
	window.ShowAndRun()
}
