<script>
  import { items } from "./stores.js";

  let localItems;

  const unsubscribe = items.subscribe(items => {
    localItems = items;
  });
</script>

{#if localItems.length > 0}
  <main>
    <div class="container">
      <div class="text-center">
        <h2>Downloaded ({localItems.length})</h2>
      </div>

      <ul class="list-group">
        {#each localItems as item}
          <li class="list-group-item">{item.name}</li>
        {/each}
      </ul>
    </div>
  </main>
{/if}
