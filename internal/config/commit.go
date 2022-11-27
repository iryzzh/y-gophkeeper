//go:build !go1.18

package config

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func commit() string {
	dir, err := os.Getwd()
	if err != nil {
		return "N/A"
	}

	var r *git.Repository
	r, err = git.PlainOpen(dir)
	if err != nil {
		return "N/A"
	}

	var h *plumbing.Reference
	h, err = r.Head()
	if err != nil {
		return "N/A"
	}

	return h.Hash().String()
}
