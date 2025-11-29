package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rivo/tview"
)

type ListView struct {
	App  *tview.Application
	List *tview.List
}

func NewListView(app *tview.Application, branches map[string]Branch, worktree *git.Worktree) *ListView {
	list := tview.NewList()
	list.ShowSecondaryText(false)

	index := 1
	for _, val := range branches {
		list.AddItem(val.DisplayName, "", rune(index), nil)
		index++
	}

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		selectedBranch, exists := branches[mainText]
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

	return &ListView{
		App:  app,
		List: list,
	}
}
