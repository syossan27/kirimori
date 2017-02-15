# Kirimori（切り盛り）
[![Latest Version](https://img.shields.io/github/release/syossan27/kirimori.svg)](https://github.com/syossan27/kirimori/releases)
[![CircleCI Status](https://circleci.com/gh/syossan27/kirimori.svg?style=shield)](https://circleci.com/gh/syossan27/kirimori)
[![Coverage Status](https://coveralls.io/repos/github/syossan27/kirimori/badge.svg?branch=master)](https://coveralls.io/github/syossan27/kirimori?branch=master)

kirimori - Add vim plugin Tool.

![screenshot](screenshot.gif)

## Requirement
[golang/go](https://github.com/golang/go)

[urfave/cli](https://github.com/urfave/cli)

[mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)

[Masterminds/glide](https://github.com/Masterminds/glide)

[BurntSushi/toml](https://github.com/BurntSushi/toml)

[haya14busa/go-vimlparser](https://github.com/haya14busa/go-vimlparser)

[Songmu/prompter](https://github.com/Songmu/prompter)

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
