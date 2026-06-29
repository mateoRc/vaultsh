const status = document.querySelector("#status");
const form = document.querySelector("#command-form");
const command = document.querySelector("#command");
const output = document.querySelector("#output");

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
  }
});

async function execute(line) {
  try {
    const response = await fetch("/api/exec", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ line }),
    });
    const result = await response.json();
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
