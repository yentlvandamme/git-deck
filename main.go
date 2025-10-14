package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	list := tview.NewList().ShowSecondaryText(false)
	list.SetChangedFunc(ListItemChanged)
	list.SetSelectedFunc(ListItemSelected)
	branches := GetBranchesIter()

	index := 1
	branches.ForEach(func(branch *plumbing.Reference) error {
		list.AddItem(branch.Name().Short(), "", rune(index), nil)
		index++
		return nil
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		printError(err)
	}
}

func ListItemChanged(index int, mainText string, secondaryText string, shortCut rune) {}

func ListItemSelected(index int, mainText string, secondaryText string, shortCut rune) {}

func GetBranchesIter() storer.ReferenceIter {
	currentPath, err := os.Getwd()
	if err != nil {
		printError(err)
	}

	repo, err := git.PlainOpen(currentPath)
	if err != nil {
		printError(err)
	}

	branches, err := repo.Branches()
	if err != nil {
		printError(err)
	}

	return branches
}

func printError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
