# gosrcdir

Simple utility for figuring out the correct path to put a repo in.

It doesn't create anything, and it doesn't include the name of the repo
itself. The idea is that you can pass its output to `mkdir -p` like so:

```
REPO="https://github.com/doozr/gosrcdir"
REPO_DIR=$(gosrcdir $REPO)
mkdir -p "$REPO_DIR"
cd "$REPO_DIR"
git clone "$REPO"
```

This will give you a nice, go-style directory tree using your GOPATH as the
root, but for any Git repo. Not just ones supported by `go get`.

Run it like so:

```
$ gosrcdir <repo URL>
```

It supports standard format URLs, like this:

```
$ gosrcdir https://github.com/doozr/gosrcdir
/home/user/go/src/github.com/doozr
```

And the weird URL format Git uses, like this:

```
$ gosrcdir git@github.com:doozr/gosrcdir.git
/home/user/go/src/github.com/doozr
```


