<template>
  <div>
    <n-button v-if="isStopped" type="primary">启动代理池</n-button>
    <n-button v-else type="error">停止代理池</n-button>
  </div>
</template>

<script setup>
import { isStopped } from '@/composables/use-proxy-pool';
</script>
