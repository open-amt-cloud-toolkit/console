# go-htmx-tailwind-example

Example CRUD app written in Go + HTMX + Tailwind CSS

This project implements a pure dynamic web app with SPA-like features but without heavy complex Javascript or Go frameworks to keep up with.  Just HTML/CSS + Go ❤️

![screenshot](./screenshot.jpeg)

## Prequisites
- Tailwind
``` bash
# Example for macOS arm64
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 tailwindcss

# Example for win x64
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-windows-x64
chmod +x tailwindcss-windows-x64
mv tailwindcss-windows-x64 tailwindcss

# Example for linux x64
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
mv tailwindcss-linux-x64 tailwindcss
```

## Develop

```
 Choose a make command to run

  init          initialize project (make init module=github.com/user/project)
  vet           vet code
  test          run unit tests
  build         build a binary
  dockerbuild   build project into a docker container image
  start         build and run local project
  css           build tailwindcss
  css-watch     watch build tailwindcss
```
