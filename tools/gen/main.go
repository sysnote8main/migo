// Command gen generates Go source files for all Misskey API endpoints
// from the OpenAPI specification (api.json).
//
// Usage: go run tools/gen/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	specData, err := os.ReadFile("api.json")
	if err != nil {
		panic(err)
	}
	var spec Spec
	if err := json.Unmarshal(specData, &spec); err != nil {
		panic(err)
	}

	knownTypes := map[string]string{
		"Error": "any", "UserLite": "types.UserLite",
		"UserDetailedNotMeOnly": "types.UserDetailedNotMeOnly",
		"MeDetailedOnly": "types.MeDetailedOnly",
		"MeDetailed": "types.MeDetailed",
		"UserDetailed": "types.UserDetailed",
		"User": "types.UserLite",
		"Note": "types.Note", "NoteDraft": "types.Note",
		"DriveFile": "types.DriveFile", "DriveFolder": "types.DriveFolder",
		"Drive": "types.Drive", "Channel": "types.Channel",
		"Poll": "types.Poll", "Page": "types.Page",
		"Announcement": "types.Announcement",
		"EmojiDetailed": "types.EmojiDetailed",
		"ChatMessage": "types.ChatMessage", "ChatRoom": "types.ChatRoom",
		"ChatRoomInvitation": "types.ChatRoomInvitation",
		"Notification": "types.Notification",
		"AbuseReport": "types.AbuseReport",
		"Meta": "types.Meta",
	}

	pkgMap := map[string]string{
		"i": "account", "admin": "admin", "notes": "notes", "users": "users",
		"drive": "drive", "auth": "auth", "following": "following", "chat": "chat",
		"channels": "channels", "clips": "clips", "flash": "flash",
		"gallery": "gallery", "pages": "pages", "federation": "federation",
		"antennas": "antennas", "hashtags": "hashtags", "roles": "roles",
		"invite": "invite", "blocking": "blocking", "mute": "mute",
		"renote-mute": "renotemute", "announcements": "announcements",
		"ap": "ap", "app": "app", "charts": "charts",
		"bubble-game": "bubblegame", "reversi": "reversi", "sw": "sw",
		"email-address": "emailaddress", "emoji": "emoji", "emojis": "emojis",
		"endpoint": "endpoint", "export-custom-emojis": "export",
		"fetch-external-resources": "fetch", "fetch-rss": "fetch",
		"get-avatar-decorations": "get", "get-online-users-count": "get",
		"meta": "meta", "miauth": "auth", "my": "account",
		"page-push": "pages", "ping": "auth", "pinned-users": "get",
		"promo": "promo", "request-reset-password": "auth",
		"reset-db": "admin", "reset-password": "auth", "retention": "admin",
		"server-info": "get", "stats": "get", "test": "test",
		"username": "username", "v2": "admin", "verify-email": "auth",
	}

	covered := map[string]bool{
		// notes
		"/notes/create": true, "/notes/show": true, "/notes/delete": true,
		"/notes/search": true, "/notes/reactions/create": true, "/notes/reactions/delete": true,
		// users
		"/users/show": true, "/users/search": true,
		// drive
		"/drive": true, "/drive/files": true, "/drive/files/show": true,
		"/drive/folders": true, "/drive/folders/create": true, "/drive/folders/update": true,
		"/drive/folders/delete": true,
		// auth
		"/auth/session/generate": true, "/auth/session/userkey": true, "/ping": true,
		"/auth/session/show": true, "/auth/accept": true,
		// following
		"/following/create": true, "/following/delete": true, "/following/invalidate": true,
		"/following/list": true, "/following/requests/list": true,
		"/following/requests/accept": true, "/following/requests/reject": true,
		"/following/requests/cancel": true,
		// notification
		"/i/notifications": true, "/i/notifications-grouped": true,
		"/notifications/mark-all-as-read": true, "/notifications/create": true,
		"/notifications/flush": true, "/notifications/test-notification": true,
		// timeline
		"/notes/timeline": true, "/notes/local-timeline": true, "/notes/hybrid-timeline": true,
		"/notes/global-timeline": true, "/channels/timeline": true,
		// chat
		"/chat/history": true, "/chat/messages/create-to-user": true,
		"/chat/messages/create-to-room": true, "/chat/messages/delete": true,
		"/chat/messages/show": true, "/chat/messages/react": true, "/chat/messages/unreact": true,
		"/chat/messages/user-timeline": true, "/chat/messages/room-timeline": true,
		"/chat/messages/search": true, "/chat/read-all": true,
		"/chat/rooms/create": true, "/chat/rooms/update": true, "/chat/rooms/delete": true,
		"/chat/rooms/show": true, "/chat/rooms/join": true, "/chat/rooms/leave": true,
		"/chat/rooms/mute": true, "/chat/rooms/joining": true, "/chat/rooms/owned": true,
		"/chat/rooms/members": true, "/chat/rooms/invitations/create": true,
		"/chat/rooms/invitations/ignore": true, "/chat/rooms/invitations/inbox": true,
		"/chat/rooms/invitations/outbox": true,
		// admin
		"/admin/emoji/add": true, "/admin/emoji/list": true, "/admin/emoji/remove": true,
		"/admin/emoji/update": true, "/admin/emoji/set-aliases-bulk": true, "/admin/emoji/copy": true,
		"/admin/emoji/delete": true, "/admin/emoji/delete-bulk": true,
		"/admin/announcements/create": true, "/admin/announcements/list": true,
		"/admin/announcements/update": true, "/admin/announcements/delete": true,
		"/admin/accounts/create": true, "/admin/accounts/delete": true, "/admin/accounts/find-by-email": true,
		"/admin/ad/create": true, "/admin/ad/list": true, "/admin/ad/update": true, "/admin/ad/delete": true,
		"/admin/invite/create": true, "/admin/invite/list": true,
		"/admin/roles/list": true, "/admin/roles/create": true, "/admin/roles/delete": true,
		"/admin/roles/assign": true, "/admin/roles/unassign": true,
		"/admin/abuse-user-reports": true, "/admin/resolve-abuse-user-report": true,
		"/admin/suspend-user": true, "/admin/unsuspend-user": true,
		"/admin/show-user": true, "/admin/update-user-note": true,
		"/admin/meta": true, "/admin/update-meta": true,
		"/admin/drive/files": true, "/admin/drive/cleanup": true,
		"/admin/drive/clean-remote-files": true, "/admin/drive/show-file": true,
		// base
		"/i": true,
	}

	// Group endpoints by package
	epByPkg := map[string][]string{}
	for path := range spec.Paths {
		if covered[path] {
			continue
		}
		seg := strings.Split(strings.Trim(path, "/"), "/")[0]
		pkg := pkgMap[seg]
		if pkg == "" {
			pkg = "misc"
		}
		epByPkg[pkg] = append(epByPkg[pkg], path)
	}

	// Packages that have existing code (only emit methods, not Service type)
	existingPkgs := map[string]bool{
		"admin": true, "notes": true, "users": true, "auth": true,
		"drive": true, "timeline": true, "following": true, "chat": true,
		"notification": true,
	}

	// Generate each package
	for pkg, paths := range epByPkg {
		genPkg(&spec, pkg, paths, knownTypes, pkgMap, existingPkgs[pkg])
	}

	fmt.Println("\nDone! Packages generated:")
	for pkg, paths := range epByPkg {
		fmt.Printf("  %s (%d endpoints)\n", pkg, len(paths))
	}
}

