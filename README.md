# coost

coost is a cookie storage module for reading, writing and importing cookie jars.

## Importing cookies from an existing browser session

To import browser session cookies (instructions for Google Chrome, should work for any browser):

1. Launch a browser, then open Developer Tools:
   a. Select "Customize and control Google Chrome" menu item, visually represented as three
       vertical dots;
   b. Select "More Tools" menu item, then "Developer Tools";
   c. Select "Network" tab in the Developer Tools.
2. Navigate to the desired website (e.g. `example.com`);
3. Select that navigation request in the Network tab of Developer Tools;
4. Select "Headers" tab in the request preview section;
5. Navigate to the `Request Headers` section (not `Response Headers`!);
6. Find `cookie:` header and `Copy value` for that header (or select the value and copy);
7. In application that uses `coost` - run `import-cookies` command with a cookie header value: `import-cookies "key1=value1; key2=value2;"`