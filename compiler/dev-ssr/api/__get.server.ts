// Verb routes are the only http entry points for server endpoints. This one
// handles "GET /api" (the url path mirrors the directory the file sits in).
const handler = (req: any, res: any) => {
  res.json({
    site: "dev-ssr",
    method: "GET",
    time: new Date().toISOString(),
  });
};

export default handler;