func genPkg(spec *Spec, pkgName string, paths []string, kt map[string]string, pm map[string]string, existing bool) {
	dir := pkgName
	os.MkdirAll(dir, 0755)

	// Phase 1: generate body to check imports needed
	body := &strings.Builder{}

	if !existing {
		pkgLabel := pkgName
		if pkgName == "account" {
			pkgLabel = "account (/i/)"
		}
		svcName := pascal(pkgName)
		fmt.Fprintf(body, `// Service provides %s-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new %s.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

`, pkgLabel, svcName+"Service")
	}

	for _, path := range paths {
		op := spec.Paths[path].Post
		if op == nil {
			continue
		}
		mn := goMethodName(path, pm)
		rp := getReqProps(op)
		rt := getRespType(op, kt)

		// Request struct
		if len(rp) > 0 {
			sn := mn + "Request"
			fmt.Fprintf(body, "// %s represents the request body for %s.\n", sn, trimPrefix(path))
			fmt.Fprintf(body, "type %s struct {\n", sn)
			for _, p := range rp {
				gt := propType(p, kt)
				fn := goFieldName(p.name)
				tag := p.name
				if p.req {
					fmt.Fprintf(body, "\t%s %s `json:\"%s\"`\n", fn, gt, tag)
				} else {
					fmt.Fprintf(body, "\t%s %s `json:\"%s,omitempty\"`\n", fn, gt, tag)
				}
			}
			fmt.Fprintf(body, "}\n\n")
		}

		// Method
		auth := op.Security != nil && len(op.Security) > 0
		if auth {
			fmt.Fprintf(body, "// %s calls %s (authenticated).\n", mn, path)
		} else {
			fmt.Fprintf(body, "// %s calls %s.\n", mn, path)
		}

		switch {
		case !auth:
			// Public endpoints
			if len(rp) > 0 && rt != "" {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context, req *%sRequest) (%s, error) {\n\tvar resp %s\n\tif err := s.cli.Do(ctx, %q, req, &resp); err != nil {\n\t\treturn *new(%s), err\n\t}\n\treturn resp, nil\n}\n\n", mn, mn, rt, rt, path, rt)
			} else if len(rp) > 0 {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context, req *%sRequest) error {\n\treturn s.cli.Do(ctx, %q, req, nil)\n}\n\n", mn, mn, path)
			} else if rt != "" {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context) (%s, error) {\n\tvar resp %s\n\tif err := s.cli.Do(ctx, %q, nil, &resp); err != nil {\n\t\treturn *new(%s), err\n\t}\n\treturn resp, nil\n}\n\n", mn, rt, rt, path, rt)
			} else {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context) error {\n\treturn s.cli.Do(ctx, %q, nil, nil)\n}\n\n", mn, path)
			}
		default:
			// Authenticated endpoints
			if len(rp) > 0 && rt != "" {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context, req *%sRequest) (%s, error) {\n\tvar resp %s\n\tif err := s.cli.Do(ctx, %q, req, &resp); err != nil {\n\t\treturn *new(%s), err\n\t}\n\treturn resp, nil\n}\n\n", mn, mn, rt, rt, path, rt)
			} else if len(rp) > 0 {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context, req *%sRequest) error {\n\treturn s.cli.Do(ctx, %q, req, nil)\n}\n\n", mn, mn, path)
			} else if rt != "" {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context) (%s, error) {\n\tvar resp %s\n\tif err := s.cli.Do(ctx, %q, nil, &resp); err != nil {\n\t\treturn *new(%s), err\n\t}\n\treturn resp, nil\n}\n\n", mn, rt, rt, path, rt)
			} else {
				fmt.Fprintf(body, "func (s *Service) %s(ctx context.Context) error {\n\treturn s.cli.Do(ctx, %q, nil, nil)\n}\n\n", mn, path)
			}
		}
	}

	// Phase 2: write output with correct imports
	out := &strings.Builder{}
	fmt.Fprintf(out, "// Code generated by tools/gen/main.go. DO NOT EDIT.\n\n")
	fmt.Fprintf(out, "package %s\n\n", pkgName)

	bodyStr := body.String()
	if strings.Contains(bodyStr, "types.") {
		fmt.Fprintf(out, "import (\n\t\"context\"\n\t\"github.com/sysnote8main/migo/types\"\n)\n\n")
	} else {
		fmt.Fprintf(out, "import \"context\"\n\n")
	}
	fmt.Fprintf(out, "%s", bodyStr)

	os.WriteFile(filepath.Join(dir, "gen_all.go"), []byte(out.String()), 0644)
}

