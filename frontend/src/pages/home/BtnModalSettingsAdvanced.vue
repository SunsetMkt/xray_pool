<template>
  <n-button @click="showModal = true" type="primary" v-bind="$attrs" :disabled="!isStopped"> 高级设置 </n-button>
  <n-modal v-model:show="showModal">
    <div class="flex flex-col bg-white p-3" style="width: 600px; max-height: 90vh; overflow: hidden">
      <div class="font-bold text-xl">高级设置</div>

      <n-divider class="!my-3" />

      <div class="flex-1 overflow-auto">
        <div class="p-2">
          <div class="font-bold text-lg">Core 基础设置</div>
          <settings-core class="mt-2" />
        </div>

        <div class="p-2 mt-4">
          <div class="font-bold text-lg">路由设置</div>
          <settings-routing class="mt-2" />
        </div>
      </div>
    </div>
  </n-modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import SettingsCore from '@/pages/home/SettingsCore.vue';
import SettingsRouting from '@/pages/home/SettingsRouting.vue';
import { isStopped } from '@/composables/use-proxy-pool';

const showModal = ref(false);
</script>
