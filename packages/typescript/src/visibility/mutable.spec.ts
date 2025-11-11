import { mutable, type Mutable } from "./mutable";

const person = {
  name: "Alice",
  age: 30,
} as const;
type MutablePerson = Mutable<typeof person>;

const mutablePerson = mutable(person);
// mutablePerson.age = 31; // This should be allowed

/* Arrays */
const readonlyArray = [1, 2, 3] as const;
type MutableArray = Mutable<typeof readonlyArray>;
