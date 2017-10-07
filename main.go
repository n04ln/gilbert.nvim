package main

import (
	"github.com/NoahOrberg/gilbert/gist"
	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

func gilbertUp(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}
	filename, err := v.BufferName(buf)
	if err != nil {
		return err
	}
	err = gist.PostToGist("", filename)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "GilbertUp"}, gilbertUp)
		return nil
	})
}