// ---- Helpers ----

func getReqProps(op *Operation) []propInfo {
	if op.RequestBody == nil || op.RequestBody.Content == nil {
		return nil
	}
	jc, ok := op.RequestBody.Content["application/json"]
	if !ok || jc.Schema == nil {
		return nil
	}
	s := *jc.Schema
	var schemas []Schema
	if s.AllOf != nil {
		schemas = append(schemas, s.AllOf...)
	} else {
		schemas = append(schemas, s)
	}
	seen := map[string]bool{}
	var props []propInfo
	for _, sc := range schemas {
		if sc.AnyOf != nil {
			for _, a := range sc.AnyOf {
				if a.Properties != nil {
					reqSet := map[string]bool{}
					for _, r := range a.Required {
						reqSet[r] = true
					}
					for n, ps := range a.Properties {
						if seen[n] {
							continue
						}
						seen[n] = true
						props = append(props, propInfo{name: n, sch: ps, req: reqSet[n]})
					}
				}
			}
		}
		if sc.Properties != nil {
			reqSet := map[string]bool{}
			for _, r := range sc.Required {
				reqSet[r] = true
			}
			for n, ps := range sc.Properties {
				if seen[n] {
					continue
				}
				seen[n] = true
				props = append(props, propInfo{name: n, sch: ps, req: reqSet[n]})
			}
		}
	}
	return props
}

func getRespType(op *Operation, kt map[string]string) string {
	if op.Responses == nil {
		return ""
	}
	r200, ok := op.Responses["200"]
	if !ok || r200.Content == nil {
		return ""
	}
	jc, ok := r200.Content["application/json"]
	if !ok || jc.Schema == nil {
		return ""
	}
	return goType(*jc.Schema, kt)
}

