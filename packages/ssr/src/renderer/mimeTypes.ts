const contentTypeMap = {
    ".html": "text/html; charset=utf-8",
    ".css": "text/css; charset=utf-8",
    ".js": "application/javascript",
    ".png": "image/png",
    ".jpg": "image/jpeg",
    ".gif": "image/gif",
    ".svg": "image/svg+xml",
    ".json": "application/json",
    ".txt": "text/plain; charset=utf-8",
} as const satisfies Record<string, string>;

export function mimeTypeForUrl(path: string) {
    let mimeType = Object.entries(contentTypeMap).find(([fileExt]) =>
        path.endsWith(fileExt),
    ) ?? ["", ""];

    return mimeType[1];
}
