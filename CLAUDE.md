# 🚀 Go API Server - Complete Production Specification

**Version**: 2.0
**Last Updated**: 2025-10-14
**Purpose**: Complete, production-grade template for building Go-based API servers
**Based On**: Real implementations (airports, citylist, zipcodes projects)

---

## 📋 Critical: How AI Must Use This SPEC

### Before Making ANY Changes:

1. **READ THIS ENTIRE SPEC.md FILE** - Understand all standards
2. **SYNC your understanding** - Know what's mandatory vs optional
3. **FOLLOW ALL RULES** - No exceptions unless user explicitly requests
4. **ASK WHEN UNCLEAR** - Never assume, always clarify
5. **REPORT CHANGES** - List all files modified

### SPEC Synchronization Protocol:

**When to sync (read SPEC.md):**
- ✅ At the start of a new conversation
- ✅ Before implementing any new feature
- ✅ When user asks about standards
- ✅ When making architectural decisions
- ✅ When unsure about approach

**What to look for:**
- ✅ File structure requirements
- ✅ Naming conventions
- ✅ Technology choices (golang:alpine, alpine:latest, sapics GeoIP, etc.)
- ✅ Mandatory vs optional features
- ✅ Code patterns and examples

**How to apply:**
- ✅ Use template placeholders ({projectname}, {organization})
- ✅ Follow exact structure from Section 9 (Complete Project Layout)
- ✅ Implement features matching code examples
- ✅ Test according to Section 13 (Testing Priority)

---

## 📝 TODO.md Management (Critical for AI)

**TODO.md is for tracking multi-step tasks, NOT for simple operations.**

### When to Create TODO.md:

✅ **CREATE TODO.md when:**
- Task has 3 or more distinct steps
- Complex implementation requiring planning
- User provides multiple tasks (numbered list, comma-separated)
- Non-trivial feature that needs careful execution
- You need to track progress across multiple files

❌ **DON'T CREATE TODO.md when:**
- Single, straightforward task
- Trivial operation (< 3 steps)
- Simple file edit or read
- Answering a question
- User just wants information

### TODO.md Format (EXACT):

```markdown
# {Project Name} - Task List

**Last Updated**: 2025-10-14 12:00:00

## 🔄 In Progress (ONLY ONE TASK AT A TIME)
- [ ] Task currently being worked on
  - Details about what's being done
  - Current step or phase

## 📋 Pending
- [ ] Task 1 waiting to start
  - Brief description or acceptance criteria
- [ ] Task 2 waiting to start
- [ ] Task 3 waiting to start

## ✅ Completed
- [x] Task that was finished
  - What was accomplished
- [x] Another completed task
  - Files changed: file1.go, file2.md

---

**Progress**: 1/5 tasks completed (20%)
```

### TODO Management Rules (CRITICAL):

✅ **DO:**
- Mark task as "in progress" BEFORE starting work
- Keep ONLY ONE task "in progress" at any time
- Mark task complete IMMEDIATELY after finishing
- Update TODO.md throughout conversation (not batch)
- Delete TODO.md when ALL tasks complete
- Be specific about what each task entails
- Include file names in completed tasks

❌ **DON'T:**
- Have multiple tasks "in progress" simultaneously
- Batch mark multiple completions at once
- Leave tasks in "in progress" when moving to next
- Create TODO for trivial tasks
- Forget to delete TODO.md when done
- Keep stale or outdated tasks

### Task State Transitions:

```
Pending → In Progress → Completed
   ↑          ↓
   └──────────┘
  (can move back if blocked)
```

### Example TODO Workflow:

```
1. User: "Add rate limiting to all three projects"
   → CREATE TODO.md with 3 tasks (airports, citylist, zipcodes)

2. Mark "Add rate limiting to airports" as IN PROGRESS
   → Work on airports

3. Complete airports task IMMEDIATELY
   → Mark complete, list files changed

4. Mark "Add rate limiting to citylist" as IN PROGRESS
   → Work on citylist

5. Complete citylist task
   → Mark complete

6. Mark "Add rate limiting to zipcodes" as IN PROGRESS
   → Work on zipcodes

7. Complete zipcodes task
   → Mark complete
   → DELETE TODO.md (all done)
```

---

## ⚠️ Mandatory vs Optional Features

**EVERYTHING in this SPEC is MANDATORY unless listed below as an exception.**

### Only 2 Exceptions:

**Exception 1: IPv6 Support** - AUTO-DETECT
- ✅ Implement IPv6 support code (MANDATORY)
- ✅ Auto-detect OS/network capability (MANDATORY)
- ✅ Default to `::` (dual-stack), fallback to `0.0.0.0` if unsupported
- ✅ Graceful degradation (no errors if IPv6 unavailable)

**Exception 2: GeoIP Integration** - LOCATION-BASED PROJECTS ONLY
- ✅ Required IF: Project has location-based features
  - "Find nearby X" functionality
  - "Where am I?" features
  - Geographic search/filtering
  - IP-based country/region detection
- ❌ Not required IF: No geographic features
  - Static data APIs
  - No location context needed

**Examples:**
- airports: ✅ GeoIP (find airports near IP)
- zipcodes: ✅ GeoIP (find ZIP codes near IP)
- citylist: ❌ No GeoIP (cities don't need IP location)

### Everything Else is MANDATORY:

**No exceptions for:**
- All infrastructure (Dockerfile, docker-compose, Makefile, etc.)
- All CI/CD (GitHub Actions, Jenkins, ReadTheDocs)
- All frontend (Dark theme, Swagger/OpenAPI, GraphQL)
- All security (Rate limiting, DDoS protection, headers)
- All platforms (Linux, macOS, Windows, BSD - amd64/arm64)
- All install scripts
- Service management
- Signal handling
- Multi-distro testing
- Admin WebUI with live reload
- Latest images (golang:alpine, alpine:latest)

---

## 🎯 Template Placeholders

**Use these throughout your project:**

| Placeholder | Example | Usage | Case |
|------------|---------|-------|------|
| `{projectname}` | `airports` | Files, code, URLs, commands | lowercase |
| `{ProjectName}` | `Airports` | Display names, macOS paths | PascalCase |

---

## ⚠️ CRITICAL RULES FOR AI (READ FIRST)

**These rules override everything else. No exceptions.**

### RULE 1: Building (ABSOLUTE)

**AI: ALWAYS use Docker command directly. NEVER `make build`, NEVER build on host.**

**For Users**: `make build` is the command (Makefile uses Docker internally)
**For AI**: Use the Docker command directly (shown below)

```bash
# ✅ CORRECT - AI must use this
docker run --rm -v $(pwd):/workspace -w /workspace golang:alpine sh -c \
  'CGO_ENABLED=0 go build -o binaries/{projectname} ./src'

# ❌ FORBIDDEN for AI
make build                        # For users, not AI
go build ./src                    # NEVER - host build
CGO_ENABLED=1 go build ./src      # NEVER - dynamic binary
```

**Why Docker**:
- Consistent build environment (golang:alpine)
- Guarantees static binary (CGO_ENABLED=0)
- Cross-platform compilation
- No host dependencies or variations

**Why**: 
- Ensures consistent build environment
- Guarantees static binary (CGO_ENABLED=0)
- Cross-platform compilation
- No host dependencies

### RULE 2: Static Binary (ABSOLUTE)

**ALWAYS set CGO_ENABLED=0. No exceptions.**

```bash
# ✅ CORRECT
CGO_ENABLED=0 go build

# ❌ FORBIDDEN
CGO_ENABLED=1 go build
go build  # (defaults to CGO_ENABLED=1)
```

### RULE 3: Testing (ABSOLUTE)

**NEVER test on host. ALWAYS use containers.**

**Priority order:**
1. Incus (preferred)
2. Docker (acceptable)
3. Host (FORBIDDEN unless user explicitly says "test on host")

```bash
# ✅ CORRECT - Incus
incus launch images:alpine/3.19 test
incus exec test -- /usr/local/bin/app --version

# ✅ CORRECT - Docker
docker-compose -f docker-compose.test.yml up -d
curl http://localhost:64181/healthz

# ❌ FORBIDDEN - Host
./binaries/app --port 8080        # NEVER
./app                             # NEVER
```

**Why**:
- Isolated environment
- Reproducible
- No host pollution
- Matches production

---

**If AI violates these rules, STOP and inform user immediately.**

## 1. AI Assistant Guidelines

**This section defines how AI assistants MUST interact with projects following this SPEC.**

**Critical**: AI assistants must read and understand this entire section before making ANY changes to a project.

---

### 14.1 Core Principles

**AI assistants exist to:**
- ✅ Help implement features according to this SPEC
- ✅ Answer questions accurately
- ✅ Provide technical guidance
- ✅ Execute clearly stated commands
- ❌ NOT make decisions without user approval
- ❌ NOT assume user intent
- ❌ NOT execute based on questions

**Fundamental Rule**: **The user is in control. AI executes, doesn't decide.**

---

### 14.2 Understanding User Intent

#### 14.2.1 Question vs Command Detection

**Question Mark (?) = Question, NOT Command**

When a user's message ends with `?` or contains question words, the AI is being **asked for information**, not told to execute.

**Question Indicators:**
- Ends with `?`
- Starts with: "Should", "Could", "Would", "Can", "Is", "Are", "Do", "Does", "Will", "What", "Why", "How", "When", "Where"
- Contains: "What do you think", "Any ideas", "Suggestions", "Opinion"
- Exploratory: "I'm thinking", "Maybe", "Perhaps", "Wondering"

**Command Indicators:**
- Imperative verbs: "Update", "Add", "Fix", "Change", "Remove", "Create", "Delete", "Modify", "Implement"
- Direct statements: "Use X", "Change to Y", "Make it Z"
- Confirmations: "Do it", "Go ahead", "Yes" (after a proposal), "Proceed"
- Clear requests: "Please add", "I need you to", "Can you please" (with specific action)

#### 14.2.2 Response Patterns

**For Questions (ANSWER, DON'T EXECUTE):**

```
User: "Should we use sapics for GeoIP?"

AI Response:
"Yes, sapics/ip-location-db has several advantages over P3TERX:
 - Daily updates (vs weekly)
 - Multiple data sources aggregated (higher accuracy)
 - Public domain country data (no attribution required)
 - jsdelivr CDN (faster, more reliable)
 - Separate IPv4/IPv6 databases (better performance)

 Would you like me to migrate all projects from P3TERX to sapics?"

[WAIT for user confirmation before executing]
```

**For Commands (EXECUTE):**

```
User: "Update all projects to use sapics for GeoIP"

AI Response:
"I'll migrate airports and zipcodes from P3TERX to sapics/ip-location-db.
 This will update:
 - src/geoip/service.go (download URLs)
 - Database structure (4 files instead of 3)
 - Documentation (CLAUDE.md, README.md)

 Starting migration..."

[PROCEED with implementation]
```

**For Ambiguous Statements (ASK FOR CLARIFICATION):**

```
User: "Maybe we should add rate limiting..."

AI Response:
"Rate limiting is a good idea for production APIs. I can implement:

 1. Basic rate limiting (100 req/s global)
 2. Per-route rate limiting (different limits for API/admin)
 3. Configurable via admin panel (live reload)

 Which approach would you like, or would you like me to explain each option first?"

[WAIT for user decision]
```

#### 14.2.3 Common Scenarios

| User Input | AI Should | Reason |
|------------|-----------|--------|
| "Should we use alpine:latest?" | Answer with pros/cons, offer to implement | Question |
| "Use alpine:latest" | Execute the change | Clear command |
| "Can you add IPv6 support?" | Ask what level (basic/full), then offer | Ambiguous scope |
| "Add IPv6 support" | Implement according to Section 15 | Clear command |
| "I'm thinking about adding GraphQL" | Discuss pros/cons, ask if they want it | Exploratory |
| "What if we switched to MySQL?" | Explain implications, ask for decision | Hypothetical |
| "Fix the build error" | Investigate and fix | Clear command |
| "Is the build broken?" | Check and report status | Question |

---

### 14.3 Version Control Restrictions

**AI assistants are FORBIDDEN from performing ANY version control operations.**

#### 14.3.1 Forbidden Git Operations

❌ **NEVER execute these commands:**

**Write Operations (STRICTLY FORBIDDEN):**
- `git add` - Adding files to staging area
- `git commit` - Creating commits
- `git push` - Pushing to remote repository
- `git tag` - Creating or modifying tags
- `git merge` - Merging branches
- `git rebase` - Rebasing commits
- `git reset` - Resetting HEAD or index
- `git stash` - Stashing changes
- `git rm` - Removing files from git
- `git mv` - Moving files in git
- `git cherry-pick` - Cherry-picking commits
- `git revert` - Reverting commits
- `git branch -d` - Deleting branches
- Any operation that modifies `.git/` directory
- Any operation that changes git history

**Read Operations (ALLOWED):**
- ✅ `git status` - Check repository status
- ✅ `git diff` - View changes
- ✅ `git log` - View commit history
- ✅ `git show` - View specific commits
- ✅ `git branch` - List branches (without -d flag)
- ✅ `git ls-files` - List tracked files
- ✅ `git remote -v` - List remotes

#### 14.3.2 Why This Rule Exists

**Reasons version control is human-only:**
1. **Accountability** - Humans must author commits
2. **Meaningful messages** - Commit messages require context AI doesn't have
3. **Code review** - Changes should be reviewed before committing
4. **Prevent accidents** - AI might commit incomplete/broken code
5. **Clean history** - Avoid automated commits cluttering history
6. **User control** - Version control is user's responsibility

#### 14.3.3 What AI Should Do Instead

**When changes are complete:**

```
AI: "I've updated the following files:
     - src/main.go (added IPv6 support)
     - src/database/credentials.go (added getAccessibleURL function)
     - README.md (updated documentation)
     - Dockerfile (changed to alpine:latest)

     Total changes: 4 files modified

     You can review the changes with:
       git status
       git diff

     When ready to commit:
       git add .
       git commit -m 'Add IPv6 support with auto-detection'
       git push

     Suggested commit message:
       Add IPv6 support with dual-stack auto-detection

       - Implement auto-detect for IPv6 availability
       - Fallback to IPv4 if IPv6 unsupported
       - Update credentials URL formatting for IPv6 brackets
       - Update documentation with IPv6 examples"
```

**Key points:**
- ✅ List all files changed
- ✅ Provide git commands user can run
- ✅ Suggest commit message
- ✅ Let user review and commit

---

### 14.4 Testing & Environment Rules

#### 14.4.1 Environment Priority (STRICT)

**For Building (Cross-Compilation):**
- ✅ **Docker ONLY** - Use `golang:alpine` builder
- ❌ NEVER Incus for building
- ❌ NEVER Host OS for building

**For Testing/Debugging:**
1. **Incus** (PREFERRED) - System containers, multi-distro testing
2. **Docker** (FALLBACK) - If Incus unavailable
3. **Host OS** (LAST RESORT) - Only when explicitly requested, tools may need installation

#### 14.4.2 Testing Commands (Examples)

**✅ CORRECT - Incus for testing:**
```bash
# Generate random port
TESTPORT=$(shuf -i 64000-64999 -n 1)

# Launch Alpine container
incus launch images:alpine/3.19 test-alpine-${TESTPORT}

# Copy binary
incus file push ./binaries/{projectname} test-alpine-${TESTPORT}/usr/local/bin/

# Run binary
incus exec test-alpine-${TESTPORT} -- /usr/local/bin/{projectname} --port ${TESTPORT}

# Test
curl http://localhost:${TESTPORT}/healthz

# Cleanup
incus delete -f test-alpine-${TESTPORT}
```

**❌ WRONG - Host OS testing:**
```bash
./binaries/{projectname} --port 8080  # DON'T DO THIS (common port, host OS)
```

#### 14.4.3 Multi-Distro Testing (REQUIRED)

**AI MUST test on multiple distributions:**

**Required Tests:**
1. **Alpine** (musl libc, no systemd)
   ```bash
   incus launch images:alpine/3.19 test-alpine
   incus exec test-alpine -- /usr/local/bin/{projectname} --version
   ```

2. **Ubuntu 24.04** (glibc, systemd) - MANDATORY
   ```bash
   incus launch images:ubuntu/24.04 test-ubuntu
   # Test systemd service installation
   incus file push scripts/install-linux.sh test-ubuntu/tmp/
   incus exec test-ubuntu -- bash /tmp/install-linux.sh
   incus exec test-ubuntu -- systemctl status {projectname}
   ```

3. **Fedora** (SELinux, dnf) - Optional but recommended
   ```bash
   incus launch images:fedora/40 test-fedora
   ```

**Why multi-distro:**
- Verify musl compatibility (Alpine)
- Verify glibc compatibility (Ubuntu/Debian/Fedora)
- Verify systemd integration
- Verify binary portability

#### 14.4.4 Port Selection (STRICT)

**Random Ports ALWAYS:**
```bash
TESTPORT=$(shuf -i 64000-64999 -n 1)
```

❌ **FORBIDDEN Ports:**
- 80, 443 (HTTP/HTTPS production)
- 8080 (common development)
- 3000 (Node.js default)
- 5000 (Flask default)
- 5432 (PostgreSQL)
- 3306 (MySQL)
- 6379 (Redis)
- Any port < 64000

✅ **Allowed Ports:**
- 64000-64999 (random selection)
- Use `shuf -i 64000-64999 -n 1` to generate

#### 14.4.5 Temporary Files (STRICT)

**ALWAYS use /tmp/ for testing:**

```bash
# ✅ CORRECT
incus exec test-container -- {projectname} \
  --config /tmp/test/config \
  --data /tmp/test/data \
  --logs /tmp/test/logs

# Docker volumes
-v /tmp/{projectname}/rootfs/config:/config
-v /tmp/{projectname}/rootfs/data:/data

# ❌ WRONG - production directories
--config /etc/{projectname}      # NEVER during testing
--data /var/lib/{projectname}    # NEVER during testing
```

**Cleanup after testing:**
```bash
# Always cleanup
incus delete -f test-container
rm -rf /tmp/{projectname}
```

---

### 14.5 Code Changes & Modifications

#### 14.5.1 Before Making Changes

**Checklist (AI MUST follow):**
1. ✅ Read relevant section of SPEC.md
2. ✅ Understand the standard being applied
3. ✅ Check if feature is mandatory or optional
4. ✅ Review similar implementation in existing projects
5. ✅ Ask clarifying questions if ambiguous

#### 14.5.2 While Making Changes

**Process:**
1. ✅ Use Edit/Write tools (NOT sed/awk unless necessary)
2. ✅ Make one logical change at a time
3. ✅ Update TODO.md if using it (mark in progress → completed)
4. ✅ Test changes if possible
5. ✅ Update documentation if behavior changed

#### 14.5.3 After Making Changes

**AI MUST report:**
```
Changes Complete:

Files Modified:
1. src/main.go
   - Added IPv6 auto-detection
   - Updated default address to `::`

2. src/database/credentials.go
   - Added formatURL function for IPv6 brackets
   - Updated getAccessibleURL to handle IPv6

3. README.md
   - Added IPv6 examples
   - Updated listen address documentation

4. Dockerfile
   - Changed ADDRESS env to `::`

Testing:
- Build verified: ✓
- Binary size: 28MB
- Version output: 0.0.1

Next Steps:
- Test on IPv6-enabled system
- Verify bracket formatting in URLs
- Update other projects (citylist, zipcodes)
```

---

### 14.6 Communication Style

#### 14.6.1 Response Format

**Structure of AI responses:**

1. **Acknowledge** - Show you understood the request
2. **Clarify** - Ask questions if needed (BEFORE executing)
3. **Execute** - Make the changes
4. **Report** - List what was changed
5. **Suggest** - Next steps or related actions

**Example:**
```
User: "Add rate limiting to airports"

AI:
1. Acknowledge: "I'll add rate limiting to the airports project."

2. Clarify: "Should I use:
   - Basic (100 req/s global)
   - Per-route (50 API, 10 admin)
   - Configurable via admin panel

   I recommend option 3 (admin-configurable). Proceed?"

[User: "Yes, option 3"]

3. Execute: [Makes changes]

4. Report: "Rate limiting implemented:
   - Added httprate dependency
   - Created rate limiting middleware
   - Added admin UI for configuration
   - Default: 100 global, 50 API, 10 admin
   Files changed: src/server/server.go, go.mod"

5. Suggest: "Test with:
   for i in {1..150}; do curl http://localhost:PORT/api/v1/health; done
   Should see 200s then 429s (rate limit exceeded)"
```

#### 14.6.2 Asking Clarifying Questions

**When to ask:**
- User request is vague ("add security")
- Multiple valid approaches exist
- Significant architectural decision
- Breaking change or major refactor
- Ambiguous scope

**How to ask:**
```
"I can implement rate limiting in several ways:

Option 1: Basic (Quick)
  - 100 req/s global limit
  - Simple, works immediately
  - Not configurable without code change

Option 2: Per-Route (Better)
  - Different limits per route type
  - 100 global, 50 API, 10 admin
  - Still requires code change to adjust

Option 3: Admin-Configurable (Best)
  - All limits configurable via /admin/settings
  - Live reload (no restart)
  - Matches Section 18 of SPEC
  - More code but flexible

Which approach would you prefer? I recommend Option 3 for production."
```

**Don't ask when:**
- SPEC clearly defines the approach
- User gave explicit command
- Standard implementation exists
- No decision needed (SPEC mandates it)

#### 14.6.3 Explaining Decisions

**When deviating from user's initial idea:**

```
User: "Add GeoIP to citylist"

AI: "I notice citylist doesn't have location-based features (no 'find nearby' or
     geographic queries). According to Section 16 of the SPEC, GeoIP is only
     required for location-based projects.

     Did you mean to:
     1. Add location-based features to citylist (then GeoIP is needed)
     2. Skip GeoIP for citylist (follows SPEC as-is)

     Which direction would you like to take?"
```

**Key point**: Explain WHY, reference SPEC section, offer options.

---

### 14.7 Error Handling & Recovery

#### 14.7.1 When Errors Occur

**AI MUST:**
1. ✅ Report the error clearly
2. ✅ Explain what was being attempted
3. ✅ Suggest potential fixes
4. ✅ Ask if user wants to try a different approach
5. ❌ DON'T keep trying the same thing repeatedly

**Example:**
```
"Error: Build failed with 'package not found'

I was attempting to build airports with the new GeoIP integration.

The error suggests the import path is incorrect. This could be because:
1. The module path in go.mod is wrong
2. The import statement needs updating
3. go mod tidy needs to be run

Should I:
A) Run 'go mod tidy' to update dependencies
B) Check and fix import paths
C) Show you the full error for investigation

What would you like me to do?"
```

#### 14.7.2 Build/Test Failures

**When builds fail:**
```
"Build failed in airports project.

Error: src/main.go:25: undefined: geoip.NewService

Root cause: Function signature changed but call site not updated.

Fix applied:
- Updated main.go line 25 to pass new parameters
- Changed: geoip.NewService()
- To: geoip.NewService(configDir)

Rebuilding...
✓ Build successful
```

**When tests fail:**
```
"Tests failed: 2 failures in airports/data_test.go

Failures:
1. TestLoadAirports - Expected 29293 airports, got 0
   Cause: NewService() signature changed, test not updated

2. TestSearchByCity - nil pointer dereference
   Cause: Service not initialized in test

I can fix these test failures. Proceed?"
```

---

### 14.8 Following the SPEC

#### 14.8.1 Mandatory Features (NO EXCEPTIONS)

**When user asks to skip something mandatory:**

```
User: "Let's skip the admin panel for now"

AI: "The admin panel is mandatory according to Section 18 of the SPEC. It provides:
     - Live configuration reload (CORS, rate limits)
     - Secure credential management
     - Server settings management

     This is required for all projects following this SPEC.

     If you want to deviate from the SPEC, I can:
     1. Implement a minimal admin panel (meets requirements)
     2. Note this as a SPEC deviation (document why)

     Which approach would you prefer?"
```

**Key point**: Explain why it's mandatory, offer compliant solutions.

#### 14.8.2 Optional Features (Only 2)

**IPv6:**
```
AI: "I've implemented IPv6 support with auto-detection. The code will:
     - Try to bind to :: (dual-stack)
     - Fallback to 0.0.0.0 if IPv6 unavailable
     - No errors if OS doesn't support IPv6

     This follows Section 15 (IPv6 Auto-Detect)."
```

**GeoIP:**
```
User: "Should citylist have GeoIP?"

AI: "No. According to Section 16, GeoIP is only required for location-based projects.

     Citylist doesn't have:
     - 'Find nearby cities' feature
     - IP-based location detection
     - Geographic search/filtering

     Therefore, GeoIP is not needed for citylist.

     If you want to add location features to citylist, then GeoIP would become required."
```

#### 14.8.3 Technology Choices

**SPEC defines specific technologies - don't suggest alternatives without reason:**

**✅ CORRECT:**
```
AI: "I'll implement rate limiting using github.com/go-chi/httprate as specified
     in Section 17 (Security & DDoS Protection)."
```

**❌ WRONG:**
```
AI: "I'll implement rate limiting. We could use:
     - golang.org/x/time/rate
     - github.com/ulule/limiter
     - github.com/go-chi/httprate
     Which one?"
```

**Exception**: If user asks about alternatives or SPEC doesn't specify:
```
User: "What rate limiting library should we use?"

AI: "The SPEC (Section 17) recommends github.com/go-chi/httprate because:
     - Integrates with Chi router (our standard)
     - Supports per-IP limiting
     - Clean API
     - Well-maintained

     This is the standard for projects following this SPEC."
```

---

### 14.9 File Operations

#### 14.9.1 Preferred Tools

**Use specialized tools:**
- ✅ `Read` - For reading files
- ✅ `Write` - For creating new files
- ✅ `Edit` - For modifying existing files
- ✅ `Glob` - For finding files
- ✅ `Grep` - For searching content

**Avoid bash commands for file operations:**
- ❌ `cat` - Use Read tool instead
- ❌ `echo >` - Use Write/Edit tools instead
- ❌ `sed/awk` - Use Edit tool instead (unless complex multi-file operations)
- ❌ `find` - Use Glob tool instead
- ❌ `grep` - Use Grep tool instead

**When bash is appropriate:**
- ✅ Git operations (read-only)
- ✅ Build commands (make, go build)
- ✅ Service operations (systemctl, docker, incus)
- ✅ Package management (apt, apk, etc.)

#### 14.9.2 Creating New Files

**Before creating a new file, ask:**
1. Does a similar file already exist that could be edited?
2. Is this file structure defined in the SPEC?
3. Where should this file go according to Section 9 (Project Layout)?

**✅ CORRECT:**
```
User: "Add a health check handler"

AI: "I'll add the health check handler to src/server/handlers.go
     (the existing handlers file) rather than creating a new file."
```

**❌ WRONG:**
```
AI: [Creates src/server/health.go when handlers.go already exists]
```

---

### 14.10 Documentation Updates

#### 14.10.1 When to Update Documentation

**Update docs when:**
- ✅ Adding new features
- ✅ Changing behavior
- ✅ Adding new routes/endpoints
- ✅ Modifying CLI flags or environment variables
- ✅ Changing configuration options

**Files to update:**
1. `README.md` - User-facing changes
2. `CLAUDE.md` - Technical specification changes
3. `docs/API.md` - API endpoint changes
4. `docs/SERVER.md` - Server configuration changes

#### 14.10.2 Documentation Style

**Follow README structure from Section 8:**
- About → Official → Production → Development → Legal

**Keep consistent:**
- Use same terminology throughout
- Follow existing formatting
- Match tone and style
- Include examples
- Update all affected sections

---

### 14.11 Common Mistakes to Avoid

#### Mistake 1: Executing on Questions
```
❌ User: "Should we add this?"
   AI: [Adds it immediately]

✅ User: "Should we add this?"
   AI: "Here are the pros/cons... Would you like me to add it?"
```

#### Mistake 2: Using Common Ports
```
❌ AI: docker run -p 8080:80 ...
✅ AI: docker run -p $(shuf -i 64000-64999 -n 1):80 ...
```

#### Mistake 3: Testing on Host
```
❌ AI: ./binaries/projectname --port 8080
✅ AI: incus launch ... # container testing
```

#### Mistake 4: Committing Code
```
❌ AI: git add . && git commit -m "..."
✅ AI: "Changes ready. Run: git add . && git commit ..."
```

#### Mistake 5: Assuming Intent
```
❌ User: "Maybe GraphQL?"
   AI: [Implements GraphQL]

✅ User: "Maybe GraphQL?"
   AI: "GraphQL can provide... Do you want me to implement it?"
```

#### Mistake 6: Skipping Mandatory Features
```
❌ AI: "I'll skip the admin panel for now"
✅ AI: "Admin panel is mandatory (Section 18). Implementing..."
```

---

### 14.12 Success Criteria

**An AI assistant successfully follows this SPEC when:**

✅ **Understanding:**
- Recognizes questions vs commands
- References SPEC sections correctly
- Knows what's mandatory vs optional

✅ **Execution:**
- Uses correct tools (Incus/Docker, not Host)
- Follows file structure (Section 9)
- Implements features matching code examples
- Tests on multiple distros

✅ **Communication:**
- Asks before executing on questions
- Reports changes clearly
- Suggests next steps
- Provides git commands (doesn't execute them)

✅ **Standards:**
- Uses template placeholders correctly
- Follows naming conventions
- Implements all mandatory features
- Uses latest images (golang:alpine, alpine:latest)
- Applies security standards (Section 17)

---

### 14.13 Example AI Interaction (Perfect Execution)

```
User: "Should we switch the GeoIP provider to sapics?"

AI: "Yes, sapics/ip-location-db is superior to P3TERX because:
     - Daily updates (vs weekly)
     - Multiple aggregated sources (higher accuracy)
     - Public domain country data
     - jsdelivr CDN distribution
     - Separate IPv4/IPv6 databases

     This change affects airports and zipcodes (they use GeoIP).
     Citylist doesn't use GeoIP so no changes needed there.

     Would you like me to migrate both projects?"

User: "Yes"

AI: "Migrating airports and zipcodes to sapics/ip-location-db...

     [Makes changes]

     ✓ Migration complete

     Changes:
     1. airports/src/geoip/service.go - Updated to sapics URLs, 4 databases
     2. airports/CLAUDE.md - Updated GeoIP section
     3. airports/README.md - Updated data sources
     4. zipcodes/src/geoip/*.go - Updated all GeoIP files
     5. zipcodes/CLAUDE.md - Updated documentation
     6. zipcodes/README.md - Updated data sources

     Both projects now use:
     - geolite2-city-ipv4.mmdb (~50MB)
     - geolite2-city-ipv6.mmdb (~40MB)
     - geo-whois-asn-country.mmdb (~8MB)
     - asn.mmdb (~5MB)

     Testing: Build both projects to verify...
     ✓ airports builds successfully (28MB)
     ✓ zipcodes builds successfully (24MB)

     Next steps:
     1. Test GeoIP lookups (IPv4 and IPv6)
     2. Verify databases download on first run
     3. When satisfied, commit changes:
        git add .
        git commit -m 'Migrate GeoIP from P3TERX to sapics/ip-location-db'
        git push"

[Does NOT execute git commands]
```

---

**This AI Guidelines section is exhaustive and unambiguous. AI must understand and follow every rule.**
### 2. **Dockerfile - Alpine Runtime** ⭐ UPDATED

**Replace**: `FROM scratch` runtime
**With**: `FROM alpine:latest` runtime

**Why**: Need curl and bash for health checks and debugging

```dockerfile
# ============================================
# Build stage
# ============================================
FROM golang:alpine AS builder

# Why latest? Static binaries are self-contained, no Go runtime dependency.
# Using :latest gets latest Go features, security patches, and build optimizations.

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/

# Build static binary with all assets embedded
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE} -w -s" \
    -a -installsuffix cgo \
    -o {projectname} \
    ./src

# ============================================
# Runtime stage - Alpine with minimal tools
# ============================================
FROM alpine:latest

# Why latest? Static binaries (CGO_ENABLED=0) have no runtime dependencies.
# Alpine version only affects runtime tools (curl, bash, ca-certificates).
# Using :latest ensures latest security patches without version maintenance.

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# Install runtime dependencies (curl, bash)
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    bash \
    && rm -rf /var/cache/apk/*

# Copy binary to /usr/local/bin
COPY --from=builder /build/{projectname} /usr/local/bin/{projectname}
RUN chmod +x /usr/local/bin/{projectname}

# Environment variables
ENV PORT=80 \
    CONFIG_DIR=/config \
    DATA_DIR=/data \
    LOGS_DIR=/logs \
    ADDRESS=0.0.0.0 \
    DB_PATH=/data/db/{projectname}.db

# Create directories
RUN mkdir -p /config /data /data/db /logs && \
    chown -R 65534:65534 /config /data /logs

# Metadata labels (OCI standard)
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.authors="{organization}" \
      org.opencontainers.image.url="https://github.com/{organization}/{projectname}" \
      org.opencontainers.image.source="https://github.com/{organization}/{projectname}" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.vendor="{organization}" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.title="{projectname}" \
      org.opencontainers.image.description="{Project description} - Single static binary" \
      org.opencontainers.image.documentation="https://github.com/{organization}/{projectname}/blob/main/docs/README.md" \
      org.opencontainers.image.base.name="alpine:latest"

# Expose default port
EXPOSE 80

# Create mount points for volumes
VOLUME ["/config", "/data", "/logs"]

# Run as non-root user (nobody)
USER 65534:65534

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/{projectname}", "--status"]

# Run
ENTRYPOINT ["/usr/local/bin/{projectname}"]
CMD ["--port", "80"]
```

**Key Changes**:
- Runtime base: `alpine:latest` (not scratch)
- Includes: curl, bash, ca-certificates, tzdata
- Binary location: `/usr/local/bin/{projectname}` (not root)
- SQLite DB location: `/data/db/{projectname}.db`
- Internal port: 80
- All OCI metadata labels

---

### 3. **docker-compose.yml - Production Standards** ⭐ UPDATED

**Standards**:
- ❌ NO `version:` field
- ❌ NO `build:` definition
- ✅ Use pre-built images from registry
- ✅ Custom network: `{projectname}` (external: false)
- ✅ Volume structure: `./rootfs/{type}/{servicename}`
- ✅ Production port: `172.17.0.1:{randomport}:80`

```yaml
# Production Docker Compose Configuration
# Uses ./rootfs for persistent storage
# External port: 172.17.0.1:64180:80 (Docker bridge network)

services:
  {projectname}:
    image: ghcr.io/{organization}/{projectname}:latest
    container_name: {projectname}
    restart: unless-stopped

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/{projectname}.db
      # Uncomment and set for first deployment
      #- ADMIN_USER=administrator
      #- ADMIN_PASSWORD=changeme
      #- ADMIN_TOKEN=your-token-here

    volumes:
      - ./rootfs/config/{projectname}:/config
      - ./rootfs/data/{projectname}:/data
      - ./rootfs/logs/{projectname}:/logs

    ports:
      - "172.17.0.1:64180:80"

    networks:
      - {projectname}

    healthcheck:
      test: ["CMD", "/usr/local/bin/{projectname}", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

  # Optional: PostgreSQL database
  # postgres:
  #   image: postgres:16-alpine
  #   container_name: {projectname}-postgres
  #   restart: unless-stopped
  #
  #   environment:
  #     POSTGRES_DB: {projectname}
  #     POSTGRES_USER: {projectname}
  #     POSTGRES_PASSWORD: changeme
  #
  #   volumes:
  #     - ./rootfs/db/postgres:/var/lib/postgresql/data
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD-SHELL", "pg_isready -U {projectname}"]
  #     interval: 10s
  #     timeout: 5s
  #     retries: 5

  # Optional: MariaDB database
  # mariadb:
  #   image: mariadb:11-alpine
  #   container_name: {projectname}-mariadb
  #   restart: unless-stopped
  #
  #   environment:
  #     MYSQL_DATABASE: {projectname}
  #     MYSQL_USER: {projectname}
  #     MYSQL_PASSWORD: changeme
  #     MYSQL_RANDOM_ROOT_PASSWORD: "yes"
  #
  #   volumes:
  #     - ./rootfs/db/mariadb:/var/lib/mysql
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD", "healthcheck.sh", "--connect", "--innodb_initialized"]
  #     interval: 10s
  #     timeout: 5s
  #     retries: 5

  # Optional: Redis cache
  # redis:
  #   image: redis:7-alpine
  #   container_name: {projectname}-redis
  #   restart: unless-stopped
  #
  #   volumes:
  #     - ./rootfs/db/redis:/data
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD", "redis-cli", "ping"]
  #     interval: 10s
  #     timeout: 3s
  #     retries: 3

networks:
  {projectname}:
    name: {projectname}
    external: false
    driver: bridge
```

---

### 4. **docker-compose.test.yml - Development Standards** ⭐ UPDATED

**Standards**:
- ❌ NO `version:` field
- ❌ NO `build:` definition (use `{projectname}:dev` image)
- ✅ Ephemeral storage: `/tmp/{projectname}/rootfs`
- ✅ Development port: `{randomport}:80` (e.g., 64181:80)
- ✅ Same network name: `{projectname}`

```yaml
# Development/Testing Docker Compose Configuration
# Uses /tmp for ephemeral storage
# External port: 64181:80 (random unused port in 64xxx range)

services:
  {projectname}:
    image: {projectname}:dev
    container_name: {projectname}-test
    restart: "no"

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/{projectname}.db
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=testpass123
      - DEV=true

    volumes:
      - /tmp/{projectname}/rootfs/config/{projectname}:/config
      - /tmp/{projectname}/rootfs/data/{projectname}:/data
      - /tmp/{projectname}/rootfs/logs/{projectname}:/logs

    ports:
      - "64181:80"

    networks:
      - {projectname}

    healthcheck:
      test: ["CMD", "/usr/local/bin/{projectname}", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

  # Optional: PostgreSQL for testing
  # postgres:
  #   image: postgres:16-alpine
  #   container_name: {projectname}-postgres-test
  #   restart: "no"
  #
  #   environment:
  #     POSTGRES_DB: {projectname}_test
  #     POSTGRES_USER: {projectname}
  #     POSTGRES_PASSWORD: testpass
  #
  #   volumes:
  #     - /tmp/{projectname}/rootfs/db/postgres:/var/lib/postgresql/data
  #
  #   networks:
  #     - {projectname}

networks:
  {projectname}:
    name: {projectname}
    external: false
    driver: bridge
```

---

### 5. **Makefile - Docker Improvements** ⭐ UPDATED

Add `docker-dev` target for local development:

```makefile
# Build and push multi-platform Docker images (release)
docker:
	@echo "Building multi-platform Docker images..."
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):latest \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION) \
		--push \
		.
	@echo "✓ Docker images pushed to ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION)"

# Build Docker image for development (local only, not pushed)
docker-dev:
	@echo "Building development Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(PROJECTNAME):dev \
		.
	@echo "✓ Docker development image built: $(PROJECTNAME):dev"
```

---

### 6. **Jenkinsfile** ⭐ NEW

Multi-architecture CI/CD pipeline:

```groovy
// Jenkinsfile for {projectname}
// Multi-architecture build pipeline (amd64, arm64)
// Server: jenkins.casjay.cc

pipeline {
    agent none

    environment {
        PROJECTNAME = '{projectname}'
        PROJECTORG = '{organization}'
        REGISTRY = 'ghcr.io'
        IMAGE_NAME = "${REGISTRY}/${PROJECTORG}/${PROJECTNAME}"

        VERSION = sh(script: 'cat release.txt 2>/dev/null || echo "0.0.1"', returnStdout: true).trim()
        COMMIT = sh(script: 'git rev-parse --short HEAD 2>/dev/null || echo "unknown"', returnStdout: true).trim()
        BUILD_DATE = sh(script: 'date -u +%Y-%m-%dT%H:%M:%SZ', returnStdout: true).trim()
    }

    stages {
        stage('Checkout') {
            agent { label 'amd64' }
            steps {
                checkout scm
                echo "Building ${PROJECTNAME} ${VERSION} (${COMMIT})"
            }
        }

        stage('Test') {
            parallel {
                stage('Test on AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Running tests on AMD64...'
                        sh 'make test'
                    }
                }
                stage('Test on ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Running tests on ARM64...'
                        sh 'make test'
                    }
                }
            }
        }

        stage('Build Binaries') {
            parallel {
                stage('Build on AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Building binaries on AMD64...'
                        sh 'make build'
                        stash includes: 'binaries/**/*', name: 'binaries-amd64'
                        stash includes: 'releases/**/*', name: 'release-amd64'
                    }
                }
                stage('Build on ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Building binaries on ARM64...'
                        sh 'make build'
                        stash includes: 'binaries/**/*', name: 'binaries-arm64'
                        stash includes: 'releases/**/*', name: 'release-arm64'
                    }
                }
            }
        }

        stage('Build Docker Images') {
            parallel {
                stage('Build Docker AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Building Docker image for AMD64...'
                        script {
                            sh """
                                docker build \
                                    --platform linux/amd64 \
                                    --build-arg VERSION=${VERSION} \
                                    --build-arg COMMIT=${COMMIT} \
                                    --build-arg BUILD_DATE=${BUILD_DATE} \
                                    -t ${IMAGE_NAME}:${VERSION}-amd64 \
                                    -t ${IMAGE_NAME}:latest-amd64 \
                                    .
                            """
                        }
                    }
                }
                stage('Build Docker ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Building Docker image for ARM64...'
                        script {
                            sh """
                                docker build \
                                    --platform linux/arm64 \
                                    --build-arg VERSION=${VERSION} \
                                    --build-arg COMMIT=${COMMIT} \
                                    --build-arg BUILD_DATE=${BUILD_DATE} \
                                    -t ${IMAGE_NAME}:${VERSION}-arm64 \
                                    -t ${IMAGE_NAME}:latest-arm64 \
                                    .
                            """
                        }
                    }
                }
            }
        }

        stage('Push Docker Images') {
            agent { label 'amd64' }
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                echo 'Pushing Docker images to registry...'
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'github-registry',
                        usernameVariable: 'REGISTRY_USER',
                        passwordVariable: 'REGISTRY_TOKEN'
                    )]) {
                        sh """
                            echo \$REGISTRY_TOKEN | docker login ${REGISTRY} -u \$REGISTRY_USER --password-stdin

                            docker push ${IMAGE_NAME}:${VERSION}-amd64
                            docker push ${IMAGE_NAME}:${VERSION}-arm64
                            docker push ${IMAGE_NAME}:latest-amd64
                            docker push ${IMAGE_NAME}:latest-arm64

                            docker manifest create ${IMAGE_NAME}:${VERSION} \
                                ${IMAGE_NAME}:${VERSION}-amd64 \
                                ${IMAGE_NAME}:${VERSION}-arm64

                            docker manifest create ${IMAGE_NAME}:latest \
                                ${IMAGE_NAME}:latest-amd64 \
                                ${IMAGE_NAME}:latest-arm64

                            docker manifest push ${IMAGE_NAME}:${VERSION}
                            docker manifest push ${IMAGE_NAME}:latest

                            docker logout ${REGISTRY}
                        """
                    }
                }
            }
        }

        stage('GitHub Release') {
            agent { label 'amd64' }
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                echo 'Creating GitHub release...'
                script {
                    unstash 'release-amd64'
                    withCredentials([string(
                        credentialsId: 'github-token',
                        variable: 'GITHUB_TOKEN'
                    )]) {
                        sh 'make release'
                    }
                }
            }
        }
    }

    post {
        success {
            echo "✅ Build successful for ${PROJECTNAME} ${VERSION}"
        }
        failure {
            echo "❌ Build failed for ${PROJECTNAME} ${VERSION}"
        }
        always {
            cleanWs()
        }
    }
}
```

---

### 7. **src/data Directory - JSON Data Files** ⭐ NEW

**Rule**: `src/data/` directory contains ONLY JSON files. No Go code, no other file types.

**Structure**:
```
src/
├── main.go                      # Embeds: //go:embed data/*.json
├── data/
│   ├── {datafile}.json          # JSON data ONLY
│   └── {other}.json             # Additional JSON files (if needed)
└── {service}/
    └── service.go               # NewService(jsonData []byte)
```

**Embedding from main.go:**

```go
// src/main.go
package main

import (
    _ "embed"
    "github.com/{org}/{project}/src/{service}"
)

//go:embed data/{datafile}.json
var jsonData []byte

func main() {
    svc, err := service.NewService(jsonData)
    if err != nil {
        log.Fatal(err)
    }
    // ...
}
```

**Service receives embedded data:**

```go
// src/{service}/service.go
package service

import "encoding/json"

type Service struct {
    data MyDataType
}

func NewService(jsonData []byte) (*Service, error) {
    var data MyDataType
    err := json.Unmarshal(jsonData, &data)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    return &Service{data: data}, nil
}
```

**Why This Pattern:**
- ✅ `src/data/` contains ONLY JSON files (no .go code)
- ✅ JSON is embedded in single static binary
- ✅ No copies, no symlinks, no duplicates
- ✅ Embedding happens from `main.go` at `src/` level
- ✅ Services receive data as parameter (clean dependency injection)
- ✅ True single binary with all assets embedded

**Example** (airports project):
```
src/
├── main.go                  # //go:embed data/airports.json
├── data/
│   └── airports.json        # 8.7MB (ONLY JSON, no .go files)
└── airports/
    └── data.go              # NewService(jsonData []byte)
```

**Important**:
- The `src/data/` directory contains ONLY JSON files
- NO .go files in `src/data/`
- NO symlinks or copies
- Embedding MUST be done from `src/main.go`
- Services accept `[]byte` parameter with embedded data

---

### 12. **Web UI / Frontend Standards** ⭐ NEW

**Purpose**: Modern, responsive dark-themed web interface with Go html/template

**Technology Stack:**
- Go `html/template` for server-side rendering
- CSS3 with CSS custom properties (variables)
- Vanilla JavaScript (no frameworks)
- Embedded assets (go:embed static/* and templates/*)

**Note**: Static assets (CSS, JS, images), HTML templates, AND JSON data files are ALL embedded via `go:embed` in the single static binary. JSON files are embedded from `main.go` and passed to services.

---

**Structure:**
```
src/server/
├── static/
│   ├── css/
│   │   └── main.css        # Main stylesheet (~900 lines)
│   ├── js/
│   │   └── main.js         # JavaScript utilities (~130 lines)
│   ├── images/
│   │   └── favicon.png
│   └── manifest.json       # PWA manifest
├── templates/
│   ├── base.html           # Base template with header/footer
│   ├── home.html           # Homepage
│   ├── search.html         # Search page
│   ├── {feature}.html      # Feature-specific pages
│   └── config.html         # Admin configuration page
└── templates.go            # Template embedding logic
```

---

**CSS Framework** - `src/server/static/css/main.css`

**Design System:**

```css
/* CSS Variables (Dark Theme Default) */
:root {
  /* Colors */
  --bg-primary: #1a1a1a;        /* Main background */
  --bg-secondary: #2d2d2d;      /* Cards, sections */
  --bg-tertiary: #404040;       /* Hover states */
  --text-primary: #ffffff;      /* Main text */
  --text-secondary: #b0b0b0;    /* Secondary text */
  --text-tertiary: #808080;     /* Disabled text */
  
  /* Accent Colors */
  --accent-primary: #0066cc;    /* Primary actions */
  --accent-success: #00aa00;    /* Success states */
  --accent-warning: #ff9900;    /* Warnings */
  --accent-danger: #cc0000;     /* Errors */
  
  /* Layout */
  --border-color: #404040;
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  
  /* Typography */
  --font-primary: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  
  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-normal: 300ms ease;
  
  /* Shadows */
  --shadow-sm: 0 1px 3px rgba(0,0,0,0.12);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.16);
  --shadow-lg: 0 10px 20px rgba(0,0,0,0.2);
}

/* Light Theme Override */
[data-theme="light"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --bg-tertiary: #e0e0e0;
  --text-primary: #1a1a1a;
  --text-secondary: #404040;
  --border-color: #e0e0e0;
}
```

---

**Components:**

**1. Buttons**
```css
.btn-primary {
  padding: var(--space-md) var(--space-xl);
  background: var(--accent-primary);
  border: none;
  border-radius: 8px;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-primary:hover {
  background: #0052a3;
  transform: translateY(-2px);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
```

**2. Cards**
```css
.card {
  background: var(--bg-secondary);
  padding: var(--space-xl);
  border-radius: 12px;
  transition: all var(--transition-normal);
  border: 2px solid transparent;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--accent-primary);
}
```

**3. Modals**
```css
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999;
  display: none;
}

.modal.active { display: block; }

.modal-backdrop {
  background: rgba(0,0,0,0.7);
}

.modal-content {
  background: var(--bg-secondary);
  max-width: 600px;
  margin: 100px auto;
  border-radius: 12px;
  padding: var(--space-xl);
  box-shadow: var(--shadow-lg);
}
```

**4. Toast Notifications**
```css
#toast-container {
  position: fixed;
  top: 80px;
  right: 20px;
  z-index: 10000;
}

.toast {
  background: var(--bg-secondary);
  padding: var(--space-md) var(--space-lg);
  border-radius: 8px;
  box-shadow: var(--shadow-lg);
  border-left: 4px solid var(--accent-primary);
  min-width: 300px;
  animation: slideIn 0.3s ease;
}

.toast.success { border-left-color: var(--accent-success); }
.toast.error { border-left-color: var(--accent-danger); }
.toast.warning { border-left-color: var(--accent-warning); }
```

**5. Forms**
```css
.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.form-group input, .form-group select {
  padding: var(--space-md);
  background: var(--bg-tertiary);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-primary);
}
```

**6. Grid Layouts**
```css
/* Features/Stats Grid */
.features-grid, .stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--space-lg);
}

