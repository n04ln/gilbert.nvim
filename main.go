package main

import (
	"github.com/NoahOrberg/gilbert.nvim/command"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	g := &command.Gilbert{}
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "GiUpload", NArgs: "?"}, g.GilbertUpload)
		p.HandleCommand(&plugin.CommandOptions{Name: "GiLoad", NArgs: "1"}, g.GilbertLoad)
		return nil
	})
}
