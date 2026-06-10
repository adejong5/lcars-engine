<script>
  import LCARSPage from '$shared/components/LCARSPage.svelte';
  import DataCascade from '$shared/components/DataCascade.svelte';
  import NavButtons from '$shared/components/NavButtons.svelte';
  import SidePanels from '$shared/components/SidePanels.svelte';
  import { ha } from '$shared/ha.svelte.js';

  // Entities shown in the data cascade. With mock mode these are fabricated;
  // against a real HA instance, use your actual entity ids.
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
</script>

<LCARSPage banner="LCARS • TEST SITE" topPanel={{ href: '/' }}>

  {#snippet topFrame()}
    <DataCascade values={cascadeValues} />
    <NavButtons buttons={[
      { label: 'HOME', href: '/' },
      { href: '#systems' },
      {},
      {}
    ]} />
  {/snippet}

  {#snippet sidebar()}
    <SidePanels panels={[
      { label: 'HOME', href: '/' },
      { href: '#' },
      { label: 'ENGINEERING', hop: '-NCC1701' },
      {},
      { label: 'OPS' }
    ]} />
  {/snippet}

  <h1>Hello</h1>
  <h2>Welcome to LCARS &#149; Classic Theme &#149; Standard Layout</h2>
  <h3 class="font-gold">Svelte Edition</h3>
  <h4>Replace This Content With Your Own</h4>
  <p class="go-big">Live long and prosper.</p>
  <p>
    Data link:
    <span class={ha.connected ? 'font-gold' : 'font-red'}>
      {ha.connected ? 'ONLINE' : 'OFFLINE'}
    </span>
    &#149; warp core temp:
    <span class="font-gold">{ha.state('sensor.warp_core_temp') ?? '—'}</span>
  </p>

  {#snippet footer()}
    Content &copy; 2026 LCARS Creator.<br>
  {/snippet}

</LCARSPage>
