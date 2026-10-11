function initGui() {
  const shellTab = document.querySelector("#shell-tab");
  const guiTab = document.querySelector("#gui-tab");
  const shellPanel = document.querySelector("#shell-panel");
  const guiPanel = document.querySelector("#gui-panel");
  const openGui = document.querySelector("#open-gui");
  const breadcrumbs = document.querySelector("#gui-breadcrumbs");
  const location = document.querySelector("#gui-location");
  const count = document.querySelector("#gui-count");
  const guiNav = document.querySelector("#gui-nav");
  const back = document.querySelector("#gui-back");
  const home = document.querySelector("#gui-home");
  const folders = document.querySelector("#gui-folders");
  const documentView = document.querySelector("#gui-document");
  const message = document.querySelector("#gui-message");
  const shortcuts = document.querySelectorAll("[data-gui-path]");
  if (!shellTab || !guiTab || !shellPanel || !guiPanel || !openGui ||
      !breadcrumbs || !location || !count || !guiNav || !back || !home ||
      !folders || !documentView || !message) return;

  let mode = "shell";
  let currentPath = "/";
  let requestNumber = 0;
  let hasLoaded = false;
  let mermaidModule;

  const modeRow = document.querySelector(".mode-row");
  const modeSentinel = document.querySelector(".mode-sticky-sentinel");
  if (modeRow && modeSentinel) {
    const updateSticky = () => {
      modeRow.classList.toggle("is-stuck", modeSentinel.getBoundingClientRect().bottom < 0);
    };
    window.addEventListener("scroll", updateSticky, { passive: true });
    updateSticky();
  }

  function setMode(next, focus = true) {
    mode = next;
    const desktop = next === "desktop";
    shellPanel.hidden = desktop;
    guiPanel.hidden = !desktop;
    shellTab.setAttribute("aria-selected", String(!desktop));
    guiTab.setAttribute("aria-selected", String(desktop));
    shellTab.tabIndex = desktop ? -1 : 0;
    guiTab.tabIndex = desktop ? 0 : -1;
    document.querySelector(".mode-switch").classList.toggle("mode-switch--desktop", desktop);
    if (focus) (desktop ? guiTab : shellTab).focus();
    if (desktop && !hasLoaded) {
      hasLoaded = true;
      navigate("/");
    }
    try { window.sessionStorage.setItem("vaultsh-view", next); } catch { /* optional */ }
  }

  shellTab.addEventListener("click", () => setMode("shell"));
  guiTab.addEventListener("click", () => setMode("desktop"));
  openGui.addEventListener("click", () => setMode("desktop"));
  for (const tab of [shellTab, guiTab]) {
    tab.addEventListener("keydown", (event) => {
      if (event.key === "ArrowLeft" || event.key === "ArrowRight") {
        event.preventDefault();
        setMode(mode === "shell" ? "desktop" : "shell");
      }
    });
  }
  for (const shortcut of shortcuts) {
    shortcut.addEventListener("click", () => navigate(shortcut.dataset.guiPath));
  }
  back.addEventListener("click", () => navigate(parentPath(currentPath)));
  home.addEventListener("click", () => navigate("/"));

  function parentPath(path) {
    return path.slice(0, path.lastIndexOf("/")) || "/";
  }

  async function navigate(path) {
    const serial = ++requestNumber;
    message.textContent = "Opening…";
    try {
      const url = new URL("/api/browse", window.location.origin);
      url.searchParams.set("path", path);
      const response = await fetch(url, { cache: "no-store" });
      if (!response.ok) throw new Error(response.status === 429 ? "Too many requests. Try again shortly." : "This item could not be opened.");
      const item = await response.json();
      if (serial !== requestNumber) return;
      currentPath = item.path;
      location.textContent = item.path;
      guiNav.hidden = item.path === "/";
      renderBreadcrumbs(item.path);
      renderShortcuts(item.path);
      if (item.kind === "directory") renderDirectory(item);
      else renderFile(item);
      message.textContent = "";
    } catch (error) {
      if (serial === requestNumber) message.textContent = error.message || "This item could not be opened.";
    }
  }

  function renderBreadcrumbs(path) {
    breadcrumbs.replaceChildren();
    const parts = path.split("/").filter(Boolean);
    if (!parts.length) addCrumb("Home", "/");
    let soFar = "";
    for (const part of parts) {
      if (soFar) {
        const separator = document.createElement("span");
        separator.className = "breadcrumb-separator";
        separator.textContent = "/";
        breadcrumbs.append(separator);
      }
      soFar += "/" + part;
      addCrumb(part, soFar);
    }
    function addCrumb(label, target) {
      const button = document.createElement("button");
      button.type = "button";
      button.textContent = label;
      button.disabled = target === path;
      button.addEventListener("click", () => navigate(target));
      breadcrumbs.append(button);
    }
  }

  function renderShortcuts(path) {
    for (const shortcut of shortcuts) {
      const target = shortcut.dataset.guiPath;
      shortcut.classList.toggle("is-current", path === target || (target === "/cv" && path.startsWith("/cv/")) ||
        (target === "/projects" && path.startsWith("/projects/")) ||
        (target === "/docs" && path.startsWith("/docs/")));
    }
  }

  function renderDirectory(item) {
    documentView.hidden = true;
    documentView.replaceChildren();
    folders.hidden = false;
    folders.replaceChildren();
    const children = (item.children || []).slice().sort((a, b) =>
      (a.kind === b.kind ? a.name.localeCompare(b.name) : a.kind === "directory" ? -1 : 1));
    count.textContent = `${children.length} item${children.length === 1 ? "" : "s"}`;
    if (!children.length) {
      const empty = document.createElement("p");
      empty.className = "gui-empty";
      empty.textContent = "This folder is empty.";
      folders.append(empty);
    }
    for (const child of children) {
      const button = document.createElement("button");
      button.type = "button";
      button.className = "gui-item";
      const icon = document.createElement("span");
      icon.className = `gui-item-icon gui-item-icon--${child.kind}`;
      icon.setAttribute("aria-hidden", "true");
      icon.textContent = child.kind === "directory" ? "▰" : "≡";
      const name = document.createElement("strong");
      const label = child.name.replace(/\.md$/i, "").replace(/[-_]/g, " ");
      name.textContent = child.name === "cv" ? "CV" : child.name === "README.md" ? "Read me" :
        label.charAt(0).toUpperCase() + label.slice(1);
      const type = document.createElement("small");
      type.textContent = child.kind === "directory" ? "Folder" : "Document";
      button.append(icon, name, type);
      button.addEventListener("click", () => navigate(child.path));
      folders.append(button);
    }
  }

  function renderFile(item) {
    folders.hidden = true;
    folders.replaceChildren();
    documentView.hidden = false;
    documentView.replaceChildren();
    count.textContent = "DOCUMENT PREVIEW";
    const body = document.createElement("div");
    body.className = "markdown-body";
    renderMarkdown(body, item.content || "");
    documentView.append(body);
    renderDiagrams(body);
  }

  async function renderDiagrams(body) {
    const diagrams = [...body.querySelectorAll("pre.mermaid")];
    if (!diagrams.length) return;
    try {
      mermaidModule ??= import("./vendor/mermaid.js");
      const { default: mermaid } = await mermaidModule;
      mermaid.initialize({ startOnLoad: false, securityLevel: "strict", theme: "neutral", htmlLabels: false, suppressErrorRendering: true });
      for (const diagram of diagrams) {
        if (!diagram.isConnected) continue;
        try {
          await mermaid.run({ nodes: [diagram], suppressErrors: true });
        } catch {
          // The original diagram source remains readable if rendering fails.
        }
      }
    } catch {
      // The source remains visible if the diagram library is unavailable.
    }
  }

  // Build DOM nodes from Markdown instead of inserting source content as HTML.
  function renderMarkdown(parent, source) {
    const lines = source.replace(/\r\n?/g, "\n").split("\n");
    let index = 0;
    while (index < lines.length) {
      const line = lines[index];
      if (!line.trim()) { index++; continue; }
      const fence = line.match(/^\s*```([\w-]*)/);
      if (fence) {
        const code = [];
        index++;
        while (index < lines.length && !/^\s*```/.test(lines[index])) code.push(lines[index++]);
        if (index < lines.length) index++;
        const pre = document.createElement("pre");
        if (fence[1].toLowerCase() === "mermaid") {
          pre.className = "mermaid";
          pre.textContent = code.join("\n");
        } else {
          const codeNode = document.createElement("code");
          codeNode.textContent = code.join("\n");
          pre.append(codeNode);
        }
        parent.append(pre);
        continue;
      }
      const heading = line.match(/^(#{1,6})\s+(.+)$/);
      if (heading) {
        const node = document.createElement(`h${heading[1].length}`);
        appendInline(node, heading[2]);
        parent.append(node);
        index++;
        continue;
      }
      if (/^\s*(?:---+|\*\*\*+)\s*$/.test(line)) {
        parent.append(document.createElement("hr"));
        index++;
        continue;
      }
      const list = line.match(/^(\s*)([-*+] |\d+\. )(.*)$/);
      if (list) {
        const ordered = /\d/.test(list[2][0]);
        const group = document.createElement(ordered ? "ol" : "ul");
        while (index < lines.length) {
          const match = lines[index].match(/^(\s*)([-*+] |\d+\. )(.*)$/);
          if (!match || ordered !== /\d/.test(match[2][0])) break;
          const item = document.createElement("li");
          const text = [match[3]];
          index++;
          while (index < lines.length && /^\s{2,}\S/.test(lines[index]) && !/^\s*([-*+] |\d+\. )/.test(lines[index])) {
            text.push(lines[index++].trim());
          }
          appendInline(item, text.join(" "));
          group.append(item);
        }
        parent.append(group);
        continue;
      }
      if (/^>\s?/.test(line)) {
        const quote = document.createElement("blockquote");
        const text = [];
        while (index < lines.length && /^>\s?/.test(lines[index])) text.push(lines[index++].replace(/^>\s?/, ""));
        appendInline(quote, text.join(" "));
        parent.append(quote);
        continue;
      }
      const paragraph = [];
      while (index < lines.length && lines[index].trim() &&
        !/^(#{1,6}\s|\s*```|\s*(?:[-*+] |\d+\. )|>\s?)/.test(lines[index])) {
        paragraph.push(lines[index++].trim());
      }
      if (!paragraph.length) { index++; continue; }
      const node = document.createElement("p");
      appendInline(node, paragraph.join(" "));
      parent.append(node);
    }
  }

  function appendInline(parent, source) {
    const token = /(\[([^\]]+)\]\(([^)]+)\)|\*\*([^*]+)\*\*|`([^`]+)`|\*([^*]+)\*)/g;
    let offset = 0;
    for (const match of source.matchAll(token)) {
      parent.append(document.createTextNode(source.slice(offset, match.index)));
      if (match[2] !== undefined) {
        const href = match[3];
        if (/^(https?:\/\/|mailto:)/i.test(href)) {
          const link = document.createElement("a");
          link.href = href;
          link.textContent = match[2];
          if (/^https?:\/\//i.test(href)) { link.target = "_blank"; link.rel = "noreferrer noopener"; }
          parent.append(link);
        } else parent.append(document.createTextNode(match[2]));
      } else {
        const node = document.createElement(match[4] !== undefined ? "strong" : match[5] !== undefined ? "code" : "em");
        node.textContent = match[4] ?? match[5] ?? match[6];
        parent.append(node);
      }
      offset = match.index + match[0].length;
    }
    parent.append(document.createTextNode(source.slice(offset)));
  }

  try {
    if (window.sessionStorage.getItem("vaultsh-view") === "desktop") setMode("desktop", false);
  } catch { /* optional */ }
}

initGui();
