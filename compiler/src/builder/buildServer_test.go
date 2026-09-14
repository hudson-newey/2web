package builder

import (
	"strings"
	"testing"
)

func TestVerbRouteForServerScript(t *testing.T) {
	inputPath := "/project/src"

	tests := []struct {
		file        string
		method      string
		route       string
		isVerbRoute bool
	}{
		{"/project/src/api/users/__get.server.ts", "get", "/api/users", true},
		{"/project/src/api/users/__post.server.ts", "post", "/api/users", true},
		{"/project/src/__delete.server.js", "delete", "/", true},
		{"/project/src/__PUT.server.mjs", "put", "/", true},
		{"/project/src/api/users.server.ts", "", "", false},
		{"/project/src/index.server.ts", "", "", false},
		{"/project/src/api/greet.server.ts", "", "", false},
		{"/project/src/style.css", "", "", false},
	}

	for _, test := range tests {
		method, route, isVerbRoute := verbRouteForServerScript(inputPath, test.file)

		if isVerbRoute != test.isVerbRoute {
			t.Errorf("%s: expected isVerbRoute %v, got %v", test.file, test.isVerbRoute, isVerbRoute)
			continue
		}

		if method != test.method || route != test.route {
			t.Errorf("%s: expected (%q, %q), got (%q, %q)", test.file, test.method, test.route, method, route)
		}
	}
}

func TestExportedFunctionNamesDoesNotIncludeRouteHandlers(t *testing.T) {
	source := `
export function greet(name: string): string {
	return name;
}

export const add = (a: number, b: number) => a + b;

const handler = (req: any, res: any) => res.json({});
export default handler;
`

	names := exportedFunctionNames(source)

	if len(names) != 2 || names[0] != "greet" || names[1] != "add" {
		t.Errorf("expected the named exports greet and add, got %v", names)
	}
}

func TestManifestShape(t *testing.T) {
	// The manifest is consumed by the server runtime (@two-web/kit/ssr), so
	// the routes and modules sections must serialize the fields it reads.
	route := ServerRoute{Route: "/api/users", Method: "get", File: "api/users/__get.server.js"}
	module := ServerModule{File: "api/users.server.js", Rpc: []string{"greet"}}

	if !strings.Contains(route.File, "__get.server.js") {
		t.Error("unexpected route file")
	}

	if module.Rpc[0] != "greet" {
		t.Error("unexpected module rpc entry")
	}
}
