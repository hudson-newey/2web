type JsonPrimitive = string | number | boolean | null;
type JsonObject = { [key: string]: JsonObject | JsonArray | JsonPrimitive };
type JsonArray = (JsonObject | JsonPrimitive)[];

export type JsonType = JsonObject | JsonArray;
