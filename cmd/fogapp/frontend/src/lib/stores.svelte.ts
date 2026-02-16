// ── Global reactive state using Svelte 5 runes ──

import { ACTIVE_STATES } from "./types";
import type {
    DiffResult,
    Repo,
    RunEvent,
    RunSummary,
    SessionSummary,
    Settings,
} from "./types";
import {
    fetchDiff,
    fetchRepos,
    fetchRunEvents,
    fetchSessionDetail,
    fetchSessions,
    fetchSettings,
    openRunStream,
    resolveAPIBaseURL,
    resolveAPIToken,
    resolveVersion,
    setAPIBaseURL,
    setAPIToken,
} from "./api";

// ── App state ──

export type ViewName = "new" | "detail" | "settings";

class AppState {
    // Connection
    daemonStatus = $state<"connecting" | "connected" | "unavailable">("connecting");
    version = $state("–");

    // Data
    settings = $state<Settings | null>(null);
    repos = $state<Repo[]>([]);
    sessions = $state<SessionSummary[]>([]);

    // UI state
    currentView = $state<ViewName>("new");
    selectedSessionID = $state("");
    selectedRunID = $state("");
    selectedTab = $state<"timeline" | "diff" | "logs" | "stats">("timeline");
    autoFollowLatest = $state(true);

    // New UI state
    sessionMode = $state<"plan" | "build">("build");
    selectedBranch = $state("");
    chatExpanded = $state(false);

    // Detail data
    detailSession = $state<SessionSummary | null>(null);
    detailRuns = $state<RunSummary[]>([]);
    detailEvents = $state<RunEvent[]>([]);
    detailDiff = $state<DiffResult | null>(null);
    detailDiffError = $state("");

    // SSE stream
    private streamSource: EventSource | null = null;
    private streamSessionID = "";
    private streamRunID = "";

    // Polling
    private pollInterval: ReturnType<typeof setInterval> | null = null;

    // ── Derived ──

    get runningSessions(): SessionSummary[] {
        return this.sessions.filter((s) => this.isSessionRunning(s));
    }

    get completedSessions(): SessionSummary[] {
        return this.sessions.filter((s) => !this.isSessionRunning(s));
    }

    get selectedRun(): RunSummary | null {
        const found = this.detailRuns.find((r) => r.id === this.selectedRunID);
        return found ?? this.detailRuns[0] ?? null;
    }

    get latestRun(): RunSummary | null {
        return this.detailRuns[0] ?? null;
    }

    get canStop(): boolean {
        const session = this.detailSession;
        const latest = this.latestRun;
        return !!(session?.busy && latest && ACTIVE_STATES[latest.state]);
    }

    // ── Methods ──

    isSessionRunning(session: SessionSummary): boolean {
        if (session.busy) return true;
        const latest = session.latest_run;
        const runState = latest?.state ?? session.status;
        return !!ACTIVE_STATES[runState];
    }

    setView(view: ViewName): void {
        this.currentView = view;
        if (view !== "detail") {
            this.closeStream();
        }
    }

    async bootstrap(): Promise<void> {
        this.daemonStatus = "connecting";
        try {
            const baseURL = await resolveAPIBaseURL();
            setAPIBaseURL(baseURL);
            const token = await resolveAPIToken();
            setAPIToken(token);
            this.version = await resolveVersion();

            await this.refreshAll();
            this.daemonStatus = "connected";
            this.startPolling();
        } catch {
            this.daemonStatus = "unavailable";
            throw new Error("Initialization failed");
        }
    }

    async refreshAll(): Promise<void> {
        const [settings, repos, sessions] = await Promise.all([
            fetchSettings(),
            fetchRepos(),
            fetchSessions(),
        ]);
        this.settings = settings;
        this.repos = repos;
        this.sessions = sessions;

        if (this.selectedSessionID) {
            await this.loadDetail();
        }
    }

    async refreshSessions(): Promise<void> {
        this.sessions = await fetchSessions();
        if (this.selectedSessionID) {
            const found = this.sessions.some(
                (s) => s.id === this.selectedSessionID,
            );
            if (!found) {
                this.closeStream();
                this.selectedSessionID = "";
                this.selectedRunID = "";
            }
        }
    }

