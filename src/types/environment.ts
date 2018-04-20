import * as t from "io-ts";
import { NonEmptyString } from "./strings";

export const Environment = t.interface({
  ARCHIVE_PATH: NonEmptyString,
  OUTPUT: NonEmptyString,
});
export type Environment = t.TypeOf<typeof Environment>;
