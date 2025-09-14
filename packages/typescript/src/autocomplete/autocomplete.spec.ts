import type { Autocomplete } from "./autocomplete";

type PowersOfTwo =
  | "1"
  | "2"
  | "4"
  | "8"
  | "16"
  | "32"
  | "64"
  | "128"
  | "256"
  | "512"
  | "1024";

type BufferSize = Autocomplete<string, PowersOfTwo>;

const bufferSize: BufferSize = "128";
