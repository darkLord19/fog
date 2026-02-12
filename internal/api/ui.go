package api

const webUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Fog</title>
  <style>
    :root {
      --bg: #f4f0e8;
      --panel: #fffdf8;
      --ink: #1c1f20;
      --muted: #5b6768;
      --accent: #0a6a58;
      --accent-soft: #d6efe8;
      --warn: #8b2e24;
      --line: #d9d2c4;
      --mono: "JetBrains Mono", Menlo, monospace;
      --sans: "Avenir Next", "Segoe UI", sans-serif;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: var(--sans);
      color: var(--ink);
      background:
        radial-gradient(circle at 10% -20%, #d6ece8 0%, transparent 40%),
        radial-gradient(circle at 90% 0%, #efe2c6 0%, transparent 30%),
        var(--bg);
    }
    .wrap {
      max-width: 1180px;
      margin: 0 auto;
      padding: 24px 16px 42px;
    }
    h1 {
      margin: 0 0 6px;
      font-size: 28px;
      letter-spacing: 0.2px;
    }
    .sub {
      margin: 0;
      color: var(--muted);
    }
    .grid {
      display: grid;
      gap: 14px;
      grid-template-columns: repeat(12, minmax(0, 1fr));
      margin-top: 18px;
    }
    .card {
      grid-column: span 12;
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 14px;
      padding: 16px;
      box-shadow: 0 6px 18px rgba(0,0,0,0.04);
    }
    .card h2 {
      margin: 0 0 10px;
      font-size: 17px;
    }
    .kpi {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
      margin-top: 6px;
    }
    .pill {
      font-family: var(--mono);
      font-size: 12px;
      padding: 6px 10px;
      border-radius: 999px;
      border: 1px solid var(--line);
      background: #fff;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 13px;
    }
    th, td {
      text-align: left;
      border-bottom: 1px solid var(--line);
      padding: 8px 6px;
      vertical-align: top;
    }
    th {
      color: var(--muted);
      font-weight: 600;
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: 0.35px;
    }
    .mono { font-family: var(--mono); }
    .state {
      font-family: var(--mono);
      font-size: 12px;
      padding: 2px 8px;
      border-radius: 10px;
      background: #eef2ee;
      display: inline-block;
    }
    .state.active { background: var(--accent-soft); color: #004b43; }
    .state.failed { background: #f6d5cf; color: #732017; }
    .row {
      display: grid;
      grid-template-columns: 150px 1fr;
      gap: 10px;
      margin-bottom: 10px;
      align-items: center;
    }
    .row-inline {
      display: flex;
      gap: 10px;
      align-items: center;
      flex-wrap: wrap;
      margin-bottom: 10px;
    }
    label { font-size: 13px; color: var(--muted); }
    input, select, textarea {
      width: 100%;
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 9px 10px;
      background: #fff;
      color: var(--ink);
      font-family: inherit;
    }
    textarea {
      min-height: 86px;
      resize: vertical;
    }
    button {
      border: 0;
      border-radius: 10px;
      padding: 10px 14px;
      background: var(--accent);
      color: #fff;
      font-weight: 600;
      cursor: pointer;
    }
    .button-muted {
      background: #edf2f1;
      color: #1d2b2a;
      border: 1px solid var(--line);
    }
    button:disabled { opacity: 0.6; cursor: not-allowed; }
    .status {
      margin-left: 8px;
      font-size: 12px;
      color: var(--muted);
    }
    .status.error { color: var(--warn); }
    .status.ok { color: #146350; }
    .callout {
      border: 1px solid #e4d7b8;
      background: #fff4dc;
      border-radius: 10px;
      padding: 10px 12px;
      margin-bottom: 12px;
      font-size: 13px;
    }
    .hint {
      margin: 0 0 10px;
      color: var(--muted);
      font-size: 13px;
    }
    .repo-list {
      max-height: 220px;
      overflow: auto;
      border: 1px solid var(--line);
      border-radius: 10px;
      padding: 8px;
      background: #fff;
    }
    .repo-item {
      display: flex;
      gap: 8px;
      align-items: center;
      padding: 4px 2px;
      font-size: 13px;
    }
    @media (min-width: 980px) {
      .onboarding { grid-column: span 12; }
      .sessions { grid-column: span 8; }
      .control { grid-column: span 4; }
      .repos { grid-column: span 6; }
      .settings { grid-column: span 6; }
    }
  </style>
</head>
<body>
  <div class="wrap">
    <h1>Fog Control Plane</h1>
    <p class="sub">Manage sessions, repos, and local defaults.</p>

    <div class="grid">
      <section class="card onboarding" id="onboarding-card" style="display:none;">
        <h2>Onboarding</h2>
        <p class="hint">Set your GitHub PAT and default tool to start using Fog.</p>
        <form id="onboarding-form">
          <div class="row">
            <label for="onboarding-pat">GitHub PAT</label>
            <input id="onboarding-pat" name="github_pat" type="password" placeholder="ghp_...">
          </div>
          <div class="row">
            <label for="onboarding-tool">Default Tool</label>
            <select id="onboarding-tool" name="default_tool"></select>
          </div>
          <button id="onboarding-save-btn" type="submit">Complete Onboarding</button>
          <span class="status" id="onboarding-status"></span>
        </form>
      </section>

      <section class="card sessions">
        <h2>Sessions</h2>
        <div class="kpi">
          <div class="pill">Active: <span id="kpi-active">0</span></div>
          <div class="pill">Total: <span id="kpi-total">0</span></div>
          <div class="pill">Updated: <span id="kpi-updated">-</span></div>
        </div>
        <div style="overflow:auto; margin-top:12px;">
          <table>
            <thead>
              <tr>
                <th>State</th>
                <th>Repo</th>
                <th>Branch</th>
                <th>Tool</th>
                <th>PR</th>
                <th>Updated</th>
              </tr>
            </thead>
            <tbody id="session-body">
              <tr><td colspan="6">Loading sessions...</td></tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="card control">
        <h2>New Session</h2>
        <form id="create-session-form">
          <div class="row">
            <label for="session-repo">Repo</label>
            <select id="session-repo" name="repo"></select>
          </div>
          <div class="row">
            <label for="session-tool">Tool</label>
            <select id="session-tool" name="tool"></select>
          </div>
          <div class="row">
            <label for="session-model">Model</label>
            <input id="session-model" name="model" placeholder="optional">
          </div>
          <div class="row">
            <label for="session-branch">Branch Name</label>
            <input id="session-branch" name="branch_name" placeholder="optional">
          </div>
          <div class="row">
            <label for="session-autopr">Auto PR</label>
            <select id="session-autopr" name="autopr">
              <option value="false">false</option>
              <option value="true">true</option>
            </select>
          </div>
          <div class="row">
            <label for="session-prompt">Prompt</label>
            <textarea id="session-prompt" name="prompt" placeholder="Describe the task"></textarea>
          </div>
          <button id="create-session-btn" type="submit">Create Session</button>
          <span class="status" id="create-session-status"></span>
        </form>

        <h2 style="margin-top:18px;">Follow-up</h2>
        <form id="followup-form">
          <div class="row">
            <label for="followup-session">Session</label>
            <select id="followup-session" name="session_id"></select>
          </div>
          <div class="row">
            <label for="followup-prompt">Prompt</label>
            <textarea id="followup-prompt" name="prompt" placeholder="Follow-up instruction"></textarea>
          </div>
          <button id="followup-btn" type="submit">Send Follow-up</button>
          <span class="status" id="followup-status"></span>
        </form>
      </section>

      <section class="card repos">
        <h2>Repositories</h2>
        <p class="hint">Identifiers are fixed as <span class="mono">owner/repo-name</span>.</p>
        <div class="row-inline">
          <button class="button-muted" id="discover-repos-btn" type="button">Discover Accessible Repos</button>
          <button id="import-repos-btn" type="button">Import Selected</button>
          <span class="status" id="repos-status"></span>
        </div>
        <div class="repo-list" id="discover-repo-list">No discovered repos yet</div>
        <div style="margin-top:12px; overflow:auto;">
          <table>
            <thead>
              <tr>
                <th>Repo</th>
                <th>Default Branch</th>
                <th>Path</th>
              </tr>
            </thead>
            <tbody id="managed-repo-body">
              <tr><td colspan="3">Loading managed repos...</td></tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="card settings">
        <h2>Settings</h2>
        <div class="callout" id="setup-warning" style="display:none;">
          Onboarding is incomplete. Save GitHub PAT and default tool above.
        </div>
        <form id="settings-form">
          <div class="row">
            <label for="default-tool">Default Tool</label>
            <select id="default-tool" name="default_tool"></select>
          </div>
          <div class="row">
            <label for="branch-prefix">Branch Prefix</label>
            <input id="branch-prefix" name="branch_prefix" placeholder="fog">
          </div>
          <div class="row">
            <label>GitHub PAT</label>
            <div class="mono" id="pat-status">-</div>
          </div>
          <button id="save-btn" type="submit">Save Settings</button>
          <span class="status" id="save-status"></span>
        </form>
      </section>
    </div>
  </div>

  <script>
    var activeStates = { CREATED: true, SETUP: true, AI_RUNNING: true, VALIDATING: true, COMMITTED: true, PR_CREATED: true };
    var latestSessions = [];
    var discoveredRepos = [];

    function escapeHTML(value) {
      return String(value || "")
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/\"/g, "&quot;")
        .replace(/'/g, "&#39;");
    }

    function stateClass(state) {
      if (state === "FAILED") return "state failed";
      if (activeStates[state]) return "state active";
      return "state";
    }

    function formatDate(value) {
      if (!value) return "-";
      var dt = new Date(value);
      if (isNaN(dt.getTime())) return "-";
      return dt.toLocaleString();
    }

    function setStatus(id, message, cls) {
      var el = document.getElementById(id);
      el.textContent = message || "";
      el.className = "status" + (cls ? (" " + cls) : "");
    }

    async function fetchJSON(url, options) {
      var res = await fetch(url, options || {});
      if (!res.ok) {
        var txt = await res.text();
        throw new Error(txt || ("HTTP " + res.status));
      }
      return res.json();
    }

    function populateSessionSelect() {
      var select = document.getElementById("followup-session");
      if (!latestSessions.length) {
        select.innerHTML = "<option value=''>No sessions</option>";
        return;
      }
      select.innerHTML = latestSessions.map(function (s) {
        return "<option value='" + escapeHTML(s.id) + "'>" +
          escapeHTML(s.repo_name + " :: " + s.branch) + "</option>";
      }).join("");
    }

    async function refreshSessions() {
      var body = document.getElementById("session-body");
      try {
        var sessions = await fetchJSON("/api/sessions");
        latestSessions = sessions || [];
        populateSessionSelect();

        var active = latestSessions.filter(function (s) { return activeStates[s.status] || s.busy; });
        document.getElementById("kpi-active").textContent = String(active.length);
        document.getElementById("kpi-total").textContent = String(latestSessions.length);
        document.getElementById("kpi-updated").textContent = new Date().toLocaleTimeString();

        if (!latestSessions.length) {
          body.innerHTML = "<tr><td colspan='6'>No sessions yet</td></tr>";
          return;
        }

        body.innerHTML = latestSessions.slice(0, 60).map(function (s) {
          var pr = s.pr_url ? "<a href='" + escapeHTML(s.pr_url) + "' target='_blank' rel='noopener'>open</a>" : "-";
          return "<tr>" +
            "<td><span class='" + stateClass(s.status) + "'>" + escapeHTML(s.status) + (s.busy ? "*" : "") + "</span></td>" +
            "<td class='mono'>" + escapeHTML(s.repo_name) + "</td>" +
            "<td class='mono'>" + escapeHTML(s.branch) + "</td>" +
            "<td>" + escapeHTML(s.tool) + "</td>" +
            "<td>" + pr + "</td>" +
            "<td>" + escapeHTML(formatDate(s.updated_at)) + "</td>" +
            "</tr>";
        }).join("");
      } catch (err) {
        body.innerHTML = "<tr><td colspan='6'>Failed to load sessions: " + escapeHTML(err.message) + "</td></tr>";
      }
    }

    function renderManagedRepos(repos) {
      var body = document.getElementById("managed-repo-body");
      if (!repos || !repos.length) {
        body.innerHTML = "<tr><td colspan='3'>No managed repos</td></tr>";
        return;
      }
      body.innerHTML = repos.map(function (r) {
        return "<tr>" +
          "<td class='mono'>" + escapeHTML(r.name) + "</td>" +
          "<td>" + escapeHTML(r.default_branch || "-") + "</td>" +
          "<td class='mono'>" + escapeHTML(r.base_worktree_path || "-") + "</td>" +
          "</tr>";
      }).join("");
    }

    function renderDiscoveredRepos() {
      var el = document.getElementById("discover-repo-list");
      if (!discoveredRepos.length) {
        el.innerHTML = "No discovered repos yet";
        return;
      }
      el.innerHTML = discoveredRepos.map(function (r, idx) {
        return "<label class='repo-item'>" +
          "<input type='checkbox' data-repo='" + escapeHTML(r.full_name) + "' id='disco-" + idx + "'>" +
          "<span class='mono'>" + escapeHTML(r.full_name) + "</span>" +
          "<span>" + escapeHTML(r.default_branch || "-") + "</span>" +
          "</label>";
      }).join("");
    }

    async function loadManagedRepos() {
      var repos = await fetchJSON("/api/repos");
      renderManagedRepos(repos);

      var repoSelect = document.getElementById("session-repo");
      if (!repos.length) {
        repoSelect.innerHTML = "<option value=''>No repos imported</option>";
      } else {
        repoSelect.innerHTML = repos.map(function (r) {
          return "<option value='" + escapeHTML(r.name) + "'>" + escapeHTML(r.name) + "</option>";
        }).join("");
      }
    }

    async function loadSettings() {
      var settings = await fetchJSON("/api/settings");
      var select = document.getElementById("default-tool");
      var onboardingSelect = document.getElementById("onboarding-tool");
      var createToolSelect = document.getElementById("session-tool");
      var options = settings.available_tools || [];
      if (!options.length && settings.default_tool) {
        options = [settings.default_tool];
      }
      var optionsHTML = options.map(function (tool) {
        var selected = tool === settings.default_tool ? " selected" : "";
        return "<option value='" + escapeHTML(tool) + "'" + selected + ">" + escapeHTML(tool) + "</option>";
      }).join("");
      select.innerHTML = optionsHTML;
      onboardingSelect.innerHTML = optionsHTML;
      createToolSelect.innerHTML = "<option value=''>default</option>" + optionsHTML;

      if (!settings.default_tool && options.length) {
        onboardingSelect.value = options[0];
      }

      document.getElementById("branch-prefix").value = settings.branch_prefix || "fog";
      document.getElementById("pat-status").textContent = settings.has_github_token ? "configured" : "missing";
      document.getElementById("onboarding-card").style.display = settings.onboarding_required ? "block" : "none";
      document.getElementById("setup-warning").style.display = settings.onboarding_required ? "block" : "none";
    }

    async function onSaveSettings(event) {
      event.preventDefault();
      var saveBtn = document.getElementById("save-btn");
      setStatus("save-status", "", "");
      saveBtn.disabled = true;

      try {
        var payload = {
          default_tool: document.getElementById("default-tool").value,
          branch_prefix: document.getElementById("branch-prefix").value
        };
        var patInput = document.getElementById("onboarding-pat").value.trim();
        if (patInput) {
          payload.github_pat = patInput;
        }

        await fetchJSON("/api/settings", {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload)
        });
        setStatus("save-status", "Saved", "ok");
        document.getElementById("onboarding-pat").value = "";
        await loadSettings();
      } catch (err) {
        setStatus("save-status", "Save failed: " + err.message, "error");
      } finally {
        saveBtn.disabled = false;
      }
    }

    async function onCompleteOnboarding(event) {
      event.preventDefault();
      var btn = document.getElementById("onboarding-save-btn");
      setStatus("onboarding-status", "", "");
      btn.disabled = true;

      try {
        var pat = document.getElementById("onboarding-pat").value.trim();
        var tool = document.getElementById("onboarding-tool").value;
        if (!pat) throw new Error("GitHub PAT is required");
        if (!tool) throw new Error("Default tool is required");

        await fetchJSON("/api/settings", {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ github_pat: pat, default_tool: tool })
        });

        setStatus("onboarding-status", "Onboarding complete", "ok");
        document.getElementById("onboarding-pat").value = "";
        await loadSettings();
      } catch (err) {
        setStatus("onboarding-status", "Onboarding failed: " + err.message, "error");
      } finally {
        btn.disabled = false;
      }
    }

    async function onCreateSession(event) {
      event.preventDefault();
      var btn = document.getElementById("create-session-btn");
      setStatus("create-session-status", "", "");
      btn.disabled = true;

      try {
        var payload = {
          repo: document.getElementById("session-repo").value,
          prompt: document.getElementById("session-prompt").value.trim(),
          model: document.getElementById("session-model").value.trim(),
          branch_name: document.getElementById("session-branch").value.trim(),
          autopr: document.getElementById("session-autopr").value === "true",
          async: true
        };
        var tool = document.getElementById("session-tool").value;
        if (tool) payload.tool = tool;

        if (!payload.repo) throw new Error("Repo is required");
        if (!payload.prompt) throw new Error("Prompt is required");

        var out = await fetchJSON("/api/sessions", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload)
        });

        setStatus("create-session-status", "Queued session " + out.session_id, "ok");
        document.getElementById("session-prompt").value = "";
        document.getElementById("session-branch").value = "";
        await refreshSessions();
      } catch (err) {
        setStatus("create-session-status", "Create failed: " + err.message, "error");
      } finally {
        btn.disabled = false;
      }
    }

    async function onFollowUp(event) {
      event.preventDefault();
      var btn = document.getElementById("followup-btn");
      setStatus("followup-status", "", "");
      btn.disabled = true;

      try {
        var sessionID = document.getElementById("followup-session").value;
        var prompt = document.getElementById("followup-prompt").value.trim();
        if (!sessionID) throw new Error("Session is required");
        if (!prompt) throw new Error("Prompt is required");

        var out = await fetchJSON("/api/sessions/" + encodeURIComponent(sessionID) + "/runs", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ prompt: prompt, async: true })
        });

        setStatus("followup-status", "Queued run " + out.run_id, "ok");
        document.getElementById("followup-prompt").value = "";
        await refreshSessions();
      } catch (err) {
        setStatus("followup-status", "Follow-up failed: " + err.message, "error");
      } finally {
        btn.disabled = false;
      }
    }

    async function onDiscoverRepos() {
      var btn = document.getElementById("discover-repos-btn");
      setStatus("repos-status", "", "");
      btn.disabled = true;
      try {
        discoveredRepos = await fetchJSON("/api/repos/discover", { method: "POST" });
        renderDiscoveredRepos();
        setStatus("repos-status", "Discovered " + discoveredRepos.length + " repos", "ok");
      } catch (err) {
        setStatus("repos-status", "Discover failed: " + err.message, "error");
      } finally {
        btn.disabled = false;
      }
    }

    async function onImportRepos() {
      var btn = document.getElementById("import-repos-btn");
      setStatus("repos-status", "", "");
      btn.disabled = true;
      try {
        var checked = Array.prototype.slice.call(document.querySelectorAll("#discover-repo-list input[type='checkbox']:checked"));
        var repos = checked.map(function (el) { return el.getAttribute("data-repo"); });
        if (!repos.length) throw new Error("Select at least one repo");

        var out = await fetchJSON("/api/repos/import", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ repos: repos })
        });

        setStatus("repos-status", "Imported " + out.imported.length + " repos", "ok");
        await loadManagedRepos();
      } catch (err) {
        setStatus("repos-status", "Import failed: " + err.message, "error");
      } finally {
        btn.disabled = false;
      }
    }

    document.getElementById("settings-form").addEventListener("submit", onSaveSettings);
    document.getElementById("onboarding-form").addEventListener("submit", onCompleteOnboarding);
    document.getElementById("create-session-form").addEventListener("submit", onCreateSession);
    document.getElementById("followup-form").addEventListener("submit", onFollowUp);
    document.getElementById("discover-repos-btn").addEventListener("click", onDiscoverRepos);
    document.getElementById("import-repos-btn").addEventListener("click", onImportRepos);

    refreshSessions();
    loadSettings().catch(function (err) {
      setStatus("save-status", "Settings load failed: " + err.message, "error");
    });
    loadManagedRepos().catch(function (err) {
      setStatus("repos-status", "Managed repos load failed: " + err.message, "error");
    });
    setInterval(refreshSessions, 4000);
  </script>
</body>
</html>
`
