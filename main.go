package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

var filesSource map[string]string
var filesTarget map[string]string

func recursivelyScanDirectory(logBase string, pathBase string, path string, maptouse map[string]string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(files); i++ {

		logMsg := logBase + "> " + strconv.Itoa(i + 1) + "/" + strconv.Itoa(len(files)) + " | " + files[i].Name() + " "
		fmt.Print("\u001B[2K\r", logMsg)

		if files[i].IsDir() {
			recursivelyScanDirectory(logMsg, pathBase + files[i].Name() + "/", path + "/" + files[i].Name(), maptouse )
		} else {
			file, err := os.Open(path + "/" + files[i].Name())
			if err != nil {
				panic(err)
			}

			hash := sha256.New()
			if _, err := io.Copy(hash, file); err != nil {
				panic(err)
			}

			result := hex.EncodeToString(hash.Sum(nil))

			file.Close()
			maptouse[pathBase + files[i].Name()] = result
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("./compare <dir1> <dir2> <log file>")
	}

	filesSource = make(map[string]string)
	filesTarget = make(map[string]string)

	fmt.Println("Scanning first directory")
	recursivelyScanDirectory("", "/", os.Args[1], filesSource)

	fmt.Println("\nScanning second directory")
	recursivelyScanDirectory("", "/", os.Args[2], filesTarget)

	fmt.Println("Comparing first directory...")

	log, err := os.Create(os.Args[3])
	if err != nil {
		panic(err)
	}

	done := 0
	for k, v := range filesSource {
		done++
		fmt.Print("\u001B[2K\r", done, " / ", len(filesSource))

		if _, ok := filesTarget[k]; ok != true {
			log.WriteString("====\n-> File does not exist\nSource: ")
			log.WriteString(k)
			log.WriteString("\nTarget: Does not exist\n\n")
			continue
		}

		if h := filesTarget[k]; h != v {
			log.WriteString("====\n-> Hashes differ\nSource: ")
			log.WriteString(k)
			log.WriteString(" (hash: ")
			log.WriteString(v)
			log.WriteString(" )\nTarget: ")
			log.WriteString(k)
			log.WriteString(" (hash: ")
			log.WriteString(h)
			log.WriteString(" )\n\n")
			continue
		}
	}

	done = 0

	fmt.Println("\nComparing second directory")
	for k, _ := range filesTarget {
		done++
		fmt.Print("\u001B[2K\r", done, " / ", len(filesTarget))

		if _, ok := filesSource[k]; ok != true {
			log.WriteString("====\n-> File does not exist\nSource: Does not exist")
			log.WriteString("\nTarget: ")
			log.WriteString(k)
			log.WriteString("\n\n")
			continue
		}
	}
	fmt.Println("\nAll done")

	log.Close()
}
