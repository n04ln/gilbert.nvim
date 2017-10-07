package main

import (
	"github.com/NoahOrberg/gilbert.nvim/src/command"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	g := &command.Gilbert{}
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "GiUpload", NArgs: "?"}, g.GilbertUpload)
		return nil
	})
}
