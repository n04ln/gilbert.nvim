package command

import (
	"github.com/NoahOrberg/gilbert/gist"
	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

func (g *Gilbert) GilbertList(v *nvim.Nvim, args []string) error {
	gs, err := gist.ListGists()
	if err != nil {
		return err
	}

	for _, g := range gs {
		util.Echom(v, g.URL)
		for k, _ := range g.Files {
			util.Echom(v, "  "+k)
		}
	}

	return nil
}