/* Results Grid */
.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-lg);
}
```

---

**JavaScript Utilities** - `src/server/static/js/main.js`

**Features:**
- Theme toggle (dark/light with localStorage persistence)
- Toast notifications
- Modal dialogs
- API helper functions
- Mobile menu toggle
- Keyboard shortcuts (Escape to close modal)
- Utility formatters (distance, coordinates, etc.)

**Example Functions:**
```javascript
// Toast Notifications
function showToast(message, type = 'info', duration = 3000) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    container.appendChild(toast);
    setTimeout(() => toast.remove(), duration);
}

// Modal Dialogs
function showModal(title, content) {
    const modal = document.createElement('div');
    modal.className = 'modal active';
    modal.innerHTML = `
        <div class="modal-backdrop" onclick="closeModal()"></div>
        <div class="modal-content">
            <div class="modal-header">
                <h2>${title}</h2>
                <button class="modal-close" onclick="closeModal()">×</button>
            </div>
            <div class="modal-body">${content}</div>
        </div>
    `;
    document.getElementById('modal-container').appendChild(modal);
}

// API Helpers
async function apiGet(endpoint) {
    const response = await fetch(endpoint);
    const data = await response.json();
    if (!data.success) {
        throw new Error(data.error?.message || 'API request failed');
    }
    return data.data;
}

// Theme Toggle
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}
```

---

**HTML Templates** - `src/server/templates/`

**Base Template** (`base.html`):
```html
<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - {{.ProjectName}}</title>
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="icon" href="/static/images/favicon.png">
</head>
<body>
    <header id="main-header">
        <div class="header-container">
            <div class="header-left">
                <button class="mobile-menu-toggle" onclick="toggleMobileMenu()">☰</button>
                <a class="logo" href="/">{{.ProjectName}}</a>
            </div>
            <nav id="main-nav" class="header-center">
                <a href="/">Home</a>
                <a href="/search">Search</a>
                <a href="/docs">API</a>
                <a href="/stats">Stats</a>
            </nav>
            <div class="header-right">
                <button class="theme-toggle" onclick="toggleTheme()">
                    <span class="theme-icon">🌙</span>
                </button>
            </div>
        </div>
    </header>
    
    <main id="main-content">
        {{template "content" .}}
    </main>
    
    <footer id="main-footer">
        <div class="footer-container">
            <p>&copy; 2024 {{.ProjectName}} - MIT License</p>
        </div>
    </footer>
    
    <div id="toast-container"></div>
    <div id="modal-container"></div>
    <script src="/static/js/main.js"></script>
</body>
</html>
```

**Template Embedding** (`templates.go`):
```go
package server

