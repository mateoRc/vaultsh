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

  try {
    const response = await fetch("/api/exec", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ line }),
    });
    const result = await response.json();
    appendEntry(line, result.output);
  } catch {
    appendEntry(line, "request failed");
  }
});

function appendEntry(line, result) {
  const separator = output.textContent ? "\n" : "";
  output.textContent += `${separator}$ ${line}\n${result}\n`;
  output.scrollTop = output.scrollHeight;
}
