"""Local browser smoke test with a read-only content API fixture.

Run: python tests/browser_smoke.py
Requires a local Chrome or Edge binary; does not start the Go service.
"""

import json
import mimetypes
import os
from pathlib import Path
import re
import subprocess
import tempfile
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import parse_qs, urlsplit


ROOT = Path(__file__).resolve().parents[1]
WEB = ROOT / "web"
CONTENT = ROOT.parent / "lab" / "content"
EXCERPTS = (
    ("head -n 11 /projects/overview.md", "## Components"),
    ("head -n 26 /docs/architecture/overview.md", "## Runtime topology"),
)


def verify_excerpts():
    sources = (
        WEB / "index.html",
        WEB / "terminal.js",
        ROOT / "internal" / "command" / "welcome.go",
        CONTENT / "README.md",
    )
    for command, section in EXCERPTS:
        count, path = re.fullmatch(r"head -n (\d+) (/[\w/.-]+\.md)", command).groups()
        lines = (CONTENT / path.lstrip("/")).read_text(encoding="utf-8").splitlines()
        excerpt = lines[:int(count)]
        assert len(excerpt) < len(lines) and lines[int(count)] == "", command
        assert section in excerpt and excerpt.count("```mermaid") == excerpt.count("```"), command
        for source in sources:
            assert command in source.read_text(encoding="utf-8"), f"{command} missing from {source}"
    print("Command excerpts: PASS complete sections and matching suggestions")

