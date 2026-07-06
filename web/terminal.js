const status = document.querySelector("#status");
const atlasStatus = document.querySelector("#atlas-status");
const forgeStatus = document.querySelector("#forge-status");
const form = document.querySelector("#command-form");
const command = document.querySelector("#command");
const prompt = document.querySelector("#prompt");
const output = document.querySelector("#output");
const requestStatus = document.querySelector("#request-status");
const nextCommands = document.querySelector("#next-commands");
const nextCommandButton = nextCommands.querySelector("button");
const clearCommand = document.querySelector("#clear-command");
const quickCommandToggle = document.querySelector("#quick-command-toggle");
const quickCommands = document.querySelector("#quick-commands");
const actionButtons = document.querySelectorAll("[data-command]");
const submitButtons = document.querySelectorAll(
  "#run-command, #clear-command, .next-commands button, .quick-commands button",
);
const suggestions = [
  ["About Mateo", "cat /cv/about.md"],
  ["Browse experience", "tree /cv/experience"],
  ["Review skills", 'cat /cv/skills.md | grep "Languages"'],
  ["Search distributed systems", "search distributed systems"],
  ["Inspect project stack", 'search "Technology" | grep "/projects/"'],
  ["Show live dashboard", "dashboard"],
  ["Show deployment status", "deployments"],
  ["Show analytics", "metrics"],
  ["Browse projects", "tree -L 2 /projects"],
  ["Find backend experience", 'grep -in "backend" /cv/about.md'],
  [
    "Review recent experience",
    "cat /cv/experience/reversinglabs.md | head -n 8",
  ],
  ["Show current role", "whoami"],
  ["Inspect the vault root", "ls -la /"],
  ["Show recent commands", "history | tail -n 5"],
  ["Count documented skills", "wc /cv/skills.md"],
  [
    "Review Vaultsh technologies",
    'cat /projects/vaultsh.md | grep "Technology"',
  ],
  ["Search languages", 'search "Languages"'],
  ["Explore available commands", "help"],
];
const maxOutputEntries = 100;
let outputEntries = [{ welcome: output.textContent }];
let sessionId = sessionStorage.getItem("vaultsh-session") || "";
let currentDirectory = sessionId
  ? sessionStorage.getItem("vaultsh-current-directory") || "/"
  : "/";
let suggestionIndex = Number.parseInt(
  sessionStorage.getItem("vaultsh-suggestion-index") || "0",
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
renderOutput();

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
    currentDirectory = result.current_directory || currentDirectory;
    sessionStorage.setItem("vaultsh-current-directory", currentDirectory);
    updatePrompt();

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
  currentDirectory = result.current_directory || currentDirectory;
  sessionStorage.setItem("vaultsh-current-directory", currentDirectory);
  updatePrompt();

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
  sessionStorage.setItem("vaultsh-suggestion-index", String(suggestionIndex));

  nextCommandButton.textContent = label;
  nextCommandButton.dataset.command = commandLine;
  nextCommands.hidden = false;
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
      output.append(document.createTextNode(entry.welcome));
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
    result.textContent = entry.result;
    container.append(commandLine, result);
    output.append(container);
  }
  output.scrollTop = output.scrollHeight;
}
