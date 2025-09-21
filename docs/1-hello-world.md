# Build a "Hello World" application

2Web is based on the "don't pay for what you don't use" philosophy.
If you don't include any reactive code in your website, you will not receive
any JavaScript.

Therefore, a simple "Hello World" website can be written like.

```html
<!DOCTYPE html>
<html>
  <body>
    <h1>Hello World!</h1>
  </body>
</html>
```

If you compile this code in a production build, the only changes that 2web will
make is minifying the supplied code.

_TIP_: You can run a production build with the `--production` flag

We do **not** perform a production build by default, because we believe that
providing readable, code from the development build is important for debugging
your code, understanding what you code does, and integrating with other tools
such as the web browsers development tools.

While we do ship embedded developer tools via the `--dev-tools` flag,
they are not as important as other web frameworks because we have spent a lot
effort ensuring that the web browsers development tools work nicely with the
development build code.
This means that when you get an error, it is a simple browser error that will
match (almost) 1-1 with your source code.

[**Next**](./2-reactive-apps.md) (Reactive Apps)
