import { tryCatch } from "fp-ts/lib/TaskEither";
import * as fs from "fs";

const defaultOptions = { encoding: "utf-8" };

export const loadFile = (path: string) =>
  tryCatch<Error, string>(
    () =>
      new Promise((res, rej) => fs.readFile(path, defaultOptions, (e, b) => (e ? rej(e) : res(b)))),
    e => (e instanceof Error ? e : new Error(`Error loading ${path}`)),
  );
