package core

// InitCore performs core dependencies initialization sequence.
func InitCore() {
	initConfig()
	initLogger()
	initVault()
}
