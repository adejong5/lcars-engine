<script>
  import LCARSPage from '$shared/components/LCARSPage.svelte';
  import DataCascade from '$shared/components/DataCascade.svelte';
  import NavButtons from '$shared/components/NavButtons.svelte';
  import SidePanels from '$shared/components/SidePanels.svelte';
  import LCARSBar from '$shared/components/LCARSBar.svelte';
  import Accordion from '$shared/components/Accordion.svelte';
  import LCARSButton from '$shared/components/LCARSButton.svelte';
  import ButtonRow from '$shared/components/ButtonRow.svelte';
  import ImageFrame from '$shared/components/ImageFrame.svelte';
  import LCARSFrame from '$shared/components/LCARSFrame.svelte';
  import BarChart from '$shared/components/BarChart.svelte';
  import Pillbox from '$shared/components/Pillbox.svelte';
  import Gallery from '$shared/components/Gallery.svelte';
  import { ha } from '$shared/ha.svelte.js';
  import { playBeep } from '$shared/sounds.js';

  // The LCARSPage section's live layout-switch demo:
  // 'standard' | 'ultra' (both columns) | 'col1' | 'col2'
  let layoutMode = $state('standard');
  let innerWidth = $state(0);

  // Live entities for the DataCascade demo in the top frame. Mock mode
  // fabricates these; against a real HA instance use real entity ids.
  const cascadeValues = [
    'sensor.warp_core_temp',
    'sensor.impulse_power',
    'sensor.shield_integrity',
    'sensor.hull_pressure',
    'sensor.life_support_o2',
    'sensor.deflector_output',
    'sensor.antimatter_flow',
    'sensor.plasma_temp',
    'sensor.eps_load',
    { entity: 'binary_sensor.red_alert', format: v => (v === 'on' ? '0001' : '0000') },
    'sensor.dilithium_matrix',
    'sensor.aux_power'
  ];

  // Offline placeholder images for the Gallery demo.
  const placeholder = (text, color = '#F1A827') =>
    'data:image/svg+xml;utf8,' + encodeURIComponent(
      `<svg xmlns="http://www.w3.org/2000/svg" width="340" height="200">
        <rect width="100%" height="100%" fill="#000" stroke="${color}"/>
        <text x="50%" y="50%" fill="${color}" font-family="sans-serif" font-size="20"
          text-anchor="middle" dominant-baseline="middle">${text}</text>
      </svg>`
    );
</script>

<svelte:window bind:innerWidth />

