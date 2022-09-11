<template>
  <n-button @click="showModal = true" type="primary" :disabled="!isStopped"> 订阅管理 </n-button>
  <n-modal v-model:show="showModal">
    <div class="flex flex-col bg-white p-3" style="width: 600px; max-height: 90vh; overflow: hidden">
      <div class="font-bold text-xl">订阅管理</div>

      <n-divider class="!my-3" />

      <div class="flex-1 overflow-y-auto overflow-x-hidden">
        <settings-subscribe />

        <node-list-panel class="mt-2" />
      </div>
    </div>
  </n-modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import SettingsSubscribe from '@/pages/home/SettingsSubscribe.vue';
import NodeListPanel from '@/pages/home/NodeListPanel.vue';
import { isStopped } from '@/composables/use-proxy-pool';

const showModal = ref(false);
</script>
