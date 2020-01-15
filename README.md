# go-billy-afero
`go-billy-afero` is a minimalist wrapper to make afero filesystems accessible as `billy.Filesystems`. 

## installation

```
go get github.com/jonfriesen/go-billy-afero
```

## usage

```go

import (
  "github.com/jonfriesen/go-billy-afero/pkg/fsadapter"
)

...

// in this example we expect that we are loading into a git repo

// open an os level file system at ./test
aferoFS := afero.NewBasePathFs(afero.NewOsFs(), "/some/path")

billyFS := fsadapter.New(aferoFS)

// using Afero with go-git (we assume we are in a git repo already)
// get a fs where the .git is the root
gitFS, err := billyFS.Chroot(".git")
HandleErr(err)

// loads the git storage
storage := filesystem.NewStorage(gitFS, cache.NewObjectLRUDefault())

// loads the git repo
repo, err := git.Open(storage, billyFS)

// get the git head
h, err := repo.Head()
HandleErr(err)

Println(repo.Hash().String())

```
