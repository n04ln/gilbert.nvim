package command

import (
	"errors"
	"strings"

	"github.com/NoahOrberg/gilbert/gist"
	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

const noName = "NoName"

type Gilbert struct {
}

func (g *Gilbert) GilbertPatch(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	filename, err := v.BufferName(buf)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	// :GiLoadでロードされる前提にするため、スラッシュから始まっていたら弾く
	temp := strings.Split(filename, "/")
	if (temp[0] == "") && (len(temp) != 2) {
		err := errors.New("didnt open :GiLoad this buffer")
		util.Echom(v, err.Error())
		return err
	}

	filename = temp[len(temp)-1]

	lines, err := v.BufferLines(buf, 0, -1, true)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	id := temp[0]
	gi := gist.Gist{
		Files: map[string]gist.File{
			filename: gist.File{
				Content: content,
			},
		},
	}

	err = gist.PatchGist(id, gi)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	return nil
}

func (g *Gilbert) GilbertLoad(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	id := args[0]

	gi, err := gist.GetGist(id)
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

	if err := v.SetBufferName(buf, id+"/"+filename); err != nil {
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
		for i, c := range lines {
			content += string(c)
			if i < len(lines)-1 {
				content += "\n"
			}
		}

		url, err = gist.PostToGistByContent("", filename, content)
		if err != nil {
			return err
		}

		splittedURL := strings.Split(url, "/")
		id := splittedURL[len(splittedURL)-1]

		err = v.SetBufferName(buf, id+"/"+filename)
	} else {
		url, err = gist.PostToGistByFile("", filename, false)
		if err != nil {
			return err
		}
	}

	if err := util.Echom(v, url); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	return nil
}
