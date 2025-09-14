import type { Unwrap } from "../structural/unwrap";

export type UppercaseAlphaCharacter = Uppercase<LowercaseAlphaCharacter>;
export type LowercaseAlphaCharacter =
  | "a"
  | "b"
  | "c"
  | "d"
  | "e"
  | "f"
  | "g"
  | "h"
  | "i"
  | "j"
  | "k"
  | "l"
  | "m"
  | "n"
  | "o"
  | "p"
  | "q"
  | "r"
  | "s"
  | "t"
  | "u"
  | "v"
  | "w"
  | "x"
  | "y"
  | "z";

export type AlphaCharacter = UppercaseAlphaCharacter | LowercaseAlphaCharacter;

export type NumericCharacter =
  | "0"
  | "1"
  | "2"
  | "3"
  | "4"
  | "5"
  | "6"
  | "7"
  | "8"
  | "9";

export type AlphanumericCharacter = Unwrap<AlphaCharacter | NumericCharacter>;

export type WhitespaceCharacter = " " | "\t" | "\n" | "\r" | "\v" | "\f";

export type PunctuationCharacter =
  | "`"
  | "~"
  | "!"
  | "@"
  | "#"
  | "$"
  | "%"
  | "^"
  | "&"
  | "*"
  | "("
  | ")"
  | "-"
  | "_"
  | "="
  | "+"
  | "["
  | "]"
  | "{"
  | "}"
  | "|"
  | "\\"
  | ";"
  | ":"
  | "'"
  | '"'
  | ","
  | "<"
  | "."
  | ">"
  | "/"
  | "?";

export type Character = Unwrap<
  AlphanumericCharacter | PunctuationCharacter | WhitespaceCharacter
>;
