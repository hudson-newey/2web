<script setup lang="ts">
import { ref } from "vue";
const count = ref(0);
</script>

<template>
  <h1 id="outTarget">{{ count }}</h1>

  <button id="inTarget" @click="count++">Increment</button>
  <button @click="count--">Decrement</button>
</template>
