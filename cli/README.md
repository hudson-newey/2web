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
| `2web test [path]`                 |       |
| `2web database <sub_command>`      | db    |
| `2web doctor <sub_command>`        |       |

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
| vite                  |       | Add a Vite builder to the project       |
| tailwind              |       | Add a Tailwind to the project           |

### Database Command

| Command | Alias | Description                  |
| ------- | ----- | ---------------------------- |
| init    |       | Initializes a local database |
| migrate | m     | Runs a database migration    |

### Serve Command

```sh
$ 2web serve <path> <arguments>
>
```

| Command            | Alias | Description                           |
| ------------------ | ----- | ------------------------------------- |
| `--no-watch`       |       | Do not watch files for changes        |
| `--no-auto-reload` |       | Do not automatically reload dev pages |

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
| [fable](https://fable.io)                         | `.fs`                                       |
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

| Dependency                                         | Required for  |
| -------------------------------------------------- | ------------- |
| [oxlint](https://oxc.rs/docs/guide/usage/linter)   | `2web lint`   |
| [oxfmt](https://oxc.rs/docs/guide/usage/formatter) | `2web format` |
