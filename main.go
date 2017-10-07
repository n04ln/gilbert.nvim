package main

import (
	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

func gilbertUp(v *nvim.Nvim, args []string) error {
	lines := [][]byte{[]byte("gilbert"), []byte("nvim")}
	buf, _ := v.CurrentBuffer()
	return v.SetBufferLines(buf, 0, -1, true, lines)
}

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "GilbertUp"}, gilbertUp)
		return nil
	})
}
