// Package cmdio provides some basic functions for reading and
// outputting information in a terminal, such as getting user
// input, or providing progress bars.
package cmdio

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// ReadInput reads input from os.Stdin and outputs a string.
// It prints out a message beforehand, so that the end user
// knows what to put into the resultant input field.
func ReadInput(message string) string {
	if message != "" {
		fmt.Printf("%s: ", message)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

// ReadInputD wraps around ReadInput and provides a default
// if input == ""
func ReadInputD(message string, d string) string {
	m := ReadInput(message)
	if m == "" {
		return d
	}

	return m
}

// ReadInputReq wraps around ReadInput, and continues to loop
// until is not an empty string.
func ReadInputReq(message string) string {
	var i string
	for i == "" {
		i = ReadInput(message)
	}

	return i
}

// ReadInputAsArray wraps around ReadInput and turns a
// string seperated by a specific seperator into a
// string array
func ReadInputAsArray(message string, sep string) []string {
	return strings.Split(ReadInput(message+" [seperator: "+sep+"] "), sep)
}

// ReadInputAsArrayD wraps around ReadInputAsArray
// and outputs a default string array if the length
// of ReadInputAsArray is equal to zero.
func ReadInputAsArrayD(message string, sep string, d []string) []string {
	m := ReadInputAsArray(message, sep)
	if len(m) == 0 {
		return d
	}

	return m
}

// ReadInputAsBool wraps around ReadInput, and if
// returns if the input is equal to the given conditional.
func ReadInputAsBool(message string, cond string) bool {
	res := ReadInput(message + "[" + cond + "] ")
	if res == cond {
		return true
	}

	return false
}

// ReadInputAsInt wraps around ReadInput and returns
// an int based on the given input.
func ReadInputAsInt(message string) (int, error) {
	var i int
	var err error

	n := ReadInput(message)
	if n == "" {
		return 0, nil
	}

	i, err = strconv.Atoi(n)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// ReadInputAsIntD wraps around ReadInputAsInt and returns
// a default if ReadInputAsInt returns an error.
func ReadInputAsIntD(message string, d int) (int, error) {
	m, err := ReadInputAsInt(message)
	if err != nil {
		return d, err
	}

	return m, nil
}

// ReadInputAsFloat wraps around ReadInput, and returns a
// float64 from the input.
func ReadInputAsFloat(message string) (float64, error) {
	var f float64
	var err error

	n := ReadInput(message)
	if n == "" {
		return 0, nil
	}

	f, err = strconv.ParseFloat(n, 64)
	if err != nil {
		return 0, err
	}

	return float64(f), nil
}

// ReadInputAsFloatD wraps around ReadInputAsFloat, and returns
// a default if any errors occur, or if the result of
// ReadInputAsFloat is NaN.
func ReadInputAsFloatD(message string, d float64) (float64, error) {
	m, err := ReadInputAsFloat(message)
	if err != nil {
		return d, err
	}

	if math.IsNaN(m) {
		return d, nil
	}

	return m, nil
}
