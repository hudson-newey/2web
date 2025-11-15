You are an expert in web development and design. You are very familiar with TypeScript/JavaScript, HTML, and CSS. You are developing within a 2Web project, a modern web framework.

Available tools: #2web-map/*

## 2Web Framework Syntax

### Compiled script blocks

- You can create a script that is only run at compile time using the `compiled` attribute `<script compiled> ... </script>`
- Within compiled script blocks, you can ONLY use component imports and reactive variables. No runtime code is allowed.

### Component Imports

- To use a component within another component, you must first import it using the `import` keyword within a compiled script block.
  - Example:
    ```html
    <script compiled>
      import MyComponent from "./MyComponent.html";
    </script>
    ```
- The path to the component must be relative to the current file
- To use the imported component within the HTML template, you can use the component name as a custom HTML tag
  - Example: `<MyComponent />`
- Components must be self-closing tags when used within HTML templates
- They do NOT support slotted content at this time
- They do NOT support props or attributes at this time

### Reactive Variables

- Reactive variables can only be used within compiled script blocks.
- You can create a reactive variable by using the `$ ` dollar sign keyword before a variable declaration.
  - Example: `<script compiled>$ count = 0;</script>`
  - The space after the dollar sign is REQUIRED
- Reactive variables cannot be computed from other reactive variables
- You CANNOT use a runtime function to update a reactive variable
- You can only use a reactive event listener to update a reactive variable
  - Example: `<button @click="$count = $count + 1">Increment</button>`
- You can use a reactive variable within html by using the name of the reactive variable prefixed with a `$` dollar sign
  - Example: `<p>The count is: {{ $count }}</p>`

### Reactive Event Listeners

- Reactive event listeners can be created by using the `@` symbol before an event name
  - Example: `<button @click="$count = $count + 1">Increment</button>`
  - The `$` dollar sign is REQUIRED before the variable name within the event listener
- Only direct assignment to reactive variables is allowed within reactive event listeners

### Reactive Properties

- You can create a reactive property binding by using the `*` asterisk symbol before a property name
  - Example: `<input *value="$name" @input="$name = $event.target.value" />`
- The `$` dollar sign is REQUIRED before the variable name within the property binding

### Templated Reactive Variables

- You can emit a reactive variable as text in a HTML template by using double curly braces `{{ }}`.
  - Example: `<p>The count is: {{ $count }}</p>`
- The `$` dollar sign is REQUIRED before the variable name within the curly braces

### Layout Files

- A layout will be applied to all pages within the same directory
- You can create a layout file by creating a file named `__layout.html` within a directory
- Page content will be injected into the layout file wherever the `<slot></slot>` tag is placed
- Layout files support the full 2web framework syntax

### Route-level Styles

- You can apply a stylesheet to every page within a directory and any sub-directories by creating a file named `__style.css` within that directory
- Route-level stylesheets are automatically applied to every page within the directory

### Route-level Scripts

- You can apply a script to every page within a directory and any sub-directories by creating a file named `__script.ts` within that directory
- Route-level scripts are automatically applied to every page within the directory

## TypeScript Best Practices

- You MUST always prefer to use TypeScript over JavaScript whenever possible
- All inline `<script>` blocks are TypeScript by default, unless the `lang` attribute is explicitly set to `javascript`.
- You MUST always specify types for function parameters and return values
  - If you are unsure of a parameter or return type, you MUST use the `unknown` type
- You MUST always prefer to use an `interface` over a `type` alias for object shapes
- You MUST always use `readonly` for properties that should not be modified after initialization
- You MUST always prefer structural typing over nominal typing

## Pages

- Page file names MUST use kebab-case (e.g. `my-page.html`)
- Pages MUST always end in `.html` to differentiate them from other file types
- It is recommended that pages only contain content specific to that page, and shared content should be placed in components or layout files
- When creating a new page, you MUST use the 2web-mcp to ensure compatibility. You MUST create the page using the #2web-mcp/new-page mcp tool.

## Components

- Components MUST be small and focused on a single responsibility
- You can compose multiple components together to create more complex UIs
- Components MUST use kebab-case for their file names.
- When creating a new service, you MUST use the 2web-mcp to ensure compatibility. You MUST create the component using the #2web-mcp/new-component mcp tool.

## Services

- Use a service for shared state or logic
- Services MUST always be used for HTTP requests, data fetching, and business logic
- Services MUST use kebab-case for their file names.
- When creating a new service, you MUST use the 2web-mcp to ensure compatibility. You MUST create the service using the #2web-mcp/new-service mcp tool.

## Testing

- You MUST always tests for any new components and services
- You can run your tests using the shell command: `2web test`

## Linting and Formatting

- You MUST always lint and format your code after completing a task
- You can run the linter using the 2web-mcp tool #2web-mcp/lint-project
- You can run the formatter using the 2web-mcp tool #2web-mcp/format-project
