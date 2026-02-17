<script lang="ts">
    import type { SessionSummary } from "$lib/types";
    import { appState } from "$lib/stores.svelte";
    import { formatRelativeTime, truncatePrompt } from "$lib/utils";
    import { MessageSquare, GitPullRequest } from "@lucide/svelte";

    let { session }: { session: SessionSummary } = $props();

    const isActive = $derived(appState.selectedSessionID === session.id);
    const isBusy = $derived(session.busy);
    const prompt = $derived(session.latest_run?.prompt ?? session.id);
    const age = $derived(formatRelativeTime(session.updated_at));

    function select() {
        appState.selectSession(session.id);
    }
</script>

<button
    class="session-item"
    class:active={isActive}
    class:busy={isBusy}
    onclick={select}
>
    <div class="status-indicator"></div>
    <div class="item-icon">
        {#if session.pr_url}
            <GitPullRequest size={14} class="text-success-soft" />
        {:else}
            <MessageSquare size={14} opacity={isActive ? 1 : 0.6} />
        {/if}
    </div>
    <div class="item-content">
        <span class="item-title">{truncatePrompt(prompt)}</span>
        <span class="item-age">{age}</span>
    </div>
    {#if isBusy}
        <div class="busy-indicator">
            <span class="busy-dot"></span>
            <span class="busy-ping"></span>
        </div>
    {/if}
</button>

<style>
    .session-item {
        position: relative;
        display: flex;
        align-items: center;
        width: 100%;
        padding: 12px 14px;
        border-radius: var(--radius-md);
        background: transparent;
        color: var(--color-text-secondary);
        font-family: inherit;
        cursor: pointer;
        text-align: left;
        transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
        gap: 12px;
        border: 1px solid transparent;
        margin-bottom: 2px;
    }

    .session-item:hover:not(.active) {
        background: var(--color-bg-hover);
        color: var(--color-text);
        padding-left: 18px;
    }

    .session-item.active {
        background: var(--color-bg-active);
        color: var(--color-text);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
        border: 1px solid var(--color-border-strong);
    }

    .status-indicator {
        position: absolute;
        left: 0;
        top: 12px;
        bottom: 12px;
        width: 3px;
        border-radius: 0 4px 4px 0;
        background: var(--color-accent);
        opacity: 0;
        transition: all 0.2s;
        transform: scaleY(0.5);
    }

    .session-item.active .status-indicator {
        opacity: 1;
        transform: scaleY(1);
    }

    .item-icon {
        flex-shrink: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 20px;
        height: 20px;
        color: var(--color-accent);
    }

    :global(.text-success-soft) {
        color: var(--color-success);
        opacity: 0.8;
    }

    .item-content {
        display: flex;
        flex-direction: column;
        flex: 1;
        min-width: 0;
    }

    .item-title {
        font-size: 13px;
        font-weight: 600;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        line-height: 1.2;
        margin-bottom: 2px;
    }

    .item-age {
        font-size: 11px;
        color: var(--color-text-muted);
        font-weight: 500;
    }

    /* Busy State */
    .busy-indicator {
        position: relative;
        width: 10px;
        height: 10px;
        flex-shrink: 0;
    }

    .busy-dot {
        position: absolute;
        width: 6px;
        height: 6px;
        background: var(--color-accent);
        border-radius: 50%;
        top: 2px;
        left: 2px;
        z-index: 2;
    }

    .busy-ping {
        position: absolute;
        width: 100%;
        height: 100%;
        border-radius: 50%;
        background: var(--color-accent);
        opacity: 0.6;
        animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
    }

    @keyframes ping {
        75%,
        100% {
            transform: scale(2.5);
            opacity: 0;
        }
    }
</style>
