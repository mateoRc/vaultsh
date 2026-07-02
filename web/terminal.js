const status = document.querySelector("#status");
const atlasStatus = document.querySelector("#atlas-status");
const forgeStatus = document.querySelector("#forge-status");
const form = document.querySelector("#command-form");
const command = document.querySelector("#command");
const output = document.querySelector("#output");
const requestStatus = document.querySelector("#request-status");
const clearCommand = document.querySelector("#clear-command");
const quickCommandToggle = document.querySelector("#quick-command-toggle");
const quickCommands = document.querySelector("#quick-commands");
const actionButtons = document.querySelectorAll("[data-command]");
const submitButtons = document.querySelectorAll(
  "#run-command, #clear-command, .quick-commands button",
);
const maxOutputEntries = 100;
let outputEntries = [output.textContent];
let sessionId = sessionStorage.getItem("vaultsh-session") || "";
let running = false;

if (window.matchMedia("(pointer: fine)").matches) {
  focusCommand();
}

fetch("/healthz")
  .then((response) => setStatus(response.ok ? "online" : "unavailable"))
  .catch(() => setStatus("unavailable"));

refreshServiceStatus();
setInterval(refreshServiceStatus, 10000);

function setStatus(state) {
  status.textContent = state;
  status.dataset.state = state;
  status.setAttribute("aria-label", `Vaultsh ${state}`);
}

async function refreshServiceStatus() {
  try {
    const response = await fetch("/api/status");
    if (!response.ok) {
      throw new Error("status unavailable");
    }
    const services = await response.json();
    setServiceStatus(atlasStatus, "Atlas", services.atlas);
    setServiceStatus(forgeStatus, "Forge", services.forge);
  } catch {
    setServiceStatus(atlasStatus, "Atlas", false);
    setServiceStatus(forgeStatus, "Forge", false);
  }
}

function setServiceStatus(element, name, available) {
  const state = available ? "available" : "unavailable";
  element.dataset.state = available ? "online" : "unavailable";
  element.title = state;
  element.setAttribute("aria-label", `${name} ${state}`);
}

function focusCommand() {
  command.focus();
  const end = command.value.length;
  command.setSelectionRange(end, end);
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
    const response = await fetch("/api/complete", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ line, cursor, session_id: sessionId }),
    });
    if (!response.ok) {
      throw new Error(await friendlyError(response));
    }
    const result = await response.json();
    sessionId = result.session_id;
    sessionStorage.setItem("vaultsh-session", sessionId);

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
  setRequestStatus("Running…");
  const slowMessage = window.setTimeout(() => {
    setRequestStatus(
      line.startsWith("search ")
        ? "Still running—Atlas searches may take a moment…"
        : "Still running…",
    );
  }, 1200);

  try {
    const response = await fetch("/api/exec", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ line, session_id: sessionId }),
    });
    if (!response.ok) {
      throw new Error(await friendlyError(response));
    }
    const result = await response.json();
    sessionId = result.session_id;
    sessionStorage.setItem("vaultsh-session", sessionId);
    handleResult(line, result);
    setRequestStatus("");
  } catch (error) {
    appendEntry(line, error.message || "Request failed. Check your connection.");
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

function handleResult(line, result) {
  if (result.action === "clear") {
    outputEntries = [];
    renderOutput();
    return;
  }

  const details = result.verbose ? `\n[verbose] ${result.verbose}` : "";
  appendEntry(line, `${result.output}${details}`);
}

function appendEntry(line, result) {
  outputEntries.push(`$ ${line}\n${result}`);
  outputEntries = outputEntries.slice(-maxOutputEntries);
  renderOutput();
}

function renderOutput() {
  output.textContent = outputEntries.join("\n\n");
  output.scrollTop = output.scrollHeight;
}
