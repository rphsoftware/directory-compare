# directory-compare
A program to perform a full comparison of 2 directories and all files in them

This program is very simple, just recursively comparing all files in 2 directories and producing a log file.


To compile just type `go build main.go`

To run, run `./main <source directory name> <target directory name> <log file name>`

Source and target will be compared and a log will be produced containing all differences the program detects.

