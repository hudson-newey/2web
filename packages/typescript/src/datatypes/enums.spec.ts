import { expectType } from "../testing/expect";
import { createEnum } from "./enums";

const myEnum = createEnum(["A", "B", "C"]);

expectType(myEnum).toBeObject();
expectType(myEnum.A).toBeSymbol();
expectType(myEnum.B).toBeSymbol();
expectType(myEnum.C).toBeSymbol();
