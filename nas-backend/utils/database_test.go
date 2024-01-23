package utils

import "testing"

func TestDBAutoMigrate(t *testing.T) {
	InitDB()
	DBAutoMigrate(DB)
}
