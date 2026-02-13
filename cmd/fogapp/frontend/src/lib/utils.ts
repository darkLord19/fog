// ── Utility functions ──

export function formatDate(value?: string): string {
    if (!value) return "–";
    const dt = new Date(value);
    if (isNaN(dt.getTime())) return "–";
    return dt.toLocaleString();
}

export function firstPromptLine(prompt?: string): string {
    const text = (prompt ?? "").trim();
    if (!text) return "Untitled session";
    const first =
        text.split(/\r?\n/).find((l) => l.trim() !== "") ?? text;
    const trimmed = first.trim();
    return trimmed.length > 110 ? trimmed.slice(0, 110) + "…" : trimmed;
}

export function relativeTime(value?: string): string {
    if (!value) return "";
    const dt = new Date(value);
    if (isNaN(dt.getTime())) return "";
    const diff = Date.now() - dt.getTime();
    const seconds = Math.floor(diff / 1000);
    if (seconds < 60) return "just now";
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    return `${days}d ago`;
}

// Aliases used by components
export const formatRelativeTime = relativeTime;
export const truncatePrompt = firstPromptLine;
