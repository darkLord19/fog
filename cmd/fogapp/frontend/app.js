(function () {
  var activeStates = { CREATED: true, SETUP: true, AI_RUNNING: true, VALIDATING: true, COMMITTED: true, PR_CREATED: true };
  var state = {
    apiBaseURL: "http://127.0.0.1:8080",
    settings: null,
    sessions: [],
    repos: [],
    discoveredRepos: [],
    selectedSessionID: "",
    selectedRunID: "",
    timelineSession: null,
    timelineRuns: [],
    timelineEvents: []
  };

  function $(id) { return document.getElementById(id); }

  function setStatus(id, message, cls) {
    var el = $(id);
    el.textContent = message || "";
    el.className = "status" + (cls ? (" " + cls) : "");
  }

  function escapeHTML(value) {
    return String(value || "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#39;");
  }

  function formatDate(value) {
    if (!value) return "-";
    var dt = new Date(value);
    if (isNaN(dt.getTime())) return "-";
    return dt.toLocaleString();
  }

  function findRepo(name) {
    return (state.repos || []).find(function (r) {
      return String(r.name || "").trim() === String(name || "").trim();
    }) || null;
  }

  function toWebRepoURL(raw) {
    raw = String(raw || "").trim();
    if (!raw) return "";
    if (raw.indexOf("git@github.com:") === 0) {
      raw = "https://github.com/" + raw.slice("git@github.com:".length);
    }
    if (raw.indexOf("https://github.com/") === 0) {
      return raw.replace(/\.git$/i, "");
    }
    return "";
  }

  async function resolveAPIBaseURL() {
    if (window.__FOG_API_BASE_URL__) {
      return String(window.__FOG_API_BASE_URL__);
    }
    try {
      var app = window.go && window.go.main && window.go.main.desktopApp;
      if (app && typeof app.APIBaseURL === "function") {
        var base = await app.APIBaseURL();
        if (base) return String(base);
      }
    } catch (_) {}
    return state.apiBaseURL;
  }

  async function resolveVersion() {
    try {
      var app = window.go && window.go.main && window.go.main.desktopApp;
      if (app && typeof app.Version === "function") {
        var v = await app.Version();
        if (v) return String(v);
      }
    } catch (_) {}
    return "-";
  }

  function openExternal(url) {
    var app = window.go && window.go.main && window.go.main.desktopApp;
    if (app && typeof app.OpenExternal === "function") {
      app.OpenExternal(url);
      return;
    }
    window.open(url, "_blank");
  }

  async function fetchJSON(path, options) {
    var url = state.apiBaseURL + path;
    var res = await fetch(url, options || {});
    if (!res.ok) {
      var text = await res.text();
      throw new Error(text || ("HTTP " + res.status));
    }
    if (res.status === 204) return null;
    return res.json();
  }

  function renderSessionSelect() {
    var select = $("followup-session");
    if (!state.sessions.length) {
      select.innerHTML = "<option value=''>No sessions</option>";
      return;
    }
    select.innerHTML = state.sessions.map(function (s) {
      var selected = s.id === state.selectedSessionID ? " selected" : "";
      return "<option value='" + escapeHTML(s.id) + "'" + selected + ">" +
        escapeHTML(s.repo_name + " :: " + s.branch) + "</option>";
    }).join("");
  }

  function renderSessions() {
    var list = $("sessions-list");
    var active = state.sessions.filter(function (s) { return activeStates[s.status] || s.busy; });
    $("kpi-active").textContent = active.length + " active";
    $("kpi-total").textContent = state.sessions.length + " total";

    if (!state.sessions.length) {
      list.innerHTML = "<div class='session-card'>No sessions yet.</div>";
      return;
    }

    list.innerHTML = state.sessions.slice(0, 80).map(function (s) {
      var stateClass = "state";
      if (s.status === "FAILED") stateClass += " state-failed";
      else if (activeStates[s.status]) stateClass += " state-active";
      var selected = s.id === state.selectedSessionID ? " selected" : "";
      return "<article class='session-card" + selected + "' data-session-id='" + escapeHTML(s.id) + "'>" +
        "<h4>" + escapeHTML(s.repo_name) + " / " + escapeHTML(s.branch) + "</h4>" +
        "<div class='session-meta'>" +
          "<span class='" + stateClass + "'>" + escapeHTML(s.status) + (s.busy ? "*" : "") + "</span>" +
          "<span>tool: " + escapeHTML(s.tool) + "</span>" +
          "<span>updated: " + escapeHTML(formatDate(s.updated_at)) + "</span>" +
        "</div>" +
      "</article>";
    }).join("");

    Array.prototype.slice.call(document.querySelectorAll(".session-card[data-session-id]")).forEach(function (el) {
      el.addEventListener("click", function () {
        var id = el.getAttribute("data-session-id");
        selectSession(id, "");
      });
    });
  }

  function renderRepos() {
    var repoSelect = $("new-repo");
    if (!state.repos.length) {
      repoSelect.innerHTML = "<option value=''>No repos imported</option>";
    } else {
      repoSelect.innerHTML = state.repos.map(function (r) {
        return "<option value='" + escapeHTML(r.name) + "'>" + escapeHTML(r.name) + "</option>";
      }).join("");
    }

    var managed = $("managed-list");
    if (!state.repos.length) {
      managed.innerHTML = "No managed repos.";
      return;
    }
    managed.innerHTML = state.repos.map(function (r) {
      return "<div>" +
        "<strong>" + escapeHTML(r.name) + "</strong>" +
        "<div style='color:#596575;font-size:12px'>" + escapeHTML(r.base_worktree_path || "-") + "</div>" +
      "</div>";
    }).join("");
  }

  function renderDiscoveredRepos() {
    var list = $("discover-list");
    if (!state.discoveredRepos.length) {
      list.innerHTML = "No discovered repos yet.";
      return;
    }
    list.innerHTML = state.discoveredRepos.map(function (r, idx) {
      return "<label class='repo-item'>" +
        "<input type='checkbox' id='repo-" + idx + "' data-repo='" + escapeHTML(r.full_name) + "'>" +
        "<span>" + escapeHTML(r.full_name) + "</span>" +
        "<span style='color:#596575;font-size:12px'>" + escapeHTML(r.default_branch || "-") + "</span>" +
      "</label>";
    }).join("");
  }

  function renderSettings() {
    var s = state.settings || {};
    var tools = s.available_tools || [];
    if (!tools.length && s.default_tool) tools = [s.default_tool];
    var options = tools.map(function (tool) {
      var selected = tool === s.default_tool ? " selected" : "";
      return "<option value='" + escapeHTML(tool) + "'" + selected + ">" + escapeHTML(tool) + "</option>";
    }).join("");
    $("settings-tool").innerHTML = options;

    var newTool = $("new-tool");
    newTool.innerHTML = "<option value=''>default</option>" + options;
    $("settings-prefix").value = s.branch_prefix || "fog";
    $("settings-pat-status").value = s.has_github_token ? "configured" : "missing";
    $("settings-onboarding").value = s.onboarding_required ? "required" : "complete";
  }

  function renderCloudStatus(cloud) {
    $("cloud-url").value = cloud.cloud_url || "";
    if (cloud.paired) {
      $("cloud-device").value = "paired (" + (cloud.device_id || "-") + ")";
    } else {
      $("cloud-device").value = "unpaired";
    }
  }

  function renderTimelineEmpty(message) {
    $("timeline-actions").innerHTML = "";
    $("timeline-summary").textContent = message || "Select a session to inspect its runs and events.";
    $("timeline-runs").innerHTML = "No runs loaded.";
    $("timeline-events").innerHTML = "No run events loaded.";
  }

  function renderTimelineSummary(session) {
    var repo = findRepo(session.repo_name);
    var base = repo && repo.default_branch ? repo.default_branch : "main";
    $("timeline-summary").innerHTML =
      "<strong>" + escapeHTML(session.repo_name) + " / " + escapeHTML(session.branch) + "</strong>" +
      "<div style='margin-top:4px;color:#596575;font-size:12px'>" +
      "tool: " + escapeHTML(session.tool) +
      " · base: " + escapeHTML(base) +
      " · updated: " + escapeHTML(formatDate(session.updated_at)) +
      "</div>";
  }

  function renderTimelineActions(session) {
    var actions = [];
    if (session.pr_url) {
      actions.push("<button type='button' class='ghost timeline-link-btn' data-url='" + escapeHTML(session.pr_url) + "'>Open PR</button>");
    }
    var repo = findRepo(session.repo_name);
    var repoWeb = repo ? toWebRepoURL(repo.url) : "";
    if (repoWeb) {
      var branchURL = repoWeb + "/tree/" + encodeURIComponent(session.branch || "");
      actions.push("<button type='button' class='ghost timeline-link-btn' data-url='" + escapeHTML(branchURL) + "'>Open Branch</button>");
      var base = repo.default_branch || "main";
      var compareURL = repoWeb + "/compare/" + encodeURIComponent(base) + "..." + encodeURIComponent(session.branch || "") + "?expand=1";
      actions.push("<button type='button' class='ghost timeline-link-btn' data-url='" + escapeHTML(compareURL) + "'>Compare</button>");
    }
    $("timeline-actions").innerHTML = actions.join("");
    Array.prototype.slice.call(document.querySelectorAll(".timeline-link-btn")).forEach(function (btn) {
      btn.addEventListener("click", function () {
        var url = btn.getAttribute("data-url");
        if (url) openExternal(url);
      });
    });
  }

  function renderTimelineRuns(runs) {
    var body = $("timeline-runs");
    if (!runs.length) {
      body.innerHTML = "No runs found for this session.";
      return;
    }
    body.innerHTML = runs.map(function (run) {
      var selected = run.id === state.selectedRunID ? " selected" : "";
      return "<div class='timeline-run" + selected + "'>" +
        "<div>" +
          "<div><strong>" + escapeHTML(run.state) + "</strong></div>" +
          "<div style='font-size:11px;color:#596575'>" + escapeHTML(formatDate(run.updated_at || run.created_at)) + "</div>" +
        "</div>" +
        "<button type='button' class='ghost timeline-run-btn' data-run-id='" + escapeHTML(run.id) + "'>View Events</button>" +
      "</div>";
    }).join("");
    Array.prototype.slice.call(document.querySelectorAll(".timeline-run-btn")).forEach(function (btn) {
      btn.addEventListener("click", function () {
        var runID = btn.getAttribute("data-run-id");
        state.selectedRunID = runID;
        loadRunEvents().catch(function (err) {
          $("timeline-events").innerHTML = "Failed to load run events: " + escapeHTML(err.message);
        });
        renderTimelineRuns(state.timelineRuns);
      });
    });
  }

  function renderTimelineEvents(events) {
    var body = $("timeline-events");
    if (!events.length) {
      body.innerHTML = "No events for this run.";
      return;
    }
    body.innerHTML = events.map(function (ev) {
      return "<div class='timeline-event'>" +
        "<div class='timeline-event-head'>" +
          "<span>" + escapeHTML(ev.type || "-") + "</span>" +
          "<span>" + escapeHTML(formatDate(ev.ts)) + "</span>" +
        "</div>" +
        "<div class='timeline-event-body'>" + escapeHTML(ev.message || ev.data || "-") + "</div>" +
      "</div>";
    }).join("");
  }

  async function loadTimeline() {
    if (!state.selectedSessionID) {
      renderTimelineEmpty("Select a session to inspect its runs and events.");
      return;
    }
    var detail = await fetchJSON("/api/sessions/" + encodeURIComponent(state.selectedSessionID));
    state.timelineSession = detail && detail.session ? detail.session : null;
    state.timelineRuns = detail && detail.runs ? detail.runs : [];
    if (!state.timelineSession) {
      renderTimelineEmpty("Session not found.");
      return;
    }
    if (!state.selectedRunID && state.timelineRuns.length) {
      state.selectedRunID = state.timelineRuns[0].id;
    }
    var hasRun = state.timelineRuns.some(function (r) { return r.id === state.selectedRunID; });
    if (!hasRun) {
      state.selectedRunID = state.timelineRuns.length ? state.timelineRuns[0].id : "";
    }

    renderTimelineSummary(state.timelineSession);
    renderTimelineActions(state.timelineSession);
    renderTimelineRuns(state.timelineRuns);
    await loadRunEvents();
  }

  async function loadRunEvents() {
    if (!state.selectedSessionID || !state.selectedRunID) {
      renderTimelineEvents([]);
      return;
    }
    var events = await fetchJSON(
      "/api/sessions/" + encodeURIComponent(state.selectedSessionID) +
      "/runs/" + encodeURIComponent(state.selectedRunID) + "/events?limit=200"
    );
    state.timelineEvents = events || [];
    renderTimelineEvents(state.timelineEvents);
  }

  async function selectSession(sessionID, runID) {
    state.selectedSessionID = String(sessionID || "").trim();
    state.selectedRunID = String(runID || "").trim();
    renderSessionSelect();
    renderSessions();
    await loadTimeline();
  }

  async function loadAll() {
    $("daemon-badge").textContent = "fogd: connected";
    var result = await Promise.all([
      fetchJSON("/api/settings"),
      fetchJSON("/api/sessions"),
      fetchJSON("/api/repos"),
      fetchJSON("/api/cloud")
    ]);
    state.settings = result[0] || {};
    state.sessions = result[1] || [];
    state.repos = result[2] || [];
    renderSettings();
    renderSessions();
    renderSessionSelect();
    renderRepos();
    renderCloudStatus(result[3] || {});

    if (!state.selectedSessionID && state.sessions.length) {
      await selectSession(state.sessions[0].id, "");
    } else {
      await loadTimeline();
    }
  }

  async function onCreateSession(event) {
    event.preventDefault();
    var btn = $("new-submit");
    setStatus("new-status", "", "");
    btn.disabled = true;
    try {
      var payload = {
        repo: $("new-repo").value,
        prompt: $("new-prompt").value.trim(),
        model: $("new-model").value.trim(),
        branch_name: $("new-branch").value.trim(),
        autopr: $("new-autopr").value === "true",
        async: true
      };
      var tool = $("new-tool").value;
      if (tool) payload.tool = tool;
      if (!payload.repo) throw new Error("Repository is required");
      if (!payload.prompt) throw new Error("Prompt is required");

      var out = await fetchJSON("/api/sessions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
      setStatus("new-status", "Queued session " + out.session_id, "ok");
      $("new-prompt").value = "";
      $("new-branch").value = "";
      await refreshSessions();
      if (out.session_id) {
        await selectSession(out.session_id, "");
      }
    } catch (err) {
      setStatus("new-status", "Create failed: " + err.message, "error");
    } finally {
      btn.disabled = false;
    }
  }

  async function onFollowup(event) {
    event.preventDefault();
    var btn = $("followup-submit");
    setStatus("followup-status", "", "");
    btn.disabled = true;
    try {
      var sessionID = $("followup-session").value;
      var prompt = $("followup-prompt").value.trim();
      if (!sessionID) throw new Error("Session is required");
      if (!prompt) throw new Error("Prompt is required");
      var out = await fetchJSON("/api/sessions/" + encodeURIComponent(sessionID) + "/runs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ prompt: prompt, async: true })
      });
      setStatus("followup-status", "Queued run " + out.run_id, "ok");
      $("followup-prompt").value = "";
      await refreshSessions();
      await selectSession(sessionID, out.run_id || "");
    } catch (err) {
      setStatus("followup-status", "Follow-up failed: " + err.message, "error");
    } finally {
      btn.disabled = false;
    }
  }

  async function onSaveSettings(event) {
    event.preventDefault();
    var btn = $("settings-submit");
    setStatus("settings-status", "", "");
    btn.disabled = true;
    try {
      var payload = {
        default_tool: $("settings-tool").value,
        branch_prefix: $("settings-prefix").value.trim()
      };
      var pat = $("settings-pat").value.trim();
      if (pat) payload.github_pat = pat;

      await fetchJSON("/api/settings", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
      setStatus("settings-status", "Saved", "ok");
      $("settings-pat").value = "";
      state.settings = await fetchJSON("/api/settings");
      renderSettings();
    } catch (err) {
      setStatus("settings-status", "Save failed: " + err.message, "error");
    } finally {
      btn.disabled = false;
    }
  }

  async function onDiscoverRepos() {
    setStatus("repos-status", "", "");
    var btn = $("discover-btn");
    btn.disabled = true;
    try {
      state.discoveredRepos = await fetchJSON("/api/repos/discover", { method: "POST" }) || [];
      renderDiscoveredRepos();
      setStatus("repos-status", "Discovered " + state.discoveredRepos.length + " repos", "ok");
    } catch (err) {
      setStatus("repos-status", "Discover failed: " + err.message, "error");
    } finally {
      btn.disabled = false;
    }
  }

  async function onImportRepos() {
    setStatus("repos-status", "", "");
    var btn = $("import-btn");
    btn.disabled = true;
    try {
      var checked = Array.prototype.slice.call(document.querySelectorAll("#discover-list input[type='checkbox']:checked"));
      var repos = checked.map(function (el) { return el.getAttribute("data-repo"); });
      if (!repos.length) throw new Error("Select at least one repo");

      var out = await fetchJSON("/api/repos/import", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ repos: repos })
      });
      setStatus("repos-status", "Imported " + out.imported.length + " repos", "ok");
      state.repos = await fetchJSON("/api/repos");
      renderRepos();
    } catch (err) {
      setStatus("repos-status", "Import failed: " + err.message, "error");
    } finally {
      btn.disabled = false;
    }
  }

  async function onSaveCloudURL() {
    setStatus("cloud-status", "", "");
    try {
      var url = $("cloud-url").value.trim();
      if (!url) throw new Error("Cloud URL is required");
      var cloud = await fetchJSON("/api/cloud", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ cloud_url: url })
      });
      renderCloudStatus(cloud);
      setStatus("cloud-status", "Cloud URL saved", "ok");
    } catch (err) {
      setStatus("cloud-status", "Save failed: " + err.message, "error");
    }
  }

  async function onPairCloud() {
    setStatus("cloud-status", "", "");
    try {
      var code = $("cloud-code").value.trim();
      if (!code) throw new Error("Pair code is required");
      var cloud = await fetchJSON("/api/cloud/pair", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code: code })
      });
      $("cloud-code").value = "";
      renderCloudStatus(cloud);
      setStatus("cloud-status", "Pairing successful", "ok");
    } catch (err) {
      setStatus("cloud-status", "Pairing failed: " + err.message, "error");
    }
  }

  async function onUnpairCloud() {
    setStatus("cloud-status", "", "");
    try {
      var team = $("cloud-team").value.trim();
      var user = $("cloud-user").value.trim();
      if (!team || !user) throw new Error("Team ID and Slack User ID are required");
      var cloud = await fetchJSON("/api/cloud/unpair", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ team_id: team, slack_user_id: user })
      });
      renderCloudStatus(cloud);
      setStatus("cloud-status", "Unpaired", "ok");
    } catch (err) {
      setStatus("cloud-status", "Unpair failed: " + err.message, "error");
    }
  }

  async function refreshSessions() {
    state.sessions = await fetchJSON("/api/sessions");
    if (state.selectedSessionID) {
      var stillExists = state.sessions.some(function (s) { return s.id === state.selectedSessionID; });
      if (!stillExists) {
        state.selectedSessionID = state.sessions.length ? state.sessions[0].id : "";
        state.selectedRunID = "";
      }
    }
    renderSessions();
    renderSessionSelect();
    if (state.selectedSessionID) {
      await loadTimeline();
    } else {
      renderTimelineEmpty("No active sessions.");
    }
  }

  async function bootstrap() {
    $("daemon-badge").textContent = "fogd: connecting";
    state.apiBaseURL = await resolveAPIBaseURL();
    var version = await resolveVersion();
    $("version-badge").textContent = "version: " + version;

    await loadAll();
    setInterval(function () {
      refreshSessions().catch(function (err) {
        setStatus("followup-status", "Refresh failed: " + err.message, "error");
      });
    }, 4000);
  }

  $("new-session-form").addEventListener("submit", onCreateSession);
  $("followup-form").addEventListener("submit", onFollowup);
  $("settings-form").addEventListener("submit", onSaveSettings);
  $("discover-btn").addEventListener("click", onDiscoverRepos);
  $("import-btn").addEventListener("click", onImportRepos);
  $("cloud-save").addEventListener("click", onSaveCloudURL);
  $("cloud-pair").addEventListener("click", onPairCloud);
  $("cloud-unpair").addEventListener("click", onUnpairCloud);

  bootstrap().catch(function (err) {
    $("daemon-badge").textContent = "fogd: unavailable";
    setStatus("new-status", "Initialization failed: " + err.message, "error");
    renderTimelineEmpty("Initialization failed.");
  });
})();
