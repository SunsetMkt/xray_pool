<template>
  <n-button @click="showModal = true" text v-bind="$attrs" title="查看日志">
    <n-icon size="22"><document-text-outline /></n-icon>
  </n-button>
  <n-modal v-model:show="showModal">
    <div class="flex flex-col bg-white p-3" style="width: 80vw; height: 90vh; overflow: hidden">
      <div class="font-bold text-xl">日志</div>

      <n-divider class="!my-3" />

      <div class="flex-1 flex flex-col overflow-hidden">
        <log-viewer class="flex-1" :logs="logs" />
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { DocumentTextOutline } from '@vicons/ionicons5';
import LogViewer from '@/components/LogViewer.vue';
import { logs } from '@/composables/use-mqtt';

const showModal = ref(false);
</script>
