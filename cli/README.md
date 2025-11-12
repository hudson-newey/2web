# 2Web CLI

A basic utility tool that can be used to create 2web projects.

This cli tool takes inspiration from the [Ember](https://cli.emberjs.com) and
[Angular](https://angular.dev/tools/cli) cli tools that emphasize configuration
over composition.

## Installation (system wide)

The easiest way to get started with the 2Web cli is to install it as a global
npm package.

```sh
$ npm install -g @two-web/cli
>
```

Once installed, you will have access to the `2web` command.

## Commands

| Command                            | Alias |
| ---------------------------------- | ----- |
| `2web new <project_name>`          | n     |
| `2web generate <generator> <name>` | g     |
| `2web template <template>`         | t     |
| `2web install <package_name>`      | i     |
| `2web serve [path]`                | s     |
| `2web build [path]`                | b     |
| `2web lint [path]`                 | l     |
| `2web format [path]`               | f     |
| `2web database <sub_command>`      | db    |
| `2web test [path]`                 |       |
| `2web doctor <sub_command>`        |       |
| `2web cms <sub_command>`           |       |

### Generate Command

| Command     | Alias | Description                                     |
| ----------- | ----- | ----------------------------------------------- |
| component   | c     | Add a component to a 2web project               |
| directive   | d     | Add a web-component directive to a 2web project |
| service     | s     | Add a service to a 2web project                 |
| aspect      | a     | Add an aspect to a 2web project                 |
| interceptor | i     | Add an interceptor to a 2web project            |
| page        | p     | Add a page to a 2web project                    |
| guard       | g     | Adds a route guard to a 2web project            |
| model       | m     | Add a model to a 2web project                   |
| enum        | e     | Add a **global** enum to a 2web project         |
| interface   |       | Add a **global** interface to a 2web project    |
| migration   |       | Adds a database migration                       |

### Template Command

Template commands differ from the "generate" command as they can only be run
once, and do not take a name as an argument.

| Command               | Alias | Description                             |
| --------------------- | ----- | --------------------------------------- |
| server-side-rendering | ssr   | Add ssr to a 2web project               |
| database              | db    | Add a database to a 2web project        |
| load-balancer         | lb    | Add a load balancer to a 2web project   |
| sitemap               |       | Adds a sitemap.xml file                 |
| robots.txt            |       | Adds a robots.txt file to the project   |
| security.txt          |       | Adds a security.txt file to the project |
| llms.txt              |       | Adds a llms.txt file to the project     |

### Database Command

| Command | Alias | Description                  |
| ------- | ----- | ---------------------------- |
| init    |       | Initializes a local database |
| migrate | m     | Runs a database migration    |

### CMS Command

| Command | Alias | Description            |
| ------- | ----- | ---------------------- |
| add     | a     | Adds a cms source      |
| view    | v     | View remote CMS source |
| sync    | s     | Sync CMS source        |
| remove  | rm    | Remove CMS source      |

### Serve Command

2Web can use Vite or an in-built development server to serve projects locally.

Vite is recommended for larger projects, while the in-built server is sufficient
for small projects, quick prototypes, or lower skill maintainers who may not
want to install npm or node.js.

2Web will use the following logic to determine which server to use:

1. If the project has an SSR target (e.g. through `2web template ssr`), use SSR server.
2. If a or `vite.config.ts` file is present in the project root, use Vite.
   While using the Vite dev server, the cli assumes that if you wanted to use
   the 2web compiler, you would be explicitly using the Vite plugin for 2web.
   1. If Vite is installed as an npm package, use the `node_modules` version
   2. If Vite is installed globally, use the global version
   3. Otherwise, use `npx vite`
3. Otherwise, use the in-built development server.
   1. If a local `./bin/2webc` compiler binary is present, build the project
      before serving any pages.
   2. Otherwise, use the global `2webc` compiler binary to build the project
      before serving any pages.
   3. If there is no 2web compiler available, serve the static files as-is.

### Build Command

2Web prefers using Vite for building projects, but can directly call the 2web
compiler (`2webc`) if Vite is not available.

The logic for determining which build tool to use is as follows:

1. If a `vite.config.ts` file is present in the project root,
   use Vite.
   1. If Vite is installed as an npm package, use the `node_modules` version
   2. If Vite is installed globally, use the global version
   3. Otherwise, use `npx vite`
2. Otherwise, use the in-built 2web compiler.
   1. If a local `./bin/2webc` compiler binary is present, use that
   2. Otherwise, use the global `2webc` compiler binary
   3. If there is no 2web compiler available, assets are directly copied to the
      output directory without any compilation and a warning is shown.

#### CMS Sources

The 2web CLI can build websites from existing document stores.

❌ = Not working, 🔧 = Developer preview, ✅ = Production ready

| Source     | State |
| ---------- | ----- |
| Wordpress  | ❌     |
| Git/GitHub | ❌     |
| OneDrive   | 🔧     |

### Doctor Command

Automatically checks for common problems in a 2web project.

| Command              | Alias | Description                                                     |
| -------------------- | ----- | --------------------------------------------------------------- |
| check                | c     | Checks common issues for a 2web project                         |
| check-dependencies   | cd    | Checks dependencies for 2web compiler, kit & cli                |
| install-dependencies |       | Installs **all** (including optional) dependencies used by 2web |

#### Dependencies

The 2web compiler requires some dependencies to convert file formats.
This is typically only needed in edge cases for unconventional file formats
e.g. creating a web page from a `.docx` (Microsoft Word) file.

To keep the number of dependencies low, 2web does not require these dependencies
to be installed until you need some of the dependencies features.

| Dependency                                        | Required for                                |
| ------------------------------------------------- | ------------------------------------------- |
| [2webc](https://github.com/hudson-newey/2web)     |                                             |
| [pandoc](https://pandoc.org)                      | `.tex`, `.docx`, `.doc`, `.odt`             |
| [dart-sass](https://sass-lang.com)                | `.sass`, `.scss`                            |
| [less](https://lesscss.org)                       | `.less`                                     |
| [rustc](https://rust-lang.org)                    | `.rs`                                       |
| [emcc](https://emscripten.org)                    | `.c`, `.h`, `.cpp`, `.c++`, `.hpp`, `.h++`  |
| [gleam](https://gleam.run)                        | `.gleam`                                    |
| [fable](https://fable.io)                         | `.fs`                                       |
| [.NET](https://dotnet.microsoft.com)              | `.cs`, `.fs`, `.vb`                         |
| [ffmpeg](https://ffmpeg.org)                      | Optimizing images/videos                    |
| [docker](https://www.docker.com)                  | Database, load balancer & deployment images |
| [docker-compose](https://docs.docker.com/compose) |                                             |
| [node](https://nodejs.org)                        |                                             |
| [npm](https://docs.npmjs.com)                     |                                             |
| [rclone](https://rclone.org)                      |                                             |
| [git](https://git-scm.com)                        |                                             |

#### Optional Dependencies

These dependencies can be temporarily downloaded when needed (e.g. through npx)
although, it is recommended to install these dependencies globally if used
frequently, and in the projects `node_modules/` for larger projects.

| Dependency                                       | Required for  |
| ------------------------------------------------ | ------------- |
| [Vite](https://vite.dev)                         | `2web serve`  |
| [oxlint](https://oxc.rs/docs/guide/usage/linter) | `2web lint`   |
| [prettier](https://prettier.io/)                 | `2web format` |
