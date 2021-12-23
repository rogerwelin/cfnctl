package utils

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

// ReadFile is a utility function for reading files
func ReadFile(path string) ([]byte, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// ReturnRandom returns a ranomized string with the length based on input
func ReturnRandom(value int) string {
	stringArr := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"}
	newString := ""

	for i := 0; i <= value; i++ {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		randIndex := r1.Intn(len(stringArr))
		newString = newString + stringArr[randIndex]
	}
	return newString
}

// TrimFileSuffix trims the suffix on the file
func TrimFileSuffix(path string) string {
	file := filepath.Base(path)
	return strings.TrimSuffix(file, filepath.Ext(file))
}
