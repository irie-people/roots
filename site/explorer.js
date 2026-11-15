// The parts explorer: a tree of the stories (stories.js), a canvas that
// draws the picked one in the palette, mode, size, width and ground asked
// for, and four panels: its controls, its code, its usage from
// components.md, and the contrast of everything it writes. The address
// keeps the story (#session-card/needs-you), so a link shows the same.
(function () {
  'use strict';
  const D = window.IRIE_DESIGN, STORIES = window.ROOTS_STORIES;
  const USAGE = JSON.parse(document.getElementById('ex-usage').textContent);
  const $ = (id) => document.getElementById(id);
  const esc = (s) => String(s).replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));
  const slug = (s) => s.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
  const keep = {
    get(k, d) { try { return localStorage.getItem('roots.ex.' + k) || d; } catch (e) { return d; } },
    set(k, v) { try { localStorage.setItem('roots.ex.' + k, v); } catch (e) { /* the page works without */ } },
  };
  const view = { mode: keep.get('mode', 'both'), density: keep.get('density', 'pointer'), width: keep.get('width', 'full'), ground: keep.get('ground', 'sidebar') };
  // ?tab=code (or usage, contrast) opens that panel first.
  const asked = new URLSearchParams(location.search).get('tab');
  let story = STORIES[0], variant = story.stories[0], args = {}, panel = ['controls', 'code', 'usage', 'contrast'].includes(asked) ? asked : 'controls';

  // ---- Colour: WCAG 2.2's ratio, and OKLCH for the agents' hues -------
  const rgb = (h) => [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16));
  const lin = (c) => { c /= 255; return c <= 0.04045 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4); };
  const lumOf = (r, g, b) => 0.2126 * lin(r) + 0.7152 * lin(g) + 0.0722 * lin(b);
  const ratioOf = (a, b) => (Math.max(a, b) + 0.05) / (Math.min(a, b) + 0.05);
  const ratio = (x, y) => ratioOf(lumOf(...rgb(x)), lumOf(...rgb(y)));
  function oklch(h) {
    const [r, g, b] = rgb(h).map(lin);
    const l = Math.cbrt(0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b);
    const m = Math.cbrt(0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b);
    const s = Math.cbrt(0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b);
    const A = 1.9779984951 * l - 2.4285922050 * m + 0.4505937099 * s;
    const B = 0.0259040371 * l + 0.7827717662 * m - 0.8086757660 * s;
    return { l: 0.2104542553 * l + 0.7936177850 * m - 0.0040720468 * s, h: (Math.atan2(B, A) * 180 / Math.PI + 360) % 360 };
  }
  const A = D.agents;
  const hueOf = (set, slot) => {
    const step = (360 - 2 * A.needs_gap) / A.slots;
    return (oklch(set.needs).h + A.needs_gap + step / 2 + slot * step) % 360;
  };
  function agentVars(set, dark, slot) {
    const pick = (p) => (dark ? p.dark : p.light), h = hueOf(set, slot).toFixed(1);
    const card = Math.min(1, Math.max(0, oklch(set.card).l + pick(A.wash.lift)));
    return {
      '--a-mark': `oklch(${pick(A.mark.lightness)} ${A.mark.chroma} ${h})`,
      '--a-ink': `oklch(${pick(A.ink.lightness)} ${A.mark.chroma} ${h})`,
      '--a-wash': `oklch(${card.toFixed(3)} ${pick(A.wash.chroma)} ${h})`,
    };
  }

  // ---- The tree ----------------------------------------------------------
  function tree(filter) {
    const f = (filter || '').trim().toLowerCase();
    const groups = {};
    for (const s of STORIES) {
      if (f && !s.title.toLowerCase().includes(f) && !s.stories.some((v) => v.name.toLowerCase().includes(f))) continue;
      (groups[s.group] = groups[s.group] || []).push(s);
    }
    $('ex-tree').innerHTML = Object.entries(groups).map(([g, list]) => `<p class="nav-group">${esc(g)}</p><ul>${list.map((s) => `
      <li><a href="#${s.id}"${s === story ? ' aria-current="true"' : ''}>${esc(s.title)}</a>${s === story ? `<ul>${s.stories.map((v) => `<li><a href="#${s.id}/${slug(v.name)}"${v === variant ? ' aria-current="page"' : ''}>${esc(v.name)}</a></li>`).join('')}</ul>` : ''}</li>`).join('')}</ul>`).join('') || '<p class="ex-dim">No part has that name.</p>';
  }

  // ---- The canvas ------------------------------------------------------------
  function palette() { return D.palettes.find((p) => p.id === (document.documentElement.dataset.palette || D.default)) || D.palettes[0]; }
  function draw() {
    const p = palette();
    const modes = view.mode === 'both' ? ['light', 'dark'] : [view.mode];
    $('ex-canvas').className = 'ex-canvas' + (view.width === 'phone' ? ' is-phone' : '');
    $('ex-canvas').innerHTML = modes.map((m) => `<div class="ex-frame ir" data-palette="${p.id}" data-mode="${m}"${view.density === 'touch' ? ' data-density="touch"' : ''} style="background:var(--${view.ground})"><span class="ex-frame-label">${esc(p.name)}, ${m}</span><div class="ex-frame-body"></div></div>`).join('');
    $('ex-canvas').querySelectorAll('.ex-frame').forEach((frame) => {
      const dark = frame.dataset.mode === 'dark', set = dark ? p.dark : p.light;
      const ctx = { D, set, dark, ratio, hue: (slot) => hueOf(set, slot) };
      frame.querySelector('.ex-frame-body').innerHTML = story.render(args, ctx);
      frame.querySelectorAll('[data-slot]').forEach((el) => {
        const v = agentVars(set, dark, +el.dataset.slot);
        for (const k in v) el.style.setProperty(k, v[k]);
      });
    });
    IrieRoots.wire($('ex-canvas'));
    if (panel === 'code' || panel === 'contrast') showPanel();
  }

  // ---- The panels --------------------------------------------------------------
  function controlsPanel() {
    const keys = Object.keys(story.controls || {});
    if (!keys.length) return '<p class="ex-dim">This one has nothing to change but the palette, mode, size, width and ground above.</p>';
    return '<form class="ex-controls" id="ex-form">' + keys.map((k) => {
      const c = story.controls[k], label = esc(c.label || k[0].toUpperCase() + k.slice(1)), v = args[k];
      let input;
      switch (c.type) {
        case 'select': input = `<select class="ir-field" name="${k}">${c.options.map((o) => `<option${o === v ? ' selected' : ''}>${esc(o)}</option>`).join('')}</select>`; break;
        case 'boolean': input = `<input type="checkbox" name="${k}"${v ? ' checked' : ''}>`; break;
        case 'number': input = `<input class="ir-field" type="number" name="${k}" value="${esc(v)}" min="${c.min ?? ''}" max="${c.max ?? ''}">`; break;
        default: input = `<input class="ir-field" type="text" name="${k}" value="${esc(v)}">`;
      }
      return `<label class="ex-control${c.type === 'boolean' ? ' is-check' : ''}"><span>${label}</span>${input}</label>`;
    }).join('') + '</form><p class="ex-dim">Each change draws the part again; the address keeps the story, not these changes.</p>';
  }
  function codePanel() {
    const html = story.render(args, { D, set: palette().light, dark: false, ratio, hue: (s) => hueOf(palette().light, s) });
    return `<div class="ex-code-head"><span class="ex-dim">The part's HTML, with web/tokens.css and web/components.css loaded, inside an element of class <code>ir</code>.</span><button class="ir-button" type="button" id="ex-copy">Copy</button></div><pre><code class="language-html">${esc(html.trim())}</code></pre>`;
  }
  function usagePanel() {
    const u = story.usage && USAGE[story.usage];
    const more = u ? `<p><a href="../components.html#${u.id}">${esc(story.usage)} in Components</a></p>` : story.docs ? `<p><a href="${story.docs}">More in the docs</a></p>` : '';
    return (u ? u.html : `<p>${esc(story.blurb || '')}</p>`) + more;
  }
  // Every element that writes words, with its colour and the ground under
  // it, after transparency, measured as the screen shows them.
  const probe = document.createElement('canvas').getContext('2d', { willReadFrequently: true });
  function rgba(css) {
    probe.clearRect(0, 0, 1, 1);
    probe.fillStyle = '#000';
    probe.fillStyle = css;
    probe.fillRect(0, 0, 1, 1);
    const d = probe.getImageData(0, 0, 1, 1).data;
    return [d[0], d[1], d[2], d[3] / 255];
  }
  function groundOf(el) {
    const layers = [];
    for (let e = el; e; e = e.parentElement) {
      const cs = getComputedStyle(e);
      let c = cs.backgroundColor;
      const img = cs.backgroundImage;
      if (img && img !== 'none') {
        const m = img.match(/(rgba?\([^)]*\)|oklch\([^)]*\)|color\([^)]*\)|#[0-9a-f]{3,8})/i);
        if (m) c = m[1];
      }
      const v = rgba(c);
      if (v[3] > 0) layers.push(v);
      if (v[3] >= 1) break;
    }
    let out = [255, 255, 255];
    for (const [r, g, b, a] of layers.reverse()) out = [r * a + out[0] * (1 - a), g * a + out[1] * (1 - a), b * a + out[2] * (1 - a)];
    return out;
  }
  const hex = ([r, g, b]) => '#' + [r, g, b].map((x) => Math.round(x).toString(16).padStart(2, '0')).join('').toUpperCase();
  function contrastPanel() {
    const rows = [];
    $('ex-canvas').querySelectorAll('.ex-frame').forEach((frame) => {
      const seen = new Set();
      frame.querySelectorAll('.ex-frame-body *').forEach((el) => {
        const own = [...el.childNodes].some((n) => n.nodeType === 3 && n.textContent.trim()) || el.matches('.ir-mark, textarea, input');
        if (!own || el.closest('[hidden]') || getComputedStyle(el).visibility === 'hidden' || getComputedStyle(el).display === 'none') return;
        const cs = getComputedStyle(el);
        const fg = rgba(cs.color), bg = groundOf(el);
        const shown = [fg[0] * fg[3] + bg[0] * (1 - fg[3]), fg[1] * fg[3] + bg[1] * (1 - fg[3]), fg[2] * fg[3] + bg[2] * (1 - fg[3])];
        const r = ratioOf(lumOf(...shown), lumOf(...bg));
        const big = parseFloat(cs.fontSize) >= 24 || (parseFloat(cs.fontSize) >= 18.66 && +cs.fontWeight >= 700);
        const need = el.matches('.ir-mark') || big ? 3 : 4.5;
        const words = el.matches('.ir-mark') ? `the ${el.dataset.state} mark` : (el.value || el.placeholder || [...el.childNodes].filter((n) => n.nodeType === 3).map((n) => n.textContent.trim()).join(' ')).slice(0, 40);
        const key = words + hex(shown) + hex(bg);
        if (seen.has(key)) return;
        seen.add(key);
        rows.push(`<tr><td>${esc(frame.dataset.mode)}</td><td>${esc(words)}</td><td><span class="sw" style="background:${hex(shown)}"></span><code>${hex(shown)}</code></td><td><span class="sw" style="background:${hex(bg)}"></span><code>${hex(bg)}</code></td><td class="${r >= need ? 'ok' : 'low'}">${r.toFixed(2)}:1</td><td>${need}:1</td></tr>`);
      });
    });
    return `<p class="ex-dim">Every piece of text and every mark, measured as drawn (WCAG 2.2: text 4.5:1, large text and marks 3:1). A low figure on a disabled or decorative piece may be right; anything else is a fault.</p><div class="table"><table><thead><tr><th>Mode</th><th>What</th><th>Colour</th><th>Ground</th><th>Contrast</th><th>Needs</th></tr></thead><tbody>${rows.join('') || '<tr><td colspan="6">Nothing to measure.</td></tr>'}</tbody></table></div>`;
  }
  function showPanel() {
    const html = { controls: controlsPanel, code: codePanel, usage: usagePanel, contrast: contrastPanel }[panel]();
    $('ex-panel').innerHTML = html;
    if (panel === 'controls' && $('ex-form')) {
      $('ex-form').addEventListener('input', (e) => {
        const el = e.target, c = story.controls[el.name];
        args[el.name] = c.type === 'boolean' ? el.checked : c.type === 'number' ? +el.value : el.value;
        draw();
      });
      $('ex-form').addEventListener('submit', (e) => e.preventDefault());
    }
    if (panel === 'code') {
      $('ex-copy').addEventListener('click', () => {
        const text = $('ex-panel').querySelector('code').textContent;
        (navigator.clipboard ? navigator.clipboard.writeText(text) : Promise.reject()).then(
          () => { $('ex-copy').textContent = 'Copied'; },
          () => { $('ex-copy').textContent = 'Select and copy'; });
      });
    }
  }

  // ---- Going to a story --------------------------------------------------------
  function go() {
    const [id, v] = location.hash.slice(1).split('/');
    story = STORIES.find((s) => s.id === id) || STORIES[0];
    variant = story.stories.find((x) => slug(x.name) === v) || story.stories[0];
    args = Object.assign({}, story.args, variant.args || {});
    $('ex-title').textContent = story.title;
    $('ex-crumb').textContent = story.group + (story.stories.length > 1 ? ' · ' + variant.name : '');
    $('ex-blurb').textContent = story.blurb || '';
    $('ex-blurb').hidden = !story.blurb;
    document.title = story.title + ' · Explorer · Roots';
    tree($('side').querySelector('.ex-find').value);
    draw();
    showPanel();
  }

  // ---- The toolbar -----------------------------------------------------------------
  for (const k of ['mode', 'density', 'width', 'ground']) {
    const seg = $('ex-' + k);
    const mark = () => seg.querySelectorAll('button').forEach((b) => b.setAttribute('aria-pressed', String(b.dataset.v === view[k])));
    mark();
    seg.addEventListener('click', (e) => {
      const b = e.target.closest('button');
      if (!b) return;
      view[k] = b.dataset.v;
      keep.set(k, view[k]);
      mark();
      draw();
    });
  }
  $('side').querySelector('.ex-find').addEventListener('input', (e) => tree(e.target.value));
  document.querySelector('.ex-tabs').addEventListener('ir-pick', (e) => {
    panel = e.target.dataset.panel;
    showPanel();
  });
  // The site's own palette picker changes <html>; the canvas follows it.
  new MutationObserver(draw).observe(document.documentElement, { attributes: true, attributeFilter: ['data-palette', 'data-mode'] });
  addEventListener('hashchange', go);
  IrieRoots.wire(document.querySelector('.ex-tabs'));
  document.querySelectorAll('.ex-tabs [role="tab"]').forEach((t) => {
    const on = t.dataset.panel === panel;
    t.setAttribute('aria-selected', String(on));
    t.tabIndex = on ? 0 : -1;
  });
  go();
})();
