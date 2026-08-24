# Scenario

**Feature**: help mode documents list CLI surface

```
Caller -> MainWith(["-h"|"--help"]) -> usage on stdout; exit 0
```

## Steps

1. Leaf sets help args only (no home seed required for usage text).
