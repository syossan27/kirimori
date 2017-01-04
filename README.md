# Kirimori（切り盛り）
kirimori - Add vim plugin Tool.

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
Success: Create setting file.
```

- Edit setting file.

```~/.kirimori.toml
# VimrcPath = "~/.vimrc"

# ManagerType = "NeoBundle"
# ManagerType = "Vundle"
# ManagerType = "dein.vim"
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
