package osdetector

import (
	"errors"
	"github.com/spf13/afero"
	"github.com/go-extras/osdetector/internal/fileutils"
	"io/ioutil"
	"regexp"
	"strings"
)

type Distro struct {
	Os         string
	BasedOn    string
	Dist       string
	PseudoName string
	Rev        string
}

const (
	OsNameLinux = "Linux"
)

const (
	UnknownDistr = "unknown"
	UnknownRev   = "unknown"
)

const (
	OsBasedOnRedHat      = "RedHat"
	OsBasedOnSuSe        = "SuSe"
	OsBasedOnMandrake    = "Mandrake"
	OsBasedOnDebian      = "Debian"
	OsBasedOnContainerOS = "ContainerOS"
)

var (
	REOSRedHatDist       = regexp.MustCompile(`(.*) release.*`)
	REOSRedHatPseudoName = regexp.MustCompile(`\((.*)\)$`)
	REOSRedHatRev        = regexp.MustCompile(`.* release (.*) \(.*`)

	// REOSAmazonLinux2Dist = regexp.MustCompile(`VERSION_ID="(.*)(\n|\z)"`)

	REOSSuSePseudoName = regexp.MustCompile(`CODENAME = (.*)(\n|\z)`)
	REOSSuSeRev        = regexp.MustCompile(`VERSION = (.*)(\n|\z)`)

	REOSDebianDist       = regexp.MustCompile(`DISTRIB_ID=(.*)(\n|\z)`)
	REOSDebianRev        = regexp.MustCompile(`DISTRIB_RELEASE=(.*)(\n|\z)`)
	REOSDebianPseudoName = regexp.MustCompile(`DISTRIB_CODENAME=(.*)(\n|\z)`)
)

type OsDetector struct {
	fs afero.Fs
}

func NewOsDetector(fs afero.Fs) *OsDetector {
	return &OsDetector{
		fs: fs,
	}
}

func (osd *OsDetector) GetRedHatDistro() (*Distro, error) {
	distro := &Distro{
		Os:      OsNameLinux,
		BasedOn: OsBasedOnRedHat,
	}

	file, err := osd.fs.Open("/etc/redhat-release")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	redHatRelease := strings.TrimSpace(string(data))

	distro.Dist = REOSRedHatDist.ReplaceAllString(redHatRelease, `$1`)
	m := REOSRedHatPseudoName.FindStringSubmatch(redHatRelease)
	if len(m) > 0 {
		distro.PseudoName = m[1]
	} else {
		distro.PseudoName = UnknownDistr
	}

	m = REOSRedHatRev.FindStringSubmatch(redHatRelease)
	if len(m) > 0 {
		distro.Rev = m[1]
	} else {
		distro.Rev = UnknownRev
	}

	return distro, nil
}

func (osd *OsDetector) GetContainerOSorAmazonLinuxDistro() (*Distro, error) {
	distro := &Distro{
		Os: OsNameLinux,
	}

	data, err := afero.Afero{Fs: osd.fs}.ReadFile("/etc/os-release")
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "ID=cos" {
			distro.BasedOn = OsBasedOnContainerOS
			return distro, nil
		}
		if strings.HasPrefix(line, "PRETTY_NAME=\"Amazon Linux 2") {
			distro.BasedOn = OsBasedOnRedHat
			distro.Dist = "Amazon Linux"
			distro.Rev = "2"
			return distro, nil
		}
	}

	return distro, nil
}

func (osd *OsDetector) GetSuSeDistro() (*Distro, error) {
	distro := &Distro{
		Os:      OsNameLinux,
		BasedOn: OsBasedOnSuSe,
	}

	file, err := osd.fs.Open("/etc/SuSe-release")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	suseRelease := string(data)
	suSeInfo := strings.Split(suseRelease, "\n")
	distro.PseudoName = suSeInfo[0]

	m := REOSSuSePseudoName.FindStringSubmatch(suseRelease)
	if len(m) > 0 {
		distro.PseudoName = m[1]
	} else {
		distro.PseudoName = UnknownDistr
	}

	m = REOSSuSeRev.FindStringSubmatch(suseRelease)
	if len(m) > 0 {
		distro.Rev = m[1]
	} else {
		distro.Rev = UnknownRev
	}

	return distro, nil
}

func (osd *OsDetector) GetDebianDistro() (*Distro, error) {
	distro := &Distro{
		Os:      OsNameLinux,
		BasedOn: OsBasedOnDebian,
	}

	file, err := osd.fs.Open("/etc/lsb-release")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	debianRelease := string(data)

	m := REOSDebianDist.FindStringSubmatch(debianRelease)
	if len(m) > 0 {
		distro.Dist = m[1]
	} else {
		distro.Dist = UnknownDistr
	}

	m = REOSDebianPseudoName.FindStringSubmatch(debianRelease)
	if len(m) > 0 {
		distro.PseudoName = m[1]
	} else {
		distro.PseudoName = UnknownDistr
	}

	m = REOSDebianRev.FindStringSubmatch(debianRelease)
	if len(m) > 0 {
		distro.Rev = m[1]
	} else {
		distro.Rev = UnknownRev
	}

	return distro, nil
}

func (osd *OsDetector) GetLinuxDistro() (*Distro, error) {
	// RedHat
	if exists, err := fileutils.FileExists(osd.fs, "/etc/redhat-release"); exists {
		return osd.GetRedHatDistro()
	} else if err != nil {
		return nil, err
	}

	// SuSe
	if exists, err := fileutils.FileExists(osd.fs, "/etc/SuSE-release"); exists {
		return osd.GetSuSeDistro()
	} else if err != nil {
		return nil, err
	}

	// Mandrake
	if exists, err := fileutils.FileExists(osd.fs, "/etc/mandrake-release"); exists {
		panic("TODO: not supported yet")
	} else if err != nil {
		return nil, err
	}

	// Debian
	if exists, err := fileutils.FileExists(osd.fs, "/etc/lsb-release"); exists {
		return osd.GetDebianDistro()
	} else if err != nil {
		return nil, err
	}

	// ContainerOS / Amazon EC2 (should be the last!)
	if exists, err := fileutils.FileExists(osd.fs, "/etc/os-release"); exists {
		return osd.GetContainerOSorAmazonLinuxDistro()
	} else if err != nil {
		return nil, err
	}

	return nil, errors.New("unsupported distro")
}

func (osd *OsDetector) GetOSDistro(os string) (*Distro, error) {
	if os == "linux" {
		return osd.GetLinuxDistro()
	}

	return nil, errors.New("unsupported distro: " + os)
}
