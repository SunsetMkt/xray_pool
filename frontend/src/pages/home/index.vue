<template>
  <div>
    <header class="border-1 p-2">XrayPool 运行状态：Running</header>

    <settings-xray-pool class="border-1 p-2 mt-2" />

    <load-balance-panel class="border-1 p-2 mt-2" />

    <div class="mt-2">
      <btn-modal-settings-subscribe />
      <btn-modal-settings-advanced class="ml-2" />
    </div>
  </div>
</template>

<script setup lang="ts">
import SettingsXrayPool from '@/pages/home/SettingsXrayPool.vue';
import LoadBalancePanel from '@/pages/home/LoadBalancePanel.vue';
import { useProxyPool } from '@/composables/use-proxy-pool';
import BtnModalSettingsSubscribe from '@/pages/home/BtnModalSettingsSubscribe.vue';
import BtnModalSettingsAdvanced from '@/pages/home/BtnModalSettingsAdvanced.vue';

useProxyPool();
</script>
