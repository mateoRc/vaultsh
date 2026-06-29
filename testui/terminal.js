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

  try {
    const response = await fetch("/api/exec", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ line: command.value }),
    });
    const result = await response.json();
    output.textContent = result.output;
  } catch {
    output.textContent = "request failed";
  }
});
