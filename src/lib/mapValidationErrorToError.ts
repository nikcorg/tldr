import * as t from "io-ts";
import { formatValidationError } from "io-ts-reporters";

export const mapValidationErrorToError = (e: Error | t.ValidationError[]) => {
  if (e instanceof Error) {
    return e;
  }
  const description = e
    .map(formatValidationError)
    .map(err => JSON.stringify(err))
    .join("\n\n");

  return new Error(`Validation error:\n\n${description}`);
};
