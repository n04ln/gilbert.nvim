package command

import (
	"github.com/NoahOrberg/gilbert/gist"
	"github.com/neovim/go-client/nvim"
)

type Gilbert struct {
}

func (g *Gilbert) GilbertUpload(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	filename, err := v.BufferName(buf)
	if err != nil {
		return err
	}

	var url string
	url = "Missing"
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
