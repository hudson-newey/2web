# Control Flow

Note that control flow is currently only available for static conditions that
can be evaluated at compile time.
Future work is planned to support runtime control flow with the same syntax.

Syntax is also subject to change.

## If Conditions

```html
{% if false <h1>Hidden Content</h1> %}
```

## For Loops

```html
<ul>
  {% for 1,2,3 <li>item: {{&value}}</li> %}
</ul>
```

[**Next**](./8-code-splitting.md) (Code Splitting)
