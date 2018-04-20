import { identity } from "fp-ts/lib/function";
import * as t from "io-ts";
import { failure, success } from "io-ts";
import { createOptionFromNullable } from "io-ts-types";

const isNonEmptyString: t.Is<string> = (x: t.mixed): x is string =>
  typeof x === "string" && x.trim().length > 0;
const validateNonEmptyString: t.Validate<t.mixed, string> = (x: t.mixed, c: t.Context) =>
  isNonEmptyString(x) ? success(x) : failure(x, c);

export const NonEmptyString = new t.Type(
  "NonEmptyString",
  isNonEmptyString,
  validateNonEmptyString,
  identity,
);

const looksLikeHref = /^https?:\/\/.*/;
const isUrlString: t.Is<string> = (x: t.mixed): x is string =>
  typeof x === "string" && looksLikeHref.test(x);
const validateUrlString: t.Validate<t.mixed, string> = (x: t.mixed, c: t.Context) =>
  isUrlString(x) ? success(x) : failure(x, c);

export const UrlString = new t.Type("URLString", isUrlString, validateUrlString, identity);

export const OptionalUrlString = createOptionFromNullable(UrlString);
export type OptionalUrlString = t.TypeOf<typeof OptionalUrlString>;
