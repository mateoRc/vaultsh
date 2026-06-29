const status = document.querySelector("#status");
const form = document.querySelector("#command-form");
const command = document.querySelector("#command");
const output = document.querySelector("#output");
let sessionId = sessionStorage.getItem("vaultsh-session") || "";

fetch("/healthz")
  .then((response) => {
    status.textContent = response.ok ? "online" : "unavailable";
  })
  .catch(() => {
    status.textContent = "unavailable";
  });

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const line = command.value;
  if (!line.trim()) {
    command.focus();
    return;
  }

  command.value = "";
  command.focus();

  await execute(line);
});

document.addEventListener("keydown", async (event) => {
  if (event.ctrlKey && event.key.toLowerCase() === "l") {
    event.preventDefault();
    await execute("clear");
    command.focus();
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
