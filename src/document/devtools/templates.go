package devtools

func devtoolsHtmlSource() string {
	return `
    <div class="__2_devtools_container">
      <nav>
        <button>Selector</button>
        <button>Reactivity Graph</button>
        <button>Logs</button>
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
      }
    </style>
  `
}
