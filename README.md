# Kirimori（切り盛り）
kirimori - Add vim plugin Tool.

![screenshot](screenshot.gif)

## Requirement
[golang/go](https://github.com/golang/go)

[urfave/cli](https://github.com/urfave/cli)

[mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)

[Masterminds/glide](https://github.com/Masterminds/glide)

[BurntSushi/toml](https://github.com/BurntSushi/toml)

## Install
````
$ go get github.com/syossan27/kirimori
````

## Usage

- Create setting file.

```
$ kirimori init
Type your .vimrc path. (default: ~/.vimrc)
> ~/.vimrc
Choose a your vim bundle plugin. (default: 1)
	1) Vundle
	2) NeoBundle
	3) dein.vim
	4) vim-plug
Type number > 1
Success: Create setting file.
```

- Edit setting file.

Open config file by any editor (default: vim)

```
$ kirimori config
```

- Add vim plugin

```
$ kirimori add Shougo/neocomplete
Success: Add plugin.
```

- Remove vim plugin

```
$ kirimori remove Shougo/neocomplete
Success: Remove plugin.
```

- List vim plugin

```
$ kirimori list
Shougo/neocomplete
Shougo/unite.vim
Shougo/neosnippet.vim
```

## Contributing

You're most welcomed!💓
Welcome pull request and issues.✨

## License

[MIT](https://github.com/tcnksm/tool/blob/master/LICENCE)

## Author

[syossan27](https://github.com/syossan27)
