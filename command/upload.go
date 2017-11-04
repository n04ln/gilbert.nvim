package command

import (
	"errors"
	"strings"

	wsconfig "github.com/NoahOrberg/gilbert/config"
	"github.com/NoahOrberg/gilbert/gist"
	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

func (g *Gilbert) GilbertUpload(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	filepath, err := v.BufferName(buf)
	if err != nil {
		return err
	}
	var fileExists bool

	var url string
	var filename string
	url = "Missing"
	if filepath == "" {
		fileExists = false
		if len(args) > 0 {
			filename = args[0]
			filepath = args[0]
		} else {
			filepath = noName
			filename = noName
		}
	} else {
		fileExists = true
		temp := strings.Split(filepath, "/")
		filename = temp[len(temp)-1]
	}

	lines, err := v.BufferLines(buf, 0, -1, true)
	if err != nil {
		return err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	if !fileExists && content == "" {
		return errors.New("this is empty buffer")
	}

	url, err = gist.PostToGistByContent("", filename, content)
	if err != nil {
		return err
	}

	util.Echom(v, url)

	if err := checkAndOpenGist(v, url); err != nil {
		return err
	}

	if err := checkAndCopyGistURL(v, url); err != nil {
		return err
	}

	splittedURL := strings.Split(url, "/")
	id := splittedURL[len(splittedURL)-1]
	if err := setGistIDFromBufferID(v, buf, id); err != nil {
		return err
	}

	if !fileExists {
		if err := v.SetBufferName(buf, id+"/"+filename); err != nil {
			return err
		}

		if err := saveFileInCurrentBufferWithName(v, wsconfig.GetConfig().Workspace+"/"+id+"/"+filename); err != nil {
			return err
		}
	}

	return nil
}