SMOKE = r"""
<pre id="smoke-results">RUNNING</pre><script>
(async () => {
  const results = [];
  const check = (condition, label) => { if (!condition) throw Error(label); results.push(label); };
  const wait = async (test) => {
    for (let i = 0; i < 80; i++) {
      if (test()) return;
      await new Promise(resolve => setTimeout(resolve, 50));
    }
    throw Error('timed out waiting for UI: ' + document.querySelector('#output').textContent.slice(-600));
  };
  try {
    await wait(() => document.querySelector('#output').textContent.includes('Explore projects'));
    check(!document.querySelector('#shell-panel').hidden, 'terminal default');
    check(document.querySelector('.workspace > .workspace-header .mode-switch') &&
          getComputedStyle(document.querySelector('.workspace-header')).backgroundColor === 'rgb(247, 200, 67)' &&
          document.querySelector('.intro-hint').textContent.includes('Explorer tab'), 'tabs inside window header');
    check(document.activeElement.id !== 'command', 'no automatic input focus');
    const examples = [...document.querySelectorAll('#output .terminal-example')];
    check(examples.length === 3 && examples.every(button => button.tagName === 'BUTTON'), 'clickable welcome examples');
    check(examples[0].textContent === 'head -n 11 /projects/overview.md' &&
          examples[1].textContent === 'head -n 26 /docs/architecture/overview.md', 'welcome excerpts');
    check([...document.querySelectorAll('#quick-commands [data-command]')].slice(0, 2)
      .map(button => button.dataset.command).join('|') === examples.slice(0, 2).map(button => button.textContent).join('|'), 'shortcut consistency');
    examples[0].click();
    await wait(() => [...document.querySelectorAll('#output .output-result')].at(-1)?.textContent.includes('**Lab** owns Docker Compose'));
    check(![...document.querySelectorAll('#output .output-result')].at(-1).textContent.includes('## Design decisions'), 'project excerpt boundary');
    document.querySelectorAll('#output .terminal-example')[1].click();
    await wait(() => [...document.querySelectorAll('#output .output-result')].at(-1)?.textContent.includes('inside the private Compose networks.'));
    check(![...document.querySelectorAll('#output .output-result')].at(-1).textContent.includes('## Component ownership'), 'architecture excerpt boundary');
    document.querySelector('#gui-tab').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/projects');
    await wait(() => document.querySelector('[data-gui-path="/docs"]'));
    check(document.querySelector('#gui-tree-location').textContent.includes('/projects') &&
          document.querySelector('[data-gui-path="/projects"]').getAttribute('aria-current') === 'page', 'project tree and current folder');
    check(document.querySelector('.gui-item strong').textContent === 'Overview', 'overview first');
    check(getComputedStyle(document.querySelector('.desktop-main')).display === 'grid', 'desktop layout');
    document.querySelector('.gui-item').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/projects/overview.md');
    check(document.querySelector('[data-gui-path="/projects/overview.md"]').getAttribute('aria-current') === 'page' &&
          document.querySelector('#gui-tree-location').textContent.includes('/projects/overview.md'), 'current file in tree');
    check(location.search.includes('overview.md'), 'shareable URL');
    [...document.querySelectorAll('.markdown-body a')].find(a => a.textContent === 'Architecture').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/docs/architecture/overview.md');
    check(document.querySelector('.markdown-body h1').textContent === 'Backend Lab Architecture', 'relative link');
    history.back();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/projects/overview.md');
    check(document.querySelector('.markdown-body h1').textContent.includes('Backend Lab'), 'browser Back');
    history.forward();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/docs/architecture/overview.md');
    check(document.querySelector('#gui-tab').getAttribute('aria-selected') === 'true', 'browser Forward and tab state');
    document.querySelector('[data-gui-path="/docs"]').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/docs');
    await wait(() => document.querySelector('[data-gui-path="/docs/architecture"]'));
    check(document.querySelector('[data-gui-path="/docs"]').getAttribute('aria-current') === 'page', 'tree folder navigation');
    document.querySelector('[data-gui-path="/docs"]').previousElementSibling.click();
    check(!document.querySelector('[data-gui-path="/docs/architecture"]'), 'tree folder collapse');
    document.querySelector('[data-gui-path="/docs"]').previousElementSibling.click();
    await wait(() => document.querySelector('[data-gui-path="/docs/architecture"]'));
    document.querySelector('[data-gui-path="/docs/architecture"]').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/docs/architecture');
    await wait(() => document.querySelector('[data-gui-path="/docs/architecture/vaultsh.md"]'));
    document.querySelector('[data-gui-path="/docs/architecture/vaultsh.md"]').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/docs/architecture/vaultsh.md');
    [...document.querySelectorAll('.markdown-body a')].find(a => a.textContent === 'browse API').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/docs/api.md');
    check(location.hash === '#browse-content' && document.activeElement.id === 'browse-content', 'relative fragment');
    await fetch('/api/fail-next');
    document.querySelector('[data-gui-path="/projects"]').click();
    await wait(() => !document.querySelector('#gui-retry').hidden);
    check(document.querySelector('#gui-message').textContent.includes('unavailable'), 'browse error');
    document.querySelector('#gui-retry').click();
    await wait(() => document.querySelector('#gui-breadcrumbs').dataset.path === '/projects');
    check(document.querySelector('#gui-retry').hidden, 'browse retry');
    document.querySelector('#shell-tab').click();
    const input = document.querySelector('#command');
    input.focus(); input.value = 'cat /projects/ov'; input.setSelectionRange(input.value.length, input.value.length);
    const tab = () => { const event = new KeyboardEvent('keydown', {key: 'Tab', bubbles: true, cancelable: true}); input.dispatchEvent(event); return event.defaultPrevented; };
    check(tab(), 'first Tab completes');
    await wait(() => input.value.includes('overview.md'));
    check(!tab(), 'second Tab leaves input');
    const shift = new KeyboardEvent('keydown', {key: 'Tab', shiftKey: true, bubbles: true, cancelable: true});
    input.dispatchEvent(shift);
    check(!shift.defaultPrevented, 'Shift+Tab moves normally');
    input.dispatchEvent(new KeyboardEvent('keydown', {key: 'Escape', bubbles: true}));
    check(document.activeElement.id === 'quick-command-toggle', 'Escape leaves input');
    await fetch('/api/fail-next-complete');
    input.focus(); input.value = 'cat /projects/ov2'; input.dispatchEvent(new Event('input'));
    check(tab(), 'completion request starts');
    await wait(() => document.querySelector('#request-status').textContent.includes('Edit to retry'));
    check(!tab(), 'completion failure permits exit');
    await fetch('/api/fail-next-exec');
    input.value = 'help';
    document.querySelector('#command-form').dispatchEvent(new Event('submit', {bubbles: true, cancelable: true}));
    await wait(() => document.querySelector('#output').textContent.includes('Service unavailable'));
    check(!input.disabled, 'execution failure unlocks input');
    await wait(() => document.activeElement === input);
    check(document.activeElement === input, 'input focus after failed command');
    input.value = 'help';
    document.querySelector('#command-form').dispatchEvent(new Event('submit', {bubbles: true, cancelable: true}));
    await wait(() => [...document.querySelectorAll('#output .output-result')].at(-1)?.textContent === 'ok');
    check(!input.disabled, 'manual command recovery');
    await wait(() => document.activeElement === input);
    check(document.activeElement === input, 'input focus after successful command');
    input.value = 'slow';
    document.querySelector('#command-form').dispatchEvent(new Event('submit', {bubbles: true, cancelable: true}));
    document.querySelector('#quick-command-toggle').focus();
    await wait(() => !input.disabled && [...document.querySelectorAll('#output .output-command')].at(-1)?.textContent.includes('slow'));
    check(document.activeElement.id === 'quick-command-toggle', 'deliberate focus move preserved');
    document.querySelector('#smoke-results').textContent = 'PASS ' + results.join(', ');
  } catch (error) {
    document.querySelector('#smoke-results').textContent = 'FAIL ' + error.message + ' after ' + results.join(', ');
  }
})();
</script>
"""

