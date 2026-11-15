// The behaviour of Roots' web parts that CSS alone cannot give
// (components.md): the tab row's keys, its picked tab, and its overflow
// chevrons (components §3.13). Load it after components.css, as a classic
// script, then call IrieRoots.wire() after drawing parts, with the element
// they were drawn into, so their tabs name their panels. A picked tab sends
// an ir-pick event, which bubbles, for the page to show its panel.
(function () {
  'use strict';
  const still = matchMedia('(prefers-reduced-motion: reduce)');
  const glide = () => (still.matches ? 'instant' : 'smooth');
  let ids = 0;

  // Each side where tabs hide gets a fade and a chevron.
  function syncTabRow(list) {
    const at = list.scrollLeft, max = list.scrollWidth - list.clientWidth;
    const start = at > 1, end = at < max - 1;
    list.classList.toggle('is-more-start', start);
    list.classList.toggle('is-more-end', end);
    const row = list.closest('.ir-tabs-row');
    if (!row) return;
    const back = row.querySelector('.ir-tabs-more.is-start'), on = row.querySelector('.ir-tabs-more.is-end');
    if (back) back.hidden = !start;
    if (on) on.hidden = !end;
  }

  // Tabs drawn anew get new ids, for their panel to name.
  function wire(root) {
    root = root || document;
    root.querySelectorAll('.ir-tabs').forEach((tabs) => {
      const panel = tabs.querySelector('[role="tabpanel"]');
      if (!panel) return;
      panel.id = 'ir-tabpanel-' + ++ids;
      tabs.querySelectorAll('[role="tab"]').forEach((tab) => {
        tab.id = 'ir-tab-' + ++ids;
        tab.setAttribute('aria-controls', panel.id);
        if (tab.getAttribute('aria-selected') === 'true') panel.setAttribute('aria-labelledby', tab.id);
      });
    });
    root.querySelectorAll('.ir-tablist').forEach(syncTabRow);
  }

  // A click or an arrow key picks a tab, which comes into view clear of
  // the fade.
  function pickTab(tab, focus) {
    const tabs = [...tab.closest('[role="tablist"]').querySelectorAll('[role="tab"]')];
    for (const t of tabs) {
      const on = t === tab;
      t.setAttribute('aria-selected', String(on));
      t.tabIndex = on ? 0 : -1;
    }
    const panel = document.getElementById(tab.getAttribute('aria-controls'));
    if (panel) {
      panel.setAttribute('aria-labelledby', tab.id);
      // A demo's tab says its panel's words; a page's own tab leaves the
      // panel to the page, which hears the pick as an ir-pick event.
      if ('says' in tab.dataset) panel.textContent = tab.dataset.says;
    }
    if (focus) tab.focus({ preventScroll: true });
    tab.scrollIntoView({ block: 'nearest', inline: 'nearest', behavior: glide() });
    tab.dispatchEvent(new CustomEvent('ir-pick', { bubbles: true }));
  }

  document.addEventListener('click', (e) => {
    const tab = e.target instanceof Element ? e.target.closest('.ir-tab') : null;
    if (tab) { pickTab(tab, false); return; }
    const more = e.target instanceof Element ? e.target.closest('.ir-tabs-more') : null;
    if (!more) return;
    const list = more.closest('.ir-tabs-row').querySelector('.ir-tablist');
    list.scrollBy({ left: (more.classList.contains('is-start') ? -0.8 : 0.8) * list.clientWidth, behavior: glide() });
  });
  document.addEventListener('keydown', (e) => {
    const tab = e.target instanceof Element ? e.target.closest('.ir-tab') : null;
    if (!tab) return;
    const tabs = [...tab.closest('[role="tablist"]').querySelectorAll('[role="tab"]')];
    const i = tabs.indexOf(tab);
    const to = { ArrowRight: i + 1, ArrowLeft: i - 1, Home: 0, End: tabs.length - 1 }[e.key];
    if (to === undefined) return;
    e.preventDefault();
    pickTab(tabs[(to + tabs.length) % tabs.length], true);
  });
  // A row's scroll events don't bubble, so they are caught on their way down.
  document.addEventListener('scroll', (e) => {
    if (e.target instanceof Element && e.target.matches('.ir-tablist')) syncTabRow(e.target);
  }, true);
  addEventListener('resize', () => document.querySelectorAll('.ir-tablist').forEach(syncTabRow));

  window.IrieRoots = { wire };
})();
