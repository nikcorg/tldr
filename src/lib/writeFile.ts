import { tryCatch } from "fp-ts/lib/TaskEither";
import * as fs from "fs";

const defaultOptions = { encoding: "utf-8" };

export const writeFile = (path: string, contents: string) =>
  tryCatch(
    () =>
      new Promise((res, rej) =>
        fs.writeFile(path, contents, defaultOptions, e => (e ? rej(e) : res())),
      ),
    e => (e instanceof Error ? e : new Error(`Error writing ${path}`)),
  );
