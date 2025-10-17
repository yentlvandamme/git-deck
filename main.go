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
	repo, err := GetRepo()
	if err != nil {
		printError(err)
	}

	branches, err := GetBranchesIter(repo)
	if err != nil {
		printError(err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		printError(err)
	}

	branchesMap := make(map[string]*plumbing.Reference)
	app := tview.NewApplication()
	list := tview.NewList().ShowSecondaryText(false)
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		selectedBranch, exists := branchesMap[mainText]
		if !exists {
			return
		}

		checkoutOpts := git.CheckoutOptions{
			Hash:                      plumbing.ZeroHash,
			Branch:                    selectedBranch.Name(),
			Create:                    false,
			Force:                     false,
			Keep:                      true,
			SparseCheckoutDirectories: make([]string, 0),
		}
		worktree.Checkout(&checkoutOpts)
		app.Stop()
	})

	index := 1
	branches.ForEach(func(branch *plumbing.Reference) error {
		list.AddItem(branch.Name().Short(), "", rune(index), nil)
		branchesMap[branch.Name().Short()] = branch
		index++
		return nil
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		printError(err)
	}
}

func GetRepo() (*git.Repository, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(currentPath)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func GetBranchesIter(repo *git.Repository) (storer.ReferenceIter, error) {
	branches, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func printError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
