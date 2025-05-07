package database

import "testing"

func TestSaveData(t *testing.T) {
	WriteFileAtomic("test.txt", []byte("test1"))
}