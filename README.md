# Go Choose Your License
**Go Choose Your License** ("go-choose-license") is a Go CLI app that helps you choose an open-source license for your next big project.

## Why This Exists
As an open-source developer, I know how many open-source licenses there are. 
I have used a multitude of different licenses based on my projects.
I remember when I was first starting out, I had no idea what license to use for my project,
and while ways to help choose a license did exist, nobody had a simple way for someone to 
choose a license for my next project. 
So, I decided to build my own, which is why we now have `go-choose-license`.

## How This Program Works
This program is simple, but genius.

## Quick Start
We have released downloadable binaries for macOS, Windows, and Linux for both `amd64` and `arm64` archetectures. 
You can run them by unzipping the executable and running it in the terminal in the unzipped directory
```bash
./go-choose-license
```
> [!WARNING]
> macOS Users will need to use the following sequence in order to open the binary since I don't pay for the Apple Developer Certificate.
> ```bash
> xattr -d com.apple.quarantine ./go-choose-license
> ```

Want to verify that the binaries have an appropriate attestation? You can do that using the GitHub CLI?
You can find information on how to verify the binaries on the [Releases](https://github.com/Bloomy52/go-choose-license/releases) Page.

## Installation
Don't want to download the binaries from GitHub Releases, I get it. Sometimes it is more fun to just build the software
from source. Or just download it via `go install`. Your choice.

In order to install the software, you need to have `Go` installed.
You can find instructions on downloading `Go` in the Installation Documentation [INSTALL.md](docs/INSTALL.md)

The following instructions assume you have installed `Go` and have it set up correctly.

### Using `go install`
You can install `go-choose-license` by running the following command in the terminal:
```bash
go install github.com/Bloomy52/go-choose-license@latest
```
Then you can run the app using the following command:
```bash
go-choose-license
```
> [!NOTE]
> You will need to make sure that your `$GOPATH/bin` has been added to `PATH` using the following commands found in the [INSTALL.md](docs/INSTALL.md) file.


### Building From Source
Want to build the package from source? We have a solution! You can follow the commands below to build the package from source. 
1. Clone the GitHub Repository & cd into the directory
```bash
git clone https://github.com/Bloomy52/go-choose-license.git
cd go-choose-license
```
2. Use `go build` to compile the app
```bash
go build -o go-choose-license .
```
3. Run the app
```bash
./go-choose-license
```

## License
This project is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for full license text.