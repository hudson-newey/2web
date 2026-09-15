// Server modules (any ".server.ts" file that isn't a verb route) compile so
// that their named exports are callable over rpc from compiled blocks.
export function getUsers(): string[] {
  return ["ada", "grace", "linus"];
}

export function getUserCount(): number {
  return 3;
}

export function greetVisitor(name: string): string {
  return "Hello, " + (name || "visitor") + "!";
}

// Remote functions can return html content (rendered with [[ ]] in the page).
export function getUserTable(filter: string): string {
  const users = ["ada", "grace", "linus", "server rendered row"];

  const rows = users
    .filter((user) => filter === "all" || user.includes(filter))
    .map((user) => "<tr><td>" + user + "</td></tr>")
    .join("");

  return "<table><tbody>" + rows + "</tbody></table>";
}
