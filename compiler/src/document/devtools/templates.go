package devtools

func devtoolsHtmlSource() string {
	return `
    <div class="__2_devtools_container">
      <nav>
        <button class="__2_devtools_button" title="Element Selector">
          <img src="https://unpkg.com/lucide-static@latest/icons/square-dashed-mouse-pointer.svg" />
        </button>

        <button class="__2_devtools_button" title="Reactivity Graph">
          <img src="https://unpkg.com/lucide-static@latest/icons/waypoints.svg" />
        </button>

        <button class="__2_devtools_button" title="Compiler Logs">
          <img src="https://unpkg.com/lucide-static@latest/icons/logs.svg" />
        </button>
      </nav>
    </div>

    <style>
      .__2_devtools_container {
        position: fixed;
        bottom: 0px;
        left: 50%;
        transform: translate(-50%, 50%);

        padding: 1em;
        padding-left: 3em;
        padding-right: 3em;

        background-color: rgb(20, 10, 50);
        border-top-left-radius: 1rem;
        border-top-right-radius: 1rem;

        transition: transform 0.2s ease-in-out, opacity 0.2s ease-in-out;

        &:hover {
          transform: translate(-50%, 0%);
        }

        .__2_devtools_button {
          padding: 0.2rem;
          padding-left: 1.2rem;
          padding-right: 1.2rem;

          margin-left: 0.1rem;
          margin-right: 0.1rem;

          border: none;
          border-radius: 0.25rem;

          background: rgba(230, 240, 240, 1);

          &:active {
            background: rgba(230, 240, 240, 0.8);
          }
        }
      }
    </style>
  `
}
