<script lang="ts">
  import { appState } from "$lib/stores.svelte";
  import SessionItem from "./SessionItem.svelte";
  import { slide, fade } from "svelte/transition";
  import {
    Sparkles,
    Settings,
    History,
    Activity,
    ChevronDown,
    Plus,
  } from "@lucide/svelte";

  let historyCollapsed = $state(false);

  function showNew() {
    appState.currentView = "new";
    appState.selectedSessionID = "";
  }

  function showSettings() {
    appState.currentView = "settings";
    appState.selectedSessionID = "";
  }

  function toggleHistory() {
    historyCollapsed = !historyCollapsed;
  }
</script>

<aside class="sidebar glass">
    <div class="sidebar-header">
        <div class="brand">
            <div class="logo-orb">
                <Sparkles size={14} />
            </div>
            <span class="brand-name">Fog</span>
        </div>
    <button id="show-new" class="btn btn-primary new-chat-btn" onclick={showNew}>
      <Plus size={16} strokeWidth={3} />
      <span>New Chat</span>
    </button>
  </div>

  <div class="sidebar-content">
    <!-- Active Sessions -->
    {#if appState.runningSessions.length > 0}
        <div class="section">
            <div class="section-label-container">
              <Activity size={10} />
              <span class="section-label">Active now</span>
            </div>
        <div id="running-sessions" class="session-list">
          {#each appState.runningSessions as session (session.id)}
            <SessionItem {session} />
          {/each}
        </div>
      </div>
    {/if}

    <!-- History -->
    {#if appState.completedSessions.length > 0}
      <div class="section">
        <button class="section-header" onclick={toggleHistory}>
          <div class="label-with-icon">
            <History size={10} />
            <span class="section-label">Past sessions</span>
          </div>
          <ChevronDown
            size={12}
            class="collapse-icon {historyCollapsed ? 'collapsed' : ''}"
          />
        </button>
        {#if !historyCollapsed}
          <div
            id="completed-sessions"
            class="session-list"
            transition:slide={{ duration: 300 }}
          >
            {#each appState.completedSessions as session (session.id)}
              <SessionItem {session} />
            {/each}
          </div>
        {/if}
      </div>
    {/if}

    {#if appState.sessions.length === 0}
      <div class="empty-state" in:fade>
        <div class="empty-circle">
          <Sparkles size={24} />
        </div>
        <p>No activity yet</p>
        <p class="sub">Your AI journey starts here</p>
      </div>
    {/if}
  </div>

  <div class="sidebar-footer">
    <button id="show-settings" class="footer-nav-item" onclick={showSettings}>
      <Settings size={16} />
      <span>Settings</span>
    </button>
    <div class="status-pill" title={appState.daemonStatus}>
      <span
        class="status-dot"
        class:connected={appState.daemonStatus === "connected"}
      ></span>
      <span id="daemon-badge" class="status-text">{appState.daemonStatus}</span>
    </div>
  </div>
</aside>

<style>
  .sidebar {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 280px;
    flex-shrink: 0;
    border-right: 1px solid var(--color-border);
    z-index: 50;
  }

  .sidebar-header {
    padding: 24px 20px 16px;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .logo-orb {
    width: 32px;
    height: 32px;
    background: var(--color-accent-gradient);
    color: white;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
  }

  .brand-name {
    font-size: 18px;
    font-weight: 700;
    letter-spacing: -0.01em;
    color: var(--color-text);
  }

  .new-chat-btn {
    width: 100%;
    height: 40px;
    border-radius: var(--radius-md);
    font-weight: 700;
  }

  .sidebar-content {
    flex: 1;
    overflow-y: auto;
    padding: 12px 0;
  }

  .section {
    margin-bottom: 24px;
  }

  .section-label-container {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 20px 8px;
    color: var(--color-text-muted);
    opacity: 0.8;
  }

  .label-with-icon {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .section-header {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 20px;
    background: none;
    border: none;
    cursor: pointer;
    color: var(--color-text-muted);
    transition: color 0.2s;
  }

  .section-header:hover {
    color: var(--color-text-secondary);
  }

  .section-label {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  :global(.collapse-icon) {
    transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  }

  :global(.collapse-icon.collapsed) {
    transform: rotate(-90deg);
  }

  .session-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 4px 12px;
  }

  .empty-state {
    padding: 60px 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    color: var(--color-text-muted);
  }

  .empty-circle {
    width: 56px;
    height: 56px;
    border-radius: 50%;
    background: var(--color-bg-active);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 16px;
    border: 1px solid var(--color-border);
    color: var(--color-text-secondary);
  }

  .empty-state p {
    font-size: 14px;
    font-weight: 600;
    color: var(--color-text-secondary);
    margin-bottom: 4px;
  }

  .empty-state .sub {
    font-size: 12px;
    font-weight: 400;
    opacity: 0.7;
  }

  .sidebar-footer {
    padding: 16px 20px;
    border-top: 1px solid var(--color-border);
    background: rgba(0, 0, 0, 0.2);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .footer-nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    background: none;
    border: none;
    color: var(--color-text-secondary);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    padding: 8px 0;
  }

  .footer-nav-item:hover {
    color: var(--color-text);
  }

  .status-pill {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 20px;
    border: 1px solid var(--color-border);
  }

  .status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #4b5563;
  }

  .status-dot.connected {
    background: var(--color-success);
    box-shadow: 0 0 8px var(--color-success);
  }

  .status-text {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.02em;
    color: var(--color-text-muted);
  }
</style>