import (
    "embed"
    "html/template"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var templates *template.Template

func initTemplates() error {
    var err error
    templates, err = template.ParseFS(templateFS, "templates/*.html")
    return err
}
```

---

**Admin Interface Standards:**

**Admin Dashboard** (`/admin`):
- Clean, modern layout
- Server statistics cards
- Settings management
- Database status
- Log viewer
- System health monitoring

**Features:**
- ✅ Dark theme by default (light theme available)
- ✅ Responsive design (mobile-friendly)
- ✅ Toast notifications for feedback
- ✅ Modal dialogs for confirmations
- ✅ Real-time updates
- ✅ Form validation
- ✅ Keyboard shortcuts
- ✅ Accessible (ARIA labels)

**Color Scheme:**
- Dark Mode: `#1a1a1a` background, `#0066cc` primary
- Light Mode: `#ffffff` background, `#0066cc` primary
- Consistent with professional, modern UI/UX

---

**Responsive Design:**

```css
/* Mobile (< 720px) */
@media (max-width: 720px) {
  .mobile-menu-toggle { display: block; }
  .header-center {
    position: absolute;
    top: 70px;
    flex-direction: column;
    display: none;
  }
  .header-center.active { display: flex; }
  .hero h1 { font-size: 2rem; }
  .form-row { grid-template-columns: 1fr; }
}

/* Desktop (> 720px) */
@media (min-width: 720px) {
  #main-content { width: 90%; }
}
```

---

**Key Principles:**
- ✅ **No external dependencies** - No Bootstrap, Tailwind, React, etc.
- ✅ **Vanilla JS** - Pure JavaScript, no jQuery
- ✅ **Embedded assets** - All CSS/JS/images built into binary
- ✅ **Progressive enhancement** - Works without JavaScript
- ✅ **Dark theme first** - Default dark, supports light
- ✅ **Mobile responsive** - Works on all screen sizes
- ✅ **Fast loading** - Minimal CSS/JS, no external requests
- ✅ **Accessible** - Semantic HTML, ARIA labels

---

**File Sizes:**
- CSS: ~867 lines / ~25KB
- JavaScript: ~130 lines / ~4KB
- Total frontend: ~30KB (gzipped: ~10KB)

**Browser Support:**
- Modern browsers with CSS Grid support
- Chrome/Edge 88+
- Firefox 78+
- Safari 14+


---

### 11. **GitHub Actions Workflows** ⭐ NEW

**Purpose**: Automated builds on push and monthly schedule (split into separate workflows)

**Files**:
- `.github/workflows/release.yml` - Binary builds and GitHub releases
- `.github/workflows/docker.yml` - Docker image builds and registry push

**Triggers** (both workflows):
- Push to `main` branch
- Monthly schedule (1st of month, 3:00 AM UTC)
- Manual trigger (workflow_dispatch)

---

**Workflow 1: release.yml** - Binary Release

**Jobs:**

1. **test**
   - Runs `make test` with Go 1.23
   - Validates all tests pass before building

2. **build-and-release**
   - Reads version from `release.txt` (does NOT update it)
   - Runs `make build` for all 8 platforms
   - Deletes existing release if exists (always recreate)
   - Creates new GitHub release `{VERSION}`
   - Attaches all platform binaries
   - Uploads artifacts (90 day retention)

---

**Workflow 2: docker.yml** - Docker Release

**Jobs:**

1. **build-and-push**
   - Reads version from `release.txt`
   - Builds multi-arch Docker images (amd64, arm64)
   - Pushes to `ghcr.io/{organization}/{projectname}`
   - Tags: `latest`, `{VERSION}`, `{branch}-{sha}`
   - Uses GitHub cache for faster builds
   - Verifies images after push

**Configuration:**

```yaml
name: Build and Release

on:
  push:
    branches:
      - main
      - master
  schedule:
    - cron: '0 3 1 * *'  # Monthly on 1st at 3 AM UTC
  workflow_dispatch:

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-binaries:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Read version
        run: |
          VERSION=$(cat release.txt)
          echo "VERSION=${VERSION}" >> $GITHUB_OUTPUT
      - name: Build
        run: VERSION=$VERSION make build
      - uses: actions/upload-artifact@v4
        with:
          name: binaries
          path: releases/*

  build-docker:
    runs-on: ubuntu-latest
    permissions:
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-qemu-action@v3
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:latest
            ghcr.io/${{ github.repository }}:${{ steps.version.outputs.VERSION }}

  create-release:
    needs: build-binaries
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/download-artifact@v4
        with:
          name: binaries
          path: releases/
      - name: Delete existing release
        run: gh release delete $VERSION -y || true
      - name: Create release
        run: gh release create $VERSION ./releases/* --title "$VERSION"
```

**Important:**
- Version is read from `release.txt`, NOT modified by Actions
- Release is always recreated (deleted first if exists)
- Builds all platforms (8 binaries)
- Multi-arch Docker images (amd64, arm64)
- Runs on every push to main and monthly

---

## 🎯 Quick Reference

### Container Tags
- **Production**: `ghcr.io/{organization}/{projectname}:latest`
- **Versioned**: `ghcr.io/{organization}/{projectname}:{VERSION}`
- **Development**: `{projectname}:dev`

### Port Mappings
- **Production**: `172.17.0.1:{randomport}:80`
- **Development**: `{randomport}:80`
- **Internal**: Always `80`

### Volume Structure
```
./rootfs/
├── config/{servicename}/
├── data/{servicename}/
├── logs/{servicename}/
└── db/{dbtype}/
```

### Releases Directory Structure



**./releases/** contains files for GitHub releases:



```

./releases/

├── {projectname}-linux-amd64

├── {projectname}-linux-arm64

├── {projectname}-macos-amd64

├── {projectname}-macos-arm64

├── {projectname}-windows-amd64.exe

├── {projectname}-windows-arm64.exe

├── {projectname}-bsd-amd64

├── {projectname}-bsd-arm64

├── {projectname}-{VERSION}-src.tar.gz    # Source archive (no VCS)

└── {projectname}-{VERSION}-src.zip       # Source archive for Windows

```



**Source Archives:**

- Created with `git archive` (excludes .git, .github, etc.)

- Both tar.gz (Linux/macOS) and zip (Windows) formats

- Includes all source code, docs, configs

- No binaries, no VCS files, no build artifacts




### Binary Requirements
- ✅ Static binary (CGO_ENABLED=0)
- ✅ All assets embedded:
  - HTML templates (go:embed in server package)
  - CSS, JS, images (go:embed in server package)
  - JSON data files (go:embed in main.go)
- ✅ True single binary - no external files needed
- ✅ Location: `/usr/local/bin/{projectname}`
- ✅ Output: `./binaries/` or `./releases/`

### URL Display
- ❌ Never: localhost, 127.0.0.1, 0.0.0.0
- ✅ Priority: FQDN → hostname → public IP → fallback

### Data Loading (JSON Files)
```
src/
├── main.go                # Embeds: //go:embed data/{file}.json
├── data/
│   └── {datafile}.json    # JSON files ONLY (no .go files)
└── {service}/
    └── service.go         # NewService(jsonData []byte)
```
- ✅ `src/data/` contains ONLY JSON files
- ✅ NO Go code in `src/data/`
- ✅ Embedded in main.go, passed to services as []byte
- ✅ True single binary - JSON embedded, not external files

---

**Apply these updates to SPEC.md for production-grade specifications** 🚀

---

### 8. **README.md Structure** ⭐ UPDATED

**Order**: About → Production → Docker → API Usage → Development

```markdown
# {PROJECT_NAME}

Brief description

## About
- Project description
- Key features
- Technology stack

## Official Links
- Website (if exists)
- Documentation
- GitHub
- Docker Hub
- Issues

## Production Deployment

### Installation
- Linux (all distros, all init systems)
- macOS (launchd)
- FreeBSD/OpenBSD/NetBSD (rc.d)
- Windows (NSSM service)
- Docker

### Binary Installation
- Platform-specific downloads
- Features list
- Key capabilities

- Production deployment
- Binary installation
- Systemd service
- Environment variables

- Docker deployment
- Docker Compose (production)
- Docker Compose (development)
- Docker run examples

- API usage
- Quick examples
- Admin panel
- Documentation links

- Development

- Requirements
- Build System & Testing (Makefile targets, platforms, versioning)
- Development Mode (dev flags, debug features)
- CI/CD (GitHub Actions, Jenkins, ReadTheDocs)

## Legal
- License (MIT - LICENSE.md)
- Third-party licenses
- Credits
- Support
```

---

### 10. **ReadTheDocs Configuration** ⭐ NEW

**Purpose**: Automatic documentation hosting with Dracula theme

**Files Required:**

1. **`.readthedocs.yml`** (root):
```yaml
version: 2

mkdocs:
  configuration: docs/mkdocs.yml
  fail_on_warning: false

formats:
  - pdf
  - epub

python:
  version: "3.11"
  install:
    - requirements: docs/requirements.txt
```

2. **`docs/mkdocs.yml`**:
```yaml
site_name: {projectname} Documentation
site_description: {Project description}
site_author: {organization}
site_url: https://{projectname}.readthedocs.io

repo_name: {organization}/{projectname}
repo_url: https://github.com/{organization}/{projectname}
edit_uri: edit/main/docs/

copyright: Copyright &copy; 2024 {organization}

theme:
  name: material
  custom_dir: overrides
  palette:
    scheme: dracula
    primary: deep purple
    accent: pink
  font:
    text: Roboto
    code: Fira Code
  features:
    - navigation.instant
    - navigation.tracking
    - navigation.tabs
    - navigation.sections
    - navigation.expand
    - navigation.top
    - search.suggest
    - search.highlight
    - search.share
    - content.code.copy
    - content.code.annotate
  icon:
    repo: fontawesome/brands/github

nav:
  - Home: index.md
  - API Reference: API.md
  - Server Guide: SERVER.md
  - GitHub: https://github.com/{organization}/{projectname}

markdown_extensions:
  - admonition
  - pymdownx.details
  - pymdownx.superfences
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.inlinehilite
  - pymdownx.snippets
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.tasklist:
      custom_checkbox: true
  - pymdownx.emoji
  - toc:
      permalink: true
  - tables
  - attr_list

plugins:
  - search
  - minify:
      minify_html: true

extra_css:
  - stylesheets/dracula.css

extra_javascript:
  - javascripts/extra.js
```

3. **`docs/requirements.txt`**:
```txt
mkdocs>=1.5.0
mkdocs-material>=9.5.0
pymdown-extensions>=10.7
mkdocs-minify-plugin>=0.8.0
mike>=2.0.0
```

4. **`docs/stylesheets/dracula.css`** - Dracula color scheme overrides
5. **`docs/javascripts/extra.js`** - Optional custom JavaScript

**Documentation Structure:**
```
docs/
├── index.md              # Documentation home
├── API.md                # API reference
├── SERVER.md             # Server administration
├── mkdocs.yml            # MkDocs configuration
├── requirements.txt      # Python dependencies
├── stylesheets/
│   └── dracula.css       # Dracula theme CSS
├── javascripts/
│   └── extra.js          # Custom JavaScript
└── overrides/            # Theme overrides (optional)
```

**Dracula Color Palette:**
- Background: `#282a36`
- Current Line: `#44475a`
- Foreground: `#f8f8f2`
- Comment: `#6272a4`
- Cyan: `#8be9fd`
- Green: `#50fa7b`
- Orange: `#ffb86c`
- Pink: `#ff79c6`
- Purple: `#bd93f9`
- Red: `#ff5555`
- Yellow: `#f1fa8c`

**Local Preview:**
```bash
cd docs
pip install -r requirements.txt
mkdocs serve
# Open http://localhost:8000
```

---

### 9. **Complete Project Layout** ⭐ REFERENCE

**Standard Go API Server Project Structure** (matches citylist/airports):

```
{projectname}/
├── .claude/
│   └── settings.local.json      # Claude Code settings
├── .github/
│   └── workflows/
│       └── build.yml            # GitHub Actions (build on push & monthly)
├── .gitattributes               # Git attributes
├── .gitignore                   # Git ignore patterns
├── .readthedocs.yml             # ReadTheDocs configuration
├── CLAUDE.md                    # Project specification (this file)
├── Dockerfile                   # Alpine-based multi-stage build
├── docker-compose.yml           # Production compose
├── docker-compose.test.yml      # Development/testing compose
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Jenkinsfile                  # CI/CD pipeline (jenkins.casjay.cc)
├── LICENSE.md                   # License file
├── Makefile                     # Build system (4 targets)
├── README.md                    # User documentation
├── release.txt                  # Version tracking (auto-increment)
├── binaries/                    # Built binaries (gitignored)
│   ├── {projectname}-linux-amd64
│   ├── {projectname}-linux-arm64
│   ├── {projectname}-windows-amd64.exe
│   ├── {projectname}-windows-arm64.exe
│   ├── {projectname}-macos-amd64
│   ├── {projectname}-macos-arm64
│   ├── {projectname}-bsd-amd64
│   ├── {projectname}-bsd-arm64
│   └── {projectname}            # Host platform binary
├── releases/                    # Release artifacts (gitignored)
├── rootfs/                      # Docker volumes (gitignored)
│   ├── config/
│   │   └── {projectname}/       # Service config
│   ├── data/
│   │   └── {projectname}/       # Service data
│   ├── logs/
│   │   └── {projectname}/       # Service logs
│   └── db/                      # External databases
│       ├── postgres/
│       ├── mariadb/
│       └── redis/
├── docs/                        # API documentation (optional)
│   ├── index.md                 # Documentation home
│   ├── API.md                   # Complete API reference
│   ├── SERVER.md                # Server administration guide
│   ├── README.md                # Documentation index
│   ├── mkdocs.yml               # MkDocs configuration (Dracula theme)
│   ├── requirements.txt         # Python dependencies for RTD
│   ├── stylesheets/
│   │   └── dracula.css          # Dracula theme CSS
│   ├── javascripts/
│   │   └── extra.js             # Custom JavaScript
│   └── overrides/               # MkDocs theme overrides
├── scripts/                     # Production scripts (optional)
│   ├── install.sh               # Installation script
│   ├── backup.sh                # Backup script
│   └── uninstall.sh             # Uninstallation script
├── test/                        # Test files (optional)
│   ├── test-docker.sh           # Docker testing script
│   ├── unit/                    # Unit tests
│   ├── integration/             # Integration tests
│   └── e2e/                     # End-to-end tests
└── src/                         # Source code
    ├── data/
    │   └── {datafile}.json      # JSON data ONLY (no .go files)
    ├── {service}/
    │   └── service.go           # Service logic
    ├── auth/                    # Authentication (optional)
    │   └── auth.go
    ├── database/
    │   ├── database.go          # DB connection & schema
    │   ├── auth.go              # Admin authentication
    │   ├── credentials.go       # Credential management
    │   └── settings.go          # Settings CRUD
    ├── paths/
    │   └── paths.go             # OS-specific directory detection
    ├── scheduler/               # Task scheduler (optional)
    │   └── scheduler.go
    ├── server/
    │   ├── server.go            # Server setup & routing
    │   ├── handlers.go          # Public handlers
    │   ├── admin_handlers.go    # Admin handlers
    │   ├── auth_middleware.go   # Auth middleware
    │   ├── templates.go         # Template embedding
    │   ├── static/              # Static assets (embedded)
    │   │   ├── css/
    │   │   │   └── main.css
    │   │   ├── js/
    │   │   │   └── main.js
    │   │   ├── images/
    │   │   └── favicon.png
    │   └── templates/           # HTML templates (embedded)
    │       ├── base.html
    │       ├── home.html
    │       └── {other}.html
    ├── utils/                   # Utility functions (optional)
    │   └── network.go
    └── main.go                  # Entry point
```

**Key Directories:**

| Directory | Purpose | Contents |
|-----------|---------|----------|
| `src/data/` | JSON data files | ONLY `.json` files (no code) |
| `src/{service}/` | Service logic | Service implementation, loads from `src/data/` |
| `src/database/` | Database layer | DB connection, auth, settings |
| `src/paths/` | OS detection | Platform-specific directory paths |
| `src/server/` | HTTP server | Routing, handlers, templates, static files |
| `src/utils/` | Utilities | Helper functions (network, etc.) |
| `docs/` | Documentation | API.md, SERVER.md, README.md (optional) |
| `scripts/` | Production scripts | install.sh, backup.sh, uninstall.sh (optional) |
| `test/` | Test files | test-docker.sh, unit/, integration/, e2e/ (optional) |
| `binaries/` | Build output | Platform-specific binaries (gitignored) |
| `releases/` | Release artifacts | Files for GitHub releases (gitignored) |
| `rootfs/` | Docker volumes | Persistent storage for containers (gitignored) |

**Required Root Files:**
- ✅ `CLAUDE.md` - Project specification
- ✅ `README.md` - User documentation
- ✅ `LICENSE.md` - License file
- ✅ `Makefile` - Build system (4 targets)
- ✅ `Dockerfile` - Alpine-based container definition
- ✅ `docker-compose.yml` - Production compose
- ✅ `docker-compose.test.yml` - Development compose
- ✅ `Jenkinsfile` - Jenkins CI/CD pipeline (jenkins.casjay.cc)
- ✅ `.github/workflows/release.yml and docker.yml` - GitHub Actions (push & monthly)
- ✅ `.readthedocs.yml` - ReadTheDocs configuration
- ✅ `release.txt` - Version tracking (semantic versioning)
- ✅ `go.mod` / `go.sum` - Go modules
- ✅ `.gitignore` - Git ignore patterns
- ✅ `.gitattributes` - Git attributes (optional)

**Required Source Packages:**
- ✅ `src/main.go` - Application entry point
- ✅ `src/database/` - Database layer with auth
- ✅ `src/paths/` - OS-specific path detection
- ✅ `src/server/` - HTTP server with templates & static files
- ✅ `src/{service}/` - Primary service logic
- ✅ `src/data/` - JSON data files (if applicable)

**Optional Source Packages:**
- `src/auth/` - Additional authentication logic
- `src/utils/` - Utility functions (network helpers, etc.)
- `src/scheduler/` - Task scheduler for periodic jobs

**Optional Directories:**
- `docs/` - **All API type documentation** (API.md, SERVER.md, etc.)
  - ReadTheDocs compatible
  - MkDocs with Material theme (Dracula color scheme)
  - Includes: index.md, mkdocs.yml, requirements.txt
- `scripts/` - **All production and install scripts** (install.sh, backup.sh, etc.)
- `test/` - **All test files** (test-docker.sh, unit/, integration/, e2e/)

---

### 13. **Testing, Debugging & Build Environment Priority** ⭐ NEW

**Critical Rule**: ALWAYS use `/tmp/` for all temporary files, test data, and debugging. NEVER write test data to production directories.


**Rule**: Building: ALWAYS use Docker. Testing/Debugging: Prefer Incus, fallback to Docker, last resort Host. Never test directly on the host OS.

**Environment Selection:**


**For Building** (make build, cross-compilation):

- ✅ **Docker** - ALWAYS use Docker (golang:alpine builder)

- ❌ Never use Incus or Host OS for builds

- Makefile already uses Docker for all builds


**For Testing/Debugging** (running services, integration tests):

- 1️⃣ **Incus** (preferred) - System containers, full OS environment

- 2️⃣ **Docker** (fallback) - If Incus unavailable

- 3️⃣ **Host OS** (last resort) - Only when containers unavailable, tools may need installation


**Priority Order:**
1. **Incus** (if available) - Preferred
2. **Docker** (if available) - Alternative containerization
3. **Host OS** - Last resort only

---

**Testing Workflow:**

**1. Incus (Preferred)**

```bash
# Build development image
make docker-dev

# Run with docker-compose (test configuration)
docker-compose -f docker-compose.test.yml up -d

# Access service
curl http://localhost:64181/healthz

# View logs
docker-compose -f docker-compose.test.yml logs -f

# Cleanup
docker-compose -f docker-compose.test.yml down
sudo rm -rf /tmp/{projectname}/rootfs
```

**Benefits:**
- ✅ Isolated environment
- ✅ Consistent across all developers
- ✅ Easy cleanup (ephemeral /tmp storage)
- ✅ No pollution of host system
- ✅ Matches production environment

---

**2. Incus (Alternative)**

If Docker is not available, use Incus containers:

```bash
# Launch container
incus launch images:alpine/3.19 {projectname}-test

# Copy binary
incus file push ./binaries/{projectname} {projectname}-test/usr/local/bin/

# Execute
incus exec {projectname}-test -- /usr/local/bin/{projectname} --port 8080

# Shell access
incus shell {projectname}-test

# Cleanup
incus delete -f {projectname}-test
```

**Benefits:**
- ✅ System container (not Docker)
- ✅ LXD/LXC based
- ✅ Full OS environment
- ✅ Easy snapshot/restore

---

**3. Host OS (Last Resort)**

Only use when containers are unavailable:

```bash
# Build for host
make build

# Run directly
./binaries/{projectname}  --port 8080

# With custom directories (avoid polluting system)
./binaries/{projectname} \
  --config /tmp/{projectname}/config \
  --data /tmp/{projectname}/data \
  --logs /tmp/{projectname}/logs \
  --port 8080

# Cleanup
rm -rf /tmp/{projectname}
```

**⚠️ Warnings:**
- May pollute host system directories
- Different from production environment
- Harder to cleanup
- Not reproducible across systems

---

**Make Targets for Testing:**

```makefile
# Run tests in Docker (always)
test:
	@docker run --rm -v $$(pwd):/workspace -w /workspace golang:alpine \
		sh -c 'go test -v -race -timeout 5m ./...'

# Build development Docker image
docker-dev:
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		-t $(PROJECTNAME):dev \
		.

# Test with docker-compose
test-docker:
	@docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for service..."
	@timeout 30 bash -c 'until curl -sf http://localhost:64181/healthz; do sleep 1; done'
	@echo "✓ Service is running"
	@docker-compose -f docker-compose.test.yml logs
	@docker-compose -f docker-compose.test.yml down
	@sudo rm -rf /tmp/{projectname}/rootfs

# Test with Incus (if Docker not available)
test-incus:
	@incus launch images:alpine/3.19 $(PROJECTNAME)-test
	@incus file push ./binaries/$(PROJECTNAME) $(PROJECTNAME)-test/usr/local/bin/
	@incus exec $(PROJECTNAME)-test -- /usr/local/bin/$(PROJECTNAME) --status
	@incus delete -f $(PROJECTNAME)-test
```

---

**CI/CD Testing:**

**GitHub Actions:**
- Always runs in Docker containers (ubuntu-latest)
- Uses `make test` which runs tests in Docker

**Jenkins:**
- Runs on dedicated agents (amd64, arm64)
- Each agent runs in isolated environment
- Uses Docker for builds

**Local Development:**
- Developers MUST use Docker or Incus
- Host OS testing only when troubleshooting container issues
- Document when host testing is used (unusual case)

---

**Best Practices:**

✅ **DO:**
- Always test in Docker first
- Use docker-compose.test.yml for integration tests
- Use ephemeral storage (/tmp) for test data
- Clean up after tests
- Document test environment used

❌ **DON'T:**
- Test directly on host OS without containers
- Leave test data in system directories
- Mix development and production configs
- Skip cleanup after tests
- Assume host environment matches production


### Testing Environment Priority
1. **Incus** (preferred) - `incus launch images:alpine/3.19 {project}-test`
2. **Docker** (alternative) - `docker-compose -f docker-compose.test.yml up -d`
3. **Host OS** (last resort) - `./binaries/{projectname} `

Always use isolated environments. Never test directly on host OS unless containers unavailable.


---

## 2. URL Display Standards

**Rule**: Never show `localhost`, `127.0.0.1`, or `0.0.0.0` to users. Always display the most relevant accessible URL.

**Priority**:
1. **FQDN** (if hostname resolves)
2. **Public IP** (outbound IP)
3. **Hostname** (if available)
4. **Fallback** (`<your-host>`)

**Implementation** - Add to `src/database/credentials.go`:

```go
// getAccessibleURL returns the most relevant URL for accessing the server
// Priority: FQDN > hostname > public IP > fallback
// NEVER shows localhost, 127.0.0.1, or 0.0.0.0
func getAccessibleURL(port string) string {
	// Try to get hostname
	hostname, err := os.Hostname()
	if err == nil && hostname != "" && hostname != "localhost" {
		// Try to resolve hostname to see if it's a valid FQDN
		if addrs, err := net.LookupHost(hostname); err == nil && len(addrs) > 0 {
			return fmt.Sprintf("http://%s:%s", hostname, port)
		}
	}

	// Try to get outbound IP (most likely accessible IP)
	if ip := getOutboundIP(); ip != "" {
		return fmt.Sprintf("http://%s:%s", ip, port)
	}

	// Fallback to hostname if we have one
	if hostname != "" && hostname != "localhost" {
		return fmt.Sprintf("http://%s:%s", hostname, port)
	}

	// Last resort: use a generic message
	return fmt.Sprintf("http://<your-host>:%s", port)
}

// getOutboundIP gets the preferred outbound IP of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
```

**Update** `SaveCredentialsToFile` signature to accept port:
```go
func SaveCredentialsToFile(creds *AdminCredentials, configDir, port string) error {
	serverURL := getAccessibleURL(port)

	content := fmt.Sprintf(`{PROJECT_NAME} - ADMIN CREDENTIALS
========================================
WEB UI LOGIN:
  URL:      %s/admin
  Username: %s

API ACCESS:
  URL:      %s/api/v1/admin
  Header:   Authorization: Bearer %s

CREDENTIALS:
  Username: %s
  Password: %s
  Token:    %s

Created: %s
========================================`, serverURL, creds.Username, serverURL, creds.Token,
    creds.Username, creds.Password, creds.Token, time.Now().Format("2006-01-02 15:04:05"))
	//...
}
```

**In main.go** - Save credentials AFTER port is determined:
```go
// Port resolution first
if port == "" {
	// ... port determination logic
}

// THEN save credentials with actual port
if creds.Token != "" {
	if err := database.SaveCredentialsToFile(creds, configDir, port); err != nil {
		log.Printf("Warning: Failed to save credentials file: %v", err)
	} else {
		log.Printf("⚠️  Access URL: %s", getAccessibleURL(port))
	}
}
```

---

## 3. Dockerfile

**Replace**: `FROM scratch` runtime
**With**: `FROM alpine:latest` runtime

**Why**: Need curl and bash for health checks and debugging

```dockerfile
# ============================================
# Build stage
# ============================================
FROM golang:alpine AS builder

# Why latest? Static binaries are self-contained, no Go runtime dependency.
# Using :latest gets latest Go features, security patches, and build optimizations.

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/

# Build static binary with all assets embedded
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE} -w -s" \
    -a -installsuffix cgo \
    -o {projectname} \
    ./src

# ============================================
# Runtime stage - Alpine with minimal tools
# ============================================
FROM alpine:latest

# Why latest? Static binaries (CGO_ENABLED=0) have no runtime dependencies.
# Alpine version only affects runtime tools (curl, bash, ca-certificates).
# Using :latest ensures latest security patches without version maintenance.

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# Install runtime dependencies (curl, bash)
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    bash \
    && rm -rf /var/cache/apk/*

# Copy binary to /usr/local/bin
COPY --from=builder /build/{projectname} /usr/local/bin/{projectname}
RUN chmod +x /usr/local/bin/{projectname}

# Environment variables
ENV PORT=80 \
    CONFIG_DIR=/config \
    DATA_DIR=/data \
    LOGS_DIR=/logs \
    ADDRESS=0.0.0.0 \
    DB_PATH=/data/db/{projectname}.db

# Create directories
RUN mkdir -p /config /data /data/db /logs && \
    chown -R 65534:65534 /config /data /logs

# Metadata labels (OCI standard)
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.authors="{organization}" \
      org.opencontainers.image.url="https://github.com/{organization}/{projectname}" \
      org.opencontainers.image.source="https://github.com/{organization}/{projectname}" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.vendor="{organization}" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.title="{projectname}" \
      org.opencontainers.image.description="{Project description} - Single static binary" \
      org.opencontainers.image.documentation="https://github.com/{organization}/{projectname}/blob/main/docs/README.md" \
      org.opencontainers.image.base.name="alpine:latest"

# Expose default port
EXPOSE 80

# Create mount points for volumes
VOLUME ["/config", "/data", "/logs"]

# Run as non-root user (nobody)
USER 65534:65534

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/{projectname}", "--status"]

# Run
ENTRYPOINT ["/usr/local/bin/{projectname}"]
CMD ["--port", "80"]
```

**Key Changes**:
- Runtime base: `alpine:latest` (not scratch)
- Includes: curl, bash, ca-certificates, tzdata
- Binary location: `/usr/local/bin/{projectname}` (not root)
- SQLite DB location: `/data/db/{projectname}.db`
- Internal port: 80
- All OCI metadata labels

---

## 4. docker-compose.yml

**Standards**:
- ❌ NO `version:` field
- ❌ NO `build:` definition
- ✅ Use pre-built images from registry
- ✅ Custom network: `{projectname}` (external: false)
- ✅ Volume structure: `./rootfs/{type}/{servicename}`
- ✅ Production port: `172.17.0.1:{randomport}:80`

```yaml
# Production Docker Compose Configuration
# Uses ./rootfs for persistent storage
# External port: 172.17.0.1:64180:80 (Docker bridge network)

services:
  {projectname}:
    image: ghcr.io/{organization}/{projectname}:latest
    container_name: {projectname}
    restart: unless-stopped

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/{projectname}.db
      # Uncomment and set for first deployment
      #- ADMIN_USER=administrator
      #- ADMIN_PASSWORD=changeme
      #- ADMIN_TOKEN=your-token-here

    volumes:
      - ./rootfs/config/{projectname}:/config
      - ./rootfs/data/{projectname}:/data
      - ./rootfs/logs/{projectname}:/logs

    ports:
      - "172.17.0.1:64180:80"

    networks:
      - {projectname}

    healthcheck:
      test: ["CMD", "/usr/local/bin/{projectname}", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

  # Optional: PostgreSQL database
  # postgres:
  #   image: postgres:16-alpine
  #   container_name: {projectname}-postgres
  #   restart: unless-stopped
  #
  #   environment:
  #     POSTGRES_DB: {projectname}
  #     POSTGRES_USER: {projectname}
  #     POSTGRES_PASSWORD: changeme
  #
  #   volumes:
  #     - ./rootfs/db/postgres:/var/lib/postgresql/data
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD-SHELL", "pg_isready -U {projectname}"]
  #     interval: 10s
  #     timeout: 5s
  #     retries: 5

  # Optional: MariaDB database
  # mariadb:
  #   image: mariadb:11-alpine
  #   container_name: {projectname}-mariadb
  #   restart: unless-stopped
  #
  #   environment:
  #     MYSQL_DATABASE: {projectname}
  #     MYSQL_USER: {projectname}
  #     MYSQL_PASSWORD: changeme
  #     MYSQL_RANDOM_ROOT_PASSWORD: "yes"
  #
  #   volumes:
  #     - ./rootfs/db/mariadb:/var/lib/mysql
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD", "healthcheck.sh", "--connect", "--innodb_initialized"]
  #     interval: 10s
  #     timeout: 5s
  #     retries: 5

  # Optional: Redis cache
  # redis:
  #   image: redis:7-alpine
  #   container_name: {projectname}-redis
  #   restart: unless-stopped
  #
  #   volumes:
  #     - ./rootfs/db/redis:/data
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD", "redis-cli", "ping"]
  #     interval: 10s
  #     timeout: 3s
  #     retries: 3

networks:
  {projectname}:
    name: {projectname}
    external: false
    driver: bridge
```

---

## 5. docker-compose.test.yml

**Standards**:
- ❌ NO `version:` field
- ❌ NO `build:` definition (use `{projectname}:dev` image)
- ✅ Ephemeral storage: `/tmp/{projectname}/rootfs`
- ✅ Development port: `{randomport}:80` (e.g., 64181:80)
- ✅ Same network name: `{projectname}`

```yaml
# Development/Testing Docker Compose Configuration
# Uses /tmp for ephemeral storage
# External port: 64181:80 (random unused port in 64xxx range)

services:
  {projectname}:
    image: {projectname}:dev
    container_name: {projectname}-test
    restart: "no"

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/{projectname}.db
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=testpass123
      - DEV=true

    volumes:
      - /tmp/{projectname}/rootfs/config/{projectname}:/config
      - /tmp/{projectname}/rootfs/data/{projectname}:/data
      - /tmp/{projectname}/rootfs/logs/{projectname}:/logs

    ports:
      - "64181:80"

    networks:
      - {projectname}

    healthcheck:
      test: ["CMD", "/usr/local/bin/{projectname}", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

  # Optional: PostgreSQL for testing
  # postgres:
  #   image: postgres:16-alpine
  #   container_name: {projectname}-postgres-test
  #   restart: "no"
  #
  #   environment:
  #     POSTGRES_DB: {projectname}_test
  #     POSTGRES_USER: {projectname}
  #     POSTGRES_PASSWORD: testpass
  #
  #   volumes:
  #     - /tmp/{projectname}/rootfs/db/postgres:/var/lib/postgresql/data
  #
  #   networks:
  #     - {projectname}

networks:
  {projectname}:
    name: {projectname}
    external: false
    driver: bridge
```

---

## 6. Makefile

Add `docker-dev` target for local development:

```makefile
# Build and push multi-platform Docker images (release)
docker:
	@echo "Building multi-platform Docker images..."
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):latest \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION) \
		--push \
		.
	@echo "✓ Docker images pushed to ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION)"

# Build Docker image for development (local only, not pushed)
docker-dev:
	@echo "Building development Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(PROJECTNAME):dev \
		.
	@echo "✓ Docker development image built: $(PROJECTNAME):dev"
```

---

## 7. Jenkinsfile

Multi-architecture CI/CD pipeline:

```groovy
// Jenkinsfile for {projectname}
// Multi-architecture build pipeline (amd64, arm64)
// Server: jenkins.casjay.cc

pipeline {
    agent none

    environment {
        PROJECTNAME = '{projectname}'
        PROJECTORG = '{organization}'
        REGISTRY = 'ghcr.io'
        IMAGE_NAME = "${REGISTRY}/${PROJECTORG}/${PROJECTNAME}"

        VERSION = sh(script: 'cat release.txt 2>/dev/null || echo "0.0.1"', returnStdout: true).trim()
        COMMIT = sh(script: 'git rev-parse --short HEAD 2>/dev/null || echo "unknown"', returnStdout: true).trim()
        BUILD_DATE = sh(script: 'date -u +%Y-%m-%dT%H:%M:%SZ', returnStdout: true).trim()
    }

    stages {
        stage('Checkout') {
            agent { label 'amd64' }
            steps {
                checkout scm
                echo "Building ${PROJECTNAME} ${VERSION} (${COMMIT})"
            }
        }

        stage('Test') {
            parallel {
                stage('Test on AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Running tests on AMD64...'
                        sh 'make test'
                    }
                }
                stage('Test on ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Running tests on ARM64...'
                        sh 'make test'
                    }
                }
            }
        }

        stage('Build Binaries') {
            parallel {
                stage('Build on AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Building binaries on AMD64...'
                        sh 'make build'
                        stash includes: 'binaries/**/*', name: 'binaries-amd64'
                        stash includes: 'releases/**/*', name: 'release-amd64'
                    }
                }
                stage('Build on ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Building binaries on ARM64...'
                        sh 'make build'
                        stash includes: 'binaries/**/*', name: 'binaries-arm64'
                        stash includes: 'releases/**/*', name: 'release-arm64'
                    }
                }
            }
        }

        stage('Build Docker Images') {
            parallel {
                stage('Build Docker AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Building Docker image for AMD64...'
                        script {
                            sh """
                                docker build \
                                    --platform linux/amd64 \
                                    --build-arg VERSION=${VERSION} \
                                    --build-arg COMMIT=${COMMIT} \
                                    --build-arg BUILD_DATE=${BUILD_DATE} \
                                    -t ${IMAGE_NAME}:${VERSION}-amd64 \
                                    -t ${IMAGE_NAME}:latest-amd64 \
                                    .
                            """
                        }
                    }
                }
                stage('Build Docker ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Building Docker image for ARM64...'
                        script {
                            sh """
                                docker build \
                                    --platform linux/arm64 \
                                    --build-arg VERSION=${VERSION} \
                                    --build-arg COMMIT=${COMMIT} \
                                    --build-arg BUILD_DATE=${BUILD_DATE} \
                                    -t ${IMAGE_NAME}:${VERSION}-arm64 \
                                    -t ${IMAGE_NAME}:latest-arm64 \
                                    .
                            """
                        }
                    }
                }
            }
        }

        stage('Push Docker Images') {
            agent { label 'amd64' }
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                echo 'Pushing Docker images to registry...'
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'github-registry',
                        usernameVariable: 'REGISTRY_USER',
                        passwordVariable: 'REGISTRY_TOKEN'
                    )]) {
                        sh """
                            echo \$REGISTRY_TOKEN | docker login ${REGISTRY} -u \$REGISTRY_USER --password-stdin

                            docker push ${IMAGE_NAME}:${VERSION}-amd64
                            docker push ${IMAGE_NAME}:${VERSION}-arm64
                            docker push ${IMAGE_NAME}:latest-amd64
                            docker push ${IMAGE_NAME}:latest-arm64

                            docker manifest create ${IMAGE_NAME}:${VERSION} \
                                ${IMAGE_NAME}:${VERSION}-amd64 \
                                ${IMAGE_NAME}:${VERSION}-arm64

                            docker manifest create ${IMAGE_NAME}:latest \
                                ${IMAGE_NAME}:latest-amd64 \
                                ${IMAGE_NAME}:latest-arm64

                            docker manifest push ${IMAGE_NAME}:${VERSION}
                            docker manifest push ${IMAGE_NAME}:latest

                            docker logout ${REGISTRY}
                        """
                    }
                }
            }
        }

        stage('GitHub Release') {
            agent { label 'amd64' }
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                echo 'Creating GitHub release...'
                script {
                    unstash 'release-amd64'
                    withCredentials([string(
                        credentialsId: 'github-token',
                        variable: 'GITHUB_TOKEN'
                    )]) {
                        sh 'make release'
                    }
                }
            }
        }
    }

    post {
        success {
            echo "✅ Build successful for ${PROJECTNAME} ${VERSION}"
        }
        failure {
            echo "❌ Build failed for ${PROJECTNAME} ${VERSION}"
        }
        always {
            cleanWs()
        }
    }
}
```

---

## 8. src/data Directory

**Rule**: `src/data/` directory contains ONLY JSON files. No Go code, no other file types.

**Structure**:
```
src/
├── main.go                      # Embeds: //go:embed data/*.json
├── data/
│   ├── {datafile}.json          # JSON data ONLY
│   └── {other}.json             # Additional JSON files (if needed)
└── {service}/
    └── service.go               # NewService(jsonData []byte)
