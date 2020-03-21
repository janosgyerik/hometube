<script>
  import { files } from "./stores.js";

  let items;

  const unsubscribe = files.subscribe(files => {
    items = files;
  });
</script>

{#if items.length > 0}
  <main>
    <div class="container">
      <div class="text-center">
        <h2>Downloaded files</h2>
      </div>

      <ul class="list-group">
        {#each items as f}
          <li class="list-group-item">{f.filename}</li>
        {/each}
      </ul>
    </div>
  </main>
{/if}
