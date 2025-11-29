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

func NewListView(app *tview.Application) (*ListView, error) {
	// This is a lot of Git stuff, all depending on the same information (repo).
	// We might as well just define this Git logic on a struct, so we don't have to pass
	// data like "repo" around everywhere.
	repo, err := GetRepo()
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	branchesMap, err := GetBranchesMap(repo)
	if err != nil {
		return nil, err
	}

	listView := ListView{
		App:      app,
		List:     tview.NewList(),
		Branches: branchesMap,
		Worktree: worktree,
	}

	listView.SetInitialListBranches()
	listView.List.ShowSecondaryText(false)
	listView.List.SetInputCapture(listView.SetListInputCaptures)
	listView.List.SetSelectedFunc(listView.SetListSelectedHandler)

	return &listView, nil
}

func (lv *ListView) SetInitialListBranches() {
	index := 1
	for _, val := range lv.Branches {
		lv.List.AddItem(val.DisplayName, "", rune(index), nil)
		index++
	}
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
