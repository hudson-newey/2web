// Verb routes are the only http entry points for server endpoints. This one
// handles "GET /api/greetings".
const handler = (req: any, res: any) => {
  res.json({ endpoint: "greetings", method: "GET" });
};

export default handler;
