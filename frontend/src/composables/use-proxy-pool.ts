import { computed, onMounted, reactive } from 'vue';
import type { ApiResponseProxyList } from '@/interfaces/proxy-pool';
import ProxyPoolApi from '@/apis/ProxyPoolApi';
import { settingsState } from '@/composables/use-settings';

export interface ProxyPoolState {
  proxyList: ApiResponseProxyList | null;
}

export const proxyPoolState = reactive<ProxyPoolState>({
  proxyList: null,
});

export const isStarting = computed(() => proxyPoolState.proxyList?.status === 'starting');
export const isRunning = computed(() => proxyPoolState.proxyList?.status === 'running');
export const isStopped = computed(() => proxyPoolState.proxyList?.status === 'stopped');

export const getProxyList = async () => {
  const [res] = await ProxyPoolApi.getStatus();
  if (res === null) return;
  proxyPoolState.proxyList = res;
};

export const startProxyPool = async () => {
  const [, err] = await ProxyPoolApi.start({ target_site_url: settingsState.model?.test_url });
  if (err !== null) {
    window.$message.error(err.message);
  }
  window.$message.success('启动成功');
};

export const stopProxyPool = async () => {
  const [, err] = await ProxyPoolApi.stop();
  if (err !== null) {
    window.$message.error(err.message);
  }
  window.$message.success('停止成功');
};

export const useProxyPool = () => {
  onMounted(getProxyList);
};
