package command

import (
	"errors"
	"os"
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

	wsconf := wsconfig.GetConfig()

	if !fileExists {
		if err := v.SetBufferName(buf, wsconf.Workspace+"/"+id+"/"+filename); err != nil {
			return err
		}

		dir := wsconf.Workspace + "/" + id
		if _, err := os.Stat(dir); err != nil {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		if err := saveFileInCurrentBufferWithName(v, dir+"/"+filename); err != nil {
			return err
		}
		// NOTE: I don't know best practice ;(
		//       re-open file in `~/.gilbert/<gist_id>/` because to active syntax highlight
		if err := openFileInCurrentBuffer(v, dir+"/"+filename); err != nil {
			return err
		}
		// NOTE: patch file because maybe(e.g. Go-file) file is formatted when it is saved.
		if err := reuploadFile(v, buf, id, filename); err != nil {
			return err
		}
	}

	util.Echom(v, url)

	return nil
}

func reuploadFile(v *nvim.Nvim, buf nvim.Buffer, id, filename string) error {
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

	files := map[string]gist.File{}
	files[filename] = gist.File{
		Content: content,
	}

	gi := gist.Gist{
		Files: files,
	}

	_, err = gist.PatchGist(id, gi)
	if err != nil {
		return err
	}

	return nil
}
