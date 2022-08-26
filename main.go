package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadFile() (string, *os.File, error) {
	var AppendFile *os.File
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-F" {
			FileContents, err := os.ReadFile(string(os.Args[i+1]))
			if err != nil {
				panic(err)
			}
			AppendFile, err := os.OpenFile(string(os.Args[i+1]), os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}
			return string(FileContents), AppendFile, nil
		}
	}
	return "", AppendFile, errors.New("LoadFile: Specify a file with < -F ExampleFile >")
}

func SELECT(FileContents []string, table, index string) ([]string, error) {
	for i := 0; i < len(FileContents); i++ {
		if strings.Split(FileContents[i], " ")[0] == "TABLE" && strings.Split(FileContents[i], " ")[1] == table {
			values, err := strconv.Atoi(strings.Split(FileContents[i], " ")[2])
			if err != nil {
				panic(err)
			}
			for x := i + 1; x < len(FileContents); x++ {
				Line := strings.Split(FileContents[x], " ")
				if Line[0] == ";" {
					break
				}
				if Line[0] == index {
					return Line[1 : 1+values], nil
				}
			}
			break
		}
	}
	return []string{}, errors.New("SELECT: Cant find " + index + " in " + table)
}

func DELETE(FileContents []string, AppendFile *os.File, table, index string) {
	var (
		offset int
	)

	if err := os.Truncate(AppendFile.Name(), 0); err != nil {
		panic(err)
	}

	for i := 0; i < len(FileContents); i++ {
		AppendFile.WriteString(FileContents[i] + "\n")
		fmt.Println(FileContents[i])
		if strings.Split(FileContents[i], " ")[0] == "TABLE" && strings.Split(FileContents[i], " ")[1] == table {
			for x := i + 1; x < len(FileContents); x++ {
				Line := strings.Split(FileContents[x], " ")
				if Line[0] != index {
					AppendFile.WriteString(strings.Join(Line, " ") + "\n")
					offset = x
					fmt.Println(strings.Join(Line, " "))
					fmt.Println(offset)
				}
			}
			break
		}
	}
}

func ADD(FileContents []string, AppendFile *os.File, table, index string, values []string) {
	var (
		FoundTable bool = false
	)

	if err := os.Truncate(AppendFile.Name(), 0); err != nil {
		panic(err)
	}

	for i := 0; i < len(FileContents); i++ {
		if len(strings.Split(FileContents[i], " ")) == 3 {
			if strings.Split(FileContents[i], " ")[0] == "TABLE" && strings.Split(FileContents[i], " ")[1] == table {
				ValueAmount, err := strconv.Atoi(strings.Split(FileContents[i], " ")[2])
				if err != nil {
					panic(err)
				}
				if len(values) != ValueAmount {
					panic(errors.New(fmt.Sprintf("Given %d values when the table wants %d value/s", len(values), ValueAmount)))
				}
				FoundTable = true
			}
		}
		if FoundTable == true {
			AppendFile.WriteString(fmt.Sprintf("%s%s", FileContents[i], "\n"))
			AppendFile.WriteString(fmt.Sprintf("%s %s%s", index, strings.Join(values, " "), "\n"))
			FoundTable = false
		} else {
			AppendFile.WriteString(fmt.Sprintf("%s%s", FileContents[i], "\n"))
		}
	}
}

func main() {
	FileContents, AppendFile, err := LoadFile()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-SELECT" {
			r, err := SELECT(strings.Split(string(FileContents), "\n"), os.Args[i+1], os.Args[i+2])
			if err != nil {
				panic(err)
			}
			fmt.Println(r)
		}
		if os.Args[i] == "-DELETE" {
			DELETE(strings.Split(string(FileContents), "\n"), AppendFile, os.Args[i+1], os.Args[i+2])
		}
		if os.Args[i] == "-ADD" {
			ADD(strings.Split(string(FileContents), "\n"), AppendFile, os.Args[i+1], os.Args[i+2], os.Args[i+3:])
		}
	}
}
