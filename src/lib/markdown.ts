import { stripIndent } from "common-tags";
import { Archive } from "../types/Archive";
import { Entry } from "../types/Entry";

export const entryToMd = (entry: Entry) => stripIndent`
- [${entry.unread ? " " : "x"}] ${entry.title.concat("  ")}
  ${entry.url}
`;

export const archiveToMd = (archive: Archive) => stripIndent`
### ${archive.title}

${archive.entries.map(entryToMd).join("\n")}`;
