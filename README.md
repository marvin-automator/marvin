# Marvin
# Installing

I'm putting a process in place to create builds for Mac and Linux that you'll be able to download from the [releases page](https://github.com/marvin-automator/marvin/releases). In the meantime, use the development instructions below to build it from source. The 

# Development
- If you don't have it already, [install Go](https://golang.org).
- If you didn't set a different workspace by setting the `$GOPATH`
  variable, run `export GOPATH=~/go` to make it available in subsequent steps.
- Set up [Mage](https://magefile.org/), a Make equivalent for Go

        go get -u -d github.com/magefile/mage
        cd $GOPATH/src/github.com/magefile/mage
        go run bootstrap.go

- Clone this repo and cd into it. The location is important:

        git clone https://github.com/marvin-automator/marvin.git $GOPATH/src/github.com/marvin-automator/marvin
        cd $GOPATH/src/github.com/marvin-automator/marvin

- Run `mage setup` To set up dependencies.

