package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	listView, err := NewListView(app)
	if err != nil {
		printError(err)
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(listView.List, true).Run(); err != nil {
		printError(err)
	}
}

func printError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
