package osdetector

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// source: https://stackoverflow.com/questions/44363911/detect-windows-version-in-go
func (osd *OsDetector) GetWindowsDistro() (*Distro, error) {
	result := &Distro{
		Os:         "Windows",
		BasedOn:    "Windows NT",
		Dist:       "Windows",
		PseudoName: "",
	}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer k.Close()

	//cv, _, err := k.GetStringValue("CurrentVersion")
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Printf("CurrentVersion: %s\n", cv)

	pn , _, err := k.GetStringValue("ProductName")
	if err != nil {
		return nil, err
	}
	result.PseudoName = pn

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return nil, err
	}
	result.Rev = fmt.Sprintf("%d.", maj)

	min, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		return nil, err
	}
	result.Rev += fmt.Sprintf("%d.", min)

	cb, _, err := k.GetStringValue("CurrentBuild")
	if err != nil {
		return nil, err
	}
	result.Rev += cb

	return result, nil
}