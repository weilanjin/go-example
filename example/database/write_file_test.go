package database

import "testing"

func TestSaveData(t *testing.T) {
	SaveData3("test.txt", []byte("test1"))
}
