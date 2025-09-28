import { useState } from "react";

export default function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <h1 id="outTarget">{count}</h1>

      <button id="inTarget" onClick={() => setCount(count + 1)}>Increment</button>
      <button onClick={() => setCount(count - 1)}>Decrement</button>
    </>
  );
}