```

**Embedding from main.go:**

```go
// src/main.go
package main

import (
    _ "embed"
    "github.com/{org}/{project}/src/{service}"
)

//go:embed data/{datafile}.json
var jsonData []byte

func main() {
    svc, err := service.NewService(jsonData)
    if err != nil {
        log.Fatal(err)
    }
    // ...
}
```

**Service receives embedded data:**

```go
// src/{service}/service.go
package service

import "encoding/json"

type Service struct {
    data MyDataType
}

func NewService(jsonData []byte) (*Service, error) {
    var data MyDataType
    err := json.Unmarshal(jsonData, &data)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    return &Service{data: data}, nil
}
```

**Why This Pattern:**
- ✅ `src/data/` contains ONLY JSON files (no .go code)
- ✅ JSON is embedded in single static binary
- ✅ No copies, no symlinks, no duplicates
- ✅ Embedding happens from `main.go` at `src/` level
- ✅ Services receive data as parameter (clean dependency injection)
- ✅ True single binary with all assets embedded

**Example** (airports project):
```
src/
├── main.go                  # //go:embed data/airports.json
├── data/
│   └── airports.json        # 8.7MB (ONLY JSON, no .go files)
└── airports/
    └── data.go              # NewService(jsonData []byte)
```

**Important**:
- The `src/data/` directory contains ONLY JSON files
- NO .go files in `src/data/`
- NO symlinks or copies
- Embedding MUST be done from `src/main.go`
- Services accept `[]byte` parameter with embedded data

---

## 9. README.md Structure

**Order**: About → Production → Docker → API Usage → Development

```markdown
# {PROJECT_NAME}

Brief description

## About
- Project description
- Key features
- Technology stack

## Official Links
- Website (if exists)
- Documentation
- GitHub
- Docker Hub
- Issues

## Production Deployment

### Installation
- Linux (all distros, all init systems)
- macOS (launchd)
- FreeBSD/OpenBSD/NetBSD (rc.d)
- Windows (NSSM service)
- Docker

### Binary Installation
- Platform-specific downloads
- Features list
- Key capabilities

- Production deployment
- Binary installation
- Systemd service
- Environment variables

- Docker deployment
- Docker Compose (production)
- Docker Compose (development)
- Docker run examples

- API usage
- Quick examples
- Admin panel
- Documentation links

- Development

- Requirements
- Build System & Testing (Makefile targets, platforms, versioning)
- Development Mode (dev flags, debug features)
- CI/CD (GitHub Actions, Jenkins, ReadTheDocs)

## Legal
- License (MIT - LICENSE.md)
- Third-party licenses
- Credits
- Support
```

---

## 10. Complete Project Layout

**Standard Go API Server Project Structure** (matches citylist/airports):

```
{projectname}/
├── .claude/
│   └── settings.local.json      # Claude Code settings
├── .github/
│   └── workflows/
│       └── build.yml            # GitHub Actions (build on push & monthly)
├── .gitattributes               # Git attributes
├── .gitignore                   # Git ignore patterns
├── .readthedocs.yml             # ReadTheDocs configuration
├── CLAUDE.md                    # Project specification (this file)
├── Dockerfile                   # Alpine-based multi-stage build
├── docker-compose.yml           # Production compose
├── docker-compose.test.yml      # Development/testing compose
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Jenkinsfile                  # CI/CD pipeline (jenkins.casjay.cc)
├── LICENSE.md                   # License file
├── Makefile                     # Build system (4 targets)
├── README.md                    # User documentation
├── release.txt                  # Version tracking (auto-increment)
├── binaries/                    # Built binaries (gitignored)
│   ├── {projectname}-linux-amd64
│   ├── {projectname}-linux-arm64
│   ├── {projectname}-windows-amd64.exe
│   ├── {projectname}-windows-arm64.exe
│   ├── {projectname}-macos-amd64
│   ├── {projectname}-macos-arm64
│   ├── {projectname}-bsd-amd64
│   ├── {projectname}-bsd-arm64
│   └── {projectname}            # Host platform binary
├── releases/                    # Release artifacts (gitignored)
├── rootfs/                      # Docker volumes (gitignored)
│   ├── config/
│   │   └── {projectname}/       # Service config
│   ├── data/
│   │   └── {projectname}/       # Service data
│   ├── logs/
│   │   └── {projectname}/       # Service logs
│   └── db/                      # External databases
│       ├── postgres/
│       ├── mariadb/
│       └── redis/
├── docs/                        # API documentation (optional)
│   ├── index.md                 # Documentation home
│   ├── API.md                   # Complete API reference
│   ├── SERVER.md                # Server administration guide
│   ├── README.md                # Documentation index
│   ├── mkdocs.yml               # MkDocs configuration (Dracula theme)
│   ├── requirements.txt         # Python dependencies for RTD
│   ├── stylesheets/
│   │   └── dracula.css          # Dracula theme CSS
│   ├── javascripts/
│   │   └── extra.js             # Custom JavaScript
│   └── overrides/               # MkDocs theme overrides
├── scripts/                     # Production scripts (optional)
│   ├── install.sh               # Installation script
│   ├── backup.sh                # Backup script
│   └── uninstall.sh             # Uninstallation script
├── test/                        # Test files (optional)
│   ├── test-docker.sh           # Docker testing script
│   ├── unit/                    # Unit tests
│   ├── integration/             # Integration tests
│   └── e2e/                     # End-to-end tests
└── src/                         # Source code
    ├── data/
    │   └── {datafile}.json      # JSON data ONLY (no .go files)
    ├── {service}/
    │   └── service.go           # Service logic
    ├── auth/                    # Authentication (optional)
    │   └── auth.go
    ├── database/
    │   ├── database.go          # DB connection & schema
    │   ├── auth.go              # Admin authentication
    │   ├── credentials.go       # Credential management
    │   └── settings.go          # Settings CRUD
    ├── paths/
    │   └── paths.go             # OS-specific directory detection
    ├── scheduler/               # Task scheduler (optional)
    │   └── scheduler.go
    ├── server/
    │   ├── server.go            # Server setup & routing
    │   ├── handlers.go          # Public handlers
    │   ├── admin_handlers.go    # Admin handlers
    │   ├── auth_middleware.go   # Auth middleware
    │   ├── templates.go         # Template embedding
    │   ├── static/              # Static assets (embedded)
    │   │   ├── css/
    │   │   │   └── main.css
    │   │   ├── js/
    │   │   │   └── main.js
    │   │   ├── images/
    │   │   └── favicon.png
    │   └── templates/           # HTML templates (embedded)
    │       ├── base.html
    │       ├── home.html
    │       └── {other}.html
    ├── utils/                   # Utility functions (optional)
    │   └── network.go
    └── main.go                  # Entry point
```

**Key Directories:**

| Directory | Purpose | Contents |
|-----------|---------|----------|
| `src/data/` | JSON data files | ONLY `.json` files (no code) |
| `src/{service}/` | Service logic | Service implementation, loads from `src/data/` |
| `src/database/` | Database layer | DB connection, auth, settings |
| `src/paths/` | OS detection | Platform-specific directory paths |
| `src/server/` | HTTP server | Routing, handlers, templates, static files |
| `src/utils/` | Utilities | Helper functions (network, etc.) |
| `docs/` | Documentation | API.md, SERVER.md, README.md (optional) |
| `scripts/` | Production scripts | install.sh, backup.sh, uninstall.sh (optional) |
| `test/` | Test files | test-docker.sh, unit/, integration/, e2e/ (optional) |
| `binaries/` | Build output | Platform-specific binaries (gitignored) |
| `releases/` | Release artifacts | Files for GitHub releases (gitignored) |
| `rootfs/` | Docker volumes | Persistent storage for containers (gitignored) |

**Required Root Files:**
- ✅ `CLAUDE.md` - Project specification
- ✅ `README.md` - User documentation
- ✅ `LICENSE.md` - License file
- ✅ `Makefile` - Build system (4 targets)
- ✅ `Dockerfile` - Alpine-based container definition
- ✅ `docker-compose.yml` - Production compose
- ✅ `docker-compose.test.yml` - Development compose
- ✅ `Jenkinsfile` - Jenkins CI/CD pipeline (jenkins.casjay.cc)
- ✅ `.github/workflows/release.yml and docker.yml` - GitHub Actions (push & monthly)
- ✅ `.readthedocs.yml` - ReadTheDocs configuration
- ✅ `release.txt` - Version tracking (semantic versioning)
- ✅ `go.mod` / `go.sum` - Go modules
- ✅ `.gitignore` - Git ignore patterns
- ✅ `.gitattributes` - Git attributes (optional)

**Required Source Packages:**
- ✅ `src/main.go` - Application entry point
- ✅ `src/database/` - Database layer with auth
- ✅ `src/paths/` - OS-specific path detection
- ✅ `src/server/` - HTTP server with templates & static files
- ✅ `src/{service}/` - Primary service logic
- ✅ `src/data/` - JSON data files (if applicable)

**Optional Source Packages:**
- `src/auth/` - Additional authentication logic
- `src/utils/` - Utility functions (network helpers, etc.)
- `src/scheduler/` - Task scheduler for periodic jobs

**Optional Directories:**
- `docs/` - **All API type documentation** (API.md, SERVER.md, etc.)
  - ReadTheDocs compatible
  - MkDocs with Material theme (Dracula color scheme)
  - Includes: index.md, mkdocs.yml, requirements.txt
- `scripts/` - **All production and install scripts** (install.sh, backup.sh, etc.)
- `test/` - **All test files** (test-docker.sh, unit/, integration/, e2e/)

---

## 11. ReadTheDocs

**Purpose**: Automatic documentation hosting with Dracula theme

**Files Required:**

1. **`.readthedocs.yml`** (root):
```yaml
version: 2

mkdocs:
  configuration: docs/mkdocs.yml
  fail_on_warning: false

formats:
  - pdf
  - epub

python:
  version: "3.11"
  install:
    - requirements: docs/requirements.txt
```

2. **`docs/mkdocs.yml`**:
```yaml
site_name: {projectname} Documentation
site_description: {Project description}
site_author: {organization}
site_url: https://{projectname}.readthedocs.io

repo_name: {organization}/{projectname}
repo_url: https://github.com/{organization}/{projectname}
edit_uri: edit/main/docs/

copyright: Copyright &copy; 2024 {organization}

theme:
  name: material
  custom_dir: overrides
  palette:
    scheme: dracula
    primary: deep purple
    accent: pink
  font:
    text: Roboto
    code: Fira Code
  features:
    - navigation.instant
    - navigation.tracking
    - navigation.tabs
    - navigation.sections
    - navigation.expand
    - navigation.top
    - search.suggest
    - search.highlight
    - search.share
    - content.code.copy
    - content.code.annotate
  icon:
    repo: fontawesome/brands/github

nav:
  - Home: index.md
  - API Reference: API.md
  - Server Guide: SERVER.md
  - GitHub: https://github.com/{organization}/{projectname}

markdown_extensions:
  - admonition
  - pymdownx.details
  - pymdownx.superfences
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.inlinehilite
  - pymdownx.snippets
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.tasklist:
      custom_checkbox: true
  - pymdownx.emoji
  - toc:
      permalink: true
  - tables
  - attr_list

plugins:
  - search
  - minify:
      minify_html: true

extra_css:
  - stylesheets/dracula.css

extra_javascript:
  - javascripts/extra.js
```

3. **`docs/requirements.txt`**:
```txt
mkdocs>=1.5.0
mkdocs-material>=9.5.0
pymdown-extensions>=10.7
mkdocs-minify-plugin>=0.8.0
mike>=2.0.0
```

4. **`docs/stylesheets/dracula.css`** - Dracula color scheme overrides
5. **`docs/javascripts/extra.js`** - Optional custom JavaScript

**Documentation Structure:**
```
docs/
├── index.md              # Documentation home
├── API.md                # API reference
├── SERVER.md             # Server administration
├── mkdocs.yml            # MkDocs configuration
├── requirements.txt      # Python dependencies
├── stylesheets/
│   └── dracula.css       # Dracula theme CSS
├── javascripts/
│   └── extra.js          # Custom JavaScript
└── overrides/            # Theme overrides (optional)
```

**Dracula Color Palette:**
- Background: `#282a36`
- Current Line: `#44475a`
- Foreground: `#f8f8f2`
- Comment: `#6272a4`
- Cyan: `#8be9fd`
- Green: `#50fa7b`
- Orange: `#ffb86c`
- Pink: `#ff79c6`
- Purple: `#bd93f9`
- Red: `#ff5555`
- Yellow: `#f1fa8c`

**Local Preview:**
```bash
cd docs
pip install -r requirements.txt
mkdocs serve
# Open http://localhost:8000
```

---

## 12. GitHub Actions

**Purpose**: Automated builds on push and monthly schedule (split into separate workflows)

**Files**:
- `.github/workflows/release.yml` - Binary builds and GitHub releases
- `.github/workflows/docker.yml` - Docker image builds and registry push

**Triggers** (both workflows):
- Push to `main` branch
- Monthly schedule (1st of month, 3:00 AM UTC)
- Manual trigger (workflow_dispatch)

---

**Workflow 1: release.yml** - Binary Release

**Jobs:**

1. **test**
   - Runs `make test` with Go 1.23
   - Validates all tests pass before building

2. **build-and-release**
   - Reads version from `release.txt` (does NOT update it)
   - Runs `make build` for all 8 platforms
   - Deletes existing release if exists (always recreate)
   - Creates new GitHub release `{VERSION}`
   - Attaches all platform binaries
   - Uploads artifacts (90 day retention)

---

**Workflow 2: docker.yml** - Docker Release

**Jobs:**

1. **build-and-push**
   - Reads version from `release.txt`
   - Builds multi-arch Docker images (amd64, arm64)
   - Pushes to `ghcr.io/{organization}/{projectname}`
   - Tags: `latest`, `{VERSION}`, `{branch}-{sha}`
   - Uses GitHub cache for faster builds
   - Verifies images after push

**Configuration:**

```yaml
name: Build and Release

on:
  push:
    branches:
      - main
      - master
  schedule:
    - cron: '0 3 1 * *'  # Monthly on 1st at 3 AM UTC
  workflow_dispatch:

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-binaries:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Read version
        run: |
          VERSION=$(cat release.txt)
          echo "VERSION=${VERSION}" >> $GITHUB_OUTPUT
      - name: Build
        run: VERSION=$VERSION make build
      - uses: actions/upload-artifact@v4
        with:
          name: binaries
          path: releases/*

  build-docker:
    runs-on: ubuntu-latest
    permissions:
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-qemu-action@v3
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:latest
            ghcr.io/${{ github.repository }}:${{ steps.version.outputs.VERSION }}

  create-release:
    needs: build-binaries
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/download-artifact@v4
        with:
          name: binaries
          path: releases/
      - name: Delete existing release
        run: gh release delete $VERSION -y || true
      - name: Create release
        run: gh release create $VERSION ./releases/* --title "$VERSION"
```

**Important:**
- Version is read from `release.txt`, NOT modified by Actions
- Release is always recreated (deleted first if exists)
- Builds all platforms (8 binaries)
- Multi-arch Docker images (amd64, arm64)
- Runs on every push to main and monthly

---

## 🎯 Quick Reference

### Container Tags
- **Production**: `ghcr.io/{organization}/{projectname}:latest`
- **Versioned**: `ghcr.io/{organization}/{projectname}:{VERSION}`
- **Development**: `{projectname}:dev`

### Port Mappings
- **Production**: `172.17.0.1:{randomport}:80`
- **Development**: `{randomport}:80`
- **Internal**: Always `80`

### Volume Structure
```
./rootfs/
├── config/{servicename}/
├── data/{servicename}/
├── logs/{servicename}/
└── db/{dbtype}/
```

### Releases Directory Structure



**./releases/** contains files for GitHub releases:



```

./releases/

├── {projectname}-linux-amd64

├── {projectname}-linux-arm64

├── {projectname}-macos-amd64

├── {projectname}-macos-arm64

├── {projectname}-windows-amd64.exe

├── {projectname}-windows-arm64.exe

├── {projectname}-bsd-amd64

├── {projectname}-bsd-arm64

├── {projectname}-{VERSION}-src.tar.gz    # Source archive (no VCS)

└── {projectname}-{VERSION}-src.zip       # Source archive for Windows

```



**Source Archives:**

- Created with `git archive` (excludes .git, .github, etc.)

- Both tar.gz (Linux/macOS) and zip (Windows) formats

- Includes all source code, docs, configs

- No binaries, no VCS files, no build artifacts




### Binary Requirements
- ✅ Static binary (CGO_ENABLED=0)
- ✅ All assets embedded:
  - HTML templates (go:embed in server package)
  - CSS, JS, images (go:embed in server package)
  - JSON data files (go:embed in main.go)
- ✅ True single binary - no external files needed
- ✅ Location: `/usr/local/bin/{projectname}`
- ✅ Output: `./binaries/` or `./releases/`

### URL Display
- ❌ Never: localhost, 127.0.0.1, 0.0.0.0
- ✅ Priority: FQDN → hostname → public IP → fallback

### Data Loading (JSON Files)
```
src/
├── main.go                # Embeds: //go:embed data/{file}.json
├── data/
│   └── {datafile}.json    # JSON files ONLY (no .go files)
└── {service}/
    └── service.go         # NewService(jsonData []byte)
```
- ✅ `src/data/` contains ONLY JSON files
- ✅ NO Go code in `src/data/`
- ✅ Embedded in main.go, passed to services as []byte
- ✅ True single binary - JSON embedded, not external files

---

**Apply these updates to SPEC.md for production-grade specifications** 🚀

---

## 13. Web UI / Frontend

**Purpose**: Modern, responsive dark-themed web interface with Go html/template

**Technology Stack:**
- Go `html/template` for server-side rendering
- CSS3 with CSS custom properties (variables)
- Vanilla JavaScript (no frameworks)
- Embedded assets (go:embed static/* and templates/*)

**Note**: Static assets (CSS, JS, images), HTML templates, AND JSON data files are ALL embedded via `go:embed` in the single static binary. JSON files are embedded from `main.go` and passed to services.

---

**Structure:**
```
src/server/
├── static/
│   ├── css/
│   │   └── main.css        # Main stylesheet (~900 lines)
│   ├── js/
│   │   └── main.js         # JavaScript utilities (~130 lines)
│   ├── images/
│   │   └── favicon.png
│   └── manifest.json       # PWA manifest
├── templates/
│   ├── base.html           # Base template with header/footer
│   ├── home.html           # Homepage
│   ├── search.html         # Search page
│   ├── {feature}.html      # Feature-specific pages
│   └── config.html         # Admin configuration page
└── templates.go            # Template embedding logic
```

---

**CSS Framework** - `src/server/static/css/main.css`

**Design System:**

```css
/* CSS Variables (Dark Theme Default) */
:root {
  /* Colors */
  --bg-primary: #1a1a1a;        /* Main background */
  --bg-secondary: #2d2d2d;      /* Cards, sections */
  --bg-tertiary: #404040;       /* Hover states */
  --text-primary: #ffffff;      /* Main text */
  --text-secondary: #b0b0b0;    /* Secondary text */
  --text-tertiary: #808080;     /* Disabled text */
  
  /* Accent Colors */
  --accent-primary: #0066cc;    /* Primary actions */
  --accent-success: #00aa00;    /* Success states */
  --accent-warning: #ff9900;    /* Warnings */
  --accent-danger: #cc0000;     /* Errors */
  
  /* Layout */
  --border-color: #404040;
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  
  /* Typography */
  --font-primary: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  
  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-normal: 300ms ease;
  
  /* Shadows */
  --shadow-sm: 0 1px 3px rgba(0,0,0,0.12);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.16);
  --shadow-lg: 0 10px 20px rgba(0,0,0,0.2);
}

/* Light Theme Override */
[data-theme="light"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --bg-tertiary: #e0e0e0;
  --text-primary: #1a1a1a;
  --text-secondary: #404040;
  --border-color: #e0e0e0;
}
```

---

**Components:**

**1. Buttons**
```css
.btn-primary {
  padding: var(--space-md) var(--space-xl);
  background: var(--accent-primary);
  border: none;
  border-radius: 8px;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-primary:hover {
  background: #0052a3;
  transform: translateY(-2px);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
```

**2. Cards**
```css
.card {
  background: var(--bg-secondary);
  padding: var(--space-xl);
  border-radius: 12px;
  transition: all var(--transition-normal);
  border: 2px solid transparent;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--accent-primary);
}
```

**3. Modals**
```css
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999;
  display: none;
}

.modal.active { display: block; }

.modal-backdrop {
  background: rgba(0,0,0,0.7);
}

.modal-content {
  background: var(--bg-secondary);
  max-width: 600px;
  margin: 100px auto;
  border-radius: 12px;
  padding: var(--space-xl);
  box-shadow: var(--shadow-lg);
}
```

**4. Toast Notifications**
```css
#toast-container {
  position: fixed;
  top: 80px;
  right: 20px;
  z-index: 10000;
}

.toast {
  background: var(--bg-secondary);
  padding: var(--space-md) var(--space-lg);
  border-radius: 8px;
  box-shadow: var(--shadow-lg);
  border-left: 4px solid var(--accent-primary);
  min-width: 300px;
  animation: slideIn 0.3s ease;
}

.toast.success { border-left-color: var(--accent-success); }
.toast.error { border-left-color: var(--accent-danger); }
.toast.warning { border-left-color: var(--accent-warning); }
```

**5. Forms**
```css
.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.form-group input, .form-group select {
  padding: var(--space-md);
  background: var(--bg-tertiary);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-primary);
}
```

**6. Grid Layouts**
```css
/* Features/Stats Grid */
.features-grid, .stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--space-lg);
}

/* Results Grid */
.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-lg);
}
```

---

**JavaScript Utilities** - `src/server/static/js/main.js`

**Features:**
- Theme toggle (dark/light with localStorage persistence)
- Toast notifications
- Modal dialogs
- API helper functions
- Mobile menu toggle
- Keyboard shortcuts (Escape to close modal)
- Utility formatters (distance, coordinates, etc.)

**Example Functions:**
```javascript
// Toast Notifications
function showToast(message, type = 'info', duration = 3000) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    container.appendChild(toast);
    setTimeout(() => toast.remove(), duration);
}

// Modal Dialogs
function showModal(title, content) {
    const modal = document.createElement('div');
    modal.className = 'modal active';
    modal.innerHTML = `
        <div class="modal-backdrop" onclick="closeModal()"></div>
        <div class="modal-content">
            <div class="modal-header">
                <h2>${title}</h2>
                <button class="modal-close" onclick="closeModal()">×</button>
            </div>
            <div class="modal-body">${content}</div>
        </div>
    `;
    document.getElementById('modal-container').appendChild(modal);
}

// API Helpers
async function apiGet(endpoint) {
    const response = await fetch(endpoint);
    const data = await response.json();
    if (!data.success) {
        throw new Error(data.error?.message || 'API request failed');
    }
    return data.data;
}

// Theme Toggle
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}
```

---

**HTML Templates** - `src/server/templates/`

**Base Template** (`base.html`):
```html
<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - {{.ProjectName}}</title>
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="icon" href="/static/images/favicon.png">
</head>
<body>
    <header id="main-header">
        <div class="header-container">
            <div class="header-left">
                <button class="mobile-menu-toggle" onclick="toggleMobileMenu()">☰</button>
                <a class="logo" href="/">{{.ProjectName}}</a>
            </div>
            <nav id="main-nav" class="header-center">
                <a href="/">Home</a>
                <a href="/search">Search</a>
                <a href="/docs">API</a>
                <a href="/stats">Stats</a>
            </nav>
            <div class="header-right">
                <button class="theme-toggle" onclick="toggleTheme()">
                    <span class="theme-icon">🌙</span>
                </button>
            </div>
        </div>
    </header>
    
    <main id="main-content">
        {{template "content" .}}
    </main>
    
    <footer id="main-footer">
        <div class="footer-container">
            <p>&copy; 2024 {{.ProjectName}} - MIT License</p>
        </div>
    </footer>
    
    <div id="toast-container"></div>
    <div id="modal-container"></div>
    <script src="/static/js/main.js"></script>
</body>
</html>
```

**Template Embedding** (`templates.go`):
```go
package server

