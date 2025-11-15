import { hydrate, prerender as ssr } from "preact-iso";
import { useState } from "preact/hooks";

export function App() {
	const [countValue, setCount] = useState(0);

  return (
		<>
			<h1 id="outTarget">{ countValue }</h1>
			<button id="inTarget" onClick={() => setCount(countValue + 1)}>Increment</button>
			<button onClick={() => setCount(countValue - 1)}>Decrement</button>
		</>
  );
}

if (typeof window !== "undefined") {
  hydrate(<App />, document.getElementById("app"));
}

export async function prerender(data) {
  return await ssr(<App {...data} />);
}
