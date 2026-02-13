<script lang="ts">
    import { appState } from "$lib/stores.svelte";
    import { formatRelativeTime } from "$lib/utils";
    import { fade } from "svelte/transition";
    import {
        Cpu,
        Hash,
        Clock,
        Zap,
        Shield,
        BarChart3,
    } from "@lucide/svelte";

    const runs = $derived(appState.detailRuns ?? []);
    const latestRun = $derived(runs[0]);
    const session = $derived(appState.detailSession);

    const metrics = $derived([
        {
            label: "Runs",
            value: String(runs.length),
            sub: "Total runs in this session",
            icon: BarChart3,
            color: "var(--color-accent)",
        },
        {
            label: "Current Status",
            value: session?.status ?? "–",
            sub: "Session status",
            icon: Shield,
            color: "#34d399",
        },
        {
            label: "Tool",
            value: session?.tool ?? "–",
            sub: "AI tool adapter",
            icon: Cpu,
            color: "#818cf8",
        },
        {
            label: "Latest Phase",
            value: latestRun?.state.replace("AI_", "") ?? "–",
            sub: "Most recent run state",
            icon: Zap,
            color: "#fbbf24",
        },
        {
            label: "Last Updated",
            value: session ? formatRelativeTime(session.updated_at) : "–",
            sub: "Relative update time",
            icon: Clock,
            color: "#94a3b8",
        },
        {
            label: "Session ID",
            value: session ? session.id.substring(0, 8) : "–",
            sub: "Short identifier",
            icon: Hash,
            color: "var(--color-text-muted)",
        },
    ]);
</script>

<div class="stats-module" in:fade>
    <div class="metrics-grid">
        {#each metrics as metric, i}
            {@const Icon = metric.icon}
            <div class="metric-card glass" style="animation-delay: {i * 50}ms">
                <div class="metric-visual">
                    <div
                        class="metric-ring"
                        style="--ring-color: {metric.color}"
                    >
                        <Icon size={20} />
                    </div>
                    <div
                        class="metric-glow"
                        style="background: {metric.color}"
                    ></div>
                </div>

                <div class="metric-info">
                    <span class="m-label">{metric.label}</span>
                    <h3 class="m-value">{metric.value}</h3>
                    <p class="m-sub">{metric.sub}</p>
                </div>

                <div class="metric-decoration">
                    <div class="deco-dots"></div>
                </div>
            </div>
        {/each}
    </div>
</div>

<style>
    .stats-module {
        height: 100%;
    }

    .metrics-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 20px;
        padding-bottom: 40px;
    }

    .metric-card {
        position: relative;
        padding: 24px;
        border-radius: 24px;
        border: 1px solid var(--color-border);
        display: flex;
        align-items: flex-start;
        gap: 20px;
        overflow: hidden;
        transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
        background: rgba(255, 255, 255, 0.02);
    }

    .metric-card:hover {
        transform: translateY(-5px);
        border-color: var(--color-border-strong);
        background: rgba(255, 255, 255, 0.05);
        box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
    }

    .metric-visual {
        position: relative;
        flex-shrink: 0;
    }

    .metric-ring {
        width: 48px;
        height: 48px;
        border-radius: 14px;
        background: var(--color-bg-active);
        border: 1px solid var(--color-border);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--ring-color);
        position: relative;
        z-index: 2;
        transition: all 0.3s;
    }

    .metric-card:hover .metric-ring {
        border-color: var(--ring-color);
        transform: rotate(-5deg) scale(1.1);
    }

    .metric-glow {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 30px;
        height: 30px;
        border-radius: 50%;
        filter: blur(20px);
        opacity: 0.2;
        z-index: 1;
        transition: opacity 0.3s;
    }

    .metric-card:hover .metric-glow {
        opacity: 0.5;
    }

    .metric-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 2px;
        min-width: 0;
    }

    .m-label {
        font-size: 11px;
        font-weight: 800;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--color-text-muted);
    }

    .m-value {
        font-size: 20px;
        font-weight: 800;
        color: var(--color-text);
        margin: 4px 0 2px;
        letter-spacing: -0.02em;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .m-sub {
        font-size: 12px;
        color: var(--color-text-secondary);
        opacity: 0.7;
    }

    .metric-decoration {
        position: absolute;
        right: 12px;
        top: 12px;
        opacity: 0.1;
    }

    .deco-dots {
        width: 20px;
        height: 20px;
        background-image: radial-gradient(white 1px, transparent 1px);
        background-size: 4px 4px;
    }

    .metric-card:hover .deco-dots {
        animation: pulse 1s infinite;
    }

    @keyframes pulse {
        0%,
        100% {
            opacity: 0.2;
        }
        50% {
            opacity: 0.5;
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
