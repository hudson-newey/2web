<script lang="ts">
  let count = $state(0);
</script>

<h1>{ count }</h1>

<button onclick={() => count++}>Increment</button>
<button onclick={() => count++}>Decrement</button>
