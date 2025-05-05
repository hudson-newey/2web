<script setup lang="ts">
import { ref } from "vue";
const count = ref(0);
</script>

<template>
  <h1>{{ count }}</h1>

  <button @click="count++">Increment</button>
  <button @click="count--">Decrement</button>
</template>
