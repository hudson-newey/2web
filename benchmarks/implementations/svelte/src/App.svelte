<script lang="ts">
  let count = $state(0);
</script>

<h1 id="outTarget">{ count }</h1>

<button id="inTarget" onclick={() => count++}>Increment</button>
<button onclick={() => count--}>Decrement</button>
