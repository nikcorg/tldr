import { array } from "fp-ts/lib/Array";
import { IO } from "fp-ts/lib/IO";
import { TaskEither, taskEither } from "fp-ts/lib/TaskEither";
import { sequence } from "fp-ts/lib/Traversable";
import * as path from "path";
import { loadFile } from "./lib/loadFile";
import { mapValidationErrorToError } from "./lib/mapValidationErrorToError";
import { archiveToMd } from "./lib/markdown";
import { readDir } from "./lib/readDir";
import { safeloadYaml } from "./lib/safeloadYaml";
import { writeFile } from "./lib/writeFile";
import { Archive } from "./types/Archive";
import { Environment } from "./types/environment";

const environmentIO = new IO(() =>
  Environment.decode(process.env).getOrElseL(e => {
    throw mapValidationErrorToError(e);
  }),
);

const isYamlFile = (x: string) => x.endsWith(".yaml");
const descSortCmp = (a: string, b: string) => (b > a ? 1 : -1);
const taskEitherSeq = <L, T>(ts: Array<TaskEither<L, T>>) => sequence(taskEither, array)(ts);

// tslint:disable-next-line:no-expression-statement
environmentIO
  .map(env => ({
    archivePath: path.resolve(env.ARCHIVE_PATH),
    output: path.resolve(env.OUTPUT),
  }))
  .map(({ archivePath, output }) =>
    readDir(archivePath)
      .map(files => files.filter(isYamlFile).sort(descSortCmp))
      .chain(files =>
        taskEitherSeq<Error, string>(
          files
            .map(fileName => path.join(archivePath, fileName))
            .map(loadFile)
            .map(yamlTask => yamlTask.chain(safeloadYaml(Archive)).map(archiveToMd)),
        ),
      )
      .map(mdDocs => mdDocs.join("\n\n"))
      .chain(mdDoc => writeFile(output, mdDoc))
      .run(),
  )
  .run()
  .then(
    r => console.log("done", r.getOrElse(undefined as any)),
    e => console.error("oops", e.getOrElse(undefined)),
  );
