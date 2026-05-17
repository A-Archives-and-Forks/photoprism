## JS/Go Code Comment Rules

Doc comments for packages and exported identifiers must be complete sentences that begin with the name of the thing being described and end with a period.

A doc comment is **required** for every function (including unexported helpers), as well as for every non-trivial Vue `methods:` / `computed:` / watcher:
- Keep comments **compact** and default to one line for WHAT in the format `// Name does X.`. Skip trivial getters (`isOpen: () => this.open`).
- Add up to three follow-up lines (`// …`) only when the WHY is non-obvious: a hidden invariant, a workaround that would otherwise be undone by a future cleanup, a contract a reader can't infer from the code.
- For short examples in comments, indent code instead of using backticks.
- If readers can derive the WHY from the function body or a nearby line, leave it out. Multi-paragraph comments belong in specifications, package README files, or issue descriptions — never in the source itself.

> **Don't include in code comments:** Issue / PR numbers, "previously…" history, alternatives considered, what the function used to do, references to old commits, names of subsequent reviewers, or any narrative that names the change rather than the steady-state behavior. That context belongs in commit messages, specs, or handover notes.
