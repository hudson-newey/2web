import type { FunctionType } from "./functions";

interface Person {
  name: string;
  age: number;
}

type GetPerson = FunctionType<[name: string, age: number], Person>;

const createPerson: GetPerson = (name, age) => {
  return { name, age };
};

createPerson("John Doe", 30);
