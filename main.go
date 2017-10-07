package main

import (
	"github.com/NoahOrberg/gilbert/gist"
	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

func gilbertUpload(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	filename, err := v.BufferName(buf)
	if err != nil {
		return err
	}

	var url string
	if filename == "" {
		if len(args) > 0 {
			filename = args[0]
		} else {
			filename = "NoName"
		}
		lines, err := v.BufferLines(buf, 0, -1, true)
		if err != nil {
			return err
		}

		var content string
		for _, c := range lines {
			content += string(c)
			content += "\n"
		}

		url, err = gist.PostToGistByContent("", filename, content)
		if err != nil {
			return err
		}
	} else {
		url, err = gist.PostToGistByFile("", filename, false)
		if err != nil {
			return err
		}
	}

	v.WriteOut(url)

	return nil
}

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "GiUpload", NArgs: "*"}, gilbertUpload)
		return nil
	})
}
