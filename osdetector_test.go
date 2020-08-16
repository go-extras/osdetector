package osdetector

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OsDetectSuite struct {
	suite.Suite
}

func (s *OsDetectSuite) TestAmazonLinux2() {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/etc", 0755)
	e := afero.WriteFile(fs, "/etc/os-release", []byte(`NAME="Amazon Linux"
VERSION="2"
ID="amzn"
ID_LIKE="centos rhel fedora"
VERSION_ID="2"
PRETTY_NAME="Amazon Linux 2"
ANSI_COLOR="0;33"
CPE_NAME="cpe:2.3:o:amazon:amazon_linux:2"
HOME_URL="https://amazonlinux.com/"
`), 0644)
	_ = e

	osd := NewOsDetector(fs)
	res, err := osd.GetOSDistro("linux")
	s.Nil(err)
	s.Equal("Linux", res.Os)
	s.Equal("RedHat", res.BasedOn)
	s.Equal("Amazon Linux", res.Dist)
	s.Equal("2", res.Rev)
	s.Equal("", res.PseudoName)
}

func (s *OsDetectSuite) TestCento610() {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/etc", 0755)
	_ = afero.WriteFile(fs, "/etc/redhat-release", []byte(`CentOS release 6.10 (Final)
`), 0644)

	osd := NewOsDetector(fs)
	res, err := osd.GetOSDistro("linux")
	s.Nil(err)
	s.Equal("Linux", res.Os)
	s.Equal("RedHat", res.BasedOn)
	s.Equal("CentOS", res.Dist)
	s.Equal("6.10", res.Rev)
	s.Equal("Final", res.PseudoName)
}

func (s *OsDetectSuite) TestUbuntu1804() {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/etc", 0755)
	_ = afero.WriteFile(fs, "/etc/lsb-release", []byte(`DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=18.04
DISTRIB_CODENAME=bionic
DISTRIB_DESCRIPTION="Ubuntu 18.04.3 LTS"
`), 0644)

	osd := NewOsDetector(fs)
	res, err := osd.GetOSDistro("linux")
	s.Nil(err)
	s.Equal("Linux", res.Os)
	s.Equal("Debian", res.BasedOn)
	s.Equal("Ubuntu", res.Dist)
	s.Equal("18.04", res.Rev)
	s.Equal("bionic", res.PseudoName)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(OsDetectSuite))
}
