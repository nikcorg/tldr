import { task } from "fp-ts/lib/Task";
import { left, right } from "fp-ts/lib/TaskEither";
import * as t from "io-ts";
import * as yaml from "js-yaml";
import { mapValidationErrorToError } from "./mapValidationErrorToError";

export const safeloadYaml = <A, O>(type: t.Type<A, O>) => (yamlSrc: string) => {
  try {
    return type
      .decode(yaml.safeLoad(yamlSrc, { schema: yaml.JSON_SCHEMA }))
      .fold(
        e => left<Error, A>(task.of(mapValidationErrorToError(e))),
        x => right<Error, A>(task.of(x)),
      );
  } catch (e) {
    return left<Error, A>(task.of(e));
  }
};
