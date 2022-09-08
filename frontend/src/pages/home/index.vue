<template>
  <div>
    <header class="border-1 p-2 flex flex-row justify-between items-center">
      <div>
        XrayPool 运行状态：
        <span v-if="isRunning" class="text-green-500">运行中</span>
        <span v-else-if="isStarting" class="text-gray-500">正在启动...</span>
        <span v-else-if="isStopped" class="text-red-500">已停止</span>
        <span v-else>未知</span>
      </div>
      <div>
        <proxy-pool-operations />
      </div>
    </header>

    <settings-xray-pool class="border-1 p-2 mt-2" />

    <div class="mt-2">
      <btn-modal-settings-subscribe />
      <btn-modal-settings-advanced class="ml-2" />
    </div>

    <load-balance-panel v-if="isRunning" class="border-1 p-2 mt-2" />
  </div>
</template>

<script setup lang="ts">
import SettingsXrayPool from '@/pages/home/SettingsXrayPool.vue';
import LoadBalancePanel from '@/pages/home/LoadBalancePanel.vue';
import { isRunning, isStarting, isStopped, useProxyPool } from '@/composables/use-proxy-pool';
import BtnModalSettingsSubscribe from '@/pages/home/BtnModalSettingsSubscribe.vue';
import BtnModalSettingsAdvanced from '@/pages/home/BtnModalSettingsAdvanced.vue';
import ProxyPoolOperations from '@/pages/home/ProxyPoolOperations.vue';
import { useSettings } from '@/composables/use-settings';

useProxyPool();
useSettings();
</script>
