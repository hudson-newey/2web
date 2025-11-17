import { html, css } from "../named-strings/index.ts";

export const content = html`
  <div>Hello World</div>

  <button
    style="${css`
      color: red;
      background-color: white;
      padding: 0.5rem;
      font-size: 1.5rem;
      border-radius: 0.5rem;
      border: 2px solid red;
    `}"
    onclick="alert('Button clicked!')"
  >
    Click Me
  </button>
`;
