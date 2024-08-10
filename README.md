# gowind

gowind simply gives you access to tailwind without needing NodeJS. It usess standalone cli that is provided
by TailwindCSS themselfs.

## usage

To use TailWindCSS that is provided by this project, you simply run `gowind` that will redirect
execution directly to the tailwind standalone cli, so you can run any commands that are supported.

### init

To initialize tailwind you can run this command

```sh
  gowind init [--postcss]
```

You can also provide the `--postcss` tag to use it with tailwind

### watching

To start watching files for style changes you can run following

```sh
  gowind -i 'path/to/input/style' -o 'path/to/output/style'
```

## installation

to install gowind you can run the following command

```sh
  go install github.com/maatko/gowind@latest
```
