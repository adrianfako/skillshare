# Trusted bundled references

Declare intentional shared files in `~/.config/skillshare/trusted-references.json`:

```json
[{"root":"/absolute/repository/subtree", "target":"docs/shared.md"}]
```

The skill and the explicitly named regular file must both be inside `root`.
Other paths stay subject to the ordinary checks. A rooted filesystem handle
prevents symlink escapes. An absent or invalid file grants no trust. Review
additions as trust decisions. `SKILLSHARE_TRUSTED_REFERENCES` selects an
alternate JSON file, useful for tests.

A trusted template file may contain the literal example `[PR #123](url)`;
that placeholder is not treated as a file. Other missing links still warn.
