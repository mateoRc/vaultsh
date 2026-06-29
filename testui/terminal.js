const status = document.querySelector("#status");

fetch("/healthz")
  .then((response) => {
    status.textContent = response.ok ? "online" : "unavailable";
  })
  .catch(() => {
    status.textContent = "unavailable";
  });
