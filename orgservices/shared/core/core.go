package core

// Init performs core dependencies initialization sequence.
func Init() {
	initConfig()
	initLogger()
	initMongoDB()
	initVault()
}
