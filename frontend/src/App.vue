<template>
  <n-message-provider>
    <naive-ui-global-message-handler />
  </n-message-provider>
  <RouterView />
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router';
import NaiveUiGlobalMessageHandler from '@/components/NaiveUiGlobalMessageHandler.vue';
</script>
