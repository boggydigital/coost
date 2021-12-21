# coost

coost is a persistently stored cookies module that can support apps and services authentication
needs. coost requires minimal configuration for consumers and support (manual) import from the
browser cookies.

## Using coost

Adding coost module to your Go app: `go get github.com/boggydigital/coost`

After that, here is how coost is commonly used:

1) Create new (persistent) cookie jar;
2) Get *http.Client from that (persistent) cookie jar;
3) (use this client anywhere in the app where you'd expect cookies will be used);
4) Defer storing cookies to update any new values we got from the server.

More details on each step are provided below.

### Creating new (persistent) cookie jar

In order to create a new (persistent) cookie jar you need to provide two pieces of information:

1) Slice of hosts for the jar (e.g. "example.com");
2) Directory to store the jar, "" value representing process working directory.

Clients need to call `jar, err := coost.NewJar(hosts, directory)` and would get a pointer to a jar
object that can be used to get an `http.Client` instance as well as `Store` that jar.

### Getting and using http.Client

Following jar initialization, clients can call `jar.NewClient()` convenience method to get a pointer
to an `http.Client` instance using that jar. Please note, that this convenience method sets some
default timeouts for the `http.Client` instance.

Alternatively, clients can use the jar they got as a result of `NewJar` method in their
own `http.Client` instance.

### Storing cookie updates from the server

During application lifecycle, cookie jar can get updated with new cookie values from the server. To
persist those values jar provides `Store` method.

Typically, clients would add `defer jar.Store()` after initializing the jar. Note: this would ignore
any errors that can happen during saving cookies file to the storage.

An extended form of the same call that would allow handling errors would be:

```go
defer func (jar coost.PersistentCookieJar) {
    if err := jar.Store(); err != nil {
    //handle Store error
    }
}(jar)
```

## cookies.json file

Please note that cookies file created by coost is not encrypted or obfuscated in any way.

cookies.json is a trivial JSON structure that stores cookies by hosts:

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
  authentication related request in order to get response cookies;
- Make sure to `Store` the jar when meaningful changes happen (e.g. when the authentication flow
  completes).

### Copying session cookies from an existing browser session

In certain scenarios it might not be practical to implement a full authentication flow - for example
for headless services that don't expect or allow user input.

To support those scenarios coost provides (relatively) easy way to import existing browser session
cookies. Here is how to do that (using Google Chrome as the browser, though it should generally work
with any browser):

1) Open a browser, then open Developer Tools:
    1) Select "Customize and control Google Chrome" menu item, visually represented as three
       vertical dots;
    2) Select "More Tools" menu item, then "Developer Tools";
    3) Select "Network" tab in the Developer Tools.
2) Navigate to the desired website (e.g. `example.com`);
3) Select that navigation request in the Network tab of Developer Tools;
4) Select "Headers" tab in the request preview section;
5) Navigate to the `Request Headers` section (make sure that's not `Response Headers`!);
6) Find `cookie:` header and `Copy value` for that header;
7) Open (or create) `cookies.json` for the client application, create a new section for the host in
   question (e.g. `example.com`);
8) In that host section add new key `cookie-header` and paste the copied value in quotes:

```json
{
  "example.com": {
    "cookie-header": "<paste-copied-value-here>"
  }
}
```

9) Repeat for any additional hosts that client application might require.

Some tips you can use to verify `cookie-header` has been imported correctly:

1) Cookie values are presented as `; ` separated key value pairs, where key and value are separated
   with `=`. Example:

```
COOKIE1=cookie-value1; COOKIE2=cookie-value2
```

2) When encountering `cookie-header` value, coost would split the contents into individual key-value
   pairs and remove original `cookie-header` entry;
3) Split `cookie-header` values are stored as separate values after `Store` is invoked, so
   effectively `cookie-header` is only imported, never preserved. If you see it after
   calling `Store` - something's wrong.

## General troubleshooting

coost has been designed as a thin-wrapper on top of Go `http.CookieJar` and has two main points of
interest:

### Restoring cookies into the jar

Expectations and common problems to check for:

1) File is located where the application expects it to be. When using "" make sure the file is in
   the current app working directory at the moment of `NewJar` call. When using specific directory
   make sure the file is present at this location and is accessible to the app;
2) Only the hosts specified for the `NewJar` call are used to populate cookie jar, regardless of
   what hosts are present in `cookie.json`;
3) `http.Client` with the jar is used for requests.

### Storing cookie jar

Expectations and common problems to check for:

1) File would be stored relative to the current app working directory at the time of the invocation.
   When using specific directory make sure it's accessible to the app;
2) Only the hosts specified at the `NewJar` call will be updated with new cookie values.

### Using nod for further debugging

coost has been tactically instrumented with `github.com/boggydigital/nod` logging calls, and clients
might find enabling nod logging output to be helpful. Please check `github.com/boggydigital/nod` for
more details.

