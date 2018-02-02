# Kirimori（切り盛り）
[![Latest Version](https://img.shields.io/github/release/syossan27/kirimori.svg)](https://github.com/syossan27/kirimori/releases)
[![CircleCI Status](https://circleci.com/gh/syossan27/kirimori.svg?style=shield)](https://circleci.com/gh/syossan27/kirimori)
[![Coverage Status](https://coveralls.io/repos/github/syossan27/kirimori/badge.svg?branch=master)](https://coveralls.io/github/syossan27/kirimori?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/syossan27/kirimori)](https://goreportcard.com/report/github.com/syossan27/kirimori)

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

or download file from releases page.

https://github.com/syossan27/kirimori/releases

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

- Search vim plugin

```
$ kirimori search Shougo
Total: 52

Plugin Name                        Short Description
shougo/unite.vim                   🐉  Unite and create user interfaces
shougo/vimproc.vim                 Interactive command execution in Vim.
shougo/neocomplete.vim             Next generation completion framework after neocomplcache
shougo/neosnippet-snippets         The standard snippets repository for neosnippet
shougo/neomru.vim                  MRU plugin includes unite.vim/denite.nvim MRU sources
shougo/neosnippet.vim              neo-snippet plugin
shougo/neocomplcache.vim           Ultimate auto-completion system for Vim.
shougo/vimshell.vim                🐚  Powerful shell implemented by vim.
shougo/unite-outline               outline source for unite.vim
shougo/vimfiler.vim                📁  Powerful file explorer implemented by Vim script
shougo/deoplete.nvim               🌠  Dark powered asynchronous completion framework for neovim/Vim8
shougo/neoyank.vim                 Saves yank history includes unite.vim history/yank source.
shougo/neoinclude.vim              Include completion framework for neocomplete/deoplete
shougo/context_filetype.vim        Context filetype library for Vim script
shougo/unite-help                  help source for unite.vim
shougo/vinarise                    Ultimate hex editing system with Vim
shougo/vinarise.vim                Ultimate hex editing system with Vim
shougo/unite-session               unite.vim session source
shougo/junkfile.vim                Create temporary file for memo, testing, ...
shougo/echodoc                     Print documents in echo area.
```

## Contributing

You're most welcomed!💓

Welcome pull request and issues.✨

## License

[MIT](https://github.com/tcnksm/tool/blob/master/LICENCE)

## Author

[syossan27](https://github.com/syossan27)
