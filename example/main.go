package main

import (
	"fmt"
	"github.com/go-extras/osdetector"
	"github.com/spf13/afero"
	"os"
	"runtime"
)

func main() {
	fs := afero.NewOsFs()
	osd := osdetector.NewOsDetector(fs)
	distro, err := osd.GetOSDistro(runtime.GOOS)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("%s / %s / %s / %s / %s\n", distro.BasedOn, distro.Dist, distro.Os, distro.PseudoName, distro.Rev)
}
