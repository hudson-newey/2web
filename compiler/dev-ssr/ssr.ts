import { runServer } from "@two-web/kit/ssr";

// Serves the compiled dev-ssr website: compiled html documents are server
// side rendered, compiled assets are served statically, and the compiled
// server routes (.server.ts files) are mounted.
//
// The compiler dev workflow for this website lives in .air-ssr.toml:
//
//	air -c .air-ssr.toml
//
// The website compiles into ./dist-ssr/ (and ./dist-ssr-server/ for the
// server routes) so that it doesn't clobber the static dev/ fixture build in
// ./dist/ that the compiler test suite reads.
runServer(Number(process.env.PORT ?? 5173), {
  clientDir: process.env.TWO_WEB_CLIENT_DIR,
  serverDir: process.env.TWO_WEB_SERVER_DIR,
});
