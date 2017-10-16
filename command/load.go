package command

import (
	"strings"

	"github.com/NoahOrberg/gilbert/gist"
	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

func (g *Gilbert) GilbertLoad(v *nvim.Nvim, args []string) error {
	temp := strings.Split(args[0], "/")
	id := temp[len(temp)-1]

	gi, err := gist.GetGist(id)
	if err != nil {
		return err
	}

	for filename, file := range gi.Files {
		if err := util.NewBuffer(v); err != nil {
			return err
		}

		buf, err := v.CurrentBuffer()
		if err != nil {
			return err
		}

		var strLines []string
		strLines = strings.Split(file.Content, "\n")

		lines := make([][]byte, 0, len(strLines))
		for _, line := range strLines {
			lines = append(lines, []byte(line))
		}

		if err := setGistIDFromBufferID(v, buf, id); err != nil {
			return err
		}

		if err := v.SetBufferName(buf, id+"/"+filename); err != nil {
			return err
		}

		if err := v.SetBufferLines(buf, 0, -1, true, lines); err != nil {
			return err
		}
	}
	return nil
}