    async refreshRepos(): Promise<void> {
        this.repos = await fetchRepos();
    }

    async selectSession(
        sessionID: string,
        followLatest = true,
    ): Promise<void> {
        this.selectedSessionID = sessionID;
        if (followLatest) this.autoFollowLatest = true;
        this.setView("detail");
        await this.loadDetail();
    }

    async loadDetail(): Promise<void> {
        if (!this.selectedSessionID) {
            this.closeStream();
            this.detailSession = null;
            this.detailRuns = [];
            this.detailEvents = [];
            this.detailDiff = null;
            this.detailDiffError = "";
            return;
        }

        const detail = await fetchSessionDetail(this.selectedSessionID);
        this.detailSession = detail?.session ?? null;
        this.detailRuns = detail?.runs ?? [];

        if (!this.detailSession) {
            this.closeStream();
            this.detailEvents = [];
            this.detailDiff = null;
            this.detailDiffError = "Session not found.";
            return;
        }

        // pick run
        if (this.autoFollowLatest || !this.selectedRunID) {
            this.selectedRunID = this.detailRuns[0]?.id ?? "";
        } else {
            const exists = this.detailRuns.some(
                (r) => r.id === this.selectedRunID,
            );
            if (!exists) {
                this.selectedRunID = this.detailRuns[0]?.id ?? "";
            }
        }

        // load events
        if (this.selectedRunID) {
            this.detailEvents = await fetchRunEvents(
                this.selectedSessionID,
                this.selectedRunID,
            );
        } else {
            this.detailEvents = [];
        }

        // load diff
        try {
            this.detailDiff = await fetchDiff(this.selectedSessionID);
            this.detailDiffError = "";
        } catch (err) {
            this.detailDiff = null;
            this.detailDiffError =
                err instanceof Error ? err.message : "Failed to load diff";
        }

        this.openStream();
    }

    // ── SSE stream ──

    private closeStream(): void {
        if (this.streamSource) {
            this.streamSource.close();
            this.streamSource = null;
        }
        this.streamSessionID = "";
        this.streamRunID = "";
    }

    private openStream(): void {
        const run = this.selectedRun;
        if (!run || !this.selectedSessionID || !this.selectedRunID) {
            this.closeStream();
            return;
        }
        if (!ACTIVE_STATES[run.state]) {
            this.closeStream();
            return;
        }
        if (
            this.streamSource &&
            this.streamSessionID === this.selectedSessionID &&
            this.streamRunID === this.selectedRunID
        ) {
            return;
        }

        this.closeStream();

        const cursor = this.latestEventID();
        this.streamSource = openRunStream(
            this.selectedSessionID,
            this.selectedRunID,
            cursor,
            (event) => {
                this.appendEvent(event);
            },
            () => {
                this.closeStream();
                this.refreshSessions().catch(() => { });
            },
            () => {
                this.closeStream();
            },
        );
        this.streamSessionID = this.selectedSessionID;
        this.streamRunID = this.selectedRunID;
    }

    private latestEventID(): number {
        if (!this.detailEvents.length) return 0;
        const last = this.detailEvents[this.detailEvents.length - 1];
        return last?.id ?? 0;
    }

    private appendEvent(event: RunEvent): void {
        if (!event?.id) return;
        const exists = this.detailEvents.some((e) => e.id === event.id);
        if (exists) return;
        this.detailEvents = [...this.detailEvents, event].sort(
            (a, b) => a.id - b.id,
        );
    }

    // ── Polling ──

    private startPolling(): void {
        this.pollInterval = setInterval(() => {
            this.refreshSessions()
                .then(() => {
                    if (
                        this.selectedSessionID &&
                        this.currentView === "detail"
                    ) {
                        return this.loadDetail();
                    }
                    return undefined;
                })
                .catch(() => {
                    this.daemonStatus = "unavailable";
                });
        }, 4000);
    }

    destroy(): void {
        if (this.pollInterval) {
            clearInterval(this.pollInterval);
        }
        this.closeStream();
    }
}

export const appState = new AppState();
