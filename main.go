package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os/user"
	"sort"
	"strings"
	"time"

	"github.com/reujab/wallpaper"
)

type arguments struct {
	setLast bool
	path    string
}

type file struct {
	name    string
	modTime time.Time
}

type filesSlice []file

func main() {
	args := parseArgs()

	imgs, err := readDir(args.path)

	if err != nil {
		fmt.Println(args.path + " does not exist or parmission is denied =/")
		return
	}

	if len(imgs) == 0 {
		fmt.Println("No images found on: " + args.path)
		fmt.Println("¯\\_(ツ)_/¯")
		return
	}

	i := imgIndex(len(imgs), args.setLast)
	randImg := imgs[i].name

	err = wallpaper.SetFromFile(args.path + randImg)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Success!!! %s set as wallpaper.\n", randImg)
	}
}

func (s filesSlice) Len() int {
	return len(s)
}

func (s filesSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s filesSlice) Less(i, j int) bool {
	return s[i].modTime.Before(s[j].modTime)
}

func parseArgs() arguments {
	useLastPtr := flag.Bool("l", false, "use last image instead of a random one")
	flag.Parse()

	return arguments{
		setLast: *useLastPtr,
		path:    imgsPath(),
	}
}

func imgIndex(l int, getLastIndex bool) int {
	if getLastIndex {
		return l - 1
	}

	seed := rand.NewSource(time.Now().UnixNano())
	randWithSeed := rand.New(seed)

	return randWithSeed.Intn(l)
}

func imgsPath() string {
	wallsDir := homeDir() + "/Pictures/"

	if args := flag.Args(); len(args) > 0 {
		wallsDir = sanitizePath(args[0])
	}

	return wallsDir
}

func sanitizePath(path string) string {
	if strings.Contains(path, "~") {
		path = strings.Replace(path, "~", homeDir(), 1)
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func readDir(path string) (filesSlice, error) {
	var imgs filesSlice

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return filesSlice{}, err
	}

	for _, f := range files {
		imgs = append(imgs, file{f.Name(), f.ModTime()})
	}

	sort.Sort(imgs)

	return imgs, err
}
