package inject

// inject/dev.go
var devMode bool

func SetDevMode(dev bool) {
	devMode = dev
}

func IsDevMode() bool {
	return devMode
}
