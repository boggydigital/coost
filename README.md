# coost

coost is a persistent cookie storage module that simplifies storing and recreating cookie jar between app sessions. coost requires minimal configuration and supports (manual) importing existing session cookies from the browser.

## Using coost

Adding coost module to your Go app: `go get github.com/boggydigital/coost`

A typical coost flow in a client app looks like this:

1) Create new (persistent) cookie jar;
2) Defer storing cookies to update any new values we got from the server;
3) Get *http.Client from that (persistent) cookie jar;
4) (use this client anywhere in the app where you'd expect cookies will be used).

More details on each step are provided below.

### Create new (persistent) cookie jar

In order to create a new (persistent) cookie jar you need to provide two pieces of information:

1) Slice of hosts for the jar - e.g. `[]string {"example.com"}`;
2) Directory to store the jar - e.g. `""`, that will represent process working directory.

Clients need to call `jar, err := coost.NewJar(hosts, directory)` and would get a pointer to a jar
object that can be used to get an `http.Client` instance as well as `Store` that jar.

### Defer storing cookies to update any new values we got from the server

During application lifecycle, cookie jar can get updated with new cookie values from the server. To
persist those values jar provides a `Store` method.

Typically, clients would add `defer jar.Store()` after initializing the jar. NOTE: this would ignore
any errors that can happen during saving cookies file to the storage.

An extended form of the same call that would allow handling errors would be:

```go
defer func (jar coost.PersistentCookieJar) {
    if err := jar.Store(); err != nil {
        //handle error
    }
}(jar)
```

### Get *http.Client from that (persistent) cookie jar

Following jar initialization, clients can call `jar.NewClient()` convenience method to get a pointer
to an `http.Client` instance using that jar. Please note, that this convenience method sets some
default timeouts for the `http.Client` instance.

Alternatively, clients can use the jar they got as a result of `NewJar` method in their
own `http.Client` instance.

## cookies.json file

Please note that cookies file created by coost is not encrypted or obfuscated in any way. This is by design.

`cookies.json` is a trivial JSON structure that stores cookies by hosts:

```json
{
  "host1": {
    "cookie1": "cookie1_value",
    "cookie2": "cookie2_value"
  },
  "host2": {
    "cookie1": "cookie1_value"
  }
}
```

## Using coost to support cookie-based authentication

coost is agnostic to the authentication model used by the service. Some common scenarios using
cookies:

1) Application handles authentication and provides prompt for username, password, 2FA, etc. Server
   sets response cookies that are then persisted by coost and users won't need to re-authenticate on
   every application session;
2) Users copy session cookies from an existing browser session.

Both scenarios are detailed below.

### Supporting authentication flow in the app

Detailing full authentication flow is out of the scope for this document and most likely would be
application specific. There are few callouts that would be helpful to clients:

- Make sure you're using `http.Client` with a cookie jar initialized by coost for every
  authentication related request in order to get (potentially updated) response cookies;
- Make sure to `Store` the jar when meaningful changes happen (e.g. when the authentication flow
  completes).

### Copying session cookies from an existing browser session

In certain scenarios it might not be practical to implement a full authentication flow - for example
for headless services that don't expect or allow user input.

To support those scenarios coost provides a way to import existing browser session
cookies. Here is how to do that:

0) (all instructions below are using Google Chrome as the browser, though it should generally work
with any browser);
1) Launch a browser, then open Developer Tools:
    1) Select "Customize and control Google Chrome" menu item, visually represented as three
       vertical dots;
    2) Select "More Tools" menu item, then "Developer Tools";
    3) Select "Network" tab in the Developer Tools.
2) Navigate to the desired website (e.g. `example.com`);
3) Select that navigation request in the Network tab of Developer Tools;
4) Select "Headers" tab in the request preview section;
5) Navigate to the `Request Headers` section (not `Response Headers`!);
6) Find `cookie:` header and `Copy value` for that header;
7) Open (or create) `cookies.json` for the client application, update (or create) a new section for the host in
   question (e.g. `example.com`);
8) In that host section add new key `cookie-header` and paste the copied value in quotes:

```json
{
  "example.com": {
    "cookie-header": "<paste-copied-value-here>"
  }
}
```

9) (Repeat steps 2-8 for any additional hosts that client application might require).

Some tips you can use to verify `cookie-header` has been imported correctly:

1) Cookie values are presented as `; ` separated key value pairs, where key and value are separated
   with `=`. Example:

```text
COOKIE1=cookie-value1; COOKIE2=cookie-value2
```

2) When encountering `cookie-header` value, coost would split the contents into individual key-value
   pairs and remove original `cookie-header` entry;
3) Split `cookie-header` values are stored as separate values after `Store` method is invoked, so
   effectively `cookie-header` is only imported, never persisted. If you see it after
   calling `Store` - something is wrong.

## General troubleshooting

coost has been designed as a thin wrapper on top of Go runtime `http.CookieJar` and has two main points of
interest:

1) Restoring cookies into a jar
2) Storing cookie jar

Both steps involve working with local storage, and clients should anticipate errors typical for general storage operations. Working with cookies themselves is handled by Go language runtime.

### Restoring cookies into the jar

Expectations and common problems to check for:

1) File should be located where the application expects it to be. When using `""` make sure the file is in
   the current app working directory at the moment of `NewJar` method call. When using specific directory
   make sure the file is present at this location and is accessible to the app;
2) Only the hosts specified for the `NewJar` method call are used to populate cookie jar, regardless of
   what hosts are present in `cookie.json`. This is by design;
3) `http.Client` with the jar is used for requests. Just loading the cookies into the jar does nothing - it's only useful, when called through an `http.Client` with that jar.

### Storing cookie jar

Expectations and common problems to check for:

1) File would be stored relative to the current app working directory at the time of the invocation.
   When using specific directory make sure it's accessible to the app;
2) Only the hosts specified at the `NewJar` method call will be updated with new cookie values. This is by design.

### Using nod for further debugging

coost has been tactically instrumented with `github.com/boggydigital/nod` logging calls, and clients
might find enabling nod logging output to be helpful. Please check `github.com/boggydigital/nod` for
more details.
