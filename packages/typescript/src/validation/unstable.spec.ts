import { unstable, unwrap } from "./unstable";

const unstableValue = unstable(document.querySelector(".my-element"));

unwrap(unstableValue, (value) => {
  if (!value) {
    throw new Error("Element not found");
  }

  console.log(value.classList);
});

console.log(unstableValue);
