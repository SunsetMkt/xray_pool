import { reactive, ref } from 'vue';
import type { SubscribeItem, SubscribeNodeItem } from '@/interfaces/subscribe';
import SubscribeApi from '@/apis/SubscribeApi';
import ProxyPoolApi from '@/apis/ProxyPoolApi';

export interface SubscribeState {
  subscribeList: SubscribeItem[] | null;
  nodeList: SubscribeNodeItem[] | null;
}

export const subscribeState = reactive<SubscribeState>({
  subscribeList: null,
  nodeList: null,
});

export const getSubscribeList = async () => {
  const [res] = await SubscribeApi.getList();
  if (res === null) return;
  subscribeState.subscribeList = res.subscribe_list;
};

export const getNodeList = async () => {
  const [res] = await SubscribeApi.getNodeList();
  if (res === null) return;
  subscribeState.nodeList = res.node_info_list;
};

export const updateLoading = ref(false);
export const updateNodeList = async () => {
  updateLoading.value = true;
  const [, err] = await ProxyPoolApi.updateNodeList();
  if (err === null) {
    await getNodeList();
  } else {
    window.$message.error(err.message);
  }
  updateLoading.value = false;
};
