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

type Args struct {
	setLast bool
	path    string
}

type File struct {
	name    string
	modTime time.Time
}

type FilesSlice []File

func (s FilesSlice) Len() int {
	return len(s)
}

func (s FilesSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FilesSlice) Less(i, j int) bool {
	return s[i].modTime.Before(s[j].modTime)
}

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

	i := getImgIndex(len(imgs), args.setLast)
	randImg := imgs[i].name

	err = wallpaper.SetFromFile(args.path + randImg)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Success!!! %s set as wallpaper.\n", randImg)
	}
}

func parseArgs() Args {
	useLastPtr := flag.Bool("l", false, "use last image instead of a random one")
	flag.Parse()

	return Args{
		setLast: *useLastPtr,
		path:    getImgsPath(),
	}
}

func getImgIndex(l int, shouldGetLastIndex bool) int {
	if shouldGetLastIndex {
		return l - 1
	}

	seed := rand.NewSource(time.Now().UnixNano())
	randWithSeed := rand.New(seed)

	return randWithSeed.Intn(l)
}

func getImgsPath() string {
	wallsDir := getHomeDir() + "/Pictures/"

	if args := flag.Args(); len(args) > 0 {
		wallsDir = sanitizePath(args[0])
	}

	return wallsDir
}

func sanitizePath(path string) string {
	if strings.Contains(path, "~") {
		path = strings.Replace(path, "~", getHomeDir(), 1)
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func readDir(path string) (FilesSlice, error) {
	var imgs FilesSlice

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return FilesSlice{}, err
	}

	for _, file := range files {
		imgs = append(imgs, File{file.Name(), file.ModTime()})
	}

	sort.Sort(imgs)

	return imgs, err
}
