import { computed, reactive } from 'vue';
import { useIntervalFn } from '@vueuse/core';
import type { ApiResponseSystemStatus } from '@/interfaces/common';
import CommonApi from '@/apis/CommonApi';

export interface SystemState {
  status: ApiResponseSystemStatus | null;
}

export const systemState = reactive<SystemState>({
  status: null,
});

export const isSetup = computed(() => systemState.status?.is_setup);

export const getSystemStatus = async () => {
  const [res, err] = await CommonApi.getSystemStatus();
  if (err !== null) {
    window.$message.error(err.message);
    return;
  }
  systemState.status = res;
};

export const useSystem = () => {
  getSystemStatus();
  useIntervalFn(getSystemStatus, 1000);
};
