# gowind

![Static Badge](https://img.shields.io/badge/tailwind-css-blue?style=for-the-badge&logo=tailwindcss)


gowind simply gives you access to tailwind without needing nodejs. It usess standalone cli that is provided
by tailwindCSS themselfs.

Binaries are downloaded based on your current operating system. It should work on all major systems like **Windows**, **Linux** and **MacOS**

[tailwindcss documentation](https://tailwindcss.com/docs/installation)

## installation

to install gowind you can run the following command

```sh
  go install github.com/maatko/gowind@latest
```

## usage

```sh
  gowind <update/clean>
```

```sh
  update - downloads the latest cli binary from tailwind
  clean  - deletes the tailwind binary from go
```

After you have updated you tailwindcss cli, you can access it by running `tailwindcss` in your terminal.