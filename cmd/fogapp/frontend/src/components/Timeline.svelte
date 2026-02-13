<script lang="ts">
    import { appState } from "$lib/stores.svelte";
    import { formatRelativeTime } from "$lib/utils";
    import { fade, slide } from "svelte/transition";
    import {
        Play,
        CheckCircle2,
        AlertCircle,
        Clock,
        ChevronDown,
        Hash,
        History,
        Terminal as TerminalIcon,
    } from "@lucide/svelte";

    const runs = $derived(appState.detailRuns ?? []);

    function isTerminal(state: string) {
        switch (state.trim()) {
            case "COMPLETED":
            case "FAILED":
            case "CANCELLED":
                return true;
            default:
                return false;
        }
    }

    function isInProgress(state: string) {
        return !isTerminal(state);
    }

    function getStatusIcon(state: string) {
        if (state === "COMPLETED") return CheckCircle2;
        if (state === "FAILED" || state === "CANCELLED") return AlertCircle;
        if (isInProgress(state)) return Play;
        return History;
    }
</script>

<div class="timeline-v2">
    {#if runs.length === 0}
        <div class="empty-feed" in:fade>
            <div class="zen-orb mini">
                <Play size={20} />
            </div>
            <p>Intelligence engine idle. No runs recorded.</p>
        </div>
    {:else}
        <div class="runs-list">
            {#each runs as run, i}
                {@const Icon = getStatusIcon(run.state)}
                <div
                    class="run-entry"
                    class:active={run.id === appState.selectedRunID}
                    style="animation-delay: {i * 100}ms"
                >
                    <div class="run-side">
                        <div
                            class="status-marker"
                            class:busy={isInProgress(run.state)}
                        >
                            <Icon size={14} />
                        </div>
                        {#if i < runs.length - 1}
                            <div class="connector"></div>
                        {/if}
                    </div>

                    <div class="run-content glass card">
                        <header class="run-header">
                            <div class="run-meta">
                                <span class="run-number"
                                    ><Hash size={10} />{runs.length - i}</span
                                >
                                <span class="run-time"
                                    ><Clock size={10} />{formatRelativeTime(
                                        run.created_at,
                                    )}</span
                                >
                            </div>
                            <div
                                class="run-status"
                                class:running={isInProgress(run.state)}
                                class:error={run.state === "FAILED"}
                                class:done={run.state === "COMPLETED"}
                                class:cancelled={run.state === "CANCELLED"}
                            >
                                <span class="status-indicator-dot"></span>
                                {run.state.replace("AI_", "")}
                            </div>
                        </header>

                        <div class="run-instructions">
                            {run.prompt}
                        </div>

                        {#if run.id === appState.selectedRunID && appState.detailEvents.length > 0}
                            <div class="run-terminal" transition:slide>
                                <div class="terminal-bar">
                                    <TerminalIcon size={12} />
                                    <span>Activity Intelligence</span>
                                </div>
                                <div class="terminal-output">
                                    {#each appState.detailEvents as evt}
                                        <div class="log-line">
                                            <span class="log-ts"
                                                >{new Date(
                                                    evt.ts,
                                                ).toLocaleTimeString([], {
                                                    hour: "2-digit",
                                                    minute: "2-digit",
                                                    second: "2-digit",
                                                })}</span
                                            >
                                            <span class="log-type"
                                                >[{evt.type}]</span
                                            >
                                            <span class="log-msg">{evt.message}</span>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {:else if run.state === "DONE"}
                            <button
                                class="expand-btn"
                                onclick={() =>
                                    (appState.selectedRunID = run.id)}
                            >
                                <ChevronDown size={14} />
                                <span>Inspect Pipeline</span>
                            </button>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    .timeline-v2 {
        padding: 0 0 80px;
    }

    .empty-feed {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 80px 40px;
        color: var(--color-text-muted);
        text-align: center;
        gap: 20px;
        background: rgba(255, 255, 255, 0.02);
        border-radius: 24px;
        border: 1px dashed var(--color-border);
    }

    .zen-orb.mini {
        width: 48px;
        height: 48px;
        background: var(--color-bg-active);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-accent);
        box-shadow: 0 0 20px rgba(59, 130, 246, 0.1);
        animation: breath 3s ease-in-out infinite;
    }

    .runs-list {
        display: flex;
        flex-direction: column;
        gap: 0;
    }

    .run-entry {
        display: flex;
        gap: 24px;
        animation: slideIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) backwards;
    }

    .run-side {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 32px;
        flex-shrink: 0;
    }

    .status-marker {
        width: 32px;
        height: 32px;
        border-radius: 50%;
        background: var(--color-bg-elevated);
        border: 1px solid var(--color-border);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 2;
        color: var(--color-text-muted);
        transition: all 0.3s;
    }

    .active .status-marker {
        border-color: var(--color-accent);
        color: var(--color-accent);
        box-shadow: 0 0 12px rgba(59, 130, 246, 0.3);
        transform: scale(1.1);
    }

    .status-marker.busy {
        background: rgba(59, 130, 246, 0.1);
        color: var(--color-accent);
        border-color: var(--color-accent);
    }

    .connector {
        width: 2px;
        flex: 1;
        background: linear-gradient(
            to bottom,
            var(--color-border),
            rgba(255, 255, 255, 0.05)
        );
    }

    .run-content {
        flex: 1;
        padding: 20px 24px;
        margin-bottom: 32px;
        border-radius: 18px;
        transition: all 0.3s;
        border: 1px solid var(--color-border);
    }

    .run-entry.active .run-content {
        border-color: var(--color-accent);
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
        background: rgba(255, 255, 255, 0.03);
    }

    .run-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 16px;
    }

    .run-meta {
        display: flex;
        align-items: center;
        gap: 16px;
    }

    .run-number {
        font-family: var(--font-mono);
        font-size: 11px;
        font-weight: 700;
        color: var(--color-accent);
        display: flex;
        align-items: center;
        gap: 4px;
        opacity: 0.8;
    }

    .run-time {
        font-size: 11px;
        color: var(--color-text-muted);
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .run-status {
        font-size: 10px;
        font-weight: 800;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 4px 10px;
        background: rgba(255, 255, 255, 0.03);
        border-radius: 20px;
        color: var(--color-text-secondary);
        border: 1px solid var(--color-border);
    }

    .status-indicator-dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: currentColor;
    }

    .run-status.running {
        color: #60a5fa;
        border-color: rgba(96, 165, 250, 0.2);
    }
    .run-status.done {
        color: #34d399;
        border-color: rgba(52, 211, 153, 0.2);
    }
    .run-status.error {
        color: #f87171;
        border-color: rgba(248, 113, 113, 0.2);
    }
    .run-status.cancelled {
        color: #fbbf24;
        border-color: rgba(251, 191, 36, 0.2);
    }

    .run-instructions {
        font-size: 14px;
        line-height: 1.6;
        color: var(--color-text);
        font-weight: 500;
        word-break: break-word;
    }

    .run-terminal {
        margin-top: 20px;
        background: #05070a;
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.05);
        overflow: hidden;
        box-shadow: inset 0 2px 10px rgba(0, 0, 0, 0.5);
    }

    .terminal-bar {
        background: rgba(255, 255, 255, 0.03);
        padding: 8px 14px;
        display: flex;
        align-items: center;
        gap: 10px;
        font-size: 11px;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--color-text-muted);
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    }

    .terminal-output {
        padding: 16px;
        display: flex;
        flex-direction: column;
        gap: 6px;
        font-family: var(--font-mono);
        font-size: 11px;
        max-height: 400px;
        overflow-y: auto;
    }

    .log-line {
        display: flex;
        gap: 14px;
    }

    .log-ts {
        color: var(--color-text-muted);
        opacity: 0.4;
        flex-shrink: 0;
    }

    .log-type {
        color: var(--color-accent);
        opacity: 0.8;
        flex-shrink: 0;
    }

    .log-msg {
        color: var(--color-text-secondary);
        word-break: break-all;
    }

    .expand-btn {
        margin-top: 16px;
        display: flex;
        align-items: center;
        gap: 8px;
        background: none;
        border: none;
        color: var(--color-text-muted);
        font-size: 11px;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        cursor: pointer;
        padding: 4px 0;
        transition: color 0.2s;
    }

    .expand-btn:hover {
        color: var(--color-accent);
    }

    @keyframes breath {
        0%,
        100% {
            transform: scale(1);
            opacity: 0.8;
        }
        50% {
            transform: scale(1.05);
            opacity: 1;
        }
    }

    @keyframes slideIn {
        from {
            opacity: 0;
            transform: translateY(20px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
