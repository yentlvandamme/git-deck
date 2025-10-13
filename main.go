package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
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

	branches.ForEach(func(branch *plumbing.Reference) error {
		fmt.Println(branch.Name().Short())
		return nil
	})
}

func printError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
