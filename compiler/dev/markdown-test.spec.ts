import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";

beforeEach(async () => {
  await navigateToPage("markdown-test.html");
});

it("has no errors", () => {
  assertNoErrors();
});

// Literally just a hard coded string to assert that the output is correct
it("has the correct html content", () => {
  expect(document.body.innerHTML).toBe(`
    <div class="container">
      <div class="header">
        <span id="site-logo">
          <a href="/">
            2Web Compiler
          </a>
        </span>
      </div>
      <div class="sidebar">
        <nav class="sidebar-nav">
          <ul>
            <li>
              <a href="/compiler/dev/overview">
                Overview
              </a>
            </li>
            <li>
              <a href="/compiler/dev/building">
                Building
              </a>
            </li>
            <li>
              <a href="/compiler/dev/testing">
                Testing
              </a>
            </li>
            <li>
              <a href="/compiler/dev/deployment">
                Deployment
              </a>
            </li>
          </ul>
        </nav>
      </div>
      <main class="main-content">
        <h1>
          This is a markdown file
        </h1>
        <p>
          The following is a code block
        </p>
        <pre><code class="language-ts">const x = new Date();
console.log(x);
</code></pre>
        <p>
          <a href="https://github.com/hudson-newey/2web" target="_blank">
            Check out the GitHub link
          </a>
        </p>
        <p>
          Notice that this file does not have have html tags.
          This is because I do not believe in HTML inside of Markdown.
        </p>
        <p>
          Note that you can still use all ssg functions.
        </p>
        <pre><code class="language-html">&lt;script&gt;
  $ count = 0;
&lt;/script&gt;

&lt;h1&gt;{{ $count }}&lt;/h1&gt;
&lt;h2&gt;{{ $nonExistent }}&lt;/h2&gt;
</code></pre>
      </main>
    </div>
  

`);
});
