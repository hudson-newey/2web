const handler = (req: any, res: any) => {
  res.json({ endpoint: "greetings", method: "POST", received: req.body });
};

export default handler;
