package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rivo/tview"
)

type ListView struct {
	App      *tview.Application
	List     *tview.List
	Branches map[string]Branch
	Worktree *git.Worktree
}

func NewListView(app *tview.Application, branches map[string]Branch, worktree *git.Worktree) *ListView {
	list := tview.NewList()
	listView := ListView{
		App:      app,
		List:     list,
		Branches: branches,
		Worktree: worktree,
	}

	listView.List.ShowSecondaryText(false)

	index := 1
	for _, val := range branches {
		list.AddItem(val.DisplayName, "", rune(index), nil)
		index++
	}

	list.SetInputCapture(listView.SetListInputCaptures)
	list.SetSelectedFunc(listView.SetListSelectedHandler)

	return &listView
}

func (lv *ListView) SetListInputCaptures(event *tcell.EventKey) *tcell.EventKey {
	currentIndex := lv.List.GetCurrentItem()
	totalLength := len(lv.Branches)

	if event.Rune() == 'j' {
		if currentIndex < totalLength {
			lv.List.SetCurrentItem(currentIndex + 1)
		}
	}

	if event.Rune() == 'k' {
		if currentIndex > 0 {
			lv.List.SetCurrentItem(currentIndex - 1)
		}
	}

	return event
}

func (lv *ListView) SetListSelectedHandler(index int, mainText string, secondaryText string, shortCut rune) {
	selectedBranch, exists := lv.Branches[mainText]
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

	err := lv.Worktree.Checkout(&checkoutOpts)
	if err != nil {
		printError(err)
	}

	lv.App.Stop()
}
