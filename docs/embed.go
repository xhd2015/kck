// Package docs embeds the kck multi-topic skill tree (Shape 3).
//
// Layout:
//
//	docs/SKILL.md                 — root skill index
//	docs/<path>/TOPIC.md          — nested topics (slash paths)
//
// Wire into skillcmd.SingleSkill with RootContent=SkillMD and TreeFS.
package docs

import "embed"

// Name is the skill directory / skillcmd skill name.
const Name = "kck"

// SkillMD is the root SKILL.md content (index only; no install plumbing flags).
//
//go:embed SKILL.md
var SkillMD string

// TreeFS holds nested topic directories (path/TOPIC.md). Root SKILL.md is not
// in TreeFS so ListTreeTopics only enumerates real topics.
//
//go:embed overview
//go:embed list
//go:embed open
//go:embed snapshot
//go:embed send
//go:embed messages
//go:embed info
//go:embed status
var TreeFS embed.FS
