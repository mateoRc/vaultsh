const status = document.querySelector("#status");
const form = document.querySelector("#command-form");
const command = document.querySelector("#command");
const output = document.querySelector("#output");
const terminal = document.querySelector(".terminal");
let sessionId = sessionStorage.getItem("vaultsh-session") || "";

if (window.matchMedia("(pointer: fine)").matches) {
  focusCommand();
}

terminal.addEventListener("click", (event) => {
  if (event.target === command) {
    return;
  }
  if (event.target.closest("a, button")) {
    return;
  }
  if (window.getSelection()?.toString()) {
    return;
  }
  focusCommand();
});

fetch("/healthz")
  .then((response) => {
    setStatus(response.ok ? "online" : "unavailable");
  })
  .catch(() => {
    setStatus("unavailable");
  });

function setStatus(state) {
  status.textContent = state;
  status.dataset.state = state;
}

function focusCommand() {
  command.focus();
  const end = command.value.length;
  command.setSelectionRange(end, end);
}

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

  appendEntry(line, result.output);
}

function appendEntry(line, result) {
  const separator = output.textContent ? "\n" : "";
  output.textContent += `${separator}$ ${line}\n${result}\n`;
  output.scrollTop = output.scrollHeight;
}
