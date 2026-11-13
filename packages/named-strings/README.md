# 2Web Kit - Named Strings

Semantic JavaScript string templates compatible with [bierner.lit-html](https://marketplace.visualstudio.com/items?itemName=bierner.lit-html)
to provide syntax highlighting to your code in strings.

## Example

The following code snippets can be syntax highlighted in vscode.

```ts
const cssCode = css`
    .primary {
        color: blue;
        font-size: 16px;
    }
`;

const javascriptCode = js`
    console.log("Hello World!");
`;

const htmlCode = html`
    <button>Hello World!</button>
`;

// https://developer.mozilla.org/en-US/docs/Web/SVG/Reference/Element/line#example
const svgCode = html`
    <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
        <line x1="0" y1="80" x2="100" y2="20" stroke="black" />
    </svg>
`;
```
