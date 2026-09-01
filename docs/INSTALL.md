# Installation Information
This document details important information to prepare for usage of `go-choose-license`

## Go
This project requires Go to be able to use `go install` or `go build`.

You can install go in different ways, but I recommend using `Homebrew` (macOS & Linux Users) and `winget` (Windows Users).\
For Homebrew, make sure you have installed it. The install command is on [https://www.brew.sh/](https://www.brew.sh/). 
Winget is installed on Windows by default.

macOS & Linux Users:
```bash
brew install go
```
Windows Users (I would run in an Elevated Prompt aka `Run As Administrator`)
```powershell
winget install --id Golang.Go --source winget --silent
```

### Adding `$GOPATH/bin` to `PATH`
In order to use `go install`, you will need to add `$GOPATH/bin` to the `PATH` variable. 
Use the following commands to add `$GOPATH/bin` to your `PATH` variable based on your OS.

macOS Users
```bash
export PATH=$PATH:$(go env GOPATH)/bin > ~/.zshrc #if using zsh
```
Linux Users (or macOS users running `bash` as their default shell)
```bash
export PATH=$PATH:$(go env GOPATH)/bin > ~/.bashrc #if using bash
```
Windows Users:
```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";" + (go env GOPATH) + "\bin", "User")
```