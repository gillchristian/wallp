# wallp

Change your wallpaper from the terminal.

## Install

```bash
$ go get github.com/gillchristian/wallp
```

## Use

```bash
$ wallp path/to/wallpapers/directory
```

Use `-l` to set the latest modified image on the directory as the wallpaper, instead of a random one.

```bash
$ wallp -l path/to/wallpapers/directory
```

If you don't provide a directory it will default to `~/Pictures`.
