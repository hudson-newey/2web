# 2Web Development

## Compiler

To run the compiler development environment, you need to run 2 commands while in
the `compiler/` directory.

[Air](https://github.com/air-verse/air) allows you to make changes to either the
golang or or website (html, css, ts, etc...) source code and have the
development website automatically re-built with the new compiler version.

```sh
$ air
>
```

To serve the development website, you can simply run `pnpm dev` in the
`compiler/` directory.

```sh
$ pnpm dev
>
```

Once the above steps have been completed, you can view the development website
at <https://localhost:5173>.

Note: Development build times are not indicative of real-world build times
because there is zero caching in these development builds.
