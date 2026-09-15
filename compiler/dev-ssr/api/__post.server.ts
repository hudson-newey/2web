// Verb routes are the only http entry points for server endpoints. This one
// handles "POST /api" (the url path mirrors the directory the file sits in).
const handler = (req: any, res: any) => {
  res.json({
    site: "dev-ssr",
    method: "POST",
    received: req.body,
  });
};

export default handler;
