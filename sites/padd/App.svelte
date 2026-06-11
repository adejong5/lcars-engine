<script>
  import LCARSPage from '$shared/components/LCARSPage.svelte';
  import DataCascade from '$shared/components/DataCascade.svelte';
  import NavButtons from '$shared/components/NavButtons.svelte';
</script>

<LCARSPage theme="padd" banner="LCARS 57436.2">

  {#snippet topFrame()}
    <DataCascade theme="padd" columns={18} />
    <NavButtons />
  {/snippet}

  <h1>Hello</h1>
  <h2>Welcome to LCARS &#149; Lower Decks PADD Theme</h2>
  <h3 class="font-radioactive">Svelte Edition</h3>
  <h4>Replace This Content With Your Own</h4>
  <p class="go-big">Live long and prosper.</p>

  {#snippet footer()}
    Content &copy; 2026 LCARS Creator.<br>
  {/snippet}

</LCARSPage>
