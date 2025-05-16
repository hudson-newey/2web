# 2Web Tutorial

Beware, 2Web is HIGHLY experimental, and not at all stable.
There are quite a few known bugs.

## Build a "Hello World" application

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

## First reactive app

To be honest, while the "hello world" example was exciting from a technical
standpoint, it doesn't show the exciting reactivity features provided by 2web.

I'll follow the typical reactivity demonstration that all web frameworks do and
demonstrate a simple counter that can increment and a number on the screen.

Because 2Web maintains simplicity as one of its core principles.
I believe that if I can't show the entire website in a single mobile/terminal
window, I have failed in creating a concise framework.
No code snippets will be shown out of context.

```html
<!DOCTYPE html>
<html>
  <body>
    <script compiled>
      $ count = 0;
    </script>

    <h1 *innerText="$count"></h1>
    <button @click="++count">Increment</button>
  </body>
</html>
```

Because 2web is so efficient, I can show you **all** of the compiled code in a
single mobile screen as well.

_NOTE_: This is an dev build. optimized for readability, not performance.

```html
<!DOCTYPE html>
<html>
  <body>
    <h1 __data-2-element-0>0</h1>
    <button onclick="__2_func_0(++__2_var__0)">Increment</button>

    <script>
      let __2_var_0 = 0;
      function __2_func_0(__2_value) {
        document.querySelector("__data-2-element-0")["innerText"] = __2_value;
      }
    </script>
  </body>
</html>
```

The first thing that you will notice is that there are no runtime objects, no
virtual DOM, no zones, no proxies, and no signals.
We are able to achieve reactivity without any runtime overhead.
Changing the variable simply changes the property of the target elements.

An annoying thing about the counter demonstration is that the `$count` compiler
variable still requires a runtime counterpart because its new value is dependent
on its previous value.
If all updates to `$count` did not rely on the previous `$count` value, 2web
would be able to create reactivity without the `__2_var_0` runtime variable.

While it would be possible to remove this runtime variable and boast about a
zero variable runtime, it would greatly hurt performance.

## "Mustache like" templates

Text templating is a common
