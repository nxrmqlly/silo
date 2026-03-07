# silo
A dead simple, opinionated notes tui.


![silo demo gif](./docs/silo.gif)

## tips
1. try using a [nerd font](https://www.nerdfonts.com/)
2. make sure your terminal supports emojis (for now)
3. try running `silo help`!

## shortcuts index
1. `tab` switch between right and left pane
2. `ctrl+e` toggle autosave
3. `ctrl+s` (editor view) save file
4. `n` (sidebar) new file
5. `d` (sidebar) delete selected file
6. `r` (sidebar) rename selected file
7. `ctrl+x` render preview of selected file


## installing prebuilt binary

### Linux and MacOS

run this in your terminal:
```sh
curl -sSL https://raw.githubusercontent.com/nxrmqlly/silo/master/install.sh | sh
```
then launch:
```
silo
```

### Windows

1. download the latest .exe from [releases](https://github.com/nxrmqlly/silo/releases/latest)
2. add the .exe to your `PATH` 
3. run:
```sh
silo
```

## building from source
requirements:
- Go 1.26+ (only tested on this, older versions may work.)

```sh
git clone https://github.com/nxrmqlly/silo
cd silo
go build ./cmd/silo
```

## future ideas

[x] setup wizard
[ ] search / replace
[ ] recent files
[ ] tabs
[ ] configurable keybindings

## license

GPL v2

