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

- `2web new <project_name>` (alias: n)
- `2web generate <generator> <name>` (alias: g)
- `2web template <template>` (alias: t)
- `2web install <package_name>` (alias: i)
- `2web db <sub_command>`
- `2web doctor <sub_command>`
- `2web serve [path]`
- `2web build [path]`
- `2web lint [path]`

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

### Template Command

Template commands differ from the "generate" command as they can only be run
once, and do not take a name as an argument.

| Command               | Alias | Description                           |
| --------------------- | ----- | ------------------------------------- |
| server-side-rendering | ssr   | Add ssr to a 2web project             |
| database              | db    | Add a database to a 2web project      |
| load-balancer         | lb    | Add a load balancer to a 2web project |

### Doctor Command

Automatically checks for common problems in a 2web project.

| Command            | Alias | Description                                      |
| ------------------ | ----- | ------------------------------------------------ |
| check              | c     | Checks common issues for a 2web project          |
| check-dependencies | cd    | Checks dependencies for 2web compiler, kit & cli |

#### Dependencies

The 2web compiler requires some dependencies to convert file formats.
This is typically only needed in edge cases for unconventional file formats
e.g. creating a web page from a `.docx` (Microsoft Word) file.

To keep the number of dependencies low, 2web does not require these dependencies
to be installed until you need some of the dependencies features.

| Dependency                          | Required for                                   |
| ----------------------------------- | ---------------------------------------------- |
| [pandoc](https://pandoc.org/)       | `.tex`, `.docx`, `.doc`, `.odt`                |
| [dart-sass](https://sass-lang.com/) | `.sass`, `.scss`                               |
| [ffmpeg](https://ffmpeg.org/)       | Optimizing images/videos for production builds |
