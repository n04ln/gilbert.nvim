package main

import (
	"github.com/NoahOrberg/gilbert.nvim/src/command"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "GiUpload", NArgs: "*"}, command.GilbertUpload)
		return nil
	})
}
