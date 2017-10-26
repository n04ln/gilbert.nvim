package command

import (
	"strings"

	wsconfig "github.com/NoahOrberg/gilbert/config"
	"github.com/NoahOrberg/gilbert/gist"
	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

func (g *Gilbert) GilbertLoad(v *nvim.Nvim, args []string) error {
	temp := strings.Split(args[0], "/")
	id := temp[len(temp)-1]

	gi, err := gist.GetGistAndSave(id)
	if err != nil {
		return err
	}

	wsconf := wsconfig.GetConfig()

	for filename, _ := range gi.Files {
		if err := util.NewBuffer(v); err != nil {
			return err
		}

		buf, err := v.CurrentBuffer()
		if err != nil {
			return err
		}

		gistfile := wsconf.Workspace + "/" + id + "/" + filename

		if err := openFileInCurrentBuffer(v, gistfile); err != nil {
			return err
		}

		if err := setGistIDFromBufferID(v, buf, id); err != nil {
			return err
		}

		if err := setGistIDToGiLoaded(v, id); err != nil {
			return err
		}

		if err := clearUndo(v); err != nil {
			return err
		}
	}
	return nil
}
