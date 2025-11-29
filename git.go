package main

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Branch struct {
	DisplayName string
	Branch      *plumbing.Reference
}

func GetBranchesIter(repo *git.Repository) (storer.ReferenceIter, error) {
	branches, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	return branches, nil
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

func GetBranchesMap(repo *git.Repository) (map[string]Branch, error) {
	branchesMap := make(map[string]Branch)
	branchesIter, err := GetBranchesIter(repo)
	if err != nil {
		return branchesMap, err
	}

	currentBranchName, err := getCurrentBranchName(repo)
	if err != nil {
		return branchesMap, err
	}

	branchesIter.ForEach(func(branch *plumbing.Reference) error {
		var displayBranchName string
		branchNameShort := branch.Name().Short()
		if currentBranchName == branchNameShort {
			displayBranchName = "* " + branchNameShort
		} else {
			displayBranchName = branchNameShort
		}

		branchesMap[branchNameShort] = Branch{
			DisplayName: displayBranchName,
			Branch:      branch,
		}
		return nil
	})

	return branchesMap, nil
}

func getCurrentBranchName(repo *git.Repository) (string, error) {
	currentBranch, err := repo.Head()
	if err != nil {
		return "", err
	}

	return currentBranch.Name().Short(), nil
}
