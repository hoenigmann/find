package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type options struct {
	l     bool   // Follow symlinks.
	empty bool   // Only print files that are empty.
	name  string // Only print files matching pattern.
}

var opt *options

func main() {
	//test()
	opt = getOptions()

	nonFlagArgs := flag.Args()
	startpath := ""

	//fmt.Println("Non flag args count: ", len(nonFlagArgs))
	if len(nonFlagArgs) == 0 {
		startpath, _ = os.Getwd()
	} else if len(nonFlagArgs) == 1 {
		startpath = nonFlagArgs[0]
	}
	//fmt.Println("Current directory path: ", startpath)

	processDir(startpath, processFile, processLink)
}

func processDir(dirpath string, fileFunc func(string), linkFunc func(string)) {

	if opt.name != "" {

		return
	}

	if !opt.empty {
		fmt.Println(dirpath)
	}

	stat, err := os.Stat(dirpath)
	if err != nil {
		log.Println(err)
		return
	}
	if !stat.IsDir() {
		return
	}

	files, err := ioutil.ReadDir(dirpath)

	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		return
	}
	for _, f := range files {

		spot := f
		abspath := dirpath + "/" + spot.Name()
		if spot.IsDir() {
			processDir(abspath, processFile, processLink)
		} else if isSymLink(abspath) && opt.l {
			processLink(abspath)
		} else { //Regular file
			processFile(abspath)
		}
	}

}

func processFile(abspath string) {
	if opt.name != "" {
		return
	}
	if opt.empty {
		if isEmptyFile(abspath) {
			fmt.Println(abspath)
		}
	} else {
		fmt.Println(abspath)
	}

	// Hmm. need to combine the variety of options better.

}

func isEmptyFile(abspath string) bool {
	stat, err := os.Stat(abspath)
	if err != nil {
		log.Fatal(err)
	}

	return stat.Size() == 0
}

func processLink(abspath string) {
	// TODO: recursive symlink case. (symlink to a symlink), is readlink following all the way? probably not.
	realmaybe := symLinkPointee(abspath)
	// TODO: need to account for relative paths. realmaybe could be ../something etc.
	processFile(realmaybe)
}

func listDir(dirpath string) (paths []string) {
	stat, err := os.Stat(dirpath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		return []string{}
	}
	files, err := ioutil.ReadDir(dirpath)
	paths = make([]string, 0)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		//fmt.Println(f.Name())
		paths = append(paths, f.Name())
	}
	return paths
}

func isDir(abspath string) bool {
	stat, err := os.Stat(abspath)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func isSymLink(abspath string) bool {
	fi, err := os.Lstat(abspath)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0
}

func symLinkPointee(abspath string) string {
	real, err := os.Readlink(abspath)
	if err != nil {
		return ""
	}

	return real
}

func getOptions() *options {
	opt := new(options)
	flag.BoolVar(&opt.l, "L", false, "Follow symlinks")
	flag.BoolVar(&opt.empty, "empty", false, "Only print files that are empty")
	flag.StringVar(&opt.name, "name", "", "Only print files whose name matches the given pattern.")
	flag.Parse()

	return opt
}

func test() {
	fmt.Println(isSymLink("/Users/ahoenigmann/proj/go/src/github.com/hoenigmann/find/cmd/symltofind"))
	fmt.Println(symLinkPointee("/Users/ahoenigmann/proj/go/src/github.com/hoenigmann/find/cmd/symltofind"))
	fmt.Println(isDir(symLinkPointee("/Users/ahoenigmann/proj/go/src/github.com/hoenigmann/find/cmd/symltofind")))
	fmt.Println(isDir(symLinkPointee("/Users/ahoenigmann/proj/go/src/github.com/hoenigmann/find/cmd/sometextln")))
}
