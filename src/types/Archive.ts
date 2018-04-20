import * as t from "io-ts";
import { ArrayOfEntry } from "./Entry";
import { NonEmptyString } from "./strings";

export const Archive = t.interface({
  entries: ArrayOfEntry,
  title: NonEmptyString,
});

export type Archive = t.TypeOf<typeof Archive>;

export const ArrayOfArchive = t.array(Archive);
export type ArrayOfArchive = t.TypeOf<typeof ArrayOfArchive>;
