export function sayHello<T extends string>(name: T) {
  const greeting = `Hello ${name}` as const;

  console.log(greeting);

  return greeting
}
