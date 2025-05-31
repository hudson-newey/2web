type Context = symbol;

export function provideContext(name: string): Context {
  return Symbol(name);
}
