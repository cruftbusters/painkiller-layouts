# painkiller-layouts

## run the unit tests

```
$ go test github.com/cruftbusters/painkiller-layouts/v1
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
$ build/painkiller-layouts 8080 http://localhost:8080 # opens http server on port 8080
```
