const status = document.querySelector("#status");
const atlasStatus = document.querySelector("#atlas-status");
const forgeStatus = document.querySelector("#forge-status");
const form = document.querySelector("#command-form");
const command = document.querySelector("#command");
const output = document.querySelector("#output");
let sessionId = sessionStorage.getItem("vaultsh-session") || "";

if (window.matchMedia("(pointer: fine)").matches) {
  focusCommand();
}

fetch("/healthz")
  .then((response) => {
    setStatus(response.ok ? "online" : "unavailable");
  })
  .catch(() => {
    setStatus("unavailable");
  });

refreshServiceStatus();
setInterval(refreshServiceStatus, 10000);

function setStatus(state) {
  status.textContent = state;
  status.dataset.state = state;
}

async function refreshServiceStatus() {
  try {
    const response = await fetch("/api/status");
    const services = await response.json();
    setServiceStatus(atlasStatus, services.atlas);
    setServiceStatus(forgeStatus, services.forge);
  } catch {
    setServiceStatus(atlasStatus, false);
    setServiceStatus(forgeStatus, false);
  }
}

function setServiceStatus(element, available) {
  element.dataset.state = available ? "online" : "unavailable";
  element.title = available ? "available" : "unavailable";
}

function focusCommand() {
  command.focus();
  const end = command.value.length;
  command.setSelectionRange(end, end);
}

form.addEventListener("click", (event) => {
  if (event.target !== command) {
    focusCommand();
  }
});

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const line = command.value;
  if (!line.trim()) {
    focusCommand();
    return;
  }

  command.value = "";
  focusCommand();

  await execute(line);
});

document.addEventListener("keydown", async (event) => {
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

async function complete() {
  const line = command.value;
  const cursor = command.selectionStart;

  try {
    const response = await fetch("/api/complete", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ line, cursor, session_id: sessionId }),
    });
    const result = await response.json();
    sessionId = result.session_id;
    sessionStorage.setItem("vaultsh-session", sessionId);

    if (!result.replacement) {
      return;
    }

    command.value =
      line.slice(0, result.start) +
      result.replacement +
      line.slice(result.end);
    const nextCursor = result.start + result.replacement.length;
    command.setSelectionRange(nextCursor, nextCursor);
  } catch {
    return;
  }
}

async function execute(line) {
  try {
    const response = await fetch("/api/exec", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ line, session_id: sessionId }),
    });
    const result = await response.json();
    sessionId = result.session_id;
    sessionStorage.setItem("vaultsh-session", sessionId);
    handleResult(line, result);
  } catch {
    appendEntry(line, "request failed");
  }
}

function handleResult(line, result) {
  if (result.action === "clear") {
    output.textContent = "";
    return;
  }

  const details = result.verbose ? `\n[verbose] ${result.verbose}` : "";
  appendEntry(line, `${result.output}${details}`);
}

function appendEntry(line, result) {
  const separator = output.textContent ? "\n" : "";
  output.textContent += `${separator}$ ${line}\n${result}\n`;
  output.scrollTop = output.scrollHeight;
}
