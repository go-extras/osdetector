//+build !windows

package osdetector

import "errors"

func (osd *OsDetector) GetWindowsDistro() (*Distro, error) {
	return nil, errors.New("unsupported distro: windows")
}