DEEP = r"""
<pre id="smoke-results">RUNNING</pre><script>
(async () => {
  const result = document.querySelector('#smoke-results');
  for (let i = 0; i < 80 && document.querySelector('#gui-breadcrumbs').dataset.path !== '/cv/experience.md'; i++)
    await new Promise(resolve => setTimeout(resolve, 50));
  if (document.querySelector('#gui-breadcrumbs').dataset.path !== '/cv/experience.md' ||
      document.querySelector('#gui-panel').hidden ||
      !document.querySelector('#gui-document').textContent.includes('ReversingLabs')) {
    result.textContent = 'FAIL direct document load'; return;
  }
  if (!matchMedia('(max-width: 36rem)').matches ||
      getComputedStyle(document.querySelector('.desktop-main')).display !== 'block' ||
      document.documentElement.scrollWidth > innerWidth + 1 ||
      document.activeElement.id === 'command') {
    result.textContent = 'FAIL mobile layout'; return;
  }
  for (let i = 0; i < 80 && !document.querySelector('[data-gui-path="/cv/experience.md"]'); i++)
    await new Promise(resolve => setTimeout(resolve, 50));
  if (!document.querySelector('#gui-tree-location').textContent.includes('/cv/experience.md') ||
      !document.querySelector('[data-gui-path="/cv/experience.md"]') ||
      getComputedStyle(document.querySelector('#gui-tree-toggle')).display === 'none') {
    result.textContent = 'FAIL mobile project tree'; return;
  }
  document.querySelector('#gui-tree-toggle').click();
  if (document.querySelector('#gui-tree-toggle').getAttribute('aria-expanded') !== 'false' ||
      getComputedStyle(document.querySelector('#gui-tree')).display !== 'none') {
    result.textContent = 'FAIL mobile tree toggle'; return;
  }
  if (!sessionStorage.getItem('smoke-reloaded')) {
    sessionStorage.setItem('smoke-reloaded', 'yes'); location.reload(); return;
  }
  result.textContent = 'PASS deep link, refresh, and mobile layout ' + innerWidth + 'px viewport';
})();
</script>
"""

MOBILE_START = r"""
<pre id="smoke-results">RUNNING</pre><script>
(async () => {
  const result = document.querySelector('#smoke-results');
  for (let i = 0; i < 80 && !document.querySelector('#output').textContent.includes('Explore projects'); i++)
    await new Promise(resolve => setTimeout(resolve, 50));
  result.textContent = !document.querySelector('#shell-panel').hidden &&
    document.activeElement.id !== 'command' &&
    matchMedia('(max-width: 36rem)').matches
      ? 'PASS mobile terminal starts without input focus'
      : 'FAIL mobile terminal initial focus';
})();
</script>
"""

