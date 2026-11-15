// The site's own look: the palette and mode a visitor picks, kept in this
// browser only, set on <html> before the page is drawn (web/tokens.css
// reads data-palette and data-mode there); the menu on a narrow screen;
// and a Copy button on each block of code.
(function () {
  'use strict';
  const root = document.documentElement;
  const keep = {
    get(k) { try { return localStorage.getItem('roots.' + k); } catch (e) { return null; } },
    set(k, v) { try { localStorage.setItem('roots.' + k, v); } catch (e) { /* the page works without */ } },
  };
  function apply() {
    const p = keep.get('palette'), m = keep.get('mode');
    if (p) root.dataset.palette = p; else delete root.dataset.palette;
    if (m === 'light' || m === 'dark') root.dataset.mode = m; else delete root.dataset.mode;
  }
  apply();
  document.addEventListener('DOMContentLoaded', () => {
    const mark = () => {
      document.querySelectorAll('[data-pick-palette]').forEach((b) => b.setAttribute('aria-pressed', String(b.dataset.pickPalette === (root.dataset.palette || 'rosewater'))));
      document.querySelectorAll('[data-pick-mode]').forEach((b) => b.setAttribute('aria-pressed', String(b.dataset.pickMode === (root.dataset.mode || 'auto'))));
    };
    mark();
    document.addEventListener('click', (e) => {
      const b = e.target instanceof Element ? e.target.closest('[data-pick-palette], [data-pick-mode], .menu-button, pre .copy') : null;
      if (!b) return;
      if (b.dataset.pickPalette) { keep.set('palette', b.dataset.pickPalette); apply(); mark(); }
      else if (b.dataset.pickMode) { keep.set('mode', b.dataset.pickMode); apply(); mark(); }
      else if (b.classList.contains('menu-button')) {
        const open = b.getAttribute('aria-expanded') !== 'true';
        b.setAttribute('aria-expanded', String(open));
        document.body.classList.toggle('menu-open', open);
      } else {
        const code = b.parentElement.querySelector('code');
        (navigator.clipboard ? navigator.clipboard.writeText(code.textContent) : Promise.reject()).then(
          () => { b.textContent = 'Copied'; setTimeout(() => { b.textContent = 'Copy'; }, 1600); },
          () => { b.textContent = 'Select and copy'; });
      }
    });
    document.querySelectorAll('main pre').forEach((pre) => {
      const b = document.createElement('button');
      b.type = 'button';
      b.className = 'copy';
      b.textContent = 'Copy';
      pre.prepend(b);
    });
  });
})();
