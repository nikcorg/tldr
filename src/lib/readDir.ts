import { tryCatch } from "fp-ts/lib/TaskEither";
import * as fs from "fs";

const defaultOptions = { encoding: "utf-8" };

export const readDir = (path: string) =>
  tryCatch<Error, string[]>(
    () =>
      new Promise((res, rej) =>
        fs.readdir(path, defaultOptions, (e, r) => (e ? rej(e) : res(r as string[]))),
      ),
    e => (e instanceof Error ? e : new Error(`Error reading ${path}`)),
  );
