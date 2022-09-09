<template>
  <div>
    <n-button v-if="isStopped" type="primary" @click="startProxyPool" :loading="loading">启动代理池</n-button>
    <n-button v-else type="error" @click="stopProxyPool" :loading="loading">停止代理池</n-button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { isStopped } from '@/composables/use-proxy-pool';
import ProxyPoolApi from '@/apis/ProxyPoolApi';
import { settingsState } from '@/composables/use-settings';

const loading = ref(false);

const startProxyPool = async () => {
  loading.value = true;
  if (settingsState.settings === null) return;
  const [, err] = await ProxyPoolApi.start({
    target_site_url: settingsState.settings.test_url,
  });
  if (err) {
    window.$message.error(err.message);
    return;
  }
  window.$message.success('代理池已启动');
  loading.value = false;
};

const stopProxyPool = async () => {
  loading.value = true;
  const [, err] = await ProxyPoolApi.stop();
  if (err) {
    window.$message.error(err.message);
    return;
  }
  window.$message.success('代理池已停止');
  loading.value = false;
};
</script>