func goType(s Schema, kt map[string]string) string {
	if s.Ref != "" {
		name := strings.TrimPrefix(s.Ref, "#/components/schemas/")
		if t, ok := kt[name]; ok {
			return t
		}
		return "any"
	}
	st := s.effType()
	switch st {
	case "string":
		if s.Format == "id" || s.Format == "misskey:id" {
			return "types.ID"
		}
		return "string"
	case "integer", "number":
		return "int64"
	case "boolean":
		return "bool"
	case "array":
		if s.Items != nil {
			return "[]" + goType(*s.Items, kt)
		}
		return "[]any"
	}
	return "any"
}

type propInfo struct {
	name string
	sch  Schema
	req  bool
}

func propType(p propInfo, kt map[string]string) string {
	s := p.sch
	st := s.effType()
	r := goType(s, kt)
	// nullable types: [string, null] etc.
	if len(s.Type) > 0 || st == "null" {
		if r == "string" || r == "int64" || r == "bool" || r == "types.ID" {
			return "*" + r
		}
	}
	return r
}

func goMethodName(path string, pm map[string]string) string {
	segs := strings.Split(strings.Trim(path, "/"), "/")
	if len(segs) == 0 {
		return ""
	}
	first := segs[0]

	// /i/xxx → strip /i prefix
	if first == "i" && len(segs) > 1 {
		return "I" + joinTitle(segs[1:])
	}

	// Strip the package-prefix segment for well-known categories
	pkg := pm[first]
	if pkg != "" && len(segs) > 1 {
		// Skip the first segment
		return joinTitle(segs[1:])
	}

	// Default: title all segments
	return joinTitle(segs)
}

func joinTitle(segs []string) string {
	var parts []string
	for _, s := range segs {
		parts = append(parts, titleSeg(s))
	}
	return strings.Join(parts, "")
}

func titleSeg(s string) string {
	var parts []string
	if strings.Contains(s, "-") {
		for _, w := range strings.Split(s, "-") {
			if len(w) > 0 {
				parts = append(parts, strings.ToUpper(w[:1])+w[1:])
			}
		}
	} else if strings.Contains(s, "_") {
		for _, w := range strings.Split(s, "_") {
			if len(w) > 0 {
				parts = append(parts, strings.ToUpper(w[:1])+w[1:])
			}
		}
	} else if len(s) > 0 {
		parts = append(parts, strings.ToUpper(s[:1])+s[1:])
	}
	return strings.Join(parts, "")
}

func goFieldName(s string) string {
	var parts []string
	if strings.Contains(s, "-") {
		for _, w := range strings.Split(s, "-") {
			if len(w) > 0 {
				parts = append(parts, strings.ToUpper(w[:1])+w[1:])
			}
		}
	} else if strings.Contains(s, "_") {
		for _, w := range strings.Split(s, "_") {
			if len(w) > 0 {
				parts = append(parts, strings.ToUpper(w[:1])+w[1:])
			}
		}
	} else {
		parts = []string{strings.ToUpper(s[:1]) + s[1:]}
	}
	name := strings.Join(parts, "")
	// Go keywords
	switch name {
	case "type":
		return "Type"
	case "range":
		return "Range"
	}
	return name
}

func pascal(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func trimPrefix(path string) string {
	return strings.TrimPrefix(path, "/")
}

// ---- OpenAPI types ----

type Spec struct {
	Paths      PathsMap                 `json:"paths"`
	Components struct {
		Schemas map[string]Schema `json:"schemas"`
	} `json:"components"`
}

type PathsMap map[string]PathItem

type PathItem struct {
	Post *Operation `json:"post"`
}

type Operation struct {
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	Tags        []string              `json:"tags"`
	Security    []map[string][]string `json:"security"`
	RequestBody *struct {
		Content map[string]struct {
			Schema *Schema `json:"schema"`
		} `json:"content"`
	} `json:"requestBody"`
	Responses map[string]struct {
		Content map[string]struct {
			Schema *Schema `json:"schema"`
		} `json:"content,omitempty"`
	} `json:"responses"`
}

type Schema struct {
	Ref         string            `json:"$ref"`
	Type        json.RawMessage   `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"`
	Items       *Schema           `json:"items,omitempty"`
	Required    []string          `json:"required,omitempty"`
	AllOf       []Schema          `json:"allOf,omitempty"`
	OneOf       []Schema          `json:"oneOf,omitempty"`
	AnyOf       []Schema          `json:"anyOf,omitempty"`
	Description string            `json:"description,omitempty"`
}

func (s Schema) effType() string {
	if s.Ref != "" {
		return "ref"
	}
	if len(s.Type) == 0 {
		if s.Properties != nil {
			return "object"
		}
		return "any"
	}
	var types []string
	if err := json.Unmarshal(s.Type, &types); err == nil {
		for _, t := range types {
			if t != "null" {
				return t
			}
		}
		return types[0]
	}
	var single string
	if err := json.Unmarshal(s.Type, &single); err == nil {
		return single
	}
	return "any"
}
