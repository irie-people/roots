// The explorer's stories, as Storybook keeps them: each part with the
// controls that change it (args) and named states to start from. render()
// gives the part's HTML, made only of Roots' ir- classes; a part with
// data-slot on it takes an agent's colours, which the explorer works out.
// `usage` names the part's section in components.md.
window.ROOTS_STORIES = (function () {
  'use strict';
  const esc = (s) => String(s).replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));
  const icon = {
    pause: '<svg viewBox="0 0 16 16" aria-hidden="true"><path d="M5.5 3.5v9M10.5 3.5v9" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>',
    boost: '<svg viewBox="0 0 16 16" aria-hidden="true"><circle class="ir-fill" cx="8" cy="8" r="6.3" fill="none" stroke="currentColor" stroke-width="1.4"/><path class="ir-cut" d="M8 11.2V4.9M5.4 7.4 8 4.8l2.6 2.6" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>',
    park: '<svg viewBox="0 0 16 16" aria-hidden="true"><path d="M2.8 9.5v3.2h10.4V9.5M8 2.8v6.6M5.2 6.8 8 9.6l2.8-2.8" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>',
    more: '<svg viewBox="0 0 16 16" aria-hidden="true"><path d="M3.5 8h.01M8 8h.01M12.5 8h.01" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round"/></svg>',
    start: '<svg viewBox="0 0 16 16" aria-hidden="true"><path d="M8 12.5V3.8M4.4 7.2 8 3.6l3.6 3.6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>',
    gear: '<svg viewBox="0 0 20 20" aria-hidden="true"><circle cx="10" cy="10" r="2.6" fill="none" stroke="currentColor" stroke-width="1.6"/><path d="M10 2.8v2M10 15.2v2M2.8 10h2M15.2 10h2M4.9 4.9l1.4 1.4M13.7 13.7l1.4 1.4M4.9 15.1l1.4-1.4M13.7 6.3l1.4-1.4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/></svg>',
    chevStart: '<svg viewBox="0 0 12 12" aria-hidden="true"><path d="M7.5 2.5 4 6l3.5 3.5" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/></svg>',
    chevEnd: '<svg viewBox="0 0 12 12" aria-hidden="true"><path d="M4.5 2.5 8 6 4.5 9.5" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/></svg>',
  };
  // The states of DESIGN §5.3, with their words and their colour roles.
  const STATES = [
    ['needs', 'Needs you', 'needs'], ['stalled', 'Stalled', 'needs'], ['failed', 'Failed', 'needs'],
    ['working', 'Working', 'agent'], ['waiting', 'Waiting on helpers', 'agent'], ['starting', 'Starting', 'muted'],
    ['paused', 'Paused', 'paused'], ['pausing', 'Pausing', 'paused'], ['frozen', 'Frozen', 'frozen'],
    ['done', 'Done, to verify', 'done'], ['idle', 'Idle', 'muted'], ['exited', 'Exited', 'muted'],
  ];
  const mark = (state, slot) => `<span class="ir-mark" data-state="${state}"${slot === undefined ? '' : ` data-slot="${slot}"`}></span>`;
  const list = (s) => String(s).split(',').map((x) => x.trim()).filter(Boolean);

  return [
    {
      id: 'colour-roles', group: 'Foundations', title: 'Colour roles', docs: '../docs/colour.html',
      blurb: 'The palette in use, role by role, drawn as each is used. Text and marks give their lowest contrast on the eight surfaces.',
      args: { use: 'all' },
      controls: { use: { type: 'select', options: ['all', 'surface', 'line', 'text', 'mark', 'fill', 'on_fill', 'picture'], label: 'Drawn as' } },
      stories: [{ name: 'Every role' }, { name: 'Text', args: { use: 'text' } }, { name: 'Surfaces', args: { use: 'surface' } }],
      render(a, ctx) {
        const roles = ctx.D.roles.filter((r) => a.use === 'all' || r.use === a.use);
        const surfaces = ctx.D.roles.filter((r) => r.use === 'surface').map((r) => r.key);
        return '<div class="ex-roles">' + roles.map((r) => {
          const c = ctx.set[r.key];
          let note = r.use.replace('_', ' ');
          let fg = ctx.set.ink, bg = c;
          if (r.use === 'text' || r.use === 'mark') {
            const low = Math.min(...surfaces.map((s) => ctx.ratio(c, ctx.set[s])));
            note = `${low.toFixed(2)}:1 lowest`;
            fg = c; bg = ctx.set.card;
          } else if (r.use === 'on_fill') { fg = c; bg = ctx.set.accent_fill; note = `${ctx.ratio(c, bg).toFixed(2)}:1 on the button`; }
          else if (r.use === 'fill') { fg = ctx.set.on_accent; }
          return `<div class="ex-role"><span class="ex-chip" style="background:${bg};color:${fg}">Aa</span><b>${r.key}</b><span>${esc(r.title)} · ${c}</span><span>${note}</span></div>`;
        }).join('') + '</div>';
      },
    },
    {
      id: 'state-marks', group: 'Foundations', title: 'State marks', usage: 'State mark',
      blurb: 'Each state has a glyph and words of its own, so colour never carries it alone. Red is only for what needs you.',
      args: { words: true },
      controls: { words: { type: 'boolean', label: 'Show the words' } },
      stories: [{ name: 'Every state' }, { name: 'Glyphs only', args: { words: false } }],
      render: (a) => '<div class="ex-states">' + STATES.map(([s, w, role]) =>
        `<div class="ex-state">${mark(s, s === 'working' || s === 'waiting' ? 2 : undefined)}${a.words ? `<span>${w}</span><span class="ex-dim">${role}</span>` : ''}</div>`).join('') + '</div>',
    },
    {
      id: 'agent-colours', group: 'Foundations', title: "Agents' colours", docs: '../docs/colour.html',
      blurb: "Six hues 50° apart, clear of the needs red, at one low chroma; cards side by side take hues at least 100° apart. Shown at the library's starting lightness.",
      args: { order: 'neighbours' },
      controls: { order: { type: 'select', options: ['neighbours', 'round the wheel'], label: 'Order' } },
      stories: [{ name: 'As neighbours' }, { name: 'Round the wheel', args: { order: 'round the wheel' } }],
      render(a, ctx) {
        const slots = a.order === 'neighbours' ? [0, 3, 1, 4, 2, 5] : [0, 1, 2, 3, 4, 5];
        return '<div class="ex-stack">' + slots.map((s, n) => {
          const hue = Math.round(ctx.hue(s));
          const prev = n ? Math.round(ctx.hue(slots[n - 1])) : null;
          const apart = prev === null ? '' : Math.min(Math.abs(hue - prev), 360 - Math.abs(hue - prev)) + '° from the card above';
          return `<div class="ir-card is-working" data-slot="${s}">${mark('working', s)}<div class="ir-card-main"><div class="ir-card-top"><span class="ir-card-name">Agent colour ${s + 1}</span><span class="ir-card-slot">${hue}°</span></div><div class="ir-card-line">${apart || 'Hues spread from the needs red'}</div></div></div>`;
        }).join('') + '</div>';
      },
    },
    {
      id: 'state-mark', group: 'Marks', title: 'State mark', usage: 'State mark',
      args: { state: 'needs', words: 'Needs you' },
      controls: { state: { type: 'select', options: STATES.map((s) => s[0]) }, words: { type: 'text' } },
      stories: [{ name: 'Needs you' }, { name: 'Working', args: { state: 'working', words: 'Running the tests' } }, { name: 'Frozen', args: { state: 'frozen', words: 'Frozen' } }],
      render: (a) => `<span class="ex-state">${mark(a.state, a.state === 'working' || a.state === 'waiting' ? 1 : undefined)}<span>${esc(a.words)}</span></span>`,
    },
    {
      id: 'count', group: 'Marks', title: 'Count', usage: 'Count',
      args: { state: 'needs', n: 2 },
      controls: { state: { type: 'select', options: ['needs', 'working', 'paused', 'frozen', 'done'] }, n: { type: 'number', min: 0, max: 99 } },
      stories: [{ name: 'Needs you' }, { name: 'Working', args: { state: 'working', n: 3 } }, { name: 'None', args: { state: 'paused', n: 0 } }],
      render: (a) => `<button class="ir-count${a.state === 'needs' && a.n > 0 ? ' is-needs' : ''}${+a.n === 0 ? ' is-zero' : ''}" type="button">${mark(a.state, a.state === 'working' ? 0 : undefined)}${+a.n}</button>`,
    },
    {
      id: 'session-card', group: 'Parts', title: 'Session card', usage: 'Session card',
      blurb: 'One session in the list: its mark, name and age, what it is doing in plain words, and its helpers.',
      args: { state: 'working', name: 'irie-tests', line: 'Running the tests', age: '4m', slot: 0, helpers: true, keyboard: true, boosted: false },
      controls: {
        state: { type: 'select', options: ['quiet', 'working', 'needs', 'selected', 'hover', 'loose', 'done'] },
        name: { type: 'text' }, line: { type: 'text' }, age: { type: 'text' },
        slot: { type: 'number', min: 0, max: 5, label: 'Agent colour' },
        helpers: { type: 'boolean', label: 'A helper' }, keyboard: { type: 'boolean', label: 'List has the keyboard' }, boosted: { type: 'boolean', label: 'Boosted' },
      },
      stories: [
        { name: 'Working' },
        { name: 'Needs you', args: { state: 'needs', name: 'api-fix', line: 'Asks: may it push to main?', age: '12m', helpers: false } },
        { name: 'Selected', args: { state: 'selected', name: 'theme-fix', line: 'Editing the palettes', slot: 3, helpers: false } },
        { name: 'Under the pointer', args: { state: 'hover', name: 'site-copy', line: 'Waiting on its helpers', slot: 1, helpers: false } },
        { name: 'Boosted', args: { state: 'done', name: 'docs-pass', line: 'Finished: to verify', boosted: true, helpers: false, age: '2h' } },
        { name: 'Quiet', args: { state: 'quiet', name: 'old-branch', line: 'Idle since yesterday', age: '1d', helpers: false } },
        { name: 'Work not landed', args: { state: 'loose', name: 'old-branch', line: 'Its work has not landed', age: '1d', helpers: false } },
      ],
      render(a) {
        const s = a.state, slot = +a.slot;
        const cls = { quiet: 'is-quiet', working: 'is-working', needs: 'is-loud', selected: 'is-working is-selected' + (a.keyboard ? ' is-keyboard' : ''), hover: 'is-hover', loose: 'is-quiet is-loose', done: '' }[s];
        const m = { quiet: 'idle', working: 'working', needs: 'needs', selected: 'working', hover: 'waiting', loose: 'idle', done: 'done' }[s];
        const agent = s === 'working' || s === 'selected' || s === 'hover';
        const actions = `<span class="ir-actions"><button class="ir-act" type="button" aria-label="Pause">${icon.pause}</button><button class="ir-act${a.boosted ? ' is-on' : ''}" type="button" aria-label="${a.boosted ? 'Stop boosting' : 'Boost'}">${icon.boost}</button><button class="ir-act" type="button" aria-label="Park">${icon.park}</button><button class="ir-act" type="button" aria-label="More">${icon.more}</button></span>`;
        const boost = a.boosted && s !== 'hover' ? `<button class="ir-act is-on" type="button" aria-label="Stop boosting">${icon.boost}</button>` : '';
        return `<div class="ir-card${cls ? ' ' + cls : ''}"${agent ? ` data-slot="${slot}"` : ''}>
  ${mark(m, agent ? slot : undefined)}
  <div class="ir-card-main">
    <div class="ir-card-top"><span class="ir-card-name">${esc(a.name)}</span><span class="ir-card-slot">${boost}<span class="ir-age">${esc(a.age)}</span>${actions}</span></div>
    <div class="ir-card-line">${esc(a.line)}</div>${a.helpers ? `
    <div class="ir-card-helpers"><div class="ir-card-helper">↳ <b>Explore</b>Reading DESIGN.md</div></div>` : ''}
  </div>
</div>`;
      },
    },
    {
      id: 'quick-actions', group: 'Parts', title: 'Quick actions', usage: 'Quick actions',
      blurb: "The buttons on the card under the pointer. A toggle that is on stays, filled, as the state's one mark.",
      args: { paused: false, boosted: true },
      controls: { paused: { type: 'boolean', label: 'Paused' }, boosted: { type: 'boolean', label: 'Boosted' } },
      stories: [{ name: 'Boost on' }, { name: 'At rest', args: { boosted: false } }],
      render: (a) => `<div class="ir-card is-hover"><span class="ir-mark" data-state="${a.paused ? 'paused' : 'idle'}"></span><div class="ir-card-main"><div class="ir-card-top"><span class="ir-card-name">site-copy</span><span class="ir-card-slot"><span class="ir-actions"><button class="ir-act${a.paused ? ' is-on' : ''}" type="button" aria-label="${a.paused ? 'Go on' : 'Pause'}">${icon.pause}</button><button class="ir-act${a.boosted ? ' is-on' : ''}" type="button" aria-label="${a.boosted ? 'Stop boosting' : 'Boost'}">${icon.boost}</button><button class="ir-act" type="button" aria-label="Park">${icon.park}</button><button class="ir-act" type="button" aria-label="More">${icon.more}</button></span></span></div><div class="ir-card-line">${a.paused ? 'Paused' : 'Under the pointer'}</div></div></div>`,
    },
    {
      id: 'group-head', group: 'Parts', title: 'Group head', usage: 'Group head',
      args: { title: 'irie', count: 5, open: true, needs: 0 },
      controls: { title: { type: 'text' }, count: { type: 'number', min: 0, max: 99 }, open: { type: 'boolean', label: 'Open' }, needs: { type: 'number', min: 0, max: 9, label: 'Need you' } },
      stories: [{ name: 'Open' }, { name: 'Folded, one needs you', args: { title: 'Parked', count: 3, open: false, needs: 1 } }],
      render: (a) => `<button class="ir-group" type="button" aria-expanded="${!!a.open}"><span class="ir-group-chev">›</span><span class="ir-group-title">${esc(a.title)}</span>${!a.open && +a.needs > 0 ? `<span class="ir-group-need">· ${+a.needs} need${+a.needs === 1 ? 's' : ''} you</span>` : ''}<span class="ir-group-count">${+a.count}</span></button>`,
    },
    {
      id: 'task-bar', group: 'Parts', title: 'Task bar', usage: 'Task bar',
      args: { placeholder: 'Start a task: say what to do', text: '' },
      controls: { placeholder: { type: 'text' }, text: { type: 'text', label: 'Typed' } },
      stories: [{ name: 'Empty' }, { name: 'A task typed', args: { text: 'Make the tests pass on Linux' } }],
      render: (a) => `<div class="ir-taskbar">
  <textarea aria-label="New task" placeholder="${esc(a.placeholder)}">${esc(a.text)}</textarea>
  <button class="ir-start" type="button" aria-label="Start">${icon.start}</button>
</div>`,
    },
    {
      id: 'buttons', group: 'Parts', title: 'Buttons', usage: 'Buttons',
      blurb: 'One prominent button per view at most; every other button is soft or an icon. A button says a verb and its object.',
      args: { kind: 'prominent', label: 'Start the task' },
      controls: { kind: { type: 'select', options: ['prominent', 'soft', 'danger', 'icon', 'start'] }, label: { type: 'text' } },
      stories: [{ name: 'Prominent' }, { name: 'Soft', args: { kind: 'soft', label: 'Open its terminal' } }, { name: 'Danger', args: { kind: 'danger', label: 'Close “api-fix”' } }, { name: 'Icon', args: { kind: 'icon', label: 'Settings' } }, { name: 'Start', args: { kind: 'start', label: 'Start' } }],
      render(a) {
        if (a.kind === 'icon') return `<button class="ir-icon-button" type="button" aria-label="${esc(a.label)}">${icon.gear}</button>`;
        if (a.kind === 'start') return `<button class="ir-start" type="button" aria-label="${esc(a.label)}">${icon.start}</button>`;
        const cls = { prominent: ' is-prominent', soft: '', danger: ' is-danger' }[a.kind];
        return `<button class="ir-button${cls}" type="button">${esc(a.label)}</button>`;
      },
    },
    {
      id: 'segmented-control', group: 'Parts', title: 'Segmented control', usage: 'Segmented control',
      args: { options: 'List, Repo', picked: 0 },
      controls: { options: { type: 'text', label: 'Choices' }, picked: { type: 'number', min: 0, max: 5 } },
      stories: [{ name: 'Two views' }, { name: 'Three', args: { options: 'Light, Dark, Both', picked: 2 } }],
      render: (a) => `<div class="ir-seg" role="group" aria-label="View">${list(a.options).map((o, i) => `<button type="button" aria-pressed="${i === +a.picked}">${esc(o)}</button>`).join('')}</div>`,
    },
    {
      id: 'field', group: 'Parts', title: 'Field', usage: 'Field',
      args: { placeholder: 'Find a session', value: '' },
      controls: { placeholder: { type: 'text' }, value: { type: 'text' } },
      stories: [{ name: 'Empty' }, { name: 'Filled', args: { value: 'theme' } }],
      render: (a) => `<input class="ir-field" type="search" placeholder="${esc(a.placeholder)}" value="${esc(a.value)}" aria-label="${esc(a.placeholder)}">`,
    },
    {
      id: 'banner', group: 'Parts', title: 'Banner', usage: 'Banner and notice',
      args: { loud: false, text: 'Irie is getting ready: 2 of 3 done' },
      controls: { loud: { type: 'boolean', label: 'Needs you' }, text: { type: 'text' } },
      stories: [{ name: 'Quiet' }, { name: 'Loud', args: { loud: true, text: 'Another machine can’t be reached: its sessions wait' } }],
      render: (a) => `<div class="ir-banner${a.loud ? ' is-loud' : ''}">${esc(a.text)}</div>`,
    },
    {
      id: 'notice', group: 'Parts', title: 'Notice', usage: 'Banner and notice',
      args: { error: false, text: 'Parked “api-fix”', action: 'Undo' },
      controls: { error: { type: 'boolean', label: 'An error' }, text: { type: 'text' }, action: { type: 'text' } },
      stories: [{ name: 'With Undo' }, { name: 'An error', args: { error: true, text: 'Couldn’t start “docs”: the repo has no main branch', action: 'Try again' } }],
      render: (a) => `<div class="ir-notice${a.error ? ' is-error' : ''}">${mark(a.error ? 'failed' : 'idle')}<span class="ir-notice-text">${esc(a.text)}</span>${a.action ? `<button class="ir-notice-undo" type="button">${esc(a.action)}</button>` : ''}</div>`,
    },
    {
      id: 'details-panel', group: 'Parts', title: 'Details panel', usage: 'Details panel',
      args: { doing: 'go test ./internal/ui', repo: 'irie, on irie/theme-fix', started: '14:02, 4 minutes ago', used: '41k tokens' },
      controls: { doing: { type: 'text' }, repo: { type: 'text' }, started: { type: 'text' }, used: { type: 'text' } },
      stories: [{ name: 'A working session' }],
      render: (a) => `<div class="ir-panel">
  <dl>
    <dt>Doing</dt><dd><code>${esc(a.doing)}</code></dd>
    <dt>Repo</dt><dd>${esc(a.repo)}</dd>
    <dt>Started</dt><dd>${esc(a.started)}</dd>
    <dt>Used</dt><dd>${esc(a.used)}</dd>
  </dl>
</div>`,
    },
    {
      id: 'tab-row', group: 'Parts', title: 'Tab row', usage: 'Tab row',
      blurb: 'Folder tabs joined to their panel. A picked label keeps its width; a row that scrolls shows a fade and a chevron on the side where tabs hide.',
      args: { tabs: 'Helpers, Review, Notices', picked: 1, narrow: false },
      controls: { tabs: { type: 'text' }, picked: { type: 'number', min: 0, max: 12 }, narrow: { type: 'boolean', label: 'Narrow' } },
      stories: [
        { name: 'Three tabs' },
        { name: 'First picked', args: { tabs: 'Terminal, Browser, Code', picked: 0 } },
        { name: 'Overflowing', args: { tabs: 'Dock, Changes, Helpers, Plan, Review, Notices, Mail, Files, Find', picked: 0, narrow: true } },
      ],
      render(a) {
        const tabs = list(a.tabs), picked = Math.min(+a.picked, tabs.length - 1);
        const t = tabs.map((name, i) => `<button class="ir-tab" type="button" role="tab" aria-selected="${i === picked}" tabindex="${i === picked ? 0 : -1}" data-says="${esc(name)}: its panel."><span class="ir-tab-label" data-text="${esc(name)}">${esc(name)}</span></button>`).join('');
        return `<div class="ir-tabs${a.narrow ? ' narrow' : ''}">
  <div class="ir-tabs-row">
    <button class="ir-tabs-more is-start" type="button" aria-label="Previous tabs" tabindex="-1" hidden>${icon.chevStart}</button>
    <div class="ir-tablist" role="tablist" aria-label="Tabs">${t}</div>
    <button class="ir-tabs-more is-end" type="button" aria-label="Next tabs" tabindex="-1" hidden>${icon.chevEnd}</button>
  </div>
  <div class="ir-tabpanel" role="tabpanel">${esc(tabs[picked] || '')}: its panel.</div>
</div>`;
      },
    },
  ];
})();
