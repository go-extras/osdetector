# Go OS Detector Library

This library detects your OS Distribution version.

## Usage

```go
fs := afero.NewOsFs()
osd := osdetector.NewOsDetector(fs)
distro, err := osd.GetOSDistro(runtime.GOOS)
if err != nil {
		return err
}
log.Info("%s / %s / %s / %s / %s", distro.BasedOn, distro.Dist, distro.Os, distro.PseudoName, distro.Rev)
```

## Known limitations

At the moment it can only detect only major Linux distributions.
