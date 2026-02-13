<script lang="ts">
    import { appState } from "$lib/stores.svelte";
    import { fade } from "svelte/transition";
    import {
        Terminal,
        Search,
        Slash,
        Cpu,
        ShieldCheck,
        Monitor,
    } from "@lucide/svelte";

    const allEvents = $derived(appState.detailEvents ?? []);
    let filterText = $state("");
    const filterQuery = $derived(filterText.trim().toLowerCase());
    const filteredEvents = $derived(
        filterQuery
            ? allEvents.filter((evt) => {
                  const message = (evt.message ?? "").toLowerCase();
                  const type = (evt.type ?? "").toLowerCase();
                  return (
                      message.includes(filterQuery) ||
                      type.includes(filterQuery)
                  );
              })
            : allEvents,
    );
</script>

<div class="logs-module" in:fade>
    <div class="terminal-frame glass">
        <!-- Professional Window Bar -->
        <header class="terminal-titlebar">
            <div class="traffic-lights">
                <span class="dot close"></span>
                <span class="dot minimize"></span>
                <span class="dot maximize"></span>
            </div>

            <div class="terminal-tab">
                <Terminal size={12} class="tab-icon" />
                <span class="tab-label">Runtime Intelligence Console</span>
                <Slash size={10} class="tab-sep" />
                <span class="tab-sub">daemon-01</span>
            </div>

            <div class="terminal-actions">
                <div class="search-mini">
                    <Search size={12} />
                    <input
                        type="text"
                        bind:value={filterText}
                        placeholder="Filter output..."
                        aria-label="Filter log output"
                    />
                </div>
            </div>
        </header>

        <!-- Immersive Content Area -->
        <div class="terminal-viewport">
            {#if allEvents.length === 0}
                <div class="terminal-welcome">
                    <Monitor size={32} class="welcome-icon" />
                    <h3>Awaiting Logs</h3>
                    <p>Execute a run to stream runtime data.</p>
                </div>
            {:else}
                <div class="terminal-header-info">
                    <div class="info-pill">
                        <Cpu size={10} />
                        <span>PROCESS: FOG_DAEMON</span>
                    </div>
                    <div class="info-pill success">
                        <ShieldCheck size={10} />
                        <span>STATE: STABLE</span>
                    </div>
                </div>

                <div class="terminal-rows">
                    {#each filteredEvents as evt}
                        <div
                            class="terminal-line"
                            class:error={evt.type === "ERROR"}
                            class:warn={evt.type === "WARN"}
                        >
                            <span class="line-ts"
                                >[{new Date(evt.ts).toLocaleTimeString()}]</span
                            >
                            <span class="line-tag">[{evt.type}]</span>
                            <span class="line-content">{evt.message}</span>
                        </div>
                    {/each}
                    <div class="terminal-cursor"></div>
                </div>
            {/if}
        </div>

        <footer class="terminal-bottom-bar">
            {#if filterQuery}
                <span>{filteredEvents.length} / {allEvents.length} events</span>
            {:else}
                <span>{allEvents.length} events recorded</span>
            {/if}
            <div class="terminal-meta">
                <span>UTF-8</span>
                <span class="sep">|</span>
                <span>Zsh</span>
            </div>
        </footer>
    </div>
</div>

<style>
    .logs-module {
        height: 100%;
        display: flex;
        flex-direction: column;
    }

    .terminal-frame {
        display: flex;
        flex-direction: column;
        background: #05070a;
        border: 1px solid var(--color-border-strong);
        border-radius: 20px;
        overflow: hidden;
        height: 100%;
        box-shadow: 0 30px 60px rgba(0, 0, 0, 0.6);
    }

    .terminal-titlebar {
        display: flex;
        align-items: center;
        padding: 12px 16px;
        background: rgba(255, 255, 255, 0.03);
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        gap: 24px;
    }

    .traffic-lights {
        display: flex;
        gap: 8px;
    }

    .dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.1);
    }

    .dot.close {
        background: #ff5f56;
        box-shadow: 0 0 6px rgba(255, 95, 86, 0.3);
    }
    .dot.minimize {
        background: #ffbd2e;
    }
    .dot.maximize {
        background: #27c93f;
    }

    .terminal-tab {
        display: flex;
        align-items: center;
        gap: 8px;
        background: rgba(255, 255, 255, 0.05);
        padding: 6px 14px;
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.1);
    }

    .tab-label {
        font-size: 11px;
        font-weight: 700;
        color: var(--color-text);
    }

    .tab-sub {
        font-size: 10px;
        font-weight: 500;
        color: var(--color-text-muted);
    }

    .terminal-actions {
        flex: 1;
        display: flex;
        justify-content: flex-end;
    }

    .search-mini {
        display: flex;
        align-items: center;
        gap: 8px;
        background: rgba(0, 0, 0, 0.3);
        padding: 4px 12px;
        border-radius: 6px;
        border: 1px solid rgba(255, 255, 255, 0.05);
        width: 180px;
        color: var(--color-text-muted);
    }

    .search-mini input {
        background: none;
        border: none;
        outline: none;
        color: var(--color-text);
        font-size: 11px;
        width: 100%;
    }

    .terminal-viewport {
        flex: 1;
        overflow-y: auto;
        padding: 20px;
        font-family: var(--font-mono);
        mask-image: linear-gradient(
            to bottom,
            transparent,
            black 40px,
            black calc(100% - 40px),
            transparent
        );
    }

    .terminal-welcome {
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: var(--color-text-muted);
        text-align: center;
        gap: 16px;
    }

    .terminal-welcome h3 {
        font-size: 16px;
        font-weight: 700;
        color: var(--color-text);
        margin: 0;
    }
    .terminal-welcome p {
        font-size: 13px;
        margin: 0;
    }

    .terminal-header-info {
        display: flex;
        gap: 12px;
        margin-bottom: 24px;
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        padding-bottom: 16px;
    }

    .info-pill {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 10px;
        font-weight: 800;
        padding: 4px 8px;
        background: rgba(255, 255, 255, 0.05);
        border-radius: 4px;
        color: var(--color-text-muted);
    }

    .info-pill.success {
        color: #34d399;
        background: rgba(52, 211, 153, 0.05);
    }

    .terminal-rows {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .terminal-line {
        display: flex;
        gap: 16px;
        font-size: 12px;
        line-height: 1.5;
        padding: 2px 0;
    }

    .line-ts {
        color: #444;
        flex-shrink: 0;
    }
    .line-tag {
        color: var(--color-accent);
        flex-shrink: 0;
        opacity: 0.8;
    }
    .line-content {
        color: #d4d4d4;
        word-break: break-all;
        white-space: pre-wrap;
    }

    .terminal-line.error .line-content {
        color: #f87171;
    }
    .terminal-line.warn .line-content {
        color: #fbbf24;
    }

    .terminal-cursor {
        width: 8px;
        height: 16px;
        background: var(--color-accent);
        margin-top: 4px;
        animation: blink 1s step-end infinite;
    }

    .terminal-bottom-bar {
        padding: 8px 16px;
        background: rgba(255, 255, 255, 0.02);
        border-top: 1px solid rgba(255, 255, 255, 0.05);
        display: flex;
        align-items: center;
        justify-content: space-between;
        font-size: 10px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--color-text-muted);
    }

    .terminal-meta {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .sep {
        opacity: 0.2;
    }

    @keyframes blink {
        50% {
            opacity: 0;
        }
    }
</style>