<LCARSPage
  layout={layoutMode === 'ultra' ? 'ultra' : 'standard'}
  column1={layoutMode === 'col1'}
  column2={layoutMode === 'col2'}
  banner="LCARS • COMPONENT DOCUMENTATION"
  topPanel={{ href: '../../' }}>

  {#snippet topFrame()}
    <DataCascade values={cascadeValues} />
    <NavButtons buttons={[
      { label: 'HOME', href: '../../' },
      { label: 'PAGE', href: '#page' },
      { label: 'BTNS', href: '#buttons' },
      { label: 'HA', href: '#home-assistant' }
    ]} />
  {/snippet}

  {#snippet sidebar()}
    <SidePanels panels={[
      { label: 'HOME', href: '../../' },
      { label: 'PAGE', href: '#page', hop: '-LAYOUT' },
      { label: 'CONTENT', href: '#bars', hop: '-WIDGETS' },
      { label: 'DATA', href: '#home-assistant', hop: '-LIVE' },
      { label: 'PADD', href: '../padd/', hop: '-THEME' }
    ]} />
  {/snippet}

  <h1>Component Documentation</h1>
  <h2>LCARS Engine &#149; Svelte 5 Components</h2>
  <p class="go-big">
    Every component in <span class="code">shared/components/</span> is documented below,
    and every section demonstrates the component it describes. This page itself is the
    <span class="code">LCARSPage</span> demo &#149; the cascade and nav buttons above are
    <span class="code">DataCascade</span> and <span class="code">NavButtons</span> &#149;
    the panels on the left are <span class="code">SidePanels</span>.
  </p>
  <p>
    Data link:
    <span class={ha.connected ? 'font-gold' : 'font-red'}>
      {ha.connected ? 'ONLINE' : 'OFFLINE'}
    </span>
    &#149; without <span class="code">VITE_HA_HOST</span> in
    <span class="code">.env.local</span> all live values on this page are mock data.
  </p>

  <!-- ════════════════════════ PAGE FRAME ════════════════════════ -->

  <span id="page"></span>
  <LCARSBar title="LCARSPage — The Page Frame" />
  <p>
    <span class="code">&lt;LCARSPage&gt;</span> renders the full LCARS shell: top frame,
    elbow bars, sidebar, content area, and footer. You are looking at it right now
    (classic theme, standard layout). Pair the <span class="code">theme</span> prop with
    the matching CSS import in your site's <span class="code">main.js</span>.
  </p>
  <ul class="lcars-list">
    <li><span class="code">theme</span> — 'classic' | 'nemesis' | 'lower-decks' | 'padd'; must match the imported theme CSS (padd = lower-decks-padd.css). See the <a href="../padd/">PADD site</a> for the padd theme live.</li>
    <li><span class="code">layout</span> — 'standard' (this page) | 'ultra'; ultra adds the two decorative left columns from the original ultra templates (classic and nemesis themes only).</li>
    <li><span class="code">column1</span> / <span class="code">column2</span> — pick the ultra columns individually: <span class="code">true</span> renders the default decorative content, a snippet renders your own. Setting either enables the ultra frame on its own; <span class="code">layout="ultra"</span> is shorthand for both.</li>
    <li><span class="code">banner</span> — top banner text.</li>
    <li><span class="code">topPanel</span> — {'{ label?, href? }'} for the big top-left button; omit href for a decorative beep.</li>
    <li><span class="code">panel2</span> — {'{ label, hop }'} for the small panel under it.</li>
    <li><span class="code">floorText</span> — extra footer line (styled by the nemesis theme).</li>
    <li><span class="code">sounds</span> — false to silence all beeps in the frame.</li>
  </ul>
  <p>
    Snippets: <span class="code">topFrame</span> (cascade/nav area),
    <span class="code">sidebar</span> (defaults to <span class="code">&lt;SidePanels /&gt;</span>),
    <span class="code">column1</span> / <span class="code">column2</span> (ultra only),
    <span class="code">children</span> (this content), and
    <span class="code">footer</span> (the attribution line always stays).
  </p>
  <p>
    Try it — switch this page's layout live. Note the theme CSS is responsive: the ultra
    columns slide away below a 1680px-wide viewport and are hidden entirely below 1500px,
    so widen the window to see them.
  </p>
  <ButtonRow>
    <LCARSButton label="Standard"
      color={layoutMode === 'standard' ? 'butterscotch' : undefined}
      onclick={() => (layoutMode = 'standard')} />
    <LCARSButton label="Ultra (both)"
      color={layoutMode === 'ultra' ? 'butterscotch' : undefined}
      onclick={() => (layoutMode = 'ultra')} />
    <LCARSButton label="Column 1 only"
      color={layoutMode === 'col1' ? 'butterscotch' : undefined}
      onclick={() => (layoutMode = 'col1')} />
    <LCARSButton label="Column 2 only"
      color={layoutMode === 'col2' ? 'butterscotch' : undefined}
      onclick={() => (layoutMode = 'col2')} />
  </ButtonRow>
  {#if layoutMode !== 'standard' && innerWidth < 1680}
    <p class="font-red blink">
      An ultra layout is active, but your window is {innerWidth}px wide — the theme CSS
      hides the ultra columns below 1680px.
    </p>
  {/if}

  <LCARSBar title="DataCascade" heading="h3" />
  <p>
    The animated number cascade in the top frame of this page. Purely random by default;
    pass <span class="code">values</span> to pin live Home Assistant entities to stable
    cells — the cascade above carries twelve mock sensors.
  </p>
  <ul class="lcars-list">
    <li><span class="code">values</span> — entity_id strings or {'{ entity, format, cls }'} objects; <span class="code">format(state)</span> maps the display text, <span class="code">cls(state)</span> adds conditional classes (e.g. <span class="code">'font-red blink'</span> above a threshold); <span class="code">{'{ text }'}</span> makes a static cell (column labels, '' for blanks). Values fill column-major. Entity cells show random digits until the state arrives, '----' if unavailable.</li>
    <li><span class="code">columns</span> — total column count to pad to with random filler (default 24).</li>
    <li><span class="code">filler</span> — exact filler column count instead; 0 = values only.</li>
    <li><span class="code">theme</span> — matches the page theme; sets row pattern and wrapper class (lower-decks/padd get hidden digits and dark filler columns).</li>
    <li><span class="code">frozen</span> — disable the color-cycling animation.</li>
    <li><span class="code">pattern</span> / <span class="code">wrapperClass</span> — low-level overrides for custom CSS.</li>
  </ul>

  <LCARSBar title="NavButtons" heading="h3" />
  <p>
    The button group beside the cascade, top right. Takes
    <span class="code">buttons</span> as a count (auto-labeled '01'…) or an array of
    {'{ label?, href? }'} — the ones above link to sections of this page.
    <span class="code">sounds={'{false}'}</span> disables the beep.
  </p>

  <LCARSBar title="SidePanels" heading="h3" />
  <p>
    The colored panel stack on the left of this page. Takes
    <span class="code">panels</span> as a count (auto-numbered filler, default 7) or an
    array of {'{ label?, hop?, href? }'} — with an href the panel becomes a link, like
    the HOME and PADD panels here. <span class="code">bottomPanel</span> configures the
    block pinned at the bottom (false to omit), and <span class="code">theme="padd"</span>
    switches to the shorter padd panel cycle.
  </p>

  <!-- ════════════════════════ CONTENT COMPONENTS ════════════════════════ -->

  <span id="bars"></span>
  <LCARSBar title="LCARSBar — Divider Bars" />
  <p>
    Every section heading on this page is an <span class="code">LCARSBar</span>.
    With <span class="code">title</span> it renders the text bar; without it, a plain
    divider:
  </p>
  <LCARSBar />
  <p>
    <span class="code">heading</span> picks the tag ('h2' default, 'h3', 'h4', 'span')
    and <span class="code">end</span> right-aligns the title:
  </p>
  <LCARSBar title="heading='h4'" heading="h4" />
  <LCARSBar title="end aligned" heading="h3" end />
  <p>
    One-off colors via the <span class="code">style</span> prop, overriding
    <span class="code">--lcars-bar-color</span> and friends:
  </p>
  <LCARSBar title="custom color" heading="h3" style="--lcars-bar-color: var(--african-violet)" />

  <span id="buttons"></span>
  <LCARSBar title="LCARSButton + ButtonRow" />
  <p>
    <span class="code">LCARSButton</span> renders an <span class="code">&lt;a&gt;</span>
    when <span class="code">href</span> is set, otherwise a beep-only
    <span class="code">&lt;button&gt;</span>. <span class="code">color</span> applies any
    theme color as <span class="code">button-&#123;color&#125;</span>.
    <span class="code">ButtonRow</span> lays out a group (.buttons), with an optional
    <span class="code">justify</span> of 'center', 'flex-end', 'space-between',
    'space-around', or 'space-evenly'.
  </p>
  <ButtonRow>
    <LCARSButton label="Default" />
    <LCARSButton label="href link" href="#buttons" />
    <LCARSButton label="color='mars'" color="mars" />
    <LCARSButton label="color='butterscotch'" color="butterscotch" />
  </ButtonRow>
  <p>justify='center':</p>
  <ButtonRow justify="center">
    <LCARSButton label="One" />
    <LCARSButton label="Two" />
  </ButtonRow>
  <p>
    Variants: <span class="code">variant='sidebar'</span> for the flat sidebar style, and
    <span class="code">variant='panel'</span> with <span class="code">pan</span> (0, 4–8, 10)
    to mimic an established panel's height and color:
  </p>
  <LCARSButton variant="sidebar" label="Sidebar Variant" color="almond-creme" />
  <LCARSButton variant="panel" pan={5} label="Panel Mimic pan=5" />

  <LCARSBar title="Pillbox" />
  <p>
    Two-column grid of pill buttons from the ultra layout (colors come from the theme's
    designed nth-child slots). <span class="code">pills</span> is a count (auto-labeled
    'J-001'…) or an array of {'{ label?, href?, spacer? }'};
    <span class="code">variant={'{2}'}</span> selects the second color set, and
    <span class="code">spacer: true</span> renders an empty slot like the original
    template:
  </p>
  <div class="flexbox">
    <Pillbox />
    <Pillbox variant={2} pills={[
      { label: 'F12-22', href: '#buttons' }, { label: 'G24-22' }, { spacer: true },
      { label: 'H-07AM' }, { label: 'I50-72' }, { label: 'J5369' }
    ]} />
  </div>

  <LCARSBar title="Accordion" />
  <p>
    Collapsible section with the template's animated toggle.
    <span class="code">title</span> sets the header, <span class="code">open</span> starts
    it expanded, <span class="code">limitWidth</span> caps it at 1240px:
  </p>
  <Accordion title="Closed by default — click me" limitWidth>
    <p>Containment field stable. Dilithium matrix alignment at
      {ha.state('sensor.dilithium_matrix') ?? '—'} microns (live value).</p>
  </Accordion>
  <Accordion title="open — starts expanded" open limitWidth>
    <p>Any markup works in here, including other components.</p>
  </Accordion>

  <LCARSBar title="BarChart" />
  <p>
    Horizontal bar chart for displaying sensor values reactively. Compute
    <span class="code">items</span> with <span class="code">$derived()</span> in the parent
    so it updates when HA entities change. Place it inside an
    <span class="code">LCARSFrame</span> for a framed display.
    Each item: <span class="code">label</span>, <span class="code">value</span>,
    <span class="code">max</span>, optional <span class="code">unit</span>,
    <span class="code">color</span>, <span class="code">decimals</span>.
  </p>
  <div class="flexbox">
    <LCARSFrame>
      <BarChart items={[
        { label: 'CPU',    value: 42,  max: 100, unit: '%', color: 'var(--butterscotch)', decimals: 0 },
        { label: 'MEMORY', value: 68,  max: 100, unit: '%', color: 'var(--african-violet)', decimals: 0 },
        { label: 'GPU',    value: 12,  max: 100, unit: '%', color: 'var(--ice)', decimals: 0 },
        { label: 'DISK',   value: 41,  max: 100, unit: '%', color: 'var(--mars)', decimals: 0 },
        { label: 'NET',    value: 0.2, max: 1,   unit: ' Gb/s', color: 'var(--gold)', decimals: 2 },
      ]} />
    </LCARSFrame>
  </div>

  <LCARSBar title="LCARSFrame" />
  <p>
    Decorative framed display with animated "audio level" bars.
    <span class="code">horizontal</span> rotates the bars,
    <span class="code">lines</span> changes the bar count (default 16), and passing
    children replaces the bars with your own content in the center cell:
  </p>
  <div class="flexbox">
    <LCARSFrame />
    <LCARSFrame horizontal />
    <LCARSFrame>
      <span class="go-big">47</span>
    </LCARSFrame>
  </div>

  <LCARSBar title="ImageFrame" />
  <p>
    Titled, bordered display frame (.image-frame). Pass <span class="code">src</span> /
    <span class="code">alt</span> for an image, or children for arbitrary content;
    recolor with theme vars through the <span class="code">style</span> prop:
  </p>
  <div class="flexbox">
    <ImageFrame title="src image" src={placeholder('IMAGE')} alt="placeholder" />
    <ImageFrame title="children content" style="--title-color: var(--butterscotch)">
      <div style="width:300px;height:180px;display:flex;align-items:center;justify-content:center;">
        <span class="go-big blink">NO SIGNAL</span>
      </div>
    </ImageFrame>
  </div>

  <LCARSBar title="Gallery" />
  <p>
    Image grid (flagged experimental in the theme CSS).
    <span class="code">images</span> takes {'{ src, alt?, caption?, href? }'} — an href
    makes the image a link with a hover border — and <span class="code">thumbs</span>
    sizes entries as 340px thumbnails:
  </p>
  <Gallery thumbs images={[
    { src: placeholder('ALPHA'), caption: 'Caption text', href: '#bars' },
    { src: placeholder('BRAVO', '#cc6666'), caption: 'No link' },
    { src: placeholder('GAMMA', '#9c82c4') }
  ]} />

  <!-- ════════════════════════ DATA + SOUND ════════════════════════ -->

  <span id="home-assistant"></span>
  <LCARSBar title="Home Assistant Data" />
  <p>
    <span class="code">$shared/ha.svelte.js</span> is a reactive store over the HA
    WebSocket API. Call <span class="code">initHA(&#123;...&#125;)</span> once in your
    site's <span class="code">main.js</span> — with
    <span class="code">VITE_HA_HOST</span> / <span class="code">VITE_HA_TOKEN</span> from
    <span class="code">.env.local</span> for a real connection, or
    <span class="code">&#123; mock: true &#125;</span> for fabricated offline data (what
    this page uses unless configured).
  </p>
  <ul class="lcars-list">
    <li><span class="code">ha.connected</span> — reactive link status, currently
      <span class={ha.connected ? 'font-gold' : 'font-red'}>{ha.connected ? 'ONLINE' : 'OFFLINE'}</span></li>
    <li><span class="code">ha.state(id)</span> — live state string: warp core temp is
      <span class="font-gold">{ha.state('sensor.warp_core_temp') ?? '—'}</span> right now</li>
    <li><span class="code">ha.entity(id)</span> — full state object with attributes</li>
    <li><span class="code">ha.callService(domain, service, data, target)</span> — call any HA service</li>
  </ul>
  <ButtonRow>
    <LCARSButton label="callService demo"
      onclick={() => ha.callService('light', 'toggle', {}, { entity_id: 'light.bridge' })} />
  </ButtonRow>
  <p>
    Reactivity is per-entity: a state change re-renders only the expressions reading that
    entity. In mock mode any entity you read is lazily invented and drifts every two
    seconds — watch the values above change.
  </p>

  <LCARSBar title="Sounds" />
  <p>
    <span class="code">$shared/sounds.js</span> exports
    <span class="code">playBeep(n)</span> and
    <span class="code">playBeepAndGo(n, url)</span> (beep, then navigate) using the four
    template beeps. Frame and button components beep on their own; every interactive
    component takes <span class="code">sounds={'{false}'}</span> to opt out.
  </p>
  <ButtonRow>
    <LCARSButton label="Beep 1" sounds={false} onclick={() => playBeep(1)} />
    <LCARSButton label="Beep 2" sounds={false} onclick={() => playBeep(2)} />
    <LCARSButton label="Beep 3" sounds={false} onclick={() => playBeep(3)} />
    <LCARSButton label="Beep 4" sounds={false} onclick={() => playBeep(4)} />
  </ButtonRow>

  <!-- ════════════════════════ THEME CSS EXTRAS ════════════════════════ -->

  <LCARSBar title="Theme CSS Utilities" />
  <p>
    The theme stylesheet ships utility classes that work on plain markup — no component
    needed:
  </p>
  <ul class="lcars-list">
    <li>Standard <span class="code">lcars-list</span> bullet</li>
    <li class="bullet-orange">Bullet color variants, e.g. <span class="code">bullet-orange</span></li>
  </ul>
  <p>
    <span class="blink">blink</span> &#149;
    <span class="pulse">pulse</span> &#149;
    <span class="go-big">go-big</span> &#149;
    <span class="code">code 0x47</span> &#149;
    <span class="uppercase">uppercase</span> &#149;
    font colors like <span class="font-gold">font-gold</span>,
    <span class="font-african-violet">font-african-violet</span>,
    <span class="font-red">font-red</span>
  </p>
  <blockquote>blockquote — things are only impossible until they're not.</blockquote>
  <p>
    <span class="code">flexbox</span> wraps children in a flex row (used for the frame
    and image demos above).
  </p>

  {#snippet footer()}
    Content &copy; 2026 LCARS Creator.<br>
  {/snippet}

</LCARSPage>
