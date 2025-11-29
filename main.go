package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rivo/tview"
)

func main() {
	repo, err := GetRepo()
	if err != nil {
		printError(err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		printError(err)
	}

	// Create the application and its list component
	app := tview.NewApplication()
	list := tview.NewList().ShowSecondaryText(false)

	branchesMap, err := GetBranchesMap(repo)
	if err != nil {
		printError(err)
	}

	index := 1
	for _, val := range branchesMap {
		list.AddItem(val.DisplayName, "", rune(index), nil)
		index++

	}

	// Select handler
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		selectedBranch, exists := branchesMap[mainText]
		if !exists {
			return
		}

		checkoutOpts := git.CheckoutOptions{
			Hash:                      plumbing.ZeroHash,
			Branch:                    selectedBranch.Branch.Name(),
			Create:                    false,
			Force:                     false,
			Keep:                      true,
			SparseCheckoutDirectories: make([]string, 0),
		}

		err := worktree.Checkout(&checkoutOpts)
		if err != nil {
			printError(err)
		}

		app.Stop()
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentIndex := list.GetCurrentItem()

		if event.Rune() == 'j' {
			if currentIndex < index {
				list.SetCurrentItem(currentIndex + 1)
			}
		}

		if event.Rune() == 'k' {
			if currentIndex > 0 {
				list.SetCurrentItem(currentIndex - 1)
			}
		}

		return event
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		printError(err)
	}
}

func printError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
