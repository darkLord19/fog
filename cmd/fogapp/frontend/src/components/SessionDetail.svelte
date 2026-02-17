<script lang="ts">
    import { toast } from "svelte-sonner";
    import { appState } from "$lib/stores.svelte";
    import { firstPromptLine } from "$lib/utils";
    import {
        followUp,
        cancelSession,
        openInEditor,
        forkSession,
    } from "$lib/api";
    import Timeline from "./Timeline.svelte";
    import DiffView from "./DiffView.svelte";
    import LogsView from "./LogsView.svelte";
    import StatsView from "./StatsView.svelte";
    import { fade, slide } from "svelte/transition";
    import {
        RefreshCcw,
        GitFork,
        Code,
        StopCircle,
        ArrowUp,
        Terminal,
        FileCode,
        Activity,
        History,
        ExternalLink,
        Search,
        Layers,
        GitPullRequest,
        X,
    } from "@lucide/svelte";

    let followupPrompt = $state("");
    let submitting = $state(false);
    let activeTab = $state("timeline");

    const session = $derived(appState.detailSession);
    const runs = $derived(appState.detailRuns ?? []);
    const latestRun = $derived(runs[0]);
    const isBusy = $derived(session?.busy ?? false);
    const titleText = $derived(firstPromptLine(latestRun?.prompt));

    const tabs = [
        { id: "timeline", label: "Timeline", icon: History },
        { id: "diff", label: "Diff", icon: FileCode },
        { id: "logs", label: "Logs", icon: Terminal },
        { id: "stats", label: "Stats", icon: Activity },
    ];

    async function handleFollowup() {
        if (!followupPrompt.trim() || !session) return;
        submitting = true;
        try {
            const out = await followUp(session.id, followupPrompt.trim());
            toast.success(`Queued run ${out.run_id}`);
            followupPrompt = "";
            await appState.loadDetail();
        } catch (err) {
            toast.error(
                "Follow-up failed: " +
                    (err instanceof Error ? err.message : "Error"),
            );
        } finally {
            submitting = false;
        }
    }

    async function handleRerun() {
        if (!session || !latestRun) return;
        try {
            const out = await followUp(session.id, latestRun.prompt);
            toast.success(`Queued re-run ${out.run_id}`);
            await appState.loadDetail();
        } catch (err) {
            toast.error(
                "Re-run failed: " +
                    (err instanceof Error ? err.message : "Error"),
            );
        }
    }

    async function handleStop() {
        if (!session) return;
        try {
            await cancelSession(session.id);
            toast.success("Cancel requested");
            await appState.loadDetail();
        } catch (err) {
            toast.error(
                "Termination failed: " +
                    (err instanceof Error ? err.message : "Error"),
            );
        }
    }

    async function handleOpen() {
        if (!session) return;
        try {
            const res = await openInEditor(session.id);
            toast.success("Opened in " + (res.editor || "editor"));
        } catch (err) {
            toast.error(
                "Failed to open editor: " +
                    (err instanceof Error ? err.message : "Error"),
            );
        }
    }

    async function handleFork() {
        if (!session || !latestRun) return;
        try {
            const out = await forkSession(session.id, latestRun.prompt);
            toast.success(`Queued fork session ${out.session_id}`);
            await appState.refreshSessions();
        } catch (err) {
            toast.error(
                "Fork failed: " +
                    (err instanceof Error ? err.message : "Error"),
            );
        }
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleFollowup();
        }
    }

    function openPR() {
        if (session?.pr_url) {
            window.runtime.BrowserOpenURL(session.pr_url);
        }
    }
</script>

