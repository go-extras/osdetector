package fileutils

import (
	"github.com/spf13/afero"
	"os"
)

func FileExists(fs afero.Fs, fname string) (bool, error) {
	if _, err := fs.Stat(fname); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
