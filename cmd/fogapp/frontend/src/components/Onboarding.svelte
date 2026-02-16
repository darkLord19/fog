<script lang="ts">
    import { onMount } from "svelte";
    import { fade, fly, slide } from "svelte/transition";
    import {
        Check,
        ChevronRight,
        Loader2,
        Search,
        Github,
    } from "@lucide/svelte";
    import { appState } from "$lib/stores.svelte";
    import { updateSettings, discoverRepos, importRepos } from "$lib/api";
    import type { DiscoveredRepo } from "$lib/types";
    import { TOOL_MODELS, getModelsForTool } from "$lib/constants";
    import Dropdown from "./Dropdown.svelte";

    let step = $state(0);
    let loading = $state(false);
    let error = $state("");

    // Step 1: PAT
    let githubPat = $state("");

    // Step 2: Tools
    let selectedTool = $state("");
    let selectedModel = $state("");

    // Step 3: Repos
    let discovered = $state<DiscoveredRepo[]>([]);
    let selectedRepos = $state<string[]>([]);
    let repoSearch = $state("");

    let availableTools = $derived(
        appState.settings?.available_tools.map((t) => ({
            value: t,
            label: t,
        })) ?? [],
    );

    let availableModels = $derived(
        selectedTool
            ? getModelsForTool(selectedTool).map((m) => ({
                  value: m,
                  label: m,
              }))
            : [],
    );

    async function nextStep() {
        error = "";
        loading = true;

        try {
            if (step === 0) {
                if (!githubPat.trim()) {
                    throw new Error(
                        "Please enter a GitHub Personal Access Token",
                    );
                }
                // Save PAT
                await updateSettings({ github_pat: githubPat });
                // Proceed to tool selection
                step = 1;

                // Pre-select if available
                if (appState.settings?.default_tool) {
                    selectedTool = appState.settings.default_tool;
                }
                if (
                    selectedTool &&
                    appState.settings?.default_models?.[selectedTool]
                ) {
                    selectedModel =
                        appState.settings.default_models[selectedTool];
                }
            } else if (step === 1) {
                if (!selectedTool)
                    throw new Error("Please select a default AI tool");
                if (!selectedModel)
                    throw new Error("Please select a default model");

                // Save preferences
                await updateSettings({
                    default_tool: selectedTool,
                    default_model: selectedModel,
                    default_models: { [selectedTool]: selectedModel },
                });

                // Start discovery for next step
                step = 2;
                await loadRepos();
            }
        } catch (err) {
            error = err instanceof Error ? err.message : "An error occurred";
        } finally {
            loading = false;
        }
    }

    async function loadRepos() {
        loading = true;
        try {
            discovered = await discoverRepos();
        } catch (err) {
            error = "Failed to discover repositories";
        } finally {
            loading = false;
        }
    }

    function toggleRepo(path: string) {
        if (selectedRepos.includes(path)) {
            selectedRepos = selectedRepos.filter((p) => p !== path);
        } else {
            selectedRepos = [...selectedRepos, path];
        }
    }

    async function finish() {
        error = "";
        loading = true;
        try {
            if (selectedRepos.length > 0) {
                await importRepos(selectedRepos);
            }
            // Refresh settings to clear onboarding_required
            await appState.refreshAll();
        } catch (err) {
            console.error("Onboarding finish error:", err);
            error =
                err instanceof Error ? err.message : "Failed to complete setup";
            loading = false;
        }
    }

    let filteredRepos = $derived(
        discovered.filter((r) =>
            r.full_name.toLowerCase().includes(repoSearch.toLowerCase()),
        ),
    );
</script>