{#if session}
    <div class="session-view" in:fade={{ duration: 300 }}>
        <!-- Sophisticated Header -->
        <header class="view-header glass">
            <div class="header-top">
                <div class="header-left">
                    <div class="breadcrumb">
                        <Layers size={14} class="crumb-icon" />
                        <span class="repo-tag">{session.repo_name}</span>
                        <span class="sep">/</span>
                        <span class="session-id"
                            >{session.id.substring(0, 8)}</span
                        >
                    </div>
                    <div id="detail-title" class="session-title">
                        {titleText}
                    </div>
                </div>

                <div class="header-actions">
                    <div class="status-pill-v2" class:busy={isBusy}>
                        <div class="indicator">
                            {#if isBusy}<div class="pulse"></div>{/if}
                        </div>
                        <span class="status-text">{session.status}</span>
                    </div>

                    <div class="action-group">
                        {#if isBusy}
                            <button
                                id="detail-stop"
                                class="action-btn danger"
                                onclick={handleStop}
                                title="Stop Session"
                            >
                                <StopCircle size={16} />
                                <span>Terminate</span>
                            </button>
                        {:else}
                            <button
                                id="detail-rerun"
                                class="action-btn secondary"
                                onclick={handleRerun}
                                title="Re-run Latest"
                            >
                                <RefreshCcw size={16} />
                                <span>Re-run</span>
                            </button>
                        {/if}

                        <button
                            id="detail-fork"
                            class="action-btn secondary"
                            onclick={handleFork}
                            title="Fork Session"
                            disabled={isBusy}
                        >
                            <GitFork size={16} />
                            <span>Fork</span>
                        </button>

                        {#if session.pr_url}
                            <button
                                class="action-btn success-soft"
                                onclick={openPR}
                                title="View Pull Request"
                            >
                                <GitPullRequest size={16} />
                                <span>Open PR</span>
                                <ExternalLink
                                    size={12}
                                    class="btn-extra-icon"
                                />
                            </button>
                        {/if}

                        <button
                            id="detail-open"
                            class="action-btn primary"
                            onclick={handleOpen}
                            title="Open in IDE"
                        >
                            <Code size={16} />
                            <span>Open in IDE</span>
                            <ExternalLink size={12} class="btn-extra-icon" />
                        </button>
                    </div>
                </div>
            </div>

            <nav class="view-nav">
                {#each tabs as tab}
                    <button
                        class="nav-tab"
                        class:active={activeTab === tab.id}
                        onclick={() => (activeTab = tab.id)}
                    >
                        <tab.icon size={14} />
                        <span>{tab.label}</span>
                        {#if activeTab === tab.id}
                            <div class="tab-indicator" in:fade></div>
                        {/if}
                    </button>
                {/each}
            </nav>
        </header>

        <!-- Dynamic Content Engine -->
        <main class="view-content" class:with-followup={!isBusy}>
            <div class="content-wrapper">
                {#if activeTab === "timeline"}
                    <Timeline />
                {:else if activeTab === "diff"}
                    <DiffView />
                {:else if activeTab === "logs"}
                    <LogsView />
                {:else if activeTab === "stats"}
                    <StatsView />
                {/if}
            </div>
        </main>

        <!-- Command Input / Follow-up -->
        {#if !isBusy}
            <div class="command-center" transition:slide={{ axis: "y" }}>
                <div class="command-orb glass">
                    <div class="input-container">
                        <textarea
                            id="followup-prompt"
                            bind:value={followupPrompt}
                            onkeydown={handleKeydown}
                            disabled={submitting}
                            placeholder="Instruct the engine to continue..."
                            rows="1"
                            class="command-input"
                        ></textarea>
                    </div>
                    <button
                        id="followup-submit"
                        class="dispatch-btn"
                        disabled={submitting || !followupPrompt.trim()}
                        onclick={handleFollowup}
                    >
                        {#if submitting}
                            <div class="loader-sm"></div>
                        {:else}
                            <ArrowUp size={18} />
                        {/if}
                    </button>
                </div>

                <div class="command-hints">
                    <span><b>Shift + Enter</b> for new line</span>
                    <span class="sep">•</span>
                    <span><b>Enter</b> to dispatch</span>
                </div>
            </div>
        {/if}
    </div>
{:else}
    <div class="empty-view" in:fade>
        <div class="zen-container">
            <div class="zen-orb">
                <Search size={32} />
            </div>
            <h2>Session Intelligence</h2>
            <p>Select a session context to begin analysis</p>
        </div>
    </div>
{/if}

<style>
    .session-view {
        display: flex;
        flex-direction: column;
        height: 100%;
        background: var(--color-bg);
        position: relative;
    }

    .view-header {
        padding: 16px 24px 0;
        z-index: 10;
        border-bottom: 1px solid var(--color-border);
    }

    .header-top {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        margin-bottom: 20px;
    }

    .header-left {
        display: flex;
        flex-direction: column;
        gap: 10px;
        min-width: 0;
    }

    .breadcrumb {
        display: flex;
        align-items: center;
        gap: 10px;
        background: rgba(255, 255, 255, 0.03);
        padding: 6px 12px;
        border-radius: 10px;
        border: 1px solid var(--color-border);
    }

    .repo-tag {
        font-weight: 700;
        font-size: 13px;
        letter-spacing: -0.01em;
    }

    .sep {
        opacity: 0.3;
        font-weight: 300;
    }

    .session-id {
        font-family: var(--font-mono);
        font-size: 12px;
        color: var(--color-text-secondary);
    }

    .session-title {
        font-size: 18px;
        font-weight: 800;
        letter-spacing: -0.02em;
        line-height: 1.2;
        color: var(--color-text);
        max-width: 720px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .header-actions {
        display: flex;
        align-items: center;
        gap: 16px;
    }

    .status-pill-v2 {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 6px 12px;
        background: var(--color-bg-active);
        border: 1px solid var(--color-border-strong);
        border-radius: 20px;
    }

    .status-pill-v2.busy {
        background: rgba(59, 130, 246, 0.1);
        border-color: rgba(59, 130, 246, 0.4);
        color: var(--color-accent);
    }

    .indicator {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: var(--color-text-muted);
        position: relative;
    }

    .busy .indicator {
        background: var(--color-accent);
    }

    .pulse {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        border-radius: 50%;
        background: var(--color-accent);
        animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
    }

    .status-text {
        font-size: 12px;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .action-group {
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .action-btn {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 0 14px;
        height: 36px;
        border-radius: 10px;
        font-size: 13px;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
        border: 1px solid transparent;
        background: none;
        color: var(--color-text-secondary);
    }

    .action-btn:hover {
        background: var(--color-bg-hover);
        color: var(--color-text);
    }

    .action-btn:disabled {
        opacity: 0.35;
        cursor: not-allowed;
        pointer-events: none;
    }

    .action-btn.primary {
        background: var(--color-accent);
        color: white;
    }

    .action-btn.primary:hover {
        background: var(--color-accent-hover);
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
    }

    .action-btn.danger {
        color: var(--color-danger);
    }

    .action-btn.danger:hover {
        background: rgba(239, 68, 68, 0.1);
        border-color: rgba(239, 68, 68, 0.2);
    }

    .action-btn.success-soft {
        background: rgba(16, 185, 129, 0.1);
        color: var(--color-success);
        border-color: rgba(16, 185, 129, 0.2);
    }

    .action-btn.success-soft:hover {
        background: rgba(16, 185, 129, 0.2);
        box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2);
    }

    .view-nav {
        display: flex;
        gap: 32px;
    }

    .nav-tab {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px 0;
        background: none;
        border: none;
        color: var(--color-text-muted);
        font-size: 13px;
        font-weight: 600;
        cursor: pointer;
        position: relative;
        transition: color 0.2s;
    }

    .nav-tab:hover {
        color: var(--color-text);
    }

    .nav-tab.active {
        color: var(--color-accent);
    }

    .tab-indicator {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: var(--color-accent);
        border-radius: 2px 2px 0 0;
        box-shadow: 0 -2px 8px rgba(59, 130, 246, 0.5);
    }

    .view-content {
        flex: 1;
        overflow-y: auto;
        padding: 24px;
        position: relative;
    }

    .view-content.with-followup {
        padding-bottom: 140px;
    }

    .content-wrapper {
        max-width: 1000px;
        margin: 0;
    }

    .command-center {
        position: absolute;
        bottom: 32px;
        left: 0;
        right: 0;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: 12px;
        padding: 0 40px;
        z-index: 20;
    }

    .command-orb {
        width: 100%;
        max-width: 760px;
        padding: 10px 10px 10px 24px;
        border-radius: 20px;
        display: flex;
        align-items: center;
        gap: 16px;
        border: 1px solid var(--color-border-strong);
        background: var(--color-bg-input);
        box-shadow:
            0 20px 50px rgba(0, 0, 0, 0.5),
            0 0 0 1px rgba(255, 255, 255, 0.05);
        transition: border-color 0.2s;
    }

    .command-orb:focus-within {
        border-color: var(--color-accent);
    }

    .input-container {
        flex: 1;
    }

    .command-input {
        width: 100%;
        background: transparent;
        border: none;
        resize: none;
        color: var(--color-text);
        font-size: 15px;
        font-family: inherit;
        line-height: 1.5;
        max-height: 150px;
        outline: none;
        display: block;
        padding: 8px 0;
    }

    .command-input::placeholder {
        color: var(--color-text-muted);
    }

    .dispatch-btn {
        width: 44px;
        height: 44px;
        border-radius: 12px;
        background: var(--color-accent);
        color: white;
        border: none;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
        flex-shrink: 0;
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
    }

    .dispatch-btn:hover:not(:disabled) {
        background: var(--color-accent-hover);
        transform: translateY(-2px);
        box-shadow: 0 8px 20px rgba(59, 130, 246, 0.5);
    }

    .dispatch-btn:disabled {
        opacity: 0.3;
        transform: none;
        background: var(--color-bg-active);
        box-shadow: none;
        cursor: not-allowed;
    }

    .command-hints {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 11px;
        color: var(--color-text-muted);
        text-transform: uppercase;
        letter-spacing: 0.05em;
        font-weight: 600;
        opacity: 0.6;
    }

    .loader-sm {
        width: 18px;
        height: 18px;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
    }

    .empty-view {
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: radial-gradient(
            circle at center,
            rgba(59, 130, 246, 0.03) 0%,
            transparent 70%
        );
    }

    .zen-container {
        text-align: center;
        animation: float 6s ease-in-out infinite;
    }

    .zen-orb {
        width: 80px;
        height: 80px;
        background: var(--color-bg-active);
        border: 1px solid var(--color-border);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 24px;
        color: var(--color-text-muted);
        box-shadow: 0 0 40px rgba(0, 0, 0, 0.2);
    }

    .zen-container h2 {
        font-size: 24px;
        font-weight: 800;
        margin-bottom: 8px;
        color: var(--color-text);
        letter-spacing: -0.02em;
    }

    .zen-container p {
        color: var(--color-text-muted);
        font-size: 15px;
    }

    @keyframes float {
        0%,
        100% {
            transform: translateY(0);
        }
        50% {
            transform: translateY(-10px);
        }
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    @keyframes ping {
        75%,
        100% {
            transform: scale(2.5);
            opacity: 0;
        }
    }
</style>
