# First reactive app

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
    <h1 id="__2_element_0">0</h1>
    <button onclick="__2_func_0(++__2_var__0)">Increment</button>

    <script type="module">
      const __2_element_0 = document.getElementById("__2_element_0");
      let __2_var_0 = 0;
      function __2_func_0(__2_value) {
        __2_element_0["innerText"] = __2_value;
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

Additionally, all of the first paint content is pre-rendered as SSG content.
You can see in the template above that the initial text of the initial value of
`$count` is inlined as the text `0` in the `__2_element_0` element.

We use ids for reactivity because they can be looked up in constant time,
while finding an element in the DOM by anything that is not an id does not
perform in constant time.

Additionally, you'll note that the `__2_element_0` is cached for the elements
lifetime, meaning that we don't have to re-query the DOM for the element every
time it changes.

[**Next**](./3-templating.md) (Templating)