DIAGRAM = r"""
<pre id="smoke-results">RUNNING</pre><script>
(async () => {
  const result = document.querySelector('#smoke-results');
  for (let i = 0; i < 400 && (document.querySelector('.mermaid--compact svg')?.querySelectorAll('.node').length ?? 0) < 7; i++)
    await new Promise(resolve => setTimeout(resolve, 50));
  const svg = document.querySelector('.mermaid--compact svg');
  const wrapper = document.querySelector('.mermaid--compact');
  if (!svg || !wrapper || svg.getBoundingClientRect().width > wrapper.clientWidth + 1 ||
      svg.getBoundingClientRect().height < 100 ||
      svg.querySelectorAll('.node').length < 7) {
    result.textContent = 'FAIL command flow diagram: svg=' + !!svg + ', wrapper=' + !!wrapper +
      ', svgWidth=' + (svg?.getBoundingClientRect().width ?? -1) + ', wrapperWidth=' + (wrapper?.clientWidth ?? -1) +
      ', svgHeight=' + (svg?.getBoundingClientRect().height ?? -1) +
      ', nodes=' + (svg?.querySelectorAll('.node').length ?? -1); return;
  }
  result.textContent = 'PASS compact command flow diagram';
})();
</script>
"""


class Handler(BaseHTTPRequestHandler):
    fail_next = False
    fail_complete = False
    fail_exec = False

    def log_message(self, *_args):
        pass

    def do_GET(self):
        url = urlsplit(self.path)
        if url.path == "/api/fail-next":
            Handler.fail_next = True
            return self.send_bytes(b"ok", "text/plain")
        if url.path == "/api/fail-next-complete":
            Handler.fail_complete = True
            return self.send_bytes(b"ok", "text/plain")
        if url.path == "/api/fail-next-exec":
            Handler.fail_exec = True
            return self.send_bytes(b"ok", "text/plain")
        if url.path == "/healthz":
            return self.send_bytes(b"ok", "text/plain")
        if url.path == "/api/status":
            return self.send_json({"atlas": True, "forge": True})
        if url.path == "/api/browse":
            if Handler.fail_next:
                Handler.fail_next = False
                return self.send_error(503)
            path = parse_qs(url.query).get("path", ["/"])[0]
            target = (CONTENT / path.lstrip("/")).resolve()
            if not target.is_relative_to(CONTENT.resolve()) or not target.exists() or target.name.startswith("."):
                return self.send_error(404)
            if target.is_dir():
                children = [{"name": child.name, "kind": "directory" if child.is_dir() else "file", "path": (path.rstrip("/") + "/" + child.name)}
                            for child in target.iterdir() if not child.name.startswith(".")]
                return self.send_json({"path": path, "kind": "directory", "children": children})
            return self.send_json({"path": path, "kind": "file", "content": target.read_text(encoding="utf-8")})
        if url.path in ("/vault/", "/vault/index.html"):
            source = (WEB / "index.html").read_text(encoding="utf-8")
            smoke = parse_qs(url.query).get("smoke", [""])[0]
            if smoke:
                source = source.replace("</body>", ({"deep": DEEP, "mobile": MOBILE_START, "diagram": DIAGRAM}.get(smoke, SMOKE)) + "</body>")
            return self.send_bytes(source.encode(), "text/html")
        if url.path.startswith("/vault/"):
            path = (WEB / url.path.removeprefix("/vault/")).resolve()
            if path.is_relative_to(WEB.resolve()) and path.is_file():
                return self.send_bytes(path.read_bytes(), mimetypes.guess_type(path)[0] or "application/octet-stream")
        self.send_error(404)

    def do_POST(self):
        body = json.loads(self.rfile.read(int(self.headers.get("Content-Length", 0))))
        if self.path == "/api/complete":
            if Handler.fail_complete:
                Handler.fail_complete = False
                return self.send_error(503)
            line = body["line"]
            return self.send_json({"start": 4, "end": len(line), "replacement": "/projects/overview.md ", "candidates": ["/projects/overview.md"], "session_id": "test", "current_directory": "/"})
        if self.path == "/api/exec":
            if Handler.fail_exec:
                Handler.fail_exec = False
                return self.send_error(503)
            line = body["line"]
            if line == "slow":
                time.sleep(0.2)
            if line == "welcome":
                output = ("Explore projects: head -n 11 /projects/overview.md\n"
                          "View architecture: head -n 26 /docs/architecture/overview.md\n"
                          "Live dashboard: dashboard")
            elif match := re.fullmatch(r"head -n (\d+) (/[\w/.-]+\.md)", line):
                target = CONTENT / match[2].lstrip("/")
                output = "\n".join(target.read_text(encoding="utf-8").splitlines()[:int(match[1])])
            else:
                output = "ok"
            return self.send_json({"output": output, "exit_code": 0, "session_id": "test", "current_directory": "/"})
        self.send_error(404)

    def send_json(self, value):
        self.send_bytes(json.dumps(value).encode(), "application/json")

    def send_bytes(self, value, content_type):
        self.send_response(200)
        self.send_header("Content-Type", content_type)
        self.send_header("Content-Length", str(len(value)))
        self.end_headers()
        self.wfile.write(value)


