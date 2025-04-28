# 2Web Router

A simple router for the 2Web framework that can turn an MPA app into an ssg app.

This router is purely compile time, meaning that routing can be compiled into
ssg, and the router does not require javascript variables.

To use, simply import the `spa-router.html` where you want the router outlet to
be.

```html
<!DOCTYPE html>
<html>
  <body>
    <header>
      <nav>
        <a href="/">Home</a>
        <a href="/about">About</a>
        <a href="/contact-us">Contact Us</a>
      </nav>
    </header>

    {% include node_modules/@hudson-newey/2web-kit/spa-router.html %}
  </body>
</html>
```
