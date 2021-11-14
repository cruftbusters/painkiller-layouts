# painkiller-layouts

## run the unit tests

```
$ go test github.com/cruftbusters/painkiller-layouts/heightmap
```

## run the acceptance tests

```
$ go test github.com/cruftbusters/painkiller-layouts/acceptance
```

## run the acceptance tests against deployment

```
$ go test github.com/cruftbusters/painkiller-layouts/acceptance -overrideBaseURL=https://layouts.painkillergis.com
```

## build and run

```
$ go build
$ build/painkiller-layouts # opens http server on port 8080
```
