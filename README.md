# painkiller-gallery

## run the unit tests

```
$ go test github.com/cruftbusters/painkiller-gallery/heightmap
```

## run the acceptance tests

```
$ go test github.com/cruftbusters/painkiller-gallery/acceptance
```

## run the acceptance tests against deployment

```
$ go test github.com/cruftbusters/painkiller-gallery/acceptance -overrideBaseURL=https://gallery.painkillergis.com
```

## build and run

```
$ go build
$ build/painkiller-gallery # opens http server on port 8080
```
