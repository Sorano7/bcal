# bcal

`bcal` is an arbitrary-precision base calculator.

## Syntax

Any combination of `0-9, A-Z, a-z, _` is a valid number literal. Alphanumerics are used as digits and separators (underscore) are ignored. Up to base-`62` is supported with alphanumerical parsing.

Denote repeating decimal with `(...)`.

For arbitrary bases, use `{..., ..., ...}` to represent each digit.

The default base is 10. Use `[...]` to annotate base.

```bcal
> [12]1 / 3

0.4
```

Use `#` at the end of the expression to configure output arguments separated by comma.

- `base=<int>`: the base to render in.
- `digit=alnum|list`: render as alphanumerics or list.
- `num=decimal|rational`: render as decimal or rational.
- `prec=<int>`: precision. Set to 0 to attempt full repeat computation.

```
> [12]1B + 1 #base=10, digit=list

{2, 4}
```

Set variables with `@<ident>`. A single `@` will point to the last result. An assignment evaluates to the value being assigned.

```
> @x = 100
> @

100
```

## REPL

The REPL has a few helpful commands.

- `:q` or `:quit`: exits.
- `:env`: shows the current environment store.
- `:trun <int>`: sets the maximum output characters. Set to 0 to disable truncation.
