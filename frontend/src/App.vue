<template>
  <n-message-provider>
    <naive-ui-global-message-handler />
  </n-message-provider>

  <n-config-provider :locale="zhCN" :date-locale="dateZhCN">
    <RouterView />
  </n-config-provider>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router';
import { zhCN, dateZhCN } from 'naive-ui';
import NaiveUiGlobalMessageHandler from '@/components/NaiveUiGlobalMessageHandler.vue';
</script>
