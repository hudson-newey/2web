export function greet(name: string): string {
  return "Hello, " + name + "!";
}

export function add(left: number, right: number): number {
  return left + right;
}

const handler = (req: any, res: any) => {
  res.json({ source: "greetings.server.ts" });
};

export default handler;
