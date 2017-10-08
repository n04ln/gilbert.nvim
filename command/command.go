package command

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/NoahOrberg/gilbert/gist"
	"github.com/neovim/go-client/nvim"
)

const noName = "NoName"

type Gilbert struct {
}

func (g *Gilbert) GilbertPatch(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	filename, err := v.BufferName(buf)
	if err != nil {
		return err
	}
	temp := strings.Split(filename, "/")
	filename = temp[len(temp)-1]
	if filename == "" {
		filename = noName
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

	if len(args) != 1 {
		return errors.New("invalid argument(need only one)")
	}
	id := args[0]
	gi := gist.Gist{
		Files: map[string]gist.File{
			filename: gist.File{
				Content: content,
			},
		},
	}
	hoge, err := json.Marshal(gi)
	v.Command("echom '" + filename + "'")
	v.Command("echom '" + string(hoge) + "'")
	err = gist.PatchGist(id, gi)
	if err != nil {
		v.Command("echo '" + err.Error() + "'")
		return err
	}

	return nil
}

func (g *Gilbert) GilbertLoad(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}
	gi, err := gist.GetGist(args[0])
	if err != nil {
		return err
	}

	if len(gi.Files) != 1 {
		return errors.New("can open only one file gist")
	}

	var filename string
	var strLines []string
	for key, value := range gi.Files {
		filename = key
		strLines = strings.Split(value.Content, "\n")
	}

	lines := make([][]byte, 0, len(strLines))
	for _, line := range strLines {
		lines = append(lines, []byte(line))
	}

	if err := v.SetBufferName(buf, filename); err != nil {
		return err
	}

	return v.SetBufferLines(buf, 0, -1, true, lines)
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
			filename = noName
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

	if err := v.Command("echom '" + url + "'"); err != nil {
		return err
	}

	return nil
}
