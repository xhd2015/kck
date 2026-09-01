// Package pickupskill embeds the kck-pickup-a-session agent skill.
package pickupskill

import _ "embed"

// Name is the skill directory / skillcmd skill name.
const Name = "kck-pickup-a-session"

// SkillMD is the agent-facing SKILL.md content.
//
//go:embed SKILL.md
var SkillMD string
