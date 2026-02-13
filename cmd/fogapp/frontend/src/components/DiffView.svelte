<script lang="ts">
    import { appState } from "$lib/stores.svelte";
    import { fade } from "svelte/transition";
    import {
        FileDiff,
        GitBranch,
        ArrowRight,
        Code2,
        FileText,
        Info,
        CheckCircle2,
    } from "@lucide/svelte";

    const diff = $derived(appState.detailDiff);
    const diffError = $derived(appState.detailDiffError);
</script>

<div class="diff-module" in:fade>
    {#if diffError}
        <div class="empty-diff">
            <div class="zen-orb mini">
                <Info size={20} />
            </div>
            <p>Unable to load diff.</p>
            <p class="error-detail">{diffError}</p>
        </div>
    {:else if !diff}
        <div class="empty-diff">
            <div class="zen-orb mini">
                <Info size={20} />
            </div>
            <p>No changes detected in this run cycle.</p>
        </div>
    {:else}
        <div class="diff-display glass card">
            <header class="diff-header">
                <div class="diff-title">
                    <div class="icon-ring">
                        <FileDiff size={14} />
                    </div>
                    <span>Proposed Intelligence Changes</span>
                </div>
                <div class="branch-flow">
                    <div class="branch-pill target">
                        <GitBranch size={12} />
                        <span>{diff.base_branch}</span>
                    </div>
                    <ArrowRight size={14} class="flow-arrow" />
                    <div class="branch-pill source">
                        <GitBranch size={12} />
                        <span>{diff.branch}</span>
                    </div>
                </div>
            </header>

            {#if diff.stat}
                <div class="diff-stat-pane">
                    <div class="stat-inner">
                        <div class="stat-icon"><FileText size={12} /></div>
                        <span class="stat-text">{diff.stat}</span>
                    </div>
                    <div class="approval-tag">
                        <CheckCircle2 size={12} />
                        <span>Review Pending</span>
                    </div>
                </div>
            {/if}

            <div class="diff-content">
                {#if diff.patch}
                    <div class="patch-viewer">
                        <div class="line-numbers-glow"></div>
                        <pre class="patch-code">
                            {#each diff.patch.split("\n") as line}
                                <div
                                    class="patch-line"
                                    class:add={line.startsWith("+")}
                                    class:sub={line.startsWith("-")}
                                    class:meta={line.startsWith("@@")}>
                                    <span class="line-content">{line}</span>
                                </div>
                            {/each}
                        </pre>
                    </div>
                {:else}
                    <div class="no-patch">
                        <Code2 size={24} />
                        <p>
                            No binary or text difference available to display.
                        </p>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .diff-module {
        height: 100%;
        display: flex;
        flex-direction: column;
    }

    .empty-diff {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 60px 20px;
        color: var(--color-text-muted);
        text-align: center;
        gap: 16px;
    }

    .error-detail {
        max-width: 760px;
        width: 100%;
        font-family: var(--font-mono);
        font-size: 12px;
        color: var(--color-text-secondary);
        opacity: 0.8;
        text-align: left;
        white-space: pre-wrap;
        word-break: break-word;
        background: rgba(255, 255, 255, 0.02);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: 12px 14px;
    }

    .zen-orb.mini {
        width: 44px;
        height: 44px;
        background: var(--color-bg-active);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-text-muted);
        border: 1px solid var(--color-border);
    }

    .diff-display {
        display: flex;
        flex-direction: column;
        height: 100%;
        overflow: hidden;
        border-radius: 20px;
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
    }

    .diff-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 20px;
        border-bottom: 1px solid var(--color-border);
        background: rgba(255, 255, 255, 0.02);
    }

    .diff-title {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 13px;
        font-weight: 700;
        letter-spacing: -0.01em;
    }

    .icon-ring {
        width: 28px;
        height: 28px;
        border-radius: 8px;
        background: var(--color-bg-active);
        border: 1px solid var(--color-border);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-accent);
    }

    .branch-flow {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .branch-pill {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 4px 10px;
        border-radius: 6px;
        font-family: var(--font-mono);
        font-size: 11px;
        font-weight: 600;
        border: 1px solid var(--color-border);
    }

    .branch-pill.target {
        background: rgba(59, 130, 246, 0.05);
        color: var(--color-text-secondary);
    }
    .branch-pill.source {
        background: rgba(59, 130, 246, 0.1);
        color: var(--color-accent);
        border-color: rgba(59, 130, 246, 0.3);
    }

    .diff-stat-pane {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 10px 20px;
        background: rgba(0, 0, 0, 0.2);
        border-bottom: 1px solid var(--color-border);
    }

    .stat-inner {
        display: flex;
        align-items: center;
        gap: 10px;
    }

    .stat-icon {
        color: var(--color-text-muted);
    }
    .stat-text {
        font-size: 12px;
        font-weight: 500;
        color: var(--color-text-secondary);
    }

    .approval-tag {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 10px;
        font-weight: 800;
        text-transform: uppercase;
        color: var(--color-text-muted);
        opacity: 0.6;
    }

    .diff-content {
        flex: 1;
        overflow: hidden;
        background: #05070a;
        position: relative;
    }

    .patch-viewer {
        height: 100%;
        overflow: auto;
        padding: 20px 0;
    }

    .patch-code {
        margin: 0;
        font-family: var(--font-mono);
        font-size: 13px;
        line-height: 1.6;
    }

    .patch-line {
        padding: 0 24px;
        white-space: pre;
        transition: background 0.2s;
    }

    .patch-line:hover {
        background: rgba(255, 255, 255, 0.03);
    }

    .patch-line.add {
        background: rgba(52, 211, 153, 0.08);
        color: #6ee7b7;
        border-left: 3px solid #34d399;
    }
    .patch-line.sub {
        background: rgba(248, 113, 113, 0.08);
        color: #fca5a5;
        border-left: 3px solid #f87171;
    }
    .patch-line.meta {
        background: rgba(59, 130, 246, 0.05);
        color: var(--color-text-muted);
        font-weight: 700;
        opacity: 0.8;
    }

    .no-patch {
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: var(--color-text-muted);
        gap: 16px;
    }

    @keyframes breath {
        0%,
        100% {
            opacity: 0.8;
            transform: scale(1);
        }
        50% {
            opacity: 1;
            transform: scale(1.05);
        }
    }
</style>
