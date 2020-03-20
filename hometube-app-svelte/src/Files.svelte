<script>
  var filenames = ["example1.avi", "example2.avi"];
</script>

{#if filenames.length > 0}
  <main>
    <div class="container">
      <div class="text-center">
        <h2>Downloaded files</h2>
      </div>

      <ul class="list-group">
        {#each filenames as fn}
          <li class="list-group-item">{fn}</li>
        {/each}
      </ul>
    </div>
  </main>
{/if}