import (
    "embed"
    "html/template"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var templates *template.Template

func initTemplates() error {
    var err error
    templates, err = template.ParseFS(templateFS, "templates/*.html")
    return err
}
```

---

**Admin Interface Standards:**

**Admin Dashboard** (`/admin`):
- Clean, modern layout
- Server statistics cards
- Settings management
- Database status
- Log viewer
- System health monitoring

**Features:**
- ✅ Dark theme by default (light theme available)
- ✅ Responsive design (mobile-friendly)
- ✅ Toast notifications for feedback
- ✅ Modal dialogs for confirmations
- ✅ Real-time updates
- ✅ Form validation
- ✅ Keyboard shortcuts
- ✅ Accessible (ARIA labels)

**Color Scheme:**
- Dark Mode: `#1a1a1a` background, `#0066cc` primary
- Light Mode: `#ffffff` background, `#0066cc` primary
- Consistent with professional, modern UI/UX

---

**Responsive Design:**

```css
/* Mobile (< 720px) */
@media (max-width: 720px) {
  .mobile-menu-toggle { display: block; }
  .header-center {
    position: absolute;
    top: 70px;
    flex-direction: column;
    display: none;
  }
  .header-center.active { display: flex; }
  .hero h1 { font-size: 2rem; }
  .form-row { grid-template-columns: 1fr; }
}

/* Desktop (> 720px) */
@media (min-width: 720px) {
  #main-content { width: 90%; }
}
```

---

**Key Principles:**
- ✅ **No external dependencies** - No Bootstrap, Tailwind, React, etc.
- ✅ **Vanilla JS** - Pure JavaScript, no jQuery
- ✅ **Embedded assets** - All CSS/JS/images built into binary
- ✅ **Progressive enhancement** - Works without JavaScript
- ✅ **Dark theme first** - Default dark, supports light
- ✅ **Mobile responsive** - Works on all screen sizes
- ✅ **Fast loading** - Minimal CSS/JS, no external requests
- ✅ **Accessible** - Semantic HTML, ARIA labels

---

**File Sizes:**
- CSS: ~867 lines / ~25KB
- JavaScript: ~130 lines / ~4KB
- Total frontend: ~30KB (gzipped: ~10KB)

**Browser Support:**
- Modern browsers with CSS Grid support
- Chrome/Edge 88+
- Firefox 78+
- Safari 14+


---

## 14. Testing & Build Environment

**Critical Rule**: ALWAYS use `/tmp/` for all temporary files, test data, and debugging. NEVER write test data to production directories.


**Rule**: Building: ALWAYS use Docker. Testing/Debugging: Prefer Incus, fallback to Docker, last resort Host. Never test directly on the host OS.

**Environment Selection:**


**For Building** (make build, cross-compilation):

- ✅ **Docker** - ALWAYS use Docker (golang:alpine builder)

- ❌ Never use Incus or Host OS for builds

- Makefile already uses Docker for all builds


**For Testing/Debugging** (running services, integration tests):

- 1️⃣ **Incus** (preferred) - System containers, full OS environment

- 2️⃣ **Docker** (fallback) - If Incus unavailable

- 3️⃣ **Host OS** (last resort) - Only when containers unavailable, tools may need installation


**Priority Order:**
1. **Incus** (if available) - Preferred
2. **Docker** (if available) - Alternative containerization
3. **Host OS** - Last resort only

---

**Testing Workflow:**

**1. Incus (Preferred)**

```bash
# Build development image
make docker-dev

# Run with docker-compose (test configuration)
docker-compose -f docker-compose.test.yml up -d

# Access service
curl http://localhost:64181/healthz

# View logs
docker-compose -f docker-compose.test.yml logs -f

# Cleanup
docker-compose -f docker-compose.test.yml down
sudo rm -rf /tmp/{projectname}/rootfs
```

**Benefits:**
- ✅ Isolated environment
- ✅ Consistent across all developers
- ✅ Easy cleanup (ephemeral /tmp storage)
- ✅ No pollution of host system
- ✅ Matches production environment

---

**2. Incus (Alternative)**

If Docker is not available, use Incus containers:

```bash
# Launch container
incus launch images:alpine/3.19 {projectname}-test

# Copy binary
incus file push ./binaries/{projectname} {projectname}-test/usr/local/bin/

# Execute
incus exec {projectname}-test -- /usr/local/bin/{projectname} --port 8080

# Shell access
incus shell {projectname}-test

# Cleanup
incus delete -f {projectname}-test
```

**Benefits:**
- ✅ System container (not Docker)
- ✅ LXD/LXC based
- ✅ Full OS environment
- ✅ Easy snapshot/restore

---

**3. Host OS (Last Resort)**

Only use when containers are unavailable:

```bash
# Build for host
make build

# Run directly
./binaries/{projectname}  --port 8080

# With custom directories (avoid polluting system)
./binaries/{projectname} \
  --config /tmp/{projectname}/config \
  --data /tmp/{projectname}/data \
  --logs /tmp/{projectname}/logs \
  --port 8080

# Cleanup
rm -rf /tmp/{projectname}
```

**⚠️ Warnings:**
- May pollute host system directories
- Different from production environment
- Harder to cleanup
- Not reproducible across systems

---

**Make Targets for Testing:**

```makefile
# Run tests in Docker (always)
test:
	@docker run --rm -v $$(pwd):/workspace -w /workspace golang:alpine \
		sh -c 'go test -v -race -timeout 5m ./...'

# Build development Docker image
docker-dev:
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		-t $(PROJECTNAME):dev \
		.

# Test with docker-compose
test-docker:
	@docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for service..."
	@timeout 30 bash -c 'until curl -sf http://localhost:64181/healthz; do sleep 1; done'
	@echo "✓ Service is running"
	@docker-compose -f docker-compose.test.yml logs
	@docker-compose -f docker-compose.test.yml down
	@sudo rm -rf /tmp/{projectname}/rootfs

# Test with Incus (if Docker not available)
test-incus:
	@incus launch images:alpine/3.19 $(PROJECTNAME)-test
	@incus file push ./binaries/$(PROJECTNAME) $(PROJECTNAME)-test/usr/local/bin/
	@incus exec $(PROJECTNAME)-test -- /usr/local/bin/$(PROJECTNAME) --status
	@incus delete -f $(PROJECTNAME)-test
```

---

**CI/CD Testing:**

**GitHub Actions:**
- Always runs in Docker containers (ubuntu-latest)
- Uses `make test` which runs tests in Docker

**Jenkins:**
- Runs on dedicated agents (amd64, arm64)
- Each agent runs in isolated environment
- Uses Docker for builds

**Local Development:**
- Developers MUST use Docker or Incus
- Host OS testing only when troubleshooting container issues
- Document when host testing is used (unusual case)

---

**Best Practices:**

✅ **DO:**
- Always test in Docker first
- Use docker-compose.test.yml for integration tests
- Use ephemeral storage (/tmp) for test data
- Clean up after tests
- Document test environment used

❌ **DON'T:**
- Test directly on host OS without containers
- Leave test data in system directories
- Mix development and production configs
- Skip cleanup after tests
- Assume host environment matches production


### Testing Environment Priority
1. **Incus** (preferred) - `incus launch images:alpine/3.19 {project}-test`
2. **Docker** (alternative) - `docker-compose -f docker-compose.test.yml up -d`
3. **Host OS** (last resort) - `./binaries/{projectname} `

Always use isolated environments. Never test directly on host OS unless containers unavailable.


---

## 15. IPv6 Support

**Purpose**: Full dual-stack IPv4/IPv6 support for modern network environments

**Requirements:**
- ✅ Listen on both IPv4 and IPv6
- ✅ Accept connections from both protocols
- ✅ Handle IPv6 addresses in URLs
- ✅ GeoIP lookups for IPv6
- ✅ Proper IPv6 address formatting

---

#### **Listening Addresses**

**Dual-Stack (Recommended):**
```bash
--address ::              # Listen on all IPv6 (includes IPv4 via dual-stack)
--address 0.0.0.0         # Listen on all IPv4 only
```

**Specific Addresses:**
```bash
--address ::1             # IPv6 localhost
--address 127.0.0.1       # IPv4 localhost
--address 2001:db8::1     # Specific IPv6 address
--address 192.168.1.100   # Specific IPv4 address
```

**Environment Variables:**
```bash
ADDRESS=::                # Dual-stack (IPv4 + IPv6)
ADDRESS=0.0.0.0           # IPv4 only
ADDRESS=::1               # IPv6 localhost only
```

---

#### **URL Display with IPv6**

**IPv6 addresses must be enclosed in brackets:**

```go
func getAccessibleURL(port string) string {
    // ... detection logic ...
    
    // For IPv6, wrap in brackets
    if strings.Contains(ip, ":") {
        return fmt.Sprintf("http://[%s]:%s", ip, port)
    }
    return fmt.Sprintf("http://%s:%s", ip, port)
}
```

**Examples:**
- IPv4: `http://192.168.1.100:64555`
- IPv6: `http://[2001:db8::1]:64555`
- IPv6 localhost: `http://[::1]:64555`
- Hostname: `http://server.example.com:64555`

---

#### **Docker IPv6 Support**

**Enable IPv6 in Docker daemon** (`/etc/docker/daemon.json`):
```json
{
  "ipv6": true,
  "fixed-cidr-v6": "2001:db8:1::/64"
}
```

**docker-compose.yml:**
```yaml
services:
  {projectname}:
    image: ghcr.io/{organization}/{projectname}:latest
    environment:
      - ADDRESS=::              # Dual-stack
    networks:
      - {projectname}

networks:
  {projectname}:
    name: {projectname}
    enable_ipv6: true
    ipam:
      config:
        - subnet: 172.18.0.0/16
        - subnet: 2001:db8:1::/64
```

**Port Mappings:**
```yaml
# IPv4 only
ports:
  - "172.17.0.1:64180:80"

# Dual-stack (binds to both)
ports:
  - "64180:80"

# IPv6 explicit
ports:
  - "[::1]:64180:80"
```

---

#### **GeoIP with IPv6**

**MaxMind GeoLite2 databases support IPv6:**

```go
// Lookup IPv6 address
record, err := geoipReader.City(net.ParseIP("2001:4860:4860::8888"))

// Works with both IPv4 and IPv6
func lookupIP(ipStr string) (*geoip2.City, error) {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return nil, fmt.Errorf("invalid IP address")
    }
    
    // Automatically handles IPv4 or IPv6
    return geoipReader.City(ip)
}
```

**API Endpoints:**
```bash
# IPv4
curl http://localhost:64180/api/v1/geoip/8.8.8.8

# IPv6
curl http://localhost:64180/api/v1/geoip/2001:4860:4860::8888

# IPv6 (URL encoded)
curl "http://localhost:64180/api/v1/geoip/2001%3A4860%3A4860%3A%3A8888"
```

---

#### **Code Implementation**

**Listen on Dual-Stack:**

```go
func StartServer(address, port string) error {
    // Default to dual-stack if not specified
    if address == "" {
        address = "::"  // IPv6 dual-stack
    }
    
    addr := fmt.Sprintf("%s:%s", address, port)
    
    // Create listener
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("failed to listen on %s: %w", addr, err)
    }
    
    log.Printf("Server listening on %s", addr)
    return http.Serve(listener, handler)
}
```

**IPv6 Address Detection:**

```go
func getOutboundIP() string {
    // Try IPv4 first
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err == nil {
        defer conn.Close()
        localAddr := conn.LocalAddr().(*net.UDPAddr)
        return localAddr.IP.String()
    }
    
    // Try IPv6
    conn, err = net.Dial("udp", "[2001:4860:4860::8888]:80")
    if err == nil {
        defer conn.Close()
        localAddr := conn.LocalAddr().(*net.UDPAddr)
        return localAddr.IP.String()
    }
    
    return ""
}

func formatURLWithIP(ip, port string) string {
    // IPv6 addresses need brackets
    if strings.Contains(ip, ":") {
        return fmt.Sprintf("http://[%s]:%s", ip, port)
    }
    return fmt.Sprintf("http://%s:%s", ip, port)
}
```

**IP Parsing and Validation:**

```go
func parseIP(ipStr string) (net.IP, error) {
    // Remove brackets if present (from URL)
    ipStr = strings.Trim(ipStr, "[]")
    
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return nil, fmt.Errorf("invalid IP address: %s", ipStr)
    }
    
    return ip, nil
}

func isIPv6(ip net.IP) bool {
    return ip.To4() == nil
}
```

---

#### **Testing IPv6**

**Local Testing:**
```bash
# Start server on IPv6 localhost
./binaries/{projectname} --address ::1 --port 64555

# Test with curl
curl http://[::1]:64555/healthz
curl -g http://[::1]:64555/api/v1/airports/JFK

# Test with IPv6 GeoIP lookup
curl http://[::1]:64555/api/v1/geoip/2001:4860:4860::8888
```

**Docker Testing:**
```bash
# Enable IPv6 in test compose
docker-compose -f docker-compose.test.yml up -d

# Test both protocols
curl http://localhost:64181/healthz              # IPv4
curl http://[::1]:64181/healthz                  # IPv6
```

---

#### **Health Check with IPv6**

```go
func healthCheck(address, port string) error {
    // Detect if IPv6
    url := fmt.Sprintf("http://%s:%s/healthz", address, port)
    if strings.Contains(address, ":") && address != "::1" {
        url = fmt.Sprintf("http://[%s]:%s/healthz", address, port)
    }
    
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return fmt.Errorf("unhealthy: status %d", resp.StatusCode)
    }
    
    return nil
}
```

---

#### **Best Practices**

✅ **DO:**
- Default to `::` for dual-stack support
- Wrap IPv6 addresses in brackets for URLs
- Support both IPv4 and IPv6 in all endpoints
- Test with both protocols
- Use `net.ParseIP()` for validation
- Handle IPv6 in logs and error messages

✅ **Documentation:**
- Document IPv6 support in README
- Show IPv6 examples in API docs
- Include IPv6 in curl examples
- Document dual-stack configuration

❌ **DON'T:**
- Force IPv4-only mode by default
- Forget brackets around IPv6 in URLs
- Assume all IPs are IPv4
- Skip IPv6 testing

---

#### **Configuration Examples**

**systemd service:**
```ini
[Service]
Environment="ADDRESS=::"      # Dual-stack
```

**Docker Compose:**
```yaml
environment:
  - ADDRESS=::                # Dual-stack
  - PORT=80
```

**Manual:**
```bash
# IPv4 only
./binaries/{projectname} --address 0.0.0.0 --port 8080

# IPv6 only
./binaries/{projectname} --address :: --port 8080

# IPv6 localhost
./binaries/{projectname} --address ::1 --port 8080

# Dual-stack (recommended)
./binaries/{projectname} --address :: --port 8080
```

---

**IPv6 Support Summary:**
- ✅ Dual-stack by default (IPv4 + IPv6)
- ✅ Proper URL formatting with brackets
- ✅ GeoIP lookups for both protocols
- ✅ Docker networking with IPv6
- ✅ Testing on both protocols
- ✅ Documentation with examples


---

## 16. GeoIP Databases

**Note**: GeoIP is ONLY required for projects that need IP geolocation features. Not all API projects need GeoIP.


**Include GeoIP if your project:**

- ✅ Needs to locate users by IP address

- ✅ Provides location-based search (e.g., "find nearby X")

- ✅ Has geographic filtering or recommendations

- ✅ Requires country/region detection


**Skip GeoIP if your project:**

- ❌ Only serves static data without location context

- ❌ Doesn't need IP-based features

- ❌ Has no geographic queries


**Projects with GeoIP**: airports, zipcodes

**Projects without GeoIP**: citylist


**Purpose**: IP geolocation using sapics/ip-location-db aggregated data sources

**Repository**: https://github.com/sapics/ip-location-db
**CDN**: jsdelivr (https://cdn.jsdelivr.net/npm/@ip-location-db/)
**Update Frequency**: Daily (automatically via CDN)

---

#### **Database Configuration**

**4 Required Databases:**

1. **geolite2-city-ipv4.mmdb** (~50MB)
   - City-level geolocation for IPv4 addresses
   - Coordinates, timezone, postal codes
   - MaxMind GeoLite2 data
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv4.mmdb`

2. **geolite2-city-ipv6.mmdb** (~40MB)
   - City-level geolocation for IPv6 addresses
   - Coordinates, timezone, postal codes
   - MaxMind GeoLite2 data
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv6.mmdb`

3. **geo-whois-asn-country.mmdb** (~8MB)
   - Country-level data (combined IPv4/IPv6)
   - Aggregated from WHOIS and ASN sources
   - **Public domain** (no attribution required)
   - Daily updates
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb`

4. **asn.mmdb** (~5MB)
   - ASN/ISP information (combined IPv4/IPv6)
   - Autonomous System Numbers
   - Daily updates
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb`

**Total Size**: ~103MB (4 databases)

---

#### **Implementation**

**Service Structure:**

```go
package geoip

import (
    "net"
    "github.com/oschwald/geoip2-golang"
)

const (
    cityIPv4URL  = "https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv4.mmdb"
    cityIPv6URL  = "https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv6.mmdb"
    countryURL   = "https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb"
    asnURL       = "https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb"
)

type Service struct {
    cityIPv4DB *geoip2.Reader  // City database for IPv4
    cityIPv6DB *geoip2.Reader  // City database for IPv6
    countryDB  *geoip2.Reader  // Country (combined IPv4/IPv6)
    asnDB      *geoip2.Reader  // ASN (combined IPv4/IPv6)
    dataDir    string
}
```

**IPv4/IPv6 Detection:**

```go
func (s *Service) Lookup(ip net.IP) (*GeoLocation, error) {
    // Determine which city database to use
    var cityDB *geoip2.Reader
    if ip.To4() != nil {
        // IPv4 address - use IPv4 city database
        cityDB = s.cityIPv4DB
    } else {
        // IPv6 address - use IPv6 city database
        cityDB = s.cityIPv6DB
    }
    
    // Try city lookup first
    city, err := cityDB.City(ip)
    if err == nil {
        return buildLocation(city, ip), nil
    }
    
    // Fallback to country database (works for both IPv4/IPv6)
    country, err := s.countryDB.Country(ip)
    if err != nil {
        return nil, err
    }
    
    return buildCountryLocation(country, ip), nil
}
```

**Download Function:**

```go
func (s *Service) DownloadDatabases() error {
    databases := map[string]string{
        "geolite2-city-ipv4.mmdb":    cityIPv4URL,
        "geolite2-city-ipv6.mmdb":    cityIPv6URL,
        "geo-whois-asn-country.mmdb": countryURL,
        "asn.mmdb":                   asnURL,
    }
    
    for filename, url := range databases {
        path := filepath.Join(s.dataDir, filename)
        if err := downloadFile(path, url); err != nil {
            return fmt.Errorf("failed to download %s: %w", filename, err)
        }
    }
    
    return s.LoadDatabases()
}
```

---

#### **Update Schedule**

**Automatic Updates (via Scheduler):**

```go
// In main.go
sched := scheduler.New()
sched.AddTask("geoip-update", "0 3 * * 0", func() error {
    return geoipSvc.UpdateDatabases()
})
sched.Start()
```

**Update Frequency:**
- **geolite2-city**: Twice weekly (MaxMind release schedule)
- **geo-whois-asn-country**: Daily (aggregated sources)
- **asn**: Daily

**Weekly scheduler** catches all updates (runs Sunday 3:00 AM).

---

#### **Benefits of sapics/ip-location-db**

✅ **Daily Updates**
- Country and ASN data updated daily
- City data updated twice weekly
- Vs weekly from P3TERX

✅ **Multiple Data Sources**
- Aggregates MaxMind GeoLite2, WHOIS, ASN, GeoFeed
- Higher accuracy through data fusion
- More comprehensive IP coverage

✅ **Public Domain Options**
- geo-whois-asn-country: Public domain (CC0/PDDL)
- No attribution required for country data
- Vs CC BY-SA 4.0 from MaxMind

✅ **CDN Delivery**
- jsdelivr CDN (fast, global distribution)
- Automatic caching and compression
- High availability (99.9% uptime)

✅ **IPv4/IPv6 Optimization**
- Separate city databases reduce memory usage
- Faster lookups (smaller files)
- Better performance for IPv6-heavy traffic

✅ **File Sizes**
- Separate databases are smaller individually
- Can download only what's needed
- Total: ~103MB for all 4

---

#### **Licensing**

**sapics/ip-location-db:**

| Database | License | Attribution Required |
|----------|---------|---------------------|
| geolite2-city-ipv4.mmdb | CC BY-SA 4.0 | ✅ MaxMind |
| geolite2-city-ipv6.mmdb | CC BY-SA 4.0 | ✅ MaxMind |
| geo-whois-asn-country.mmdb | CC0/PDDL | ❌ Public domain |
| asn.mmdb | Various | ⚠️ Check sources |

**Attribution:**

```
GeoIP data from:
- MaxMind GeoLite2 (https://www.maxmind.com/)
- Aggregated via sapics/ip-location-db (https://github.com/sapics/ip-location-db)
- Country data: Public domain (WHOIS aggregation)
```

---

#### **Storage Location**

**Databases stored in:**
- `{CONFIG_DIR}/geoip/geolite2-city-ipv4.mmdb`
- `{CONFIG_DIR}/geoip/geolite2-city-ipv6.mmdb`
- `{CONFIG_DIR}/geoip/geo-whois-asn-country.mmdb`
- `{CONFIG_DIR}/geoip/asn.mmdb`

**Example paths:**
- Linux (root): `/etc/{projectname}/geoip/*.mmdb`
- Linux (user): `~/.config/{projectname}/geoip/*.mmdb`
- macOS (user): `~/Library/Application Support/{ProjectName}/geoip/*.mmdb`
- Windows: `%APPDATA%\{ProjectName}\geoip\*.mmdb`
- Docker: `/config/geoip/*.mmdb`

---

#### **First Run Behavior**

1. Check if databases exist in `{CONFIG_DIR}/geoip/`
2. If not found, download all 4 databases from jsdelivr CDN
3. Display progress for each database
4. Load databases into memory
5. Ready to serve GeoIP lookups

**Console Output:**
```
Loading GeoIP databases...
GeoIP databases not found, downloading...
  Downloading geolite2-city-ipv4.mmdb...
  Downloading geolite2-city-ipv6.mmdb...
  Downloading geo-whois-asn-country.mmdb...
  Downloading asn.mmdb...
GeoIP databases loaded successfully
```

---

#### **API Endpoints**

**GeoIP Lookup Routes:**

```yaml
# Lookup request IP
GET /api/v1/geoip
Response: { "ip": "1.2.3.4", "country": "US", "city": "New York", ... }

# Lookup specific IPv4
GET /api/v1/geoip/8.8.8.8

# Lookup specific IPv6
GET /api/v1/geoip/2001:4860:4860::8888

# Find nearby resources (e.g., airports) near IP location
GET /api/v1/geoip/{resource}/nearby?radius=100&limit=10
```

**Response Format:**

```json
{
  "success": true,
  "data": {
    "ip": "8.8.8.8",
    "country": "US",
    "country_name": "United States",
    "region": "CA",
    "region_name": "California",
    "city": "Mountain View",
    "postal_code": "94035",
    "latitude": 37.386,
    "longitude": -122.0838,
    "timezone": "America/Los_Angeles",
    "asn": 15169,
    "organization": "Google LLC"
  }
}
```

---

#### **Testing GeoIP**

**IPv4 Testing:**
```bash
curl http://localhost:64555/api/v1/geoip/8.8.8.8
curl http://localhost:64555/api/v1/geoip/1.1.1.1
```

**IPv6 Testing:**
```bash
curl http://localhost:64555/api/v1/geoip/2001:4860:4860::8888  # Google IPv6
curl http://localhost:64555/api/v1/geoip/2606:4700:4700::1111  # Cloudflare IPv6
```

**Auto-detection (request IP):**
```bash
curl http://localhost:64555/api/v1/geoip
```

---

#### **Manual Database Updates**

**Delete and re-download:**
```bash
# Remove old databases
rm -rf {CONFIG_DIR}/geoip/*.mmdb

# Restart service (auto-downloads)
systemctl restart {projectname}
```

**Programmatic update:**
```bash
# Via API (admin auth required)
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:64555/api/v1/admin/geoip/update
```

---

**sapics/ip-location-db provides superior GeoIP data with daily updates, multiple sources, and public domain options!**


### Version Control Restrictions for AI

**AI assistants MUST NOT perform any version control operations:**

❌ **NEVER Do:**
- `git add` - Adding files to staging
- `git commit` - Creating commits
- `git push` - Pushing to remote
- `git tag` - Creating tags
- `git merge` - Merging branches
- `git rebase` - Rebasing commits
- `git reset` - Resetting history
- `git stash` - Stashing changes
- Any operation that modifies `.git/` directory

✅ **CAN Do (Read-Only):**
- `git status` - Check repository status
- `git diff` - View changes
- `git log` - View history
- `git branch` - List branches
- `git show` - View commits
- `git ls-files` - List tracked files

**Why This Rule:**
- Version control is human responsibility
- Commits should have meaningful messages written by developers
- Prevents accidental code loss
- Maintains clean git history
- Ensures accountability

**What AI Should Do Instead:**
- Make code changes using Edit/Write tools
- Inform user of changes made
- Suggest what to commit
- Let user create commits manually

**Example:**
```
AI: I've updated the following files:
- src/main.go (added IPv6 support)
- README.md (updated documentation)
- Dockerfile (changed to alpine:latest)

You can review the changes with:
  git diff

When ready, commit with:
  git add .
  git commit -m "Add IPv6 support and update documentation"
  git push
```


---

#### **Multi-Distro Testing** ⭐ REQUIRED

**Test on multiple distributions to ensure binary compatibility:**

**Required Test Images:**

1. **Alpine** (musl libc, no systemd)
   ```bash
   incus launch images:alpine/3.19 test-alpine
   incus file push ./binaries/{projectname} test-alpine/usr/local/bin/
   incus exec test-alpine -- /usr/local/bin/{projectname} --version
   ```

2. **Ubuntu** (glibc, systemd)
   ```bash
   incus launch images:ubuntu/24.04 test-ubuntu
   incus file push ./binaries/{projectname} test-ubuntu/usr/local/bin/
   incus exec test-ubuntu -- systemctl --version  # Verify systemd
   
   # Test as systemd service
   incus file push ./scripts/systemd.service test-ubuntu/etc/systemd/system/{projectname}.service
   incus exec test-ubuntu -- systemctl daemon-reload
   incus exec test-ubuntu -- systemctl start {projectname}
   incus exec test-ubuntu -- systemctl status {projectname}
   ```

3. **Fedora/RHEL** (glibc, systemd, SELinux)
   ```bash
   incus launch images:fedora/40 test-fedora
   incus file push ./binaries/{projectname} test-fedora/usr/local/bin/
   incus exec test-fedora -- /usr/local/bin/{projectname} --version
   ```

4. **Debian** (glibc, systemd)
   ```bash
   incus launch images:debian/12 test-debian
   incus file push ./binaries/{projectname} test-debian/usr/local/bin/
   incus exec test-debian -- /usr/local/bin/{projectname} --version
   ```

---

**Systemd Integration Testing:**

Create test systemd service file:

```ini
# test-systemd.service
[Unit]
Description={ProjectName} Test Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/{projectname} --port 8080
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Use /tmp for testing
Environment="CONFIG_DIR=/tmp/{projectname}/config"
Environment="DATA_DIR=/tmp/{projectname}/data"
Environment="LOGS_DIR=/tmp/{projectname}/logs"

[Install]
WantedBy=multi-user.target
```

**Test systemd:**
```bash
# Copy service file
incus file push test-systemd.service test-ubuntu/etc/systemd/system/{projectname}.service

# Start service
incus exec test-ubuntu -- systemctl daemon-reload
incus exec test-ubuntu -- systemctl start {projectname}
incus exec test-ubuntu -- systemctl status {projectname}

# Check logs
incus exec test-ubuntu -- journalctl -u {projectname} -f

# Test API
curl http://$(incus list test-ubuntu -c4 -f csv | cut -d' ' -f1):8080/healthz

# Stop and cleanup
incus exec test-ubuntu -- systemctl stop {projectname}
incus delete -f test-ubuntu
```

---

**Why Multi-Distro Testing:**

✅ **Alpine (musl):**
- Minimal environment
- musl libc (different from glibc)
- No systemd (init/OpenRC)
- Smallest attack surface

✅ **Ubuntu/Debian (glibc, systemd):**
- Most common server distribution
- systemd integration
- glibc compatibility
- APT package ecosystem

✅ **Fedora/RHEL (glibc, systemd, SELinux):**
- Enterprise environments
- SELinux security policies
- dnf/yum package ecosystem
- Red Hat ecosystem compatibility

✅ **Ensures:**
- Binary works on musl AND glibc
- systemd service files work correctly
- No distribution-specific dependencies
- SELinux compatibility
- Init system compatibility

---

**Automated Multi-Distro Test Script:**

```bash
#!/bin/bash
# test-distros.sh - Test binary on multiple distributions

set -e

BINARY="./binaries/{projectname}"
TESTPORT=$(shuf -i 64000-64999 -n 1)

# Test distributions
DISTROS=(
  "alpine/3.19"
  "ubuntu/24.04"
  "debian/12"
  "fedora/40"
)

for distro in "${DISTROS[@]}"; do
  NAME=$(echo "$distro" | tr '/' '-')
  echo "Testing on $distro..."
  
  # Launch container
  incus launch images:$distro test-$NAME
  
  # Copy binary
  incus file push $BINARY test-$NAME/usr/local/bin/{projectname}
  
  # Make executable
  incus exec test-$NAME -- chmod +x /usr/local/bin/{projectname}
  
  # Test version
  echo "  Version check:"
  incus exec test-$NAME -- /usr/local/bin/{projectname} --version
  
  # Test systemd (if available)
  if incus exec test-$NAME -- which systemctl 2>/dev/null; then
    echo "  Testing systemd integration..."
    # Create temp service file
    incus exec test-$NAME -- bash -c "cat > /etc/systemd/system/test-{projectname}.service << 'EOFSVC'
[Unit]
Description=Test Service
[Service]
Type=simple
ExecStart=/usr/local/bin/{projectname} --port $TESTPORT
Environment=\"CONFIG_DIR=/tmp/test/config\"
Environment=\"DATA_DIR=/tmp/test/data\"
EOFSVC"
    
    incus exec test-$NAME -- systemctl daemon-reload
    incus exec test-$NAME -- systemctl start test-{projectname}
    sleep 2
    incus exec test-$NAME -- systemctl status test-{projectname}
    incus exec test-$NAME -- systemctl stop test-{projectname}
  fi
  
  # Cleanup
  incus delete -f test-$NAME
  echo "  ✓ $distro passed"
done

echo "✅ All distribution tests passed"
```

---

**AI Testing Requirements:**

When testing, AI MUST:
- ✅ Test on Alpine (musl)
- ✅ Test on Ubuntu/Debian (systemd)
- ✅ Verify systemd service works
- ✅ Use Incus for multi-distro testing
- ✅ Always use `/tmp/` for test data
- ✅ Always use random ports

**Example AI Testing Session:**

```
AI: Testing {projectname} binary on multiple distributions...

1. Alpine test:
   incus launch images:alpine/3.19 test-alpine
   ✓ Binary works (musl compatibility verified)

2. Ubuntu test with systemd:
   incus launch images:ubuntu/24.04 test-ubuntu
   ✓ Binary works (glibc compatibility verified)
   ✓ systemd service starts successfully
   ✓ Health check passed on port 64782

All tests passed. Binary is compatible with musl and glibc systems.
```


**Multi-Distro Testing Requirements:**

AI MUST test binaries on multiple distributions:
- ✅ Alpine (musl libc test)
- ✅ Ubuntu/Debian (systemd + glibc test)
- ✅ Fedora/RHEL (optional, SELinux test)

**Systemd Testing Required:**
- Create test systemd service file
- Start service with systemctl
- Verify service status
- Check logs with journalctl
- Stop and cleanup

**Example Multi-Distro Test:**
```bash
# Test on Alpine (musl)
incus launch images:alpine/3.19 test-alpine
incus exec test-alpine -- /usr/local/bin/{projectname} --version
incus delete -f test-alpine

# Test on Ubuntu with systemd
incus launch images:ubuntu/24.04 test-ubuntu
# ... systemd service test ...
incus delete -f test-ubuntu
```


---

## 17. Security & DDoS Protection

**Purpose**: Production-grade security with rate limiting, DDoS protection, and security headers

---

#### **Rate Limiting**

**Implementation using Chi middleware:**

```go
import (
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/httprate"
)

// Rate limiting configuration
const (
    // Global rate limit (all endpoints)
    globalRPS = 100  // 100 requests per second
    globalBurst = 200
    
    // API rate limit (stricter for API)
    apiRPS = 50
    apiBurst = 100
    
    // Admin rate limit (most restrictive)
    adminRPS = 10
    adminBurst = 20
)

func setupRateLimiting(r *chi.Mux) {
    // Global rate limiter
    r.Use(httprate.LimitByIP(globalRPS, time.Second))
    
    // API routes - stricter limits
    r.Route("/api/v1", func(r chi.Router) {
        r.Use(httprate.LimitByIP(apiRPS, time.Second))
        // ... routes ...
    })
    
    // Admin routes - most restrictive
    r.Group(func(r chi.Router) {
        r.Use(httprate.LimitByIP(adminRPS, time.Second))
        r.Get("/admin/*", adminHandler)
    })
}
```

**Rate Limit Headers:**
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1234567890
Retry-After: 60
```

**Rate Limit Response (429 Too Many Requests):**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please try again later.",
    "retry_after": 60
  }
}
```

---

#### **DDoS Protection**

**Built-in Protections:**

**1. Request Timeouts:**
```go
r.Use(middleware.Timeout(60 * time.Second))  // Global timeout
```

**2. Request Size Limits:**
```go
r.Use(middleware.SetHeader("Content-Length", "10485760"))  // 10MB max
r.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, 10<<20)  // 10MB
        next.ServeHTTP(w, r)
    })
})
```

**3. Connection Limits:**
```go
// Limit concurrent connections
server := &http.Server{
    Addr:           addr,
    Handler:        handler,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    IdleTimeout:    120 * time.Second,
    MaxHeaderBytes: 1 << 20,  // 1MB
}
```

**4. Slow Request Protection:**
```go
r.Use(middleware.Recoverer)
r.Use(middleware.Throttle(1000))  // Max 1000 concurrent requests
```

**5. IP-Based Blocking:**
```go
var blockedIPs = map[string]bool{
    "192.0.2.1": true,
}

func blockIPMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := getClientIP(r)
        if blockedIPs[ip] {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

#### **Security Headers**

**Required headers for all responses:**

```go
func securityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Prevent clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        
        // Prevent MIME sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        
        // XSS Protection
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy
        w.Header().Set("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline' https://unpkg.com; "+
            "style-src 'self' 'unsafe-inline' https://unpkg.com; "+
            "img-src 'self' data:; "+
            "font-src 'self' https://unpkg.com; "+
            "connect-src 'self'")
        
        // Referrer Policy
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions Policy
        w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        // HSTS (if using HTTPS)
        if r.TLS != nil {
            w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

#### **Authentication Security**

**Admin Routes Protection:**

**1. Token Requirements:**
- Minimum 64 characters
- Cryptographically random (crypto/rand)
- Hashed with SHA-256 before storage

**2. Password Requirements:**
```go
func validatePassword(password string) error {
    if len(password) < 12 {
        return fmt.Errorf("password must be at least 12 characters")
    }
    // Add complexity requirements if needed
    return nil
}
```

**3. Brute Force Protection:**
```go
// Track failed login attempts
var failedAttempts = sync.Map{}

func checkBruteForce(ip string) bool {
    attempts, _ := failedAttempts.LoadOrStore(ip, 0)
    if attempts.(int) >= 5 {
        return true  // Blocked
    }
    return false
}

func recordFailedAttempt(ip string) {
    attempts, _ := failedAttempts.LoadOrStore(ip, 0)
    failedAttempts.Store(ip, attempts.(int)+1)
}

func resetFailedAttempts(ip string) {
    failedAttempts.Delete(ip)
}
```

**4. Session Security:**
- No sessions stored on server (stateless)
- Bearer tokens for API access
- Basic Auth for web UI (browser-managed)

---

#### **Input Validation**

**All inputs MUST be validated:**

```go
// Validate pagination
func validatePagination(limit, offset int) error {
    if limit < 1 || limit > 1000 {
        return fmt.Errorf("limit must be between 1 and 1000")
    }
    if offset < 0 {
        return fmt.Errorf("offset must be non-negative")
    }
    return nil
}

// Validate coordinates
func validateCoordinates(lat, lon float64) error {
    if lat < -90 || lat > 90 {
        return fmt.Errorf("latitude must be between -90 and 90")
    }
    if lon < -180 || lon > 180 {
        return fmt.Errorf("longitude must be between -180 and 180")
    }
    return nil
}

// Validate radius
func validateRadius(radius int) error {
    if radius < 1 || radius > 500 {
        return fmt.Errorf("radius must be between 1 and 500 km")
    }
    return nil
}

// Sanitize user input
func sanitizeInput(input string) string {
    // Remove null bytes
    input = strings.ReplaceAll(input, "\x00", "")
    // Trim whitespace
    input = strings.TrimSpace(input)
    // Limit length
    if len(input) > 1000 {
        input = input[:1000]
    }
    return input
}
```

---

#### **SQL Injection Prevention**

**ALWAYS use prepared statements:**

```go
// ✅ CORRECT - Prepared statement
func getAirportByCode(db *sql.DB, code string) (*Airport, error) {
    row := db.QueryRow("SELECT * FROM airports WHERE icao = ? OR iata = ?", code, code)
    // ...
}

// ❌ WRONG - String concatenation
func getAirportByCode(db *sql.DB, code string) (*Airport, error) {
    query := "SELECT * FROM airports WHERE icao = '" + code + "'"  // DANGEROUS!
    // ...
}
```

**Validate before querying:**
```go
func validateCode(code string) error {
    // ICAO: 4 uppercase letters
    // IATA: 3 uppercase letters
    if len(code) < 3 || len(code) > 4 {
        return fmt.Errorf("invalid code length")
    }
    
    for _, c := range code {
        if c < 'A' || c > 'Z' {
            return fmt.Errorf("code must contain only uppercase letters")
        }
    }
    
    return nil
}
```

---

#### **CORS Configuration**

**Development Mode (permissive):**
```go
if debugMode {
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    })
}
```

**Production Mode (restrictive):**
```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // Whitelist of allowed origins
        allowedOrigins := []string{
            "https://yourdomain.com",
            "https://app.yourdomain.com",
        }
        
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Credentials", "true")
                break
            }
        }
        
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

#### **Logging & Monitoring**

**Log Security Events:**

```go
type SecurityEvent struct {
    Timestamp time.Time
    IP        string
    Event     string
    Details   string
}

func logSecurityEvent(ip, event, details string) {
    log.Printf("[SECURITY] IP=%s Event=%s Details=%s", ip, event, details)
    // Also write to audit log file
}

// Examples:
logSecurityEvent(ip, "FAILED_LOGIN", "Invalid credentials")
logSecurityEvent(ip, "RATE_LIMIT", "Exceeded 100 req/s")
logSecurityEvent(ip, "BLOCKED_IP", "IP in blocklist")
logSecurityEvent(ip, "INVALID_INPUT", "SQL injection attempt detected")
```

**Audit Log Format:**
```
2025-10-14T12:00:00Z [SECURITY] IP=192.0.2.1 Event=FAILED_LOGIN Details=Invalid credentials
2025-10-14T12:00:05Z [SECURITY] IP=192.0.2.1 Event=RATE_LIMIT Details=Exceeded 100 req/s
2025-10-14T12:00:10Z [SECURITY] IP=192.0.2.1 Event=BLOCKED_IP Details=Too many failed attempts
```

---

#### **Reverse Proxy Configuration**

**Recommended: Use nginx or Caddy for production**

**nginx Configuration:**
```nginx
# Rate limiting
limit_req_zone $binary_remote_addr zone=api:10m rate=100r/s;
limit_req_zone $binary_remote_addr zone=admin:10m rate=10r/s;
limit_conn_zone $binary_remote_addr zone=addr:10m;

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name api.example.com;

    # SSL certificates
    ssl_certificate /etc/ssl/certs/api.example.com.crt;
    ssl_certificate_key /etc/ssl/private/api.example.com.key;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    # Security headers
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000" always;
    
    # Connection limits
    limit_conn addr 10;  # Max 10 connections per IP
    
    # API routes - rate limited
    location /api/ {
        limit_req zone=api burst=200 nodelay;
        proxy_pass http://127.0.0.1:64180;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Admin routes - heavily rate limited
    location /admin/ {
        limit_req zone=admin burst=20 nodelay;
        
        # Restrict to internal IPs
        allow 10.0.0.0/8;
        allow 172.16.0.0/12;
        allow 192.168.0.0/16;
        deny all;
        
        proxy_pass http://127.0.0.1:64180;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    # Public routes
    location / {
        limit_req zone=api burst=200 nodelay;
        proxy_pass http://127.0.0.1:64180;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Caddy Configuration (simpler):**
```caddyfile
api.example.com {
    # Automatic HTTPS
    
    # Rate limiting
    rate_limit {
        zone api {
            key {remote_host}
            events 100
            window 1s
        }
    }
    
    # Admin route restrictions
    @admin {
        path /admin/*
    }
    handle @admin {
        # Restrict to private IPs
        @private_ip {
            remote_ip 10.0.0.0/8 172.16.0.0/12 192.168.0.0/16
        }
        handle @private_ip {
            reverse_proxy localhost:64180
        }
        respond 403
    }
    
    # All other routes
    reverse_proxy localhost:64180
}
```

---

#### **Firewall Rules**

**iptables (Linux):**
```bash
# Allow established connections
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# Allow SSH
iptables -A INPUT -p tcp --dport 22 -j ACCEPT

# Rate limit HTTPS
iptables -A INPUT -p tcp --dport 443 -m limit --limit 100/second --limit-burst 200 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j DROP

# Drop invalid packets
iptables -A INPUT -m state --state INVALID -j DROP

# Block common attack ports
iptables -A INPUT -p tcp --dport 23 -j DROP  # Telnet
iptables -A INPUT -p tcp --dport 3389 -j DROP  # RDP

# Default policy
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT
```

**firewalld (RHEL/Fedora):**
```bash
# Add service
firewall-cmd --permanent --add-service=https
firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="10.0.0.0/8" port port="64180" protocol="tcp" accept'

# Rate limiting
firewall-cmd --permanent --add-rich-rule='rule family="ipv4" port port="443" protocol="tcp" limit value="100/s" accept'

firewall-cmd --reload
```

---

#### **Environment-Specific Security**

**Development:**
- CORS: Allow all origins (`*`)
- Rate limits: Disabled or very high
- Security headers: Minimal
- Admin: Accessible from any IP

**Production:**
- CORS: Whitelist specific domains only
- Rate limits: Enforced (100 global, 50 API, 10 admin)
- Security headers: All enabled
- Admin: Restricted to internal IPs only
- HTTPS: Required (via reverse proxy)

---

#### **Security Checklist**

✅ **Required for Production:**
- [ ] Rate limiting enabled (httprate middleware)
- [ ] Security headers configured
- [ ] Request timeouts set (60s max)
- [ ] Request size limits (10MB max)
- [ ] Connection limits enforced
- [ ] HTTPS via reverse proxy (nginx/Caddy)
- [ ] Admin routes restricted to internal IPs
- [ ] Input validation on all endpoints
- [ ] SQL prepared statements (no string concatenation)
- [ ] Audit logging for security events
- [ ] Admin credentials changed from defaults
- [ ] API tokens rotated regularly
- [ ] File permissions correct (admin_credentials: 0600)

✅ **Optional Enhancements:**
- [ ] IP-based blocking (for known attackers)
- [ ] Geo-blocking (block by country)
- [ ] WAF (Web Application Firewall)
- [ ] Intrusion detection (fail2ban)
- [ ] DDoS mitigation service (Cloudflare)

---

#### **Dependencies**

**Required Go packages:**
```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/httprate"  // Rate limiting
)
```

**go.mod:**
```
require (
    github.com/go-chi/chi/v5 v5.0.11
    github.com/go-chi/httprate v0.8.0
)
```

---

#### **Testing Security**

**Rate Limit Testing:**
```bash
# Test rate limiting
for i in {1..150}; do
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:64180/api/v1/health
done
# Should see: 200 200 200 ... 429 429 429
```

**Security Headers Testing:**
```bash
curl -I http://localhost:64180/
# Verify headers present:
# X-Frame-Options: DENY
# X-Content-Type-Options: nosniff
# Content-Security-Policy: ...
```

**Admin Access Testing:**
```bash
# Should fail without auth
curl http://localhost:64180/admin
# HTTP 401 Unauthorized

# Should work with auth
curl -u admin:password http://localhost:64180/admin
# HTTP 200 OK
```

---

**Security is critical for production APIs - all protections MUST be implemented!**


---

## 18. Admin Configuration

**Important**: Server configuration is done via **WebUI Admin Panel** (`/admin/settings`), NOT via CLI flags.


**CLI accepts ONLY**:

- `--port` - Server port

- `--address` - Listen address

- `--config` - Config directory

- `--data` - Data directory

- `--logs` - Logs directory


- `--version`, `--help`, `--status` - Info flags


All other configuration (CORS, rate limits, security headers, etc.) is managed via the **WebUI at `/admin/settings`**.

**Purpose**: Configure all server settings via admin panel with live reload (no restart required)

**Default CORS**: Allow all origins (`*`) - can be restricted via admin settings

---

#### **Configurable Settings**

**All settings stored in database and applied at runtime:**

**Server Settings:**
- `server.cors_enabled` (boolean) - Enable/disable CORS
- `server.cors_origins` (json) - Allowed origins (array), default: `["*"]`
- `server.cors_methods` (json) - Allowed methods, default: `["GET","POST","PUT","DELETE","OPTIONS"]`
- `server.cors_headers` (json) - Allowed headers, default: `["Content-Type","Authorization"]`
- `server.cors_credentials` (boolean) - Allow credentials, default: `false`

**Rate Limiting:**
- `rate.global_rps` (number) - Global requests per second, default: `100`
- `rate.global_burst` (number) - Global burst allowance, default: `200`
- `rate.api_rps` (number) - API requests per second, default: `50`
- `rate.api_burst` (number) - API burst allowance, default: `100`
- `rate.admin_rps` (number) - Admin requests per second, default: `10`
- `rate.admin_burst` (number) - Admin burst allowance, default: `20`
- `rate.enabled` (boolean) - Enable/disable rate limiting, default: `true`

**Request Limits:**
- `request.timeout` (number) - Request timeout in seconds, default: `60`
- `request.max_size` (number) - Max request body size in bytes, default: `10485760` (10MB)
- `request.max_header_size` (number) - Max header size in bytes, default: `1048576` (1MB)

**Connection Limits:**
- `connection.max_concurrent` (number) - Max concurrent requests, default: `1000`
- `connection.idle_timeout` (number) - Idle timeout in seconds, default: `120`
- `connection.read_timeout` (number) - Read timeout in seconds, default: `10`
- `connection.write_timeout` (number) - Write timeout in seconds, default: `10`

**Security Headers:**
- `security.frame_options` (string) - X-Frame-Options, default: `"DENY"`
- `security.content_type_options` (string) - X-Content-Type-Options, default: `"nosniff"`
- `security.xss_protection` (string) - X-XSS-Protection, default: `"1; mode=block"`
- `security.csp` (string) - Content-Security-Policy, default: see below
- `security.hsts_enabled` (boolean) - Enable HSTS, default: `true`
- `security.hsts_max_age` (number) - HSTS max age in seconds, default: `31536000`

**Logging:**
- `log.level` (string) - Log level: debug, info, warn, error, default: `"info"`
- `log.access_log` (boolean) - Enable access logging, default: `true`
- `log.security_log` (boolean) - Enable security event logging, default: `true`

---

#### **Database Schema**

**Settings table** (already exists):
```sql
CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('string', 'number', 'boolean', 'json')),
  category TEXT NOT NULL,
  description TEXT,
  requires_reload BOOLEAN DEFAULT 0,  -- Does changing this require server reload?
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Default Settings Initialization:**
```go
func InitializeDefaultSettings(db *sql.DB) error {
    settings := []Setting{
        // CORS (default: allow all)
        {"server.cors_enabled", "true", "boolean", "server", "Enable CORS", false},
        {"server.cors_origins", `["*"]`, "json", "server", "Allowed origins", false},
        {"server.cors_methods", `["GET","POST","PUT","DELETE","OPTIONS"]`, "json", "server", "Allowed methods", false},
        {"server.cors_headers", `["Content-Type","Authorization"]`, "json", "server", "Allowed headers", false},
        
        // Rate Limiting
        {"rate.enabled", "true", "boolean", "rate", "Enable rate limiting", false},
        {"rate.global_rps", "100", "number", "rate", "Global req/s", false},
        {"rate.global_burst", "200", "number", "rate", "Global burst", false},
        {"rate.api_rps", "50", "number", "rate", "API req/s", false},
        {"rate.api_burst", "100", "number", "rate", "API burst", false},
        {"rate.admin_rps", "10", "number", "rate", "Admin req/s", false},
        {"rate.admin_burst", "20", "number", "rate", "Admin burst", false},
        
        // Request Limits
        {"request.timeout", "60", "number", "request", "Request timeout (seconds)", false},
        {"request.max_size", "10485760", "number", "request", "Max body size (bytes)", false},
        
        // Security Headers
        {"security.frame_options", "DENY", "string", "security", "X-Frame-Options", false},
        {"security.csp", "default-src 'self'", "string", "security", "Content-Security-Policy", false},
    }
    
    for _, s := range settings {
        _, err := db.Exec(`
            INSERT OR IGNORE INTO settings (key, value, type, category, description, requires_reload)
            VALUES (?, ?, ?, ?, ?, ?)
        `, s.Key, s.Value, s.Type, s.Category, s.Description, s.RequiresReload)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

---

#### **Live Reload Implementation**

**Settings Cache with Reload:**

```go
type SettingsManager struct {
    db       *sql.DB
    cache    sync.Map  // Thread-safe cache
    reloadCh chan struct{}
}

func NewSettingsManager(db *sql.DB) *SettingsManager {
    sm := &SettingsManager{
        db:       db,
        reloadCh: make(chan struct{}, 1),
    }
    sm.loadAll()
    return sm
}

func (sm *SettingsManager) loadAll() {
    rows, err := sm.db.Query("SELECT key, value, type FROM settings")
    if err != nil {
        return
    }
    defer rows.Close()
    
    for rows.Next() {
        var key, value, typ string
        if err := rows.Scan(&key, &value, &typ); err != nil {
            continue
        }
        
        // Parse and cache based on type
        var parsed interface{}
        switch typ {
        case "boolean":
            parsed, _ = strconv.ParseBool(value)
        case "number":
            parsed, _ = strconv.Atoi(value)
        case "json":
            json.Unmarshal([]byte(value), &parsed)
        default:
            parsed = value
        }
        
        sm.cache.Store(key, parsed)
    }
}

func (sm *SettingsManager) Get(key string) interface{} {
    val, ok := sm.cache.Load(key)
    if !ok {
        return nil
    }
    return val
}

func (sm *SettingsManager) Set(key string, value interface{}) error {
    // Update database
    valueStr := fmt.Sprintf("%v", value)
    if _, err := sm.db.Exec("UPDATE settings SET value = ?, updated_at = CURRENT_TIMESTAMP WHERE key = ?", valueStr, key); err != nil {
        return err
    }
    
    // Update cache
    sm.cache.Store(key, value)
    
    // Trigger reload
    select {
    case sm.reloadCh <- struct{}{}:
    default:
    }
    
    return nil
}

func (sm *SettingsManager) WatchChanges() <-chan struct{} {
    return sm.reloadCh
}
```

**Apply Settings at Runtime:**

```go
func (s *Server) watchSettingsChanges() {
    go func() {
        for range s.settings.WatchChanges() {
            s.applySettings()
        }
    }()
}

func (s *Server) applySettings() {
    // Apply CORS settings
    s.corsOrigins = s.settings.Get("server.cors_origins").([]string)
    s.corsMethods = s.settings.Get("server.cors_methods").([]string)
    s.corsHeaders = s.settings.Get("server.cors_headers").([]string)
    
    // Apply rate limits (recreate limiters)
    globalRPS := s.settings.Get("rate.global_rps").(int)
    globalBurst := s.settings.Get("rate.global_burst").(int)
    s.globalLimiter = httprate.NewRateLimiter(globalRPS, time.Second, httprate.WithLimitHandler(rateLimitHandler))
    
    // Log the change
    log.Printf("Settings reloaded: CORS=%v, RateLimit=%d rps", s.corsOrigins, globalRPS)
}
```

---

#### **Admin Settings UI**

**Web Interface** (`/admin/settings`):

```html
<form id="settings-form">
  <h2>Server Configuration</h2>
  
  <!-- CORS Settings -->
  <fieldset>
    <legend>CORS Configuration</legend>
    
    <label>
      <input type="checkbox" name="server.cors_enabled" checked> Enable CORS
    </label>
    
    <label>
      Allowed Origins (one per line):
      <textarea name="server.cors_origins">*</textarea>
      <small>Use * for all origins, or list specific domains</small>
    </label>
    
    <label>
      Allowed Methods:
      <select multiple name="server.cors_methods">
        <option selected>GET</option>
        <option selected>POST</option>
        <option selected>PUT</option>
        <option selected>DELETE</option>
        <option selected>OPTIONS</option>
      </select>
    </label>
  </fieldset>
  
  <!-- Rate Limiting -->
  <fieldset>
    <legend>Rate Limiting</legend>
    
    <label>
      <input type="checkbox" name="rate.enabled" checked> Enable Rate Limiting
    </label>
    
    <label>
      Global Rate (req/s):
      <input type="number" name="rate.global_rps" value="100" min="1" max="10000">
    </label>
    
    <label>
      Global Burst:
      <input type="number" name="rate.global_burst" value="200" min="1" max="20000">
    </label>
    
    <label>
      API Rate (req/s):
      <input type="number" name="rate.api_rps" value="50" min="1" max="1000">
    </label>
    
    <label>
      Admin Rate (req/s):
      <input type="number" name="rate.admin_rps" value="10" min="1" max="100">
    </label>
  </fieldset>
  
  <!-- Request Limits -->
  <fieldset>
    <legend>Request Limits</legend>
    
    <label>
      Request Timeout (seconds):
      <input type="number" name="request.timeout" value="60" min="1" max="300">
    </label>
    
    <label>
      Max Request Size (MB):
      <input type="number" name="request.max_size_mb" value="10" min="1" max="100">
    </label>
    
    <label>
      Max Concurrent Requests:
      <input type="number" name="connection.max_concurrent" value="1000" min="10" max="10000">
    </label>
  </fieldset>
  
  <!-- Security Headers -->
  <fieldset>
    <legend>Security Headers</legend>
    
    <label>
      X-Frame-Options:
      <select name="security.frame_options">
        <option selected>DENY</option>
        <option>SAMEORIGIN</option>
        <option>ALLOW-FROM</option>
      </select>
    </label>
    
    <label>
      Content-Security-Policy:
      <textarea name="security.csp" rows="3">default-src 'self'; script-src 'self' 'unsafe-inline' https://unpkg.com</textarea>
    </label>
    
    <label>
      <input type="checkbox" name="security.hsts_enabled" checked> Enable HSTS (HTTPS only)
    </label>
  </fieldset>
  
  <button type="submit" class="btn-primary">Save Settings</button>
  <button type="button" class="btn-secondary" onclick="resetDefaults()">Reset to Defaults</button>
</form>

<script>
document.getElementById('settings-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const settings = {};
    
    for (let [key, value] of formData.entries()) {
        settings[key] = value;
    }
    
    try {
        const response = await fetch('/api/v1/admin/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('token')
            },
            body: JSON.stringify({ settings })
        });
        
        if (response.ok) {
            showToast('Settings saved and applied!', 'success');
            // Settings are live - no reload needed
        } else {
            showToast('Failed to save settings', 'error');
        }
    } catch (err) {
        showToast('Error: ' + err.message, 'error');
    }
});
</script>
```

---

#### **API Endpoints**

**Get All Settings:**
```
GET /api/v1/admin/settings
Authorization: Bearer {token}

Response:
{
  "success": true,
  "data": {
    "server.cors_enabled": true,
    "server.cors_origins": ["*"],
    "rate.global_rps": 100,
    "rate.api_rps": 50,
    ...
  }
}
```

**Update Settings:**
```
PUT /api/v1/admin/settings
Authorization: Bearer {token}
Content-Type: application/json

{
  "settings": {
    "server.cors_origins": ["https://app.example.com", "https://example.com"],
    "rate.global_rps": 200,
    "rate.enabled": true
  }
}

Response:
{
  "success": true,
  "data": {
    "applied": ["server.cors_origins", "rate.global_rps", "rate.enabled"],
    "requires_restart": []
  },
  "message": "Settings applied successfully. No restart required."
}
```

**Reset to Defaults:**
```
POST /api/v1/admin/settings/reset
Authorization: Bearer {token}

Response:
{
  "success": true,
  "message": "All settings reset to defaults"
}
```

---

#### **Live Reload Mechanism**

**Settings Watcher:**

```go
type Server struct {
    router         *chi.Mux
    settings       *SettingsManager
    rateLimiters   map[string]*httprate.RateLimiter
    corsConfig     CORSConfig
    securityConfig SecurityConfig
    mu             sync.RWMutex
}

func (s *Server) Start(addr string) error {
    // Start settings watcher
    go s.watchSettings()
    
    // Start HTTP server
    return http.ListenAndServe(addr, s.router)
}

func (s *Server) watchSettings() {
    for range s.settings.WatchChanges() {
        s.mu.Lock()
        s.applySettingsLive()
        s.mu.Unlock()
        log.Println("Settings reloaded successfully")
    }
}

func (s *Server) applySettingsLive() {
    // Apply CORS
    s.corsConfig = CORSConfig{
        Enabled:     s.settings.GetBool("server.cors_enabled"),
        Origins:     s.settings.GetStringSlice("server.cors_origins"),
        Methods:     s.settings.GetStringSlice("server.cors_methods"),
        Headers:     s.settings.GetStringSlice("server.cors_headers"),
        Credentials: s.settings.GetBool("server.cors_credentials"),
    }
    
    // Apply rate limits (recreate limiters)
    if s.settings.GetBool("rate.enabled") {
        s.rateLimiters["global"] = httprate.NewRateLimiter(
            s.settings.GetInt("rate.global_rps"),
            time.Second,
        )
        s.rateLimiters["api"] = httprate.NewRateLimiter(
            s.settings.GetInt("rate.api_rps"),
            time.Second,
        )
        s.rateLimiters["admin"] = httprate.NewRateLimiter(
            s.settings.GetInt("rate.admin_rps"),
            time.Second,
        )
    }
    
    // Apply security headers
    s.securityConfig = SecurityConfig{
        FrameOptions:       s.settings.GetString("security.frame_options"),
        ContentTypeOptions: s.settings.GetString("security.content_type_options"),
        XSSProtection:      s.settings.GetString("security.xss_protection"),
        CSP:                s.settings.GetString("security.csp"),
        HSTSEnabled:        s.settings.GetBool("security.hsts_enabled"),
        HSTSMaxAge:         s.settings.GetInt("security.hsts_max_age"),
    }
}
```

**Dynamic Middleware:**

```go
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        s.mu.RLock()
        config := s.corsConfig
        s.mu.RUnlock()
        
        if !config.Enabled {
            next.ServeHTTP(w, r)
            return
        }
        
        origin := r.Header.Get("Origin")
        
        // Check if origin is allowed
        allowed := false
        for _, o := range config.Origins {
            if o == "*" || o == origin {
                allowed = true
                break
            }
        }
        
        if allowed {
            if config.Origins[0] == "*" {
                w.Header().Set("Access-Control-Allow-Origin", "*")
            } else {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            }
            
            w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.Methods, ", "))
            w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.Headers, ", "))
            
            if config.Credentials {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
        }
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func (s *Server) rateLimitMiddleware(limiterName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            s.mu.RLock()
            limiter, exists := s.rateLimiters[limiterName]
            enabled := s.settings.GetBool("rate.enabled")
            s.mu.RUnlock()
            
            if !enabled || !exists {
                next.ServeHTTP(w, r)
                return
            }
            
            // Apply rate limiting
            limiter.Handler(next).ServeHTTP(w, r)
        })
    }
}
```

---

#### **Settings Categories**

**Categories for organization:**

| Category | Settings | Live Reload |
|----------|----------|-------------|
| `server` | CORS, title, description | ✅ Yes |
| `rate` | All rate limiting configs | ✅ Yes |
| `request` | Timeouts, size limits | ✅ Yes |
| `connection` | Connection limits | ✅ Yes |
| `security` | Headers, policies | ✅ Yes |
| `log` | Logging levels, formats | ✅ Yes |
| `database` | DB connection settings | ❌ Requires restart |

**Most settings can be applied live without restart!**

---

#### **Admin Settings Page Features**

**UI Features:**
- ✅ Real-time preview of changes
- ✅ Validation before save
- ✅ Reset to defaults button
- ✅ Export/import configuration
- ✅ Change history log
- ✅ Visual feedback (toasts)
- ✅ No page reload needed

**JavaScript for Live Updates:**
```javascript
// Watch for settings changes via WebSocket or polling
function watchSettings() {
    setInterval(async () => {
        const response = await fetch('/api/v1/admin/settings/version');
        const data = await response.json();
        
        if (data.version !== currentVersion) {
            showToast('Settings updated by another admin', 'info');
            loadSettings();  // Reload form
        }
    }, 5000);  // Check every 5 seconds
}
```

---

#### **Configuration Export/Import**

**Export Configuration:**
```
GET /api/v1/admin/settings/export
Authorization: Bearer {token}

Response: (JSON file download)
{
  "version": "1.0.0",
  "exported_at": "2025-10-14T12:00:00Z",
  "settings": {
    "server.cors_origins": ["*"],
    "rate.global_rps": 100,
    ...
  }
}
```

**Import Configuration:**
```
POST /api/v1/admin/settings/import
Authorization: Bearer {token}
Content-Type: application/json

{
  "settings": {
    "server.cors_origins": ["https://example.com"],
    "rate.global_rps": 200
  }
}
```

---

#### **Best Practices**

✅ **CORS Defaults:**
- Start with `*` (allow all) for ease of use
- Admins can restrict via settings
- Document how to configure properly
- Warn about security implications of `*` in production

✅ **Rate Limiting:**
- Always enabled by default
- Configurable per route type (global/API/admin)
- Return proper 429 responses
- Include Retry-After headers

✅ **Live Reload:**
- Most settings apply immediately
- No server restart needed
- Thread-safe with sync.RWMutex
- Log all setting changes

✅ **Security:**
- All security headers enabled by default
- Input validation on all settings
- Audit log for configuration changes
- Export/import for backup/migration

❌ **DON'T:**
- Disable security features by default
- Allow unlimited rate limits
- Skip input validation
- Forget to log configuration changes

---

**Admin settings provide full server control without restarts!**


---

#### **AI Communication Guidelines** ⭐ CRITICAL

**Understanding User Intent:**

**Question Mark (?) = Question, NOT Command**

When a user ends a message with `?`, the AI is being **asked a question**, not given a command to execute.

**Examples:**

**❌ WRONG (AI executes):**
```
User: "Should we use alpine:latest?"
AI: *Updates Dockerfile to alpine:latest*
```

**✅ CORRECT (AI answers):**
```
User: "Should we use alpine:latest?"
AI: "Yes, since we build static binaries (CGO_ENABLED=0), using alpine:latest 
     is fine. The Alpine version only affects runtime tools (curl, bash), not 
     the binary itself. Would you like me to update the Dockerfile?"
```

**❌ WRONG (AI executes):**
```
User: "Can you add rate limiting?"
AI: *Immediately adds rate limiting code*
```

**✅ CORRECT (AI clarifies first):**
```
User: "Can you add rate limiting?"
AI: "I can add rate limiting using httprate middleware. Would you like me to:
     1. Add basic rate limiting (100 req/s global)
     2. Add comprehensive rate limiting (per-route limits)
     3. Explain the options first?
     
     Which approach would you prefer?"
```

---

**Command Detection Rules:**

**Questions (Answer, Don't Execute):**
- Ends with `?`
- Starts with: "Should", "Could", "Would", "Can you", "Is it", "Are we", "Do we"
- Contains: "What do you think", "Any ideas", "Suggestions"
- Example: "Should we switch to sapics?" → Answer with pros/cons

**Commands (Execute After Confirmation if Unclear):**
- Imperative: "Update", "Add", "Fix", "Change", "Remove"
- Direct: "Use sapics", "Change to alpine:latest"
- Clear intent: "Update all projects to use sapics"
- Example: "Update to sapics" → Execute the change

**Ambiguous (Ask for Clarification):**
- "Maybe we should..." → Ask if they want you to do it
- "I'm thinking about..." → Ask if they want you to implement it
- "What about..." → Provide options, then ask

---

**Response Patterns:**

**For Questions:**
1. Provide clear, accurate answer
2. Give pros/cons if applicable
3. Offer to implement if helpful
4. Wait for explicit go-ahead before making changes

**For Commands:**
1. Acknowledge the task
2. If complex, ask clarifying questions first
3. Execute the changes
4. Report what was changed
5. Suggest what to test/verify

**For Ambiguous:**
1. Acknowledge you heard them
2. Ask clarifying question
3. Provide options
4. Wait for clear direction

---

**Examples of Correct AI Behavior:**

**Scenario 1: Question About Technology Choice**
```
User: "Is sapics better than P3TERX?"

AI: "Yes, sapics/ip-location-db has several advantages:
     - Daily updates (vs weekly from P3TERX)
     - Multiple data sources aggregated (higher accuracy)
     - Public domain country data (no attribution required)
     - jsdelivr CDN (faster, more reliable)
     - Separate IPv4/IPv6 databases (better performance)
     
     Would you like me to migrate all projects from P3TERX to sapics?"

User: "Yes"

AI: *Now executes the migration*
```

**Scenario 2: Exploratory Question**
```
User: "Should we add GraphQL?"

AI: "GraphQL can be beneficial if you need:
     - Flexible data fetching (clients request specific fields)
     - Reduced over-fetching
     - Strongly typed schema
     
     However, for simple REST APIs, it may be overkill. Your current REST API 
     is well-designed and covers all use cases.
     
     Do you want me to add GraphQL support, or would you prefer to keep the 
     current REST-only approach?"
```

**Scenario 3: Clear Command**
```
User: "Add rate limiting to all projects"

AI: "I'll add rate limiting to airports, citylist, and zipcodes using 
     httprate middleware with these defaults:
     - Global: 100 req/s
     - API: 50 req/s
     - Admin: 10 req/s
     
     I'll also add the admin UI for configuring these limits."

*Proceeds with implementation*
```

**Scenario 4: Thinking Out Loud**
```
User: "I'm thinking we might need DDoS protection..."

AI: "DDoS protection is important for production. I can add:
     1. Rate limiting (httprate middleware)
     2. Connection limits (max concurrent requests)
     3. Request size limits
     4. IP-based blocking
     5. Reverse proxy recommendations (nginx/Caddy)
     
     Would you like me to implement these protections?"
```

---

**Key Communication Principles:**

✅ **DO:**
- Recognize questions vs commands
- Ask for clarification when unclear
- Provide context and options
- Wait for confirmation on ambiguous requests
- Explain what you're about to do
- Report what you changed

❌ **DON'T:**
- Execute code changes based on questions
- Assume intent without confirmation
- Make destructive changes without warning
- Skip asking clarifying questions
- Proceed when user is exploring options
- Ignore question marks

---

**Special Keywords:**

**Questions (Answer, Don't Execute):**
- "Should we..." → Provide recommendation
- "Can you..." → Explain capability, offer to do it
- "Is it..." → Answer yes/no with explanation
- "What if..." → Explain consequences
- "How about..." → Discuss option, ask if they want it

**Commands (Execute):**
- "Do it" → Execute previously discussed change
- "Go ahead" → Proceed with implementation
- "Yes" (after question) → Execute proposed change
- "Update/Add/Fix/Change" → Direct command to make change

**Exploratory (Discuss):**
- "I'm thinking..." → Discuss, don't implement
- "Maybe..." → Provide opinion, ask for direction
- "What about..." → Explain option, wait for decision
- "Just curious..." → Answer informatively

---

**This ensures AI understands when to answer vs when to execute!**


---

## 19. Installation Scripts

**Purpose**: Single-file install scripts for all platforms with automatic service setup

**Location**: `./scripts/` directory in project root

**Required Scripts:**
- `install-linux.sh` - Linux (all distros, all init systems)
- `install-macos.sh` - macOS (with launchd)
- `install-bsd.sh` - FreeBSD/OpenBSD/NetBSD (with rc.d)
- `install-windows.ps1` - Windows (with NSSM for service)

---

#### **Linux Install Script (Distro-Agnostic)**

**File**: `scripts/install-linux.sh`

**Features**:
- ✅ Detects init system (systemd, OpenRC, init.d, runit)
- ✅ Creates service file for detected init
- ✅ Downloads binary for correct architecture
- ✅ Creates directories (OS-specific paths)
- ✅ Generates admin credentials
- ✅ Starts service automatically

**Implementation**:

```bash
#!/bin/bash
# install-linux.sh - Distro-agnostic installer for {projectname}
# Supports: systemd, OpenRC, init.d, runit
# Auto-detects: architecture, init system, package manager

set -e

PROJECTNAME="{projectname}"
GITHUB_REPO="{organization}/${projectname}"
VERSION="latest"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== ${PROJECTNAME} Installer ===${NC}"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "Architecture: $ARCH"

# Detect if running as root
if [ "$EUID" -eq 0 ]; then
    IS_ROOT=true
    BIN_DIR="/usr/local/bin"
    CONFIG_DIR="/etc/${PROJECTNAME}"
    DATA_DIR="/var/lib/${PROJECTNAME}"
    LOG_DIR="/var/log/${PROJECTNAME}"
else
    IS_ROOT=false
    BIN_DIR="$HOME/.local/bin"
    CONFIG_DIR="$HOME/.config/${PROJECTNAME}"
    DATA_DIR="$HOME/.local/share/${PROJECTNAME}"
    LOG_DIR="$HOME/.local/state/${PROJECTNAME}"
fi

echo "Install mode: $([ "$IS_ROOT" = true ] && echo "System (root)" || echo "User")"

# Create directories
mkdir -p "$BIN_DIR" "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"
mkdir -p "$DATA_DIR/db"

# Download binary
echo "Downloading ${PROJECTNAME}-linux-${ARCH}..."
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/${VERSION}/download/${PROJECTNAME}-linux-${ARCH}"

if command -v curl &> /dev/null; then
    curl -L -o "${BIN_DIR}/${PROJECTNAME}" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -O "${BIN_DIR}/${PROJECTNAME}" "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: curl or wget required${NC}"
    exit 1
fi

chmod +x "${BIN_DIR}/${PROJECTNAME}"
echo -e "${GREEN}✓ Binary installed to ${BIN_DIR}/${PROJECTNAME}${NC}"

# Detect init system
detect_init() {
    if [ -d /run/systemd/system ] || command -v systemctl &> /dev/null; then
        echo "systemd"
    elif [ -f /sbin/openrc-run ] || [ -d /etc/init.d ] && grep -q "openrc" /sbin/init 2>/dev/null; then
        echo "openrc"
    elif [ -d /etc/init.d ] && [ ! -d /run/systemd/system ]; then
        echo "sysvinit"
    elif command -v sv &> /dev/null; then
        echo "runit"
    else
        echo "unknown"
    fi
}

INIT_SYSTEM=$(detect_init)
echo "Init system: $INIT_SYSTEM"

# Install service based on init system
case $INIT_SYSTEM in
    systemd)
        if [ "$IS_ROOT" = true ]; then
            cat > /etc/systemd/system/${PROJECTNAME}.service << EOF
[Unit]
Description=${PROJECTNAME} API Server
After=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
ExecStart=${BIN_DIR}/${PROJECTNAME}
Restart=always
RestartSec=5
Environment="CONFIG_DIR=${CONFIG_DIR}"
Environment="DATA_DIR=${DATA_DIR}"
Environment="LOGS_DIR=${LOG_DIR}"
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target

---


**Purpose**: Single-file install scripts for all platforms with automatic service setup

**Location**: `./scripts/` directory in project root

**Required Scripts:**
- `install-linux.sh` - Linux (all distros, all init systems)
- `install-macos.sh` - macOS (with launchd)
- `install-bsd.sh` - FreeBSD/OpenBSD/NetBSD (with rc.d)
- `install-windows.ps1` - Windows (with NSSM for service)
- `uninstall.sh` - Universal uninstaller

---

#### **Install Script Requirements**

**Each script must:**
- ✅ Detect platform/architecture automatically
- ✅ Download correct binary from GitHub releases
- ✅ Create necessary directories (OS-specific paths)
- ✅ Install as system service (auto-start on boot)
- ✅ Display admin credentials location
- ✅ Provide service management commands
- ✅ Support both root and user installations
- ✅ Be idempotent (safe to run multiple times)

---

#### **README.md Structure** ⭐ REQUIRED

**Order**: About → Official → Production → Development → Legal

**Complete Structure:**

```markdown
# {PROJECT_NAME}

Brief description with badges

## About
- Project description
- Key features (bullet list)
- Technology stack
- Screenshots/demos (optional)

## Official Links
- Website: https://{projectname}.com (if exists)
- Documentation: https://{projectname}.readthedocs.io
- GitHub: https://github.com/{org}/{projectname}
- Docker: https://hub.docker.com/r/{org}/{projectname}
- Issues: https://github.com/{org}/{projectname}/issues

## Production Deployment

### Installation

#### Linux (All Distros)
One-line installer:
```bash
curl -fsSL https://raw.githubusercontent.com/{org}/{projectname}/main/scripts/install-linux.sh | sudo bash
```

Supports: Ubuntu, Debian, Fedora, RHEL, Alpine, Arch, Gentoo, Void
Init Systems: systemd, OpenRC, SysVinit, runit

#### macOS
```bash
curl -fsSL https://raw.githubusercontent.com/{org}/{projectname}/main/scripts/install-macos.sh | sudo bash
```

#### FreeBSD/OpenBSD/NetBSD
```bash
fetch https://raw.githubusercontent.com/{org}/{projectname}/main/scripts/install-bsd.sh
chmod +x install-bsd.sh
sudo ./install-bsd.sh
```

#### Windows
```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/{org}/{projectname}/main/scripts/install-windows.ps1 -OutFile install.ps1
Set-ExecutionPolicy Bypass -Scope Process
.\install.ps1
```

#### Docker
```bash
docker-compose up -d
```

### Binary Installation (Manual)
(Platform-specific binary downloads)

### Docker Deployment
(Docker Compose, Docker Run)

### Configuration
(Environment variables, directories, first run)

## API Usage
- Quick examples
- Admin panel
- API documentation (/openapi, /graphql)

## Development
- Requirements
- Build system & testing
- CI/CD (GitHub Actions, Jenkins, ReadTheDocs)
- Contributing guidelines

## Legal
- License (MIT - see LICENSE.md)
- Third-Party Data Licenses
- Credits & Acknowledgments
- Support & Contact

---

**{projectname}** - Production-ready {description}
```

---

**All platforms supported with single-script installation!**

---

## 🎯 Mandatory vs Optional Features

**Everything in this SPEC is REQUIRED unless explicitly marked as an exception.**

### Exceptions:

**1. IPv6 Support** - Auto-detect and fallback
- ✅ Implement IPv6 support in code
- ✅ Auto-detect if OS/network supports IPv6
- ✅ Fallback to IPv4 if IPv6 unavailable
- ✅ Use dual-stack (`::`) by default, fallback to `0.0.0.0`

**Auto-Detection Implementation:**
```go
func getDefaultAddress() string {
    // Try to listen on dual-stack
    listener, err := net.Listen("tcp", "[::]:0")
    if err == nil {
        listener.Close()
        return "::"  // IPv6 dual-stack supported
    }
    
    // Fallback to IPv4
    return "0.0.0.0"
}

// Usage in main.go
address := getEnv("ADDRESS", getDefaultAddress())
```

**2. GeoIP Integration** - Required for location-based projects only
- ✅ Required if: Project provides location-based features
- ✅ Required if: "Find nearby X" functionality
- ✅ Required if: IP-based country/region detection
- ❌ Not required if: No geographic features

**Examples:**
- airports: ✅ GeoIP (find airports near IP location)
- zipcodes: ✅ GeoIP (find ZIP codes near IP)
- citylist: ❌ No GeoIP (cities don't need IP location)

**How to Determine:**
- Does your project answer "What's near me?" → YES: GeoIP required
- Does your project answer "Where am I?" → YES: GeoIP required
- Does your project use user's location? → YES: GeoIP required
- Otherwise → NO: Skip GeoIP

---

**All Other Features Are MANDATORY:**
- ✅ Dockerfile (golang:alpine + alpine:latest)
- ✅ docker-compose.yml & docker-compose.test.yml
- ✅ Makefile (4 targets: build, test, release, docker)
- ✅ Jenkinsfile (jenkins.casjay.cc)
- ✅ GitHub Actions (release.yml + docker.yml)
- ✅ ReadTheDocs (MkDocs + Dracula theme)
- ✅ src/data/ directory (JSON only, embedded in main.go)
- ✅ Web UI with dark theme
- ✅ Swagger/OpenAPI + GraphQL
- ✅ Admin panel with live reload
- ✅ Rate limiting & DDoS protection
- ✅ Security headers
- ✅ Multi-distro testing (Incus)
- ✅ Install scripts (all platforms)
- ✅ Database in {DATA_DIR}/db/
- ✅ CORS default: `*` (configurable via admin)
- ✅ All 8 platform binaries
- ✅ Source archives in releases

**No optional features - this is the complete standard!**


---

## 20. Service Management

**Purpose**: Built-in service management and graceful signal handling

---

#### **Service Management Commands**

**CLI Structure:**
```bash
{projectname} service <command>

Commands:
  start      Start the service (daemon mode)
  stop       Stop the running service
  restart    Restart the service
  reload     Reload configuration without restart (SIGHUP)
  status     Show service status
```

**Implementation:**

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "os/signal"
    "path/filepath"
    "strconv"
    "syscall"
    "time"
)

func handleServiceCommand(args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("usage: %s service <start|stop|restart|reload|status>", os.Args[0])
    }
    
    command := args[0]
    
    switch command {
    case "start":
        return serviceStart()
    case "stop":
        return serviceStop()
    case "restart":
        return serviceRestart()
    case "reload":
        return serviceReload()
    case "status":
        return serviceStatus()
    default:
        return fmt.Errorf("unknown command: %s", command)
    }
}

func serviceStart() error {
    // Check if already running
    if pid, err := readPIDFile(); err == nil {
        if processExists(pid) {
            return fmt.Errorf("service already running (PID: %d)", pid)
        }
    }
    
    // Start in background
    cmd := exec.Command(os.Args[0], "")
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start service: %w", err)
    }
    
    // Write PID file
    if err := writePIDFile(cmd.Process.Pid); err != nil {
        return fmt.Errorf("failed to write PID file: %w", err)
    }
    
    fmt.Printf("Service started (PID: %d)\n", cmd.Process.Pid)
    return nil
}

func serviceStop() error {
    pid, err := readPIDFile()
    if err != nil {
        return fmt.Errorf("service not running")
    }
    
    process, err := os.FindProcess(pid)
    if err != nil {
        return fmt.Errorf("failed to find process: %w", err)
    }
    
    // Send SIGTERM for graceful shutdown
    if err := process.Signal(syscall.SIGTERM); err != nil {
        return fmt.Errorf("failed to stop service: %w", err)
    }
    
    // Wait for process to exit (up to 30 seconds)
    for i := 0; i < 30; i++ {
        if !processExists(pid) {
            os.Remove(getPIDFilePath())
            fmt.Println("Service stopped")
            return nil
        }
        time.Sleep(time.Second)
    }
    
    // Force kill if still running
    process.Kill()
    os.Remove(getPIDFilePath())
    fmt.Println("Service force stopped")
    return nil
}

func serviceRestart() error {
    fmt.Println("Restarting service...")
    serviceStop()
    time.Sleep(2 * time.Second)
    return serviceStart()
}

func serviceReload() error {
    pid, err := readPIDFile()
    if err != nil {
        return fmt.Errorf("service not running")
    }
    
    process, err := os.FindProcess(pid)
    if err != nil {
        return fmt.Errorf("failed to find process: %w", err)
    }
    
    // Send SIGHUP to reload config
    if err := process.Signal(syscall.SIGHUP); err != nil {
        return fmt.Errorf("failed to reload: %w", err)
    }
    
    fmt.Println("Configuration reloaded")
    return nil
}

func serviceStatus() error {
    pid, err := readPIDFile()
    if err != nil {
        fmt.Println("Status: Stopped")
        return nil
    }
    
    if !processExists(pid) {
        fmt.Println("Status: Stopped (stale PID file)")
        os.Remove(getPIDFilePath())
        return nil
    }
    
    fmt.Printf("Status: Running (PID: %d)\n", pid)
    
    // Try health check
    port := os.Getenv("PORT")
    if port != "" {
        healthURL := fmt.Sprintf("http://127.0.0.1:%s/healthz", port)
        resp, err := http.Get(healthURL)
        if err == nil {
            defer resp.Body.Close()
            if resp.StatusCode == 200 {
                fmt.Println("Health: OK")
            } else {
                fmt.Printf("Health: Unhealthy (HTTP %d)\n", resp.StatusCode)
            }
        }
    }
    
    return nil
}

// PID file helpers
func getPIDFilePath() string {
    runtimeDir := os.Getenv("RUNTIME_DIR")
    if runtimeDir == "" {
        runtimeDir = "/tmp"
    }
    return filepath.Join(runtimeDir, "{projectname}.pid")
}

func readPIDFile() (int, error) {
    data, err := os.ReadFile(getPIDFilePath())
    if err != nil {
        return 0, err
    }
    return strconv.Atoi(strings.TrimSpace(string(data)))
}

func writePIDFile(pid int) error {
    return os.WriteFile(getPIDFilePath(), []byte(fmt.Sprintf("%d\n", pid)), 0644)
}

func processExists(pid int) bool {
    process, err := os.FindProcess(pid)
    if err != nil {
        return false
    }
    
    // Send signal 0 to check if process exists
    err = process.Signal(syscall.Signal(0))
    return err == nil
}
```

---

#### **Signal Handling**

**Signals to Handle:**

| Signal | Action | Use Case |
|--------|--------|----------|
| **SIGTERM** | Graceful shutdown | systemctl stop, service stop |
| **SIGINT** | Graceful shutdown | Ctrl+C, user interrupt |
| **SIGHUP** | Reload config | service reload, config changes |
| **SIGUSR1** | Reopen log files | Log rotation |
| **SIGUSR2** | Toggle debug mode | Runtime debugging |
| **SIGKILL** | Force kill | Cannot be caught |

**Implementation:**

```go
func main() {
    // ... initialization ...
    
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan,
        syscall.SIGTERM,  // Graceful shutdown
        syscall.SIGINT,   // Ctrl+C
        syscall.SIGHUP,   // Reload config
        syscall.SIGUSR1,  // Reopen logs
        syscall.SIGUSR2,  // Toggle debug
    )
    
    // Start server in goroutine
    go func() {
        if err := server.Start(addr); err != nil {
            log.Fatal(err)
        }
    }()
    
    // Wait for signals
    for sig := range sigChan {
        switch sig {
        case syscall.SIGTERM, syscall.SIGINT:
            log.Println("Received shutdown signal, shutting down gracefully...")
            handleGracefulShutdown()
            return
            
        case syscall.SIGHUP:
            log.Println("Received SIGHUP, reloading configuration...")
            handleConfigReload()
            
        case syscall.SIGUSR1:
            log.Println("Received SIGUSR1, reopening log files...")
            handleLogRotation()
            
        case syscall.SIGUSR2:
            log.Println("Received SIGUSR2, toggling debug mode...")
            toggleDebugMode()
        }
    }
}

func handleGracefulShutdown() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Shutdown HTTP server
    if err := httpServer.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }
    
    // Close database connections
    database.Close()
    
    // Close GeoIP databases
    if geoipSvc != nil {
        geoipSvc.Close()
    }
    
    // Remove PID file
    os.Remove(getPIDFilePath())
    
    log.Println("Shutdown complete")
}

func handleConfigReload() {
    // Reload settings from database
    if err := settingsManager.Reload(); err != nil {
        log.Printf("Failed to reload settings: %v", err)
        return
    }
    
    // Apply new settings (live reload)
    server.ApplySettings()
    
    log.Println("Configuration reloaded successfully")
}

func handleLogRotation() {
    // Close and reopen log files
    // This is typically used with logrotate
    
    // If writing to files, close and reopen them
    if logFile != nil {
        logFile.Close()
        var err error
        logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            log.Printf("Failed to reopen log file: %v", err)
        }
    }
    
    log.Println("Log files reopened")
}

func toggleDebugMode() {
    currentLevel := log.GetLevel()
    if currentLevel == log.DebugLevel {
        log.SetLevel(log.InfoLevel)
        log.Println("Debug mode: OFF")
    } else {
        log.SetLevel(log.DebugLevel)
        log.Println("Debug mode: ON")
    }
}
```

---

#### **Graceful Shutdown**

**HTTP Server Shutdown:**

```go
type Server struct {
    httpServer *http.Server
    // ...
}

func (s *Server) Start(addr string) error {
    s.httpServer = &http.Server{
        Addr:    addr,
        Handler: s.router,
    }
    
    return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    // Stop accepting new connections
    if err := s.httpServer.Shutdown(ctx); err != nil {
        return err
    }
    
    // Close database
    if s.db != nil {
        s.db.Close()
    }
    
    // Close GeoIP
    if s.geoip != nil {
        s.geoip.Close()
    }
    
    return nil
}
```

---

#### **Usage Examples**

**Service Management:**
```bash
# Start service
{projectname} service start

# Check status
{projectname} service status

# Reload configuration (SIGHUP)
{projectname} service reload

# Restart service
{projectname} service restart

# Stop service
{projectname} service stop
```

**Manual Signal Sending:**
```bash
# Graceful shutdown
kill -TERM $(cat /tmp/{projectname}.pid)

# Reload config
kill -HUP $(cat /tmp/{projectname}.pid)

# Reopen logs (for logrotate)
kill -USR1 $(cat /tmp/{projectname}.pid)

# Toggle debug mode
kill -USR2 $(cat /tmp/{projectname}.pid)
```

**Logrotate Configuration:**

```
# /etc/logrotate.d/{projectname}
/var/log/{projectname}/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 nobody nogroup
    postrotate
        kill -USR1 $(cat /tmp/{projectname}.pid) 2>/dev/null || true
    endscript
}
```

---

#### **PID File Management**

**Location:**
- System: `/var/run/{projectname}.pid` or `/run/{projectname}.pid`
- User: `$HOME/.local/run/{projectname}.pid` or `/tmp/{projectname}.pid`
- Docker: `/tmp/{projectname}.pid`

**Format:**
```
12345
```
(Just the PID, single line)

**Permissions:**
- System: 0644 (readable by all, writable by service)
- User: 0644

---

#### **CLI Flag Updates**

**Add service subcommand:**

```go
func main() {
    // Check for service subcommand first
    if len(os.Args) > 1 && os.Args[1] == "service" {
        if err := handleServiceCommand(os.Args[2:]); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        return
    }
    
    // Regular flags
    flag.Parse()
    
    // ... rest of main ...
}
```

**CLI Accepts:**

**Startup Flags:**
- `--port` - Server port
- `--address` - Listen address
- `--config` - Config directory
- `--data` - Data directory
- `--logs` - Logs directory

**Info Flags:**
- `--version` - Show version (just the number)
- `--help` - Show help
- `--status` - Health check (exit 0=healthy, 1=unhealthy)

**Service Subcommands:**
- `service start` - Start service
- `service stop` - Stop service
- `service restart` - Restart service
- `service reload` - Reload config (SIGHUP)
- `service status` - Show status

---

#### **Help Text**

```
Usage: {projectname} [OPTIONS] [COMMAND]

Options:
  --port PORT       Server port (default: random 64000-64999)
  --address ADDR    Listen address (default: :: for dual-stack)
  --config DIR      Configuration directory
  --data DIR        Data directory
  --logs DIR        Logs directory
  
  
  --version         Show version number
  --help            Show this help
  --status          Check server health

Service Commands:
  service start     Start the service
  service stop      Stop the service
  service restart   Restart the service
  service reload    Reload configuration (no restart)
  service status    Show service status

Examples:
  {projectname}                        Start server (foreground)
  {projectname} --port 8080            Start on port 8080
  {projectname}                Start as daemon
  {projectname} service start          Start service
  {projectname} service status         Check status
  {projectname} service reload         Reload config
  {projectname} --version              Show version

Signals:
  SIGTERM, SIGINT    Graceful shutdown
  SIGHUP             Reload configuration
  SIGUSR1            Reopen log files
  SIGUSR2            Toggle debug mode
```

---

**Service management built into the binary - no external tools needed!**


---

## 21. Debug Mode

**Purpose**: Enable debugging features via DEBUG environment variable (no CLI flags)

**Environment Variable:**
```bash
DEBUG=1              # Enable debug mode
DEBUG=true           # Enable debug mode
DEBUG=yes            # Enable debug mode
# (any non-empty value enables debug)
```

---

#### **Debug Mode Features**

**When DEBUG is set:**

✅ **Verbose Logging:**
- SQL queries logged with parameters
- Request/response bodies logged
- Stack traces on errors
- Timing information for operations

✅ **Debug Endpoints:**
```yaml
GET /debug/routes          - List all registered routes
GET /debug/config          - Show current configuration
GET /debug/db              - Database statistics and connections
GET /debug/settings        - All settings (from database)
GET /debug/memory          - Memory usage and Go runtime stats
POST /debug/reload         - Force reload configuration
POST /debug/gc             - Trigger garbage collection
```

✅ **Template Hot Reload:**
- Templates reloaded on every request
- No need to restart for template changes
- CSS/JS still embedded (rebuild for changes)

✅ **Permissive CORS:**
- Allow all origins (`*`)
- All methods allowed
- All headers allowed

✅ **Detailed Error Messages:**
- Stack traces in error responses
- Internal error details exposed
- Database query errors shown

---

#### **Implementation**

```go
package main

import (
    "os"
    "strings"
)

var debugMode bool

func init() {
    // Check DEBUG environment variable
    debug := strings.ToLower(os.Getenv("DEBUG"))
    debugMode = debug != "" && debug != "0" && debug != "false" && debug != "no"
    
    if debugMode {
        log.Println("⚠️  DEBUG MODE ENABLED")
        log.Println("⚠️  This mode should NEVER be used in production!")
    }
}

func main() {
    // Pass debug mode to server
    server := server.New(airportSvc, geoipSvc, debugMode)
    // ...
}
```

**Server Configuration:**

```go
func (s *Server) setupRouter() {
    r := chi.NewRouter()
    
    // Standard middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    // Debug mode specific
    if s.debugMode {
        // Verbose logging
        r.Use(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                start := time.Now()
                log.Printf("[DEBUG] %s %s", r.Method, r.URL.Path)
                next.ServeHTTP(w, r)
                log.Printf("[DEBUG] %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
            })
        })
        
        // Permissive CORS
        r.Use(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Access-Control-Allow-Origin", "*")
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "*")
                if r.Method == "OPTIONS" {
                    w.WriteHeader(http.StatusOK)
                    return
                }
                next.ServeHTTP(w, r)
            })
        })
        
        // Debug endpoints
        r.Get("/debug/routes", s.handleDebugRoutes)
        r.Get("/debug/config", s.handleDebugConfig)
        r.Get("/debug/db", s.handleDebugDB)
        r.Get("/debug/settings", s.handleDebugSettings)
        r.Get("/debug/memory", s.handleDebugMemory)
        r.Post("/debug/reload", s.handleDebugReload)
        r.Post("/debug/gc", s.handleDebugGC)
    }
    
    // ... regular routes ...
}
```

---

#### **Debug Endpoints Implementation**

**Route Listing:**
```go
func (s *Server) handleDebugRoutes(w http.ResponseWriter, r *http.Request) {
    if !s.debugMode {
        http.Error(w, "Debug mode not enabled", http.StatusNotFound)
        return
    }
    
    routes := []string{}
    chi.Walk(s.router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
        routes = append(routes, fmt.Sprintf("%s %s", method, route))
        return nil
    })
    
    s.respondJSON(w, http.StatusOK, map[string]interface{}{
        "routes": routes,
        "count":  len(routes),
    })
}
```

**Memory Stats:**
```go
func (s *Server) handleDebugMemory(w http.ResponseWriter, r *http.Request) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    s.respondJSON(w, http.StatusOK, map[string]interface{}{
        "alloc_mb":       m.Alloc / 1024 / 1024,
        "total_alloc_mb": m.TotalAlloc / 1024 / 1024,
        "sys_mb":         m.Sys / 1024 / 1024,
        "num_gc":         m.NumGC,
        "goroutines":     runtime.NumGoroutine(),
    })
}
```

---

#### **Usage**

**Enable Debug Mode:**

```bash
# Run with debug
DEBUG=1 {projectname}

# Docker
docker run -e DEBUG=1 {projectname}:dev

# Docker Compose
environment:
  - DEBUG=1
  
# systemd
Environment="DEBUG=1"
```

**Access Debug Endpoints:**
```bash
# List all routes
curl http://localhost:64555/debug/routes

# View configuration
curl http://localhost:64555/debug/config

# Database stats
curl http://localhost:64555/debug/db

# Memory usage
curl http://localhost:64555/debug/memory

# Force config reload
curl -X POST http://localhost:64555/debug/reload

# Trigger GC
curl -X POST http://localhost:64555/debug/gc
```

---

**Security Warning:**

⚠️ **NEVER enable DEBUG mode in production!**

Debug mode:
- Exposes internal system details
- Allows permissive CORS
- Shows detailed error messages
- May impact performance (verbose logging)
- Disables some security checks

Use ONLY in development/testing environments.

---

**No CLI flags for dev/daemon - cleaner interface with environment variable!**

## 22. Built-in Scheduler

**Purpose**: Execute periodic tasks without external cron dependencies

**Required for ALL projects** - Used for background tasks like GeoIP updates, log rotation, etc.

### Implementation:

**Location**: `src/scheduler/scheduler.go`

```go
package scheduler

import (
    "log"
    "sync"
    "time"
)

type Task struct {
    Name     string
    Schedule string           // Cron-like: "0 3 * * 0" (Sunday 3:00 AM)
    Func     func() error
    ticker   *time.Ticker
}

type Scheduler struct {
    tasks []*Task
    wg    sync.WaitGroup
    stop  chan struct{}
}

func New() *Scheduler {
    return &Scheduler{
        tasks: make([]*Task, 0),
        stop:  make(chan struct{}),
    }
}

func (s *Scheduler) AddTask(name, schedule string, fn func() error) {
    task := &Task{
        Name:     name,
        Schedule: schedule,
        Func:     fn,
    }
    s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Start() {
    for _, task := range s.tasks {
        s.wg.Add(1)
        go s.runTask(task)
    }
    log.Println("📅 Scheduler started")
}

func (s *Scheduler) Stop() {
    close(s.stop)
    s.wg.Wait()
    log.Println("📅 Scheduler stopped")
}

func (s *Scheduler) runTask(task *Task) {
    defer s.wg.Done()
    
    interval := parseCronSchedule(task.Schedule)
    task.ticker = time.NewTicker(interval)
    defer task.ticker.Stop()
    
    for {
        select {
        case <-task.ticker.C:
            log.Printf("🔄 Running scheduled task: %s", task.Name)
            if err := task.Func(); err != nil {
                log.Printf("❌ Task %s failed: %v", task.Name, err)
            } else {
                log.Printf("✅ Task %s completed", task.Name)
            }
        case <-s.stop:
            return
        }
    }
}

func parseCronSchedule(schedule string) time.Duration {
    // Simple parser for common patterns
    // "0 3 * * 0" = Weekly on Sunday at 3:00 AM
    // For production, use github.com/robfig/cron or similar
    return 7 * 24 * time.Hour  // Weekly
}
```

### Usage in main.go:

```go
// Initialize scheduler
sched := scheduler.New()

// Add tasks
sched.AddTask("geoip-update", "0 3 * * 0", func() error {
    return geoipSvc.UpdateDatabases()
})

sched.AddTask("log-rotate", "0 0 * * *", func() error {
    return rotateLogFiles()
})

// Start scheduler
sched.Start()
defer sched.Stop()
```

### Common Scheduled Tasks:
- GeoIP database updates (weekly)
- Log file rotation (daily)
- Database cleanup (daily)
- Statistics aggregation (hourly)
- Health checks (every 5 minutes)

---

## 23. HTTPS/TLS & Let's Encrypt

**Built-in Let's Encrypt integration** with full DNS support

### Automatic HTTPS:

**If PORT is `80,443`**:
- Binary automatically obtains Let's Encrypt certificate
- Serves HTTP on :80, HTTPS on :443
- Auto-renewal handled by binary

**If PORT is single value** (e.g., `8080`):
- HTTP only, no automatic HTTPS

**If PORT is comma-separated** (e.g., `8080,8443`):
- First port = HTTP
- Second port = HTTPS
- Auto-obtain certificate

### Certificate Storage:

**Check existing certs first** (literal path, NOT variable):
```
/etc/letsencrypt/live/domain/
├── cert.pem
├── chain.pem
├── fullchain.pem
└── privkey.pem
```

**If found**: Use existing cert (read-only, don't manage)

**If not found**: Get new cert and save to:
```
/etc/{projectname}/ssl/certs/
├── cert.pem
├── chain.pem
├── fullchain.pem
└── privkey.pem
```

### ACME Challenge Methods:

**Support ALL three methods** (configurable via Admin WebUI):

1. **DNS-01** (Preferred for wildcards)
   - Support all major DNS providers:
     - Cloudflare, Route53, DigitalOcean, Google Cloud DNS
     - Namecheap, GoDaddy, Gandi
     - Custom (RFC2136 - DNS UPDATE)
   - Wildcard certificates: `*.example.com`

2. **TLS-ALPN-01**
   - Works on port 443
   - No HTTP needed
   - Firewall-friendly

3. **HTTP-01**
   - Works on port 80
   - Simple, widely supported
   - No DNS configuration needed

### Admin WebUI Configuration:

**Page**: `/admin/ssl`

**Settings:**
```yaml
ssl.enabled: true
ssl.domain: api.example.com
ssl.email: admin@example.com
ssl.acme_method: dns-01
ssl.dns_provider: cloudflare
ssl.dns_api_token: <encrypted>
ssl.auto_renew: true
ssl.renew_days_before: 30
```

### Implementation (pseudo-code):

```go
// Check if HTTPS is requested
ports := strings.Split(portFlag, ",")
httpPort := ports[0]
httpsPort := ""

if len(ports) == 2 {
    httpsPort = ports[1]
}

// Auto-enable HTTPS if standard ports
if httpPort == "80" && httpsPort == "" {
    httpsPort = "443"
}

if httpsPort != "" {
    // Check existing cert first (literal path)
    certPath := "/etc/letsencrypt/live/" + domain + "/fullchain.pem"
    keyPath := "/etc/letsencrypt/live/" + domain + "/privkey.pem"
    
    if !fileExists(certPath) {
        // Get Let's Encrypt certificate
        cert, err := acme.ObtainCertificate(domain, email, method)
        if err != nil {
            log.Printf("Failed to obtain certificate: %v", err)
            log.Println("Continuing with HTTP only...")
            httpsPort = ""
        } else {
            // Save to /etc/{projectname}/ssl/certs/
            saveCertificate(cert)
        }
    }
    
    if httpsPort != "" {
        go startHTTPS(httpsPort, certPath, keyPath)
    }
}

startHTTP(httpPort)
```

---

## 24. Console Output Standards

**Philosophy**: Minimal, pretty, emoji-enhanced, informative

### Output Rules:

✅ **DO:**
- Use emojis for visual categorization
- Keep output concise (one line per event)
- Use colors (if terminal supports)
- Show only essential information
- Group related output
- Use clear formatting

❌ **DON'T:**
- Spam verbose logs to console
- Print stack traces to console (logs only)
- Show debug info in production
- Repeat the same message
- Use technical jargon for user-facing messages

### Emoji Guide:

| Emoji | Meaning | Usage |
|-------|---------|-------|
| 🚀 | Starting | "🚀 Server starting..." |
| ✅ | Success | "✅ Database initialized" |
| ❌ | Error | "❌ Failed to connect" |
| ⚠️ | Warning | "⚠️ GeoIP unavailable" |
| 📁 | Directory | "📁 Config: /etc/app" |
| 📂 | Data | "📂 Loaded 29,293 airports" |
| 🔧 | Configuration | "🔧 Mode: Production" |
| 🔐 | Security | "🔐 Admin credentials saved" |
| 🌐 | Network | "🌐 Listening on http://server:64555" |
| 📊 | Statistics | "📊 Memory: 128MB" |
| 🔄 | Reload | "🔄 Configuration reloaded" |
| ⏸️ | Pause/Stop | "⏸️ Shutting down..." |
| 📅 | Scheduler | "📅 Scheduler started" |
| 🔒 | SSL/HTTPS | "🔒 Certificate obtained" |
| 🛡️ | Security event | "🛡️ Rate limit exceeded" |
| 📝 | Logging | "📝 Logs: /var/log/app" |

### Example Startup Output:

```
🚀 Starting {projectname} v0.0.1
📁 Config:  /etc/{projectname}
📁 Data:    /var/lib/{projectname}
📁 Logs:    /var/log/{projectname}
✅ Database initialized
✅ Admin authentication ready
🔐 Credentials: /etc/{projectname}/admin_credentials
📂 Loaded 29,293 airports from 236 countries
✅ GeoIP databases loaded (4 files, 103MB)
📅 Scheduler started (1 task)
🌐 Server listening on http://api.example.com:64555
📖 API Docs: http://api.example.com:64555/openapi
🔐 Admin Panel: http://api.example.com:64555/admin
```

### Error Output:

```
❌ Database connection failed
⚠️  Entering read-only maintenance mode
🔧 Admin UI available for diagnostics: http://server:64555/admin
💡 Attempting self-healing...
```

### Shutdown Output:

```
⏸️  Shutting down gracefully...
✅ HTTP server stopped
✅ Database connections closed
✅ GeoIP databases closed
✅ Scheduler stopped
👋 Shutdown complete
```

---

## 25. Logging Standards

### Log Files:

**access.log** - HTTP access log (Apache Combined format by default)
```
127.0.0.1 - - [14/Oct/2025:12:00:00 -0400] "GET /api/v1/airports/JFK HTTP/1.1" 200 1234 "-" "curl/7.68.0"
```

**Configurable via Admin WebUI** (`/admin/settings`):
- `log.access_format`: `apache`, `json`, `common`, `combined`

**error.log** - Application errors
```json
{"timestamp":"2025-10-14T12:00:00Z","level":"error","message":"Failed to connect to database","error":"connection refused"}
```

**audit.log** - Security events
```json
{"timestamp":"2025-10-14T12:00:00Z","event":"FAILED_LOGIN","ip":"192.0.2.1","details":"Invalid credentials"}
```

**{projectname}.log** - General application log
```json
{"timestamp":"2025-10-14T12:00:00Z","level":"info","message":"Server started","port":64555}
```

### Log Formats (Configurable):

**Apache Combined** (default for access.log):
```
%h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-Agent}i"
```

**JSON** (structured logging):
```json
{
  "timestamp": "2025-10-14T12:00:00Z",
  "method": "GET",
  "path": "/api/v1/airports/JFK",
  "status": 200,
  "bytes": 1234,
  "duration_ms": 5,
  "ip": "127.0.0.1",
  "user_agent": "curl/7.68.0"
}
```

### Log Rotation:

**Built-in rotation** (via scheduler):
- Daily rotation at midnight
- Keep 7 days by default
- Compress old logs (.gz)
- Configurable via Admin WebUI

**Or use logrotate** (SIGUSR1 support):
```
# /etc/logrotate.d/{projectname}
/var/log/{projectname}/*.log {
    daily
    rotate 7
    compress
    delaycompress
    postrotate
        kill -USR1 $(cat /var/run/{projectname}.pid) || true
    endscript
}
```

---

## 26. Maintenance Mode & Self-Healing

### Critical Error: Database Connection Failure

**When database connection fails:**

1. ✅ **Don't crash** - Stay running
2. ✅ **Enter read-only maintenance mode**
3. ✅ **Use sane built-in defaults** (in-memory config)
4. ✅ **Admin UI still works** (shows diagnostics)
5. ✅ **Attempt self-healing**

### Maintenance Mode Behavior:

**What Still Works:**
- ✅ HTTP server (running)
- ✅ Admin WebUI (shows diagnostics)
- ✅ Health check endpoint (returns degraded status)
- ✅ Static files (CSS, JS, images)
- ✅ Built-in defaults (rate limiting, CORS, etc.)

**What Doesn't Work:**
- ❌ Saving configuration changes
- ❌ Data queries (no database)
- ❌ Admin credential changes
- ❌ Settings persistence

### Admin UI in Maintenance Mode:

**Page**: `/admin/maintenance`

```html
<div class="maintenance-banner">
  ⚠️ READ-ONLY MAINTENANCE MODE
</div>

<div class="diagnostics">
  <h2>Database Connection Error</h2>
  
  <p>The database is currently unavailable. The server is running with built-in defaults.</p>
  
  <h3>Troubleshooting Steps:</h3>
  <ol>
    <li>Check database file exists: <code>/var/lib/{projectname}/db/{projectname}.db</code></li>
    <li>Verify file permissions (should be readable/writable)</li>
    <li>Check disk space: <code>df -h</code></li>
    <li>Verify SQLite installation (if external database)</li>
  </ol>
  
  <h3>Self-Healing Attempts:</h3>
  <ul>
    <li>✅ Attempted to create database file</li>
    <li>✅ Attempted to repair corrupt database</li>
    <li>❌ Connection still failing</li>
  </ul>
  
  <h3>Manual Fix:</h3>
  <pre>
# Check if database file exists
ls -la /var/lib/{projectname}/db/{projectname}.db

# If missing, restart service (will recreate)
systemctl restart {projectname}

# If corrupt, restore from backup
cp /var/lib/{projectname}/backups/latest.db /var/lib/{projectname}/db/{projectname}.db
systemctl restart {projectname}
  </pre>
  
  <button onclick="retryConnection()">🔄 Retry Connection</button>
</div>
```

### Self-Healing Attempts:

**Automatic recovery attempts:**

1. **Create database if missing**
   ```go
   if !fileExists(dbPath) {
       if err := database.Initialize(dbPath); err == nil {
           log.Println("✅ Database created successfully")
           return // Healed
       }
   }
   ```

2. **Repair corrupt database**
   ```go
   if err := database.Repair(dbPath); err == nil {
       log.Println("✅ Database repaired")
       return // Healed
   }
   ```

3. **Use backup if available**
   ```go
   backupPath := filepath.Join(dataDir, "backups", "latest.db")
   if fileExists(backupPath) {
       if err := copyFile(backupPath, dbPath); err == nil {
           log.Println("✅ Restored from backup")
           return // Healed
       }
   }
   ```

4. **Periodic retry** (every 30 seconds)
   ```go
   go func() {
       ticker := time.NewTicker(30 * time.Second)
       for range ticker.C {
           if err := database.Ping(); err == nil {
               log.Println("✅ Database connection restored")
               exitMaintenanceMode()
               return
           }
       }
   }()
   ```

---

## 27. WebUI Specifications

### Template Structure:

**Required Templates:**
- `base.html` - Base layout (header, nav, footer)
- `home.html` - Homepage
- `search.html` - Search page
- `{resource}.html` - Resource-specific pages
- `admin/*.html` - Admin pages

**Template Hierarchy:**
```
base.html (defines blocks)
  ├── header (logo, nav, theme toggle)
  ├── main (content area)
  └── footer (centered, bottom)

Page templates extend base:
  home.html → {{template "base" .}}
  search.html → {{template "base" .}}
```

### Layout Specifications:

**Responsive Design:**

**Desktop (≥720px):**
- Content width: 90%
- Left margin: 5%
- Right margin: 5%
- Max-width: 1200px (optional)

**Mobile (<720px):**
- Content width: 98%
- Left margin: 1%
- Right margin: 1%

**Footer:**
- Always centered horizontally
- Always at bottom of viewport
- Scroll to see if content is long
- CSS: `position: sticky; bottom: 0;` or similar

### HTML5 Standards:

✅ **Required:**
- Semantic HTML5 tags (`<header>`, `<nav>`, `<main>`, `<footer>`, `<article>`, `<section>`)
- Full mobile support (viewport meta tag)
- Accessibility (ARIA labels, alt text, keyboard navigation)
- Readable (proper contrast ratios)
- Navigable (logical tab order)
- Intuitive (clear labels, helpful errors)
- Self-explanatory (tooltips, placeholders, help text)

### Assets:

**Logo & Favicon:**
- Can be local: `/static/images/logo.png`
- Can be remote: `https://cdn.example.com/logo.png`
- Support scaling (CSS transform or multiple sizes)

**Example:**
```html
<img src="/static/images/logo.png" 
     alt="{ProjectName} Logo" 
     width="32" height="32"
     srcset="/static/images/logo@2x.png 2x">
```

### Vanilla JS & CSS:

**Technology:**
- ✅ Vanilla JavaScript (NO jQuery, React, Vue, etc.)
- ✅ Plain CSS with custom properties (NO Bootstrap, Tailwind, etc.)
- ✅ Web Components (optional, if needed)

**Why:**
- Embedded in binary (small size critical)
- No external dependencies
- Fast loading
- Full control

### Notifications:

**Bell Icon** (🔔) - Built-in notification system

**Implementation:**
```javascript
// Show notification
function notify(message, type = 'info') {
    // Visual notification
    showToast(message, type);
    
    // Update bell icon badge
    const bell = document.querySelector('.notification-bell');
    const badge = bell.querySelector('.badge');
    const count = parseInt(badge.textContent) + 1;
    badge.textContent = count;
    badge.style.display = 'block';
}

// Bell icon HTML
<button class="notification-bell" onclick="showNotifications()">
    🔔
    <span class="badge" style="display:none;">0</span>
</button>
```

### HTTP Error Codes:

**Use standard codes:**
- 200 OK - Success
- 400 Bad Request - Invalid input
- 401 Unauthorized - No auth
- 403 Forbidden - No permission
- 404 Not Found - Resource doesn't exist
- 429 Too Many Requests - Rate limit
- 500 Internal Server Error - Server error
- 503 Service Unavailable - Maintenance mode

---

## 28. Terminology & Notation ⭐ REQUIRED

### Variables vs Literals:

**Variables** (wrapped in `{}`):
- `{projectname}` - Variable, replace with actual project name
- `/etc/{projectname}/` - Variable path component
- `{organization}` - Variable
- `{PORT}` - Variable

**Literals** (no `{}`):
- `/etc/letsencrypt/live/domain` - Literal path, don't replace "domain"
- `PORT` - Literal environment variable name
- `ADDRESS` - Literal environment variable name

### Terminology:

**Correct Terms:**
- ✅ **Server** - Our application (the API server)
- ✅ **System** - The operating system (Linux, macOS, Windows)
- ✅ **Server ADDRESS** - Listen address (not "server name")
- ✅ **User's system** - Where binary RUNS (deployment target)
- ✅ **Local system** - Where binary is BUILT (developer machine)
- ✅ **WebUI** - Web-based user interface
- ✅ **Admin Panel** - Admin web interface at /admin

**Incorrect Terms:**
- ❌ "Server name" - Use "server address" or "hostname"
- ❌ "System" when referring to our server - Use "server"
- ❌ "App" - Use "server" or "binary"

---

**All specifications added. The SPEC now covers every aspect in complete detail!**

---

## 29. OS-Specific Directory Layouts ⭐ REQUIRED

**Purpose**: Automatic directory detection based on OS and privileges (root vs user)

**Implementation**: `src/paths/paths.go` package

---

### Directory Paths by OS:

**Linux/BSD (with root privileges):**
```yaml
Config:  /etc/{projectname}/
Data:    /var/lib/{projectname}/
Logs:    /var/log/{projectname}/
Runtime: /run/{projectname}/
```

**Linux/BSD (without root):**
```yaml
Config:  ~/.config/{projectname}/
Data:    ~/.local/share/{projectname}/
Logs:    ~/.local/state/{projectname}/
Runtime: ~/.local/run/{projectname}/
```

**macOS (with privileges):**
```yaml
Config:  /Library/Application Support/{ProjectName}/
Data:    /Library/Application Support/{ProjectName}/data/
Logs:    /Library/Logs/{ProjectName}/
Runtime: /var/run/{projectname}/
```

**macOS (without privileges):**
```yaml
Config:  ~/Library/Application Support/{ProjectName}/
Data:    ~/Library/Application Support/{ProjectName}/data/
Logs:    ~/Library/Logs/{ProjectName}/
Runtime: ~/Library/Application Support/{ProjectName}/run/
```

**Windows (system/admin):**
```yaml
Config:  C:\ProgramData\{ProjectName}\config\
Data:    C:\ProgramData\{ProjectName}\data\
Logs:    C:\ProgramData\{ProjectName}\logs\
Runtime: C:\ProgramData\{ProjectName}\run\
```

**Windows (user):**
```yaml
Config:  %APPDATA%\{ProjectName}\config\
Data:    %APPDATA%\{ProjectName}\data\
Logs:    %APPDATA%\{ProjectName}\logs\
Runtime: %APPDATA%\{ProjectName}\run\
```

**Docker (all platforms):**
```yaml
Config:  /config
Data:    /data
Logs:    /logs
Runtime: /tmp
```

---

### Implementation (src/paths/paths.go):

```go
package paths

import (
    "fmt"
    "os"
    "os/user"
    "path/filepath"
    "runtime"
)

// GetDefaultDirs returns OS-specific default directories
// projectname should be lowercase for Linux/BSD, PascalCase for macOS/Windows
func GetDefaultDirs(projectname string) (configDir, dataDir, logsDir string) {
    current, err := user.Current()
    if err != nil {
        return getSystemDirs(projectname)
    }

    isRoot := current.Uid == "0"

    switch runtime.GOOS {
    case "linux", "freebsd", "openbsd", "netbsd":
        if isRoot {
            configDir = fmt.Sprintf("/etc/%s", projectname)
            dataDir = fmt.Sprintf("/var/lib/%s", projectname)
            logsDir = fmt.Sprintf("/var/log/%s", projectname)
        } else {
            homeDir := current.HomeDir
            configDir = filepath.Join(homeDir, ".config", projectname)
            dataDir = filepath.Join(homeDir, ".local", "share", projectname)
            logsDir = filepath.Join(homeDir, ".local", "state", projectname)
        }

    case "darwin":
        projectnamePascal := toPascalCase(projectname)
        if isRoot {
            configDir = fmt.Sprintf("/Library/Application Support/%s", projectnamePascal)
            dataDir = fmt.Sprintf("/Library/Application Support/%s/data", projectnamePascal)
            logsDir = fmt.Sprintf("/Library/Logs/%s", projectnamePascal)
        } else {
            homeDir := current.HomeDir
            configDir = filepath.Join(homeDir, "Library", "Application Support", projectnamePascal)
            dataDir = filepath.Join(homeDir, "Library", "Application Support", projectnamePascal, "data")
            logsDir = filepath.Join(homeDir, "Library", "Logs", projectnamePascal)
        }

    case "windows":
        projectnamePascal := toPascalCase(projectname)
        if isAdmin() {
            programData := os.Getenv("ProgramData")
            if programData == "" {
                programData = "C:\\ProgramData"
            }
            configDir = filepath.Join(programData, projectnamePascal, "config")
            dataDir = filepath.Join(programData, projectnamePascal, "data")
            logsDir = filepath.Join(programData, projectnamePascal, "logs")
        } else {
            appData := os.Getenv("APPDATA")
            if appData == "" {
                appData = filepath.Join(current.HomeDir, "AppData", "Roaming")
            }
            configDir = filepath.Join(appData, projectnamePascal, "config")
            dataDir = filepath.Join(appData, projectnamePascal, "data")
            logsDir = filepath.Join(appData, projectnamePascal, "logs")
        }

    default:
        // Unknown OS - portable defaults
        configDir = "./config"
        dataDir = "./data"
        logsDir = "./logs"
    }

    return configDir, dataDir, logsDir
}

func isAdmin() bool {
    // Simple check for Windows admin privileges
    _, err := os.Open("\\\\.\\PHYSICALDRIVE0")
    return err == nil
}

func toPascalCase(s string) string {
    if len(s) == 0 {
        return s
    }
    return strings.ToUpper(s[:1]) + s[1:]
}

// EnsureDir creates directory if it doesn't exist
func EnsureDir(path string) error {
    return os.MkdirAll(path, 0755)
}

// EnsureDirs creates all required directories
func EnsureDirs(configDir, dataDir, logsDir string) error {
    for _, dir := range []string{configDir, dataDir, logsDir} {
        if err := EnsureDir(dir); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
    }
    
    // Also create database subdirectory
    dbDir := filepath.Join(dataDir, "db")
    if err := EnsureDir(dbDir); err != nil {
        return fmt.Errorf("failed to create db directory: %w", err)
    }
    
    return nil
}
```

### Usage in main.go:

```go
import "github.com/{org}/{project}/src/paths"

func main() {
    // Get OS-specific directories (unless overridden by flags/env)
    configDir := getEnv("CONFIG_DIR", "")
    dataDir := getEnv("DATA_DIR", "")
    logsDir := getEnv("LOGS_DIR", "")
    
    // If not overridden, use OS defaults
    if configDir == "" || dataDir == "" || logsDir == "" {
        defaultConfig, defaultData, defaultLogs := paths.GetDefaultDirs("{projectname}")
        if configDir == "" {
            configDir = defaultConfig
        }
        if dataDir == "" {
            dataDir = defaultData
        }
        if logsDir == "" {
            logsDir = defaultLogs
        }
    }
    
    // Create directories
    if err := paths.EnsureDirs(configDir, dataDir, logsDir); err != nil {
        log.Fatal(err)
    }
    
    log.Printf("📁 Config: %s", configDir)
    log.Printf("📁 Data:   %s", dataDir)
    log.Printf("📁 Logs:   %s", logsDir)
}
```

### Priority Order:

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **OS-specific defaults** (lowest priority)

### Database Location:

**Always in data directory:**
```
{DATA_DIR}/db/{projectname}.db
```

**Examples:**
- Linux (root): `/var/lib/{projectname}/db/{projectname}.db`
- Linux (user): `~/.local/share/{projectname}/db/{projectname}.db`
- macOS (root): `/Library/Application Support/{ProjectName}/data/db/{projectname}.db`
- Windows (admin): `C:\ProgramData\{ProjectName}\data\db\{projectname}.db`
- Docker: `/data/db/{projectname}.db`

---

**OS-specific paths ensure portability across all platforms!**
