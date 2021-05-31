package shared

// InitCore performs core dependencies initialization sequence.
func InitCore() {
	initEnv()
	initLogger()
	initLevelDB()
}

// CloseCore performs core dependencies close sequence.
func CloseCore() {
	closeLevelDB()
}
