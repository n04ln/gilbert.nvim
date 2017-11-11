package command

import (
	"github.com/NoahOrberg/gilbert/gist"
	"github.com/neovim/go-client/nvim"
)

func (g *Gilbert) GilbertList(v *nvim.Nvim, args []string) error {
	gs, err := gist.ListGists()
	if err != nil {
		return error
	}

	for _, g := range gs {
		for k, v := range g.Files {
			v.Command("echom '" + k + "'")
		}
	}
}
