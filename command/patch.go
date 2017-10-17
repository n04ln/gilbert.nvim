package command

import (
	"strings"

	"github.com/NoahOrberg/gilbert/gist"
	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

func (g *Gilbert) GilbertPatch(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	gistID, err := getGistIDFromBufferID(v, buf)
	if err != nil {
		return err
	}

	bufs, err := getBufferIDsFromGistID(v, gistID)
	if err != nil {
		return err
	}

	files := map[string]gist.File{}

	giloadedIds, err := getGistIDByGiLoaded(v)
	if err != nil {
		return err
	}

	for _, buf := range bufs {
		filename, err := v.BufferName(buf)
		if err != nil {
			return err
		}

		temp := strings.Split(filename, "/")
		if _, ok := giloadedIds[gistID]; ok {
			if temp[0] != gistID {
				continue
			}
		}

		filename = temp[len(temp)-1]

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

		files[filename] = gist.File{
			Content: content,
		}
	}

	gi := gist.Gist{
		Files: files,
	}

	res, err := gist.PatchGist(gistID, gi)
	if err != nil {
		return err
	}

	util.Echom(v, res.HTMLURL)

	if err := checkAndOpenGist(v, res.HTMLURL); err != nil {
		return err
	}

	if err := checkAndCopyGistURL(v, res.HTMLURL); err != nil {
		return err
	}

	if _, ok := giloadedIds[gistID]; ok {
		if err := deleteBuffersOfGistID(v, gistID); err != nil {
			return err
		}
	}

	return nil
}
