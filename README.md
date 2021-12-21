# pesco

pesco is a persistently stored cookies module that can support apps and services authentication needs. pesco requires minimal configuration for consumers and support (manual) import from the browser cookies.

## Using pesco

Adding pesco module to your Go app: `go get github.com/boggydigital/pesco`

After that, here is how pesco is commonly used:

1) create new (persistent) cookie jar
2) get *http.Client from that (persistent) cookie jar
3) (use this client anywhere in the app where you'd expect cookies will be used)
4) defer storing cookies to update any new values we got from the server 

More details on each step are provided below.

### Creating new (persistent) cookie jar

TBD

### Getting and using http.Client

TBD

### Storing cookie updates from the server

TBD

## Adding existing browser cookies to pesco cookies.json file

TBD
