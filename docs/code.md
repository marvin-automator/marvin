# Conttributing Code

First of all:

![Thank you](https://media.giphy.com/media/3oz8xIsloV7zOmt81G/giphy.gif)

Thank you for helping out.

# If you're new to Go

If you've never contributed to a Go project before, **we'll help you out**.
Contributing to a Go project can be a bit confusing becaus of the way import
paths are resolved. Check out [this amazing tutorial on how to contribute to a Go repository](https://splice.com/blog/contributing-open-source-git-repositories-go/).

# How to contribute

  1. Fork this repository
  2. Run `go get github.com/marvin-automator/marvin`, and then `cd $GOPATH/src/github.com/marvin-automator/marvin
  3. Àdd your fork as a remote repository by running: `git remote add fork CLONE_URL`, where you replace CLONE_URL with the clone url of your fork.
  4. Install [dep](http://github.com/golang/dep)
  5. Run `dep ensure`
  6. Make your changes
  7. Run the tests by running `buffalo test`, and make sure the tests pass.
  8. Commit the changes, and push them to your fork: `git push fork master`
  9. Create a pull request to the master branch of this repository.