<div class="onboarding-overlay" transition:fade>
    <div class="onboarding-card" in:fly={{ y: 20, duration: 400 }}>
        <div class="progress-bar">
            <div
                class="progress-fill"
                style="width: {(step + 1) * 33.33}%"
            ></div>
        </div>

        <div class="step-content">
            {#if step === 0}
                <div class="step-header">
                    <div class="icon-circle">
                        <Github size={24} />
                    </div>
                    <h2>Connect GitHub</h2>
                    <p>
                        Enter your GitHub Personal Access Token (PAT) to access
                        your repositories and create pull requests.
                    </p>
                </div>

                <div class="input-group">
                    <label for="pat">Personal Access Token</label>
                    <input
                        id="pat"
                        type="password"
                        bind:value={githubPat}
                        placeholder="ghp_..."
                        class="text-input"
                    />
                    <p class="hint">
                        Requires <code>repo</code> and <code>read:user</code> scopes.
                    </p>
                </div>
            {:else if step === 1}
                <div class="step-header">
                    <div class="icon-circle">
                        <div class="ai-icon">✨</div>
                    </div>
                    <h2>Select AI Model</h2>
                    <p>
                        Choose your preferred AI tool and model for code
                        generation.
                    </p>
                </div>

                <div class="input-group">
                    <label for="tool">Default Tool</label>
                    <Dropdown
                        bind:value={selectedTool}
                        options={availableTools}
                        placeholder="Select Tool..."
                        class="full-width"
                    />
                </div>

                <div class="input-group">
                    <label for="model">Default Model</label>
                    {#if !selectedTool || availableModels.length > 0}
                        <Dropdown
                            bind:value={selectedModel}
                            options={availableModels}
                            placeholder={selectedTool
                                ? "Select Model..."
                                : "Select Tool First..."}
                            disabled={!selectedTool}
                            class="full-width"
                        />
                    {:else}
                        <input
                            id="model"
                            type="text"
                            bind:value={selectedModel}
                            placeholder="e.g. gpt-4"
                            class="text-input"
                        />
                    {/if}
                </div>
            {:else if step === 2}
                <div class="step-header">
                    <div class="icon-circle">
                        <Search size={24} />
                    </div>
                    <h2>Add Repositories</h2>
                    <p>Select GitHub repositories to import into Fog.</p>
                </div>

                <div class="repo-selection">
                    <div class="search-bar">
                        <Search size={16} class="search-icon" />
                        <input
                            type="text"
                            bind:value={repoSearch}
                            placeholder="Filter repositories..."
                        />
                    </div>

                    <div class="repo-list">
                        {#if loading && discovered.length === 0}
                            <div class="loading-state">
                                <Loader2 class="spin" size={24} />
                                <span>Fetching from GitHub...</span>
                            </div>
                        {:else}
                            {#each filteredRepos as repo}
                                <button
                                    class="repo-item {selectedRepos.includes(
                                        repo.full_name,
                                    )
                                        ? 'selected'
                                        : ''}"
                                    onclick={() => toggleRepo(repo.full_name)}
                                >
                                    <div class="repo-info">
                                        <span class="repo-name"
                                            >{repo.full_name}</span
                                        >
                                        <span class="repo-path"
                                            >{repo.clone_url}</span
                                        >
                                    </div>
                                    {#if selectedRepos.includes(repo.full_name)}
                                        <Check size={18} class="check-icon" />
                                    {/if}
                                </button>
                            {/each}
                        {/if}
                    </div>
                    <p class="hint-text">{selectedRepos.length} selected</p>
                </div>
            {/if}
        </div>

        {#if error}
            <div class="error-banner" transition:slide>
                {error}
            </div>
        {/if}

        <div class="actions">
            {#if step === 2}
                <button
                    class="btn btn-primary full-width"
                    onclick={finish}
                    disabled={loading}
                >
                    {#if loading}
                        <Loader2 class="spin" size={16} />
                        Finishing...
                    {:else}
                        Finish Setup
                    {/if}
                </button>
            {:else}
                <button
                    class="btn btn-primary full-width"
                    onclick={nextStep}
                    disabled={loading}
                >
                    {#if loading}
                        <Loader2 class="spin" size={16} />
                        Saving...
                    {:else}
                        Continue <ChevronRight size={16} />
                    {/if}
                </button>
            {/if}
            {#if step > 0 && !loading && step < 2}
                <button class="btn btn-ghost" onclick={() => step--}
                    >Back</button
                >
            {/if}
        </div>
    </div>
</div>

<style>
    .onboarding-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.8);
        backdrop-filter: blur(8px);
        z-index: 9999;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 20px;
    }

    .onboarding-card {
        background: var(--color-bg-elevated);
        border: 1px solid var(--color-border);
        border-radius: 16px;
        width: 100%;
        max-width: 480px;
        box-shadow: var(--shadow-xl);
        width: 100%;
        max-width: 480px;
        box-shadow: var(--shadow-xl);
        position: relative;
        display: flex;
        flex-direction: column;
        max-height: 90vh;
    }

    .progress-bar {
        height: 4px;
        background: var(--color-bg-surface);
        width: 100%;
        border-top-left-radius: 16px;
        border-top-right-radius: 16px;
        overflow: hidden; /* Keep bar fill contained */
    }

    .progress-fill {
        height: 100%;
        background: var(--color-accent);
        transition: width 0.3s ease;
    }

    .step-content {
        padding: 32px;
        /* Remove overflow-y: auto to allow dropdowns to fly out */
    }

    .step-header {
        text-align: center;
        margin-bottom: 32px;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .icon-circle {
        width: 64px;
        height: 64px;
        border-radius: 50%;
        background: var(--color-bg-surface);
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 16px;
        color: var(--color-text);
        border: 1px solid var(--color-border);
    }

    .ai-icon {
        font-size: 24px;
    }

    h2 {
        font-size: 24px;
        font-weight: 600;
        color: var(--color-text);
        margin-bottom: 8px;
    }

    p {
        font-size: 14px;
        color: var(--color-text-secondary);
        line-height: 1.5;
    }

    .input-group {
        margin-bottom: 20px;
    }

    .input-group label {
        display: block;
        font-size: 13px;
        font-weight: 500;
        color: var(--color-text);
        margin-bottom: 8px;
    }

    .text-input {
        width: 100%;
        background: var(--color-bg);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        padding: 10px 12px;
        color: var(--color-text);
        font-size: 14px;
        outline: none;
        transition: border-color 0.2s;
    }

    .text-input:focus {
        border-color: var(--color-accent);
    }

    .hint {
        font-size: 12px;
        color: var(--color-text-muted);
        margin-top: 6px;
    }

    code {
        background: var(--color-bg-surface);
        padding: 2px 4px;
        border-radius: 4px;
        font-family: var(--font-mono);
    }

    .actions {
        padding: 24px 32px;
        border-top: 1px solid var(--color-border);
        background: var(--color-bg-surface);
        display: flex;
        flex-direction: column;
        gap: 12px;
        border-bottom-left-radius: 16px;
        border-bottom-right-radius: 16px;
    }

    .btn {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 10px 16px;
        border-radius: 8px;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
        border: none;
    }

    .btn-primary {
        background: var(--color-accent);
        color: #000;
    }

    .btn-primary:hover:not(:disabled) {
        opacity: 0.9;
    }

    .btn-primary:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .btn-ghost {
        background: transparent;
        color: var(--color-text-secondary);
    }

    .btn-ghost:hover {
        color: var(--color-text);
        background: var(--color-bg-hover);
    }

    .full-width {
        width: 100%;
    }

    :global(.spin) {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    /* Repo Selection */
    .repo-selection {
        display: flex;
        flex-direction: column;
        height: 300px;
    }

    .search-bar {
        position: relative;
        margin-bottom: 12px;
    }

    :global(.search-icon) {
        position: absolute;
        left: 10px;
        top: 50%;
        transform: translateY(-50%);
        color: var(--color-text-muted);
    }

    .search-bar input {
        width: 100%;
        background: var(--color-bg-surface);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        padding: 8px 12px 8px 36px;
        color: var(--color-text);
        font-size: 13px;
        outline: none;
    }

    .search-bar input:focus {
        border-color: var(--color-accent);
    }

    .repo-list {
        flex: 1;
        overflow-y: auto;
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-bg);
    }

    .loading-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
        color: var(--color-text-muted);
        gap: 12px;
        font-size: 13px;
    }

    .repo-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        padding: 10px 12px;
        border: none;
        background: transparent;
        border-bottom: 1px solid var(--color-border);
        cursor: pointer;
        text-align: left;
        transition: background 0.1s;
    }

    .repo-item:last-child {
        border-bottom: none;
    }

    .repo-item:hover {
        background: var(--color-bg-hover);
    }

    .repo-item.selected {
        background: rgba(250, 204, 21, 0.05);
    }

    .repo-info {
        display: flex;
        flex-direction: column;
        gap: 2px;
        overflow: hidden;
    }

    .repo-name {
        font-size: 13px;
        font-weight: 500;
        color: var(--color-text);
    }

    .repo-path {
        font-size: 11px;
        color: var(--color-text-muted);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    :global(.check-icon) {
        color: var(--color-accent);
        flex-shrink: 0;
    }

    .hint-text {
        text-align: right;
        font-size: 12px;
        color: var(--color-text-muted);
        margin-top: 8px;
    }

    .error-banner {
        background: rgba(239, 68, 68, 0.1);
        color: var(--color-danger);
        padding: 10px;
        text-align: center;
        font-size: 13px;
        border-top: 1px solid rgba(239, 68, 68, 0.2);
    }
</style>
