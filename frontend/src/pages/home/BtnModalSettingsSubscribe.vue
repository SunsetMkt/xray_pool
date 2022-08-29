<template>
  <n-button @click="showModal = true" type="primary"> 订阅管理 </n-button>
  <n-modal v-model:show="showModal">
    <n-card style="width: 600px" title="订阅管理" :bordered="false" size="huge" role="dialog" aria-modal="true">
      <settings-subscribe />
    </n-card>
  </n-modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import SettingsSubscribe from '@/pages/home/SettingsSubscribe.vue';

const showModal = ref(false);
</script>
