function initTerminal() {
  const status = document.querySelector("#status");
  const atlasStatus = document.querySelector("#atlas-status");
  const forgeStatus = document.querySelector("#forge-status");
  const form = document.querySelector("#command-form");
  const command = document.querySelector("#command");
  const prompt = document.querySelector("#prompt");
  const output = document.querySelector("#output");
  const requestStatus = document.querySelector("#request-status");
  const nextCommands = document.querySelector("#next-commands");
  const nextCommandButton = nextCommands
    ? nextCommands.querySelector("button")
    : null;
  const clearCommand = document.querySelector("#clear-command");
  const quickCommandToggle = document.querySelector("#quick-command-toggle");
  const quickCommands = document.querySelector("#quick-commands");
  const actionButtons = document.querySelectorAll("[data-command]");
  const submitButtons = document.querySelectorAll(
    "#run-command, #clear-command, .next-commands button, .quick-commands button",
  );
  if (
    !status ||
    !atlasStatus ||
    !forgeStatus ||
    !form ||
    !command ||
    !prompt ||
    !output ||
    !requestStatus ||
    !nextCommands ||
    !nextCommandButton ||
    !clearCommand ||
    !quickCommandToggle ||
    !quickCommands
  ) {
    return;
  }
  window.addEventListener("error", () => {
    status.textContent = "unavailable";
    status.dataset.state = "unavailable";
    atlasStatus.dataset.state = "unavailable";
    forgeStatus.dataset.state = "unavailable";
  });
  const endpoints = Object.freeze({
    health: new URL("/healthz", window.location.origin).toString(),
    status: new URL("/api/status", window.location.origin).toString(),
    execute: new URL("/api/exec", window.location.origin).toString(),
    complete: new URL("/api/complete", window.location.origin).toString(),
  });
  const storageKeys = Object.freeze({
    session: "vaultsh-session",
    currentDirectory: "vaultsh-current-directory",
    suggestionIndex: "vaultsh-suggestion-index",
  });
  const serviceStates = Object.freeze({
    online: "online",
    unavailable: "unavailable",
  });
  const terminalLinkPattern = /\[([^\]\n]+)\]\((https?:\/\/[^\s)]+|mailto:[^\s)]+)\)/g;
  const autoWelcomeCommand = "welcome";
  const autoWelcomeDelayMilliseconds = 400;
  const autoWelcomeKeystrokeMilliseconds = 100;
  const statusTimeoutMilliseconds = 3000;
  const statusRefreshMilliseconds = 10000;
  const suggestions = [
    ["About Mateo", "cat /cv/about.md"],
    ["Browse experience", "tree /cv/experience"],
    ["Review skills", 'cat /cv/skills.md | grep "Languages"'],
    ["Search distributed systems", "search distributed systems"],
    ["Browse technologies", "search Technology"],
    ["Find Java experience", 'search Java | grep "/cv/experience/"'],
    ["Show live dashboard", "dashboard"],
    ["Show deployment status", "deployments"],
    ["Show analytics", "metrics"],
    ["Contact Mateo", "contact"],
    ["Browse projects", "tree -L 2 /projects"],
    ["Find backend experience", 'grep -in "backend" /cv/about.md'],
    [
      "Review recent experience",
      "cat /cv/experience/reversinglabs.md | head -n 8",
    ],
    ["Show current role", "whoami"],
    ["Inspect the vault root", "ls -la /"],
    [
      "Review Vaultsh technologies",
      'cat /projects/vaultsh.md | grep "Technology"',
    ],
    ["Search languages", 'search "Languages"'],
    ["Explore available commands", "help"],
  ];
  const maxOutputEntries = 100;
  let outputEntries = [{ welcome: output.textContent }];
  let sessionId = storageGet(storageKeys.session) || "";
  let currentDirectory = sessionId
    ? storageGet(storageKeys.currentDirectory) || "/"
    : "/";
  let suggestionIndex = Number.parseInt(
    storageGet(storageKeys.suggestionIndex) || "0",
    10,
  );
  if (
    !Number.isInteger(suggestionIndex) ||
    suggestionIndex < 0 ||
    suggestionIndex >= suggestions.length
  ) {
    suggestionIndex = 0;
  }
  let running = false;

  updatePrompt();
  refreshVaultStatus();
  refreshServiceStatus();
  renderOutput();

  if (window.matchMedia("(pointer: fine)").matches) {
    focusCommand();
  }
  window.setTimeout(typeWelcomeCommand, autoWelcomeDelayMilliseconds);

  setInterval(refreshVaultStatus, statusRefreshMilliseconds);
  setInterval(refreshServiceStatus, statusRefreshMilliseconds);
  window.addEventListener("offline", setUnavailableStatus);
  window.addEventListener("online", () => {
    refreshVaultStatus();
    refreshServiceStatus();
  });

  function setStatus(state) {
    status.textContent = state;
    status.dataset.state = state;
    status.setAttribute("aria-label", `Vaultsh ${state}`);
  }

  function storageGet(key) {
    try {
      return window.sessionStorage.getItem(key);
    } catch {
      return "";
    }
  }

  function storageSet(key, value) {
    try {
      window.sessionStorage.setItem(key, value);
    } catch {
      // Storage can be unavailable in private or restricted browser contexts.
    }
  }

  async function refreshVaultStatus() {
    if (!navigator.onLine) {
      setUnavailableStatus();
      return;
    }

    const fallback = window.setTimeout(
      () => setStatus(serviceStates.unavailable),
      statusTimeoutMilliseconds,
    );
    try {
      const response = await fetchWithTimeout(endpoints.health, {
        cache: "no-store",
      });
      window.clearTimeout(fallback);
      setStatus(
        response.ok ? serviceStates.online : serviceStates.unavailable,
      );
    } catch {
      window.clearTimeout(fallback);
      setStatus(serviceStates.unavailable);
    }
  }

  async function refreshServiceStatus() {
    if (!navigator.onLine) {
      setUnavailableStatus();
      return;
    }

    const fallback = window.setTimeout(() => {
      setServiceStatus(atlasStatus, "Atlas", false);
      setServiceStatus(forgeStatus, "Forge", false);
    }, statusTimeoutMilliseconds);
    try {
      const response = await fetchWithTimeout(endpoints.status);
      if (!response.ok) {
        throw new Error("status unavailable");
      }
      const services = await response.json();
      window.clearTimeout(fallback);
      setServiceStatus(atlasStatus, "Atlas", services.atlas);
      setServiceStatus(forgeStatus, "Forge", services.forge);
    } catch {
      window.clearTimeout(fallback);
      setServiceStatus(atlasStatus, "Atlas", false);
      setServiceStatus(forgeStatus, "Forge", false);
    }
  }

  async function fetchWithTimeout(url, options = {}) {
    const controller = new AbortController();
    const timeout = window.setTimeout(
      () => controller.abort(),
      statusTimeoutMilliseconds,
    );

    try {
      return await fetch(url, { ...options, signal: controller.signal });
    } finally {
      window.clearTimeout(timeout);
    }
  }

  function setUnavailableStatus() {
    setStatus(serviceStates.unavailable);
    setServiceStatus(atlasStatus, "Atlas", false);
    setServiceStatus(forgeStatus, "Forge", false);
  }

  function setServiceStatus(element, name, available) {
    const state = available ? "available" : "unavailable";
    element.dataset.state = available
      ? serviceStates.online
      : serviceStates.unavailable;
    element.title = state;
    element.setAttribute("aria-label", `${name} ${state}`);
  }

  function focusCommand() {
    command.focus();
    const end = command.value.length;
    command.setSelectionRange(end, end);
  }

  function typeWelcomeCommand() {
    if (running || command.value !== "") {
      return;
    }

    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
      command.value = autoWelcomeCommand;
      submitCommandForm();
      return;
    }

    let nextIndex = 0;
    command.readOnly = true;
    const typing = window.setInterval(() => {
      nextIndex++;
      command.value = autoWelcomeCommand.slice(0, nextIndex);

      if (nextIndex >= autoWelcomeCommand.length) {
        window.clearInterval(typing);
        command.readOnly = false;
        submitCommandForm();
      }
    }, autoWelcomeKeystrokeMilliseconds);
  }

  function submitCommandForm() {
    form.dispatchEvent(
      new Event("submit", { bubbles: true, cancelable: true }),
    );
  }

  function shellPrompt(path = currentDirectory) {
    return `mateo@vault:${path}$`;
  }

  function updatePrompt() {
    prompt.textContent = shellPrompt();
  }

  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const line = command.value;
    if (!line.trim() || running) {
      focusCommand();
      return;
    }

    command.value = "";
    await execute(line);
    focusCommand();
  });

  clearCommand.addEventListener("click", () => execute("clear"));

  quickCommandToggle.addEventListener("click", () => {
    setQuickCommandsExpanded(quickCommands.hidden);
  });

  for (const button of actionButtons) {
    button.addEventListener("click", async () => {
      document.querySelector(".terminal").scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
      await execute(button.dataset.command);
      setQuickCommandsExpanded(false);
    });
  }

  document.addEventListener("keydown", async (event) => {
    if (event.key === "Escape" && !quickCommands.hidden) {
      setQuickCommandsExpanded(false);
      quickCommandToggle.focus();
      return;
    }

    if (event.ctrlKey && event.key.toLowerCase() === "l") {
      event.preventDefault();
      await execute("clear");
      focusCommand();
      return;
    }

    if (event.key === "Tab" && document.activeElement === command) {
      event.preventDefault();
      await complete();
    }
  });

  function setQuickCommandsExpanded(expanded) {
    quickCommands.hidden = !expanded;
    quickCommandToggle.setAttribute("aria-expanded", String(expanded));
  }

  async function complete() {
    const line = command.value;
    const cursor = command.selectionStart;
    setRequestStatus("Completing…");

    try {
      const response = await fetch(endpoints.complete, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ line, cursor, session_id: sessionId }),
      });
      if (!response.ok) {
        throw new Error(await friendlyError(response));
      }
      const result = await response.json();
      updateSession(result);

      if (result.replacement) {
        command.value =
          line.slice(0, result.start) +
          result.replacement +
          line.slice(result.end);
        const nextCursor = result.start + result.replacement.length;
        command.setSelectionRange(nextCursor, nextCursor);
      }
      setRequestStatus("");
    } catch (error) {
      setRequestStatus(error.message || "Completion unavailable.");
    }
  }

  async function execute(line) {
    if (running) {
      return;
    }

    setRunning(true);
    const submittedFrom = currentDirectory;
    setRequestStatus("Running…");
    const slowMessage = window.setTimeout(() => {
      setRequestStatus(
        line.startsWith("search ")
          ? "Still running—Atlas searches may take a moment…"
          : "Still running…",
      );
    }, 1200);

    try {
      const response = await fetch(endpoints.execute, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ line, session_id: sessionId }),
      });
      if (!response.ok) {
        throw new Error(await friendlyError(response));
      }
      const result = await response.json();
      handleResult(line, result, submittedFrom);
      setRequestStatus("");
    } catch (error) {
      appendEntry(
        line,
        error.message || "Request failed. Check your connection.",
        submittedFrom,
        1,
      );
      setRequestStatus("");
    } finally {
      window.clearTimeout(slowMessage);
      setRunning(false);
    }
  }

  async function friendlyError(response) {
    if (response.status === 429) {
      return "Too many requests—try again shortly.";
    }
    if (response.status === 413) {
      return "Command is too long.";
    }
    if (response.status === 503 || response.status >= 500) {
      return "Service unavailable—try again shortly.";
    }
    return "Request could not be completed.";
  }

  function setRunning(value) {
    running = value;
    command.disabled = value;
    for (const button of submitButtons) {
      button.disabled = value;
    }
  }

  function setRequestStatus(message) {
    requestStatus.textContent = message;
  }

  function handleResult(line, result, submittedFrom) {
    updateSession(result);

    if (result.action === "clear") {
      outputEntries = [];
      nextCommands.hidden = true;
      renderOutput();
      return;
    }

    const details = result.verbose ? `\n[verbose] ${result.verbose}` : "";
    appendEntry(line, `${result.output}${details}`, submittedFrom, result.exit_code);
    suggestNext();
  }

  function suggestNext() {
    const [label, commandLine] = suggestions[suggestionIndex];
    suggestionIndex = (suggestionIndex + 1) % suggestions.length;
    storageSet(storageKeys.suggestionIndex, String(suggestionIndex));

    nextCommandButton.textContent = label;
    nextCommandButton.dataset.command = commandLine;
    nextCommands.hidden = false;
  }

  function updateSession(result) {
    sessionId = result.session_id;
    currentDirectory = result.current_directory || currentDirectory;
    storageSet(storageKeys.session, sessionId);
    storageSet(storageKeys.currentDirectory, currentDirectory);
    updatePrompt();
  }

  function appendEntry(line, result, path, exitCode) {
    outputEntries.push({ line, result, path, exitCode });
    outputEntries = outputEntries.slice(-maxOutputEntries);
    renderOutput();
  }

  function renderOutput() {
    output.replaceChildren();
    for (const entry of outputEntries) {
      if (entry.welcome !== undefined) {
        appendTerminalText(output, entry.welcome);
        continue;
      }

      const container = document.createElement("div");
      container.className = "output-entry";

      const commandLine = document.createElement("div");
      commandLine.className = "output-command";
      const entryPrompt = document.createElement("span");
      entryPrompt.className = "output-prompt";
      entryPrompt.textContent = `${shellPrompt(entry.path)} `;
      commandLine.append(entryPrompt, document.createTextNode(entry.line));

      const result = document.createElement("div");
      result.className = "output-result";
      result.dataset.exitCode = String(entry.exitCode);
      appendTerminalText(result, entry.result);
      container.append(commandLine, result);
      output.append(container);
    }
    output.scrollTop = output.scrollHeight;
  }

  function appendTerminalText(parent, text) {
    let offset = 0;
    terminalLinkPattern.lastIndex = 0;
    let match = terminalLinkPattern.exec(text);

    while (match !== null) {
      parent.append(document.createTextNode(text.slice(offset, match.index)));
      parent.append(createTerminalLink(match[1], match[2]));
      offset = match.index + match[0].length;
      match = terminalLinkPattern.exec(text);
    }

    parent.append(document.createTextNode(text.slice(offset)));
  }

  function createTerminalLink(label, href) {
    const link = document.createElement("a");
    link.href = href;
    link.textContent = label;
    if (href.startsWith("http")) {
      link.target = "_blank";
      link.rel = "noreferrer";
    }
    return link;
  }
}

initTerminal();
