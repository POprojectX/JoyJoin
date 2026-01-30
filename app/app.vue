<script setup lang="ts">
import { useCounterStore } from '../stores/counter'
const counter = useCounterStore()
</script>

<template>
  <div class="p-6">
    <button class="px-3 py-2 bg-black text-white rounded" @click="counter.inc()">
      Count: {{ counter.n }}
    </button>
  </div>
</template>
