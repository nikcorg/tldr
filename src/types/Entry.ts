import * as t from "io-ts";
import { createOptionFromNullable } from "io-ts-types";
import { NonEmptyString, OptionalUrlString, UrlString } from "./strings";

export const SimpleEntry = UrlString;
export type SimpleEntry = t.TypeOf<typeof SimpleEntry>;

export const ConciseEntry = t.interface({
  title: NonEmptyString,
  url: UrlString,
});
export type ConciseEntry = t.TypeOf<typeof ConciseEntry>;

export const RelatedEntry = t.union([SimpleEntry, ConciseEntry]);
export type RelatedEntry = t.TypeOf<typeof RelatedEntry>;

export const ArrayOfRelatedEntry = t.array(RelatedEntry);
export type ArrayOfRelatedEntry = t.TypeOf<typeof ArrayOfRelatedEntry>;

export const OptionalArrayOfRelatedEntry = createOptionFromNullable(ArrayOfRelatedEntry);
export type OptionalArrayOfRelatedEntry = t.TypeOf<typeof OptionalArrayOfRelatedEntry>;

export const Entry = t.interface({
  related: OptionalArrayOfRelatedEntry,
  source: OptionalUrlString,
  title: NonEmptyString,
  unread: t.boolean,
  url: UrlString,
});
export type Entry = t.TypeOf<typeof Entry>;

export const ArrayOfEntry = t.array(Entry);
export type ArrayOfEntry = t.TypeOf<typeof ArrayOfEntry>;