def main():
    verify_excerpts()
    browser = os.environ.get("BROWSER") or next((str(path) for path in (
        Path(r"C:\Program Files\Google\Chrome\Application\chrome.exe"),
        Path(r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe"),
    ) if path.exists()), None)
    if not browser:
        raise SystemExit("Chrome or Edge required")
    server = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    try:
        for query, expected, width in (
            ("smoke=main", "PASS terminal default", 1280),
            ("smoke=mobile", "PASS mobile terminal starts without input focus", 512),
            ("smoke=deep&doc=/cv/experience.md", "PASS deep link, refresh, and mobile layout", 512),
            ("smoke=diagram&doc=/docs/architecture/vaultsh.md", "PASS compact command flow diagram", 512),
        ):
            with tempfile.TemporaryDirectory() as profile:
                url = f"http://127.0.0.1:{server.server_port}/vault/?{query}"
                completed = subprocess.run([browser, "--headless=new", "--disable-gpu", "--no-sandbox", "--disable-software-rasterizer", "--no-first-run",
                                            f"--user-data-dir={profile}", f"--window-size={width},844",
                                            "--virtual-time-budget=30000" if "smoke=diagram" in query else "--virtual-time-budget=12000", "--dump-dom", url],
                                           capture_output=True, text=True, encoding="utf-8", errors="replace", timeout=35)
                marker = completed.stdout.split('id="smoke-results">')[-1].split("</pre>")[0]
                print(f"{width}px: {marker}")
                if completed.returncode or expected not in marker:
                    raise SystemExit(f"Browser smoke failed: {completed.stdout[-1300:]} {completed.stderr[-600:]}")
                if (screenshot_dir := os.environ.get("VAULTSH_SCREENSHOTS")) and "smoke=diagram" not in query:
                    screenshot = Path(screenshot_dir) / f"vaultsh-{width}.png"
                    screenshot.parent.mkdir(parents=True, exist_ok=True)
                    document = "/projects/overview.md" if width == 1280 else "/cv/experience.md"
                    capture_url = f"http://127.0.0.1:{server.server_port}/vault/?doc={document}"
                    subprocess.run([browser, "--headless=new", "--disable-gpu", "--no-sandbox",
                                    "--disable-software-rasterizer", "--no-first-run",
                                    f"--user-data-dir={profile}", f"--window-size={width},844",
                                    "--virtual-time-budget=4000", f"--screenshot={screenshot}", capture_url],
                                   capture_output=True, timeout=25, check=True)
        with tempfile.TemporaryDirectory() as profile:
            url = f"http://127.0.0.1:{server.server_port}/vault/"
            completed = subprocess.run([browser, "--headless=new", "--disable-gpu", "--no-sandbox",
                                        "--disable-software-rasterizer", "--disable-javascript",
                                        f"--user-data-dir={profile}", "--dump-dom", url],
                                       capture_output=True, text=True, encoding="utf-8", errors="replace", timeout=25)
            if completed.returncode or "Backend Lab connects Vaultsh's read-only virtual shell" not in completed.stdout:
                raise SystemExit("No-JavaScript fallback missing")
            print("No-JavaScript fallback: PASS")
    finally:
        server.shutdown()


if __name__ == "__main__":
    main()
