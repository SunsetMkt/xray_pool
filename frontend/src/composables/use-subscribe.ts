import { reactive } from 'vue';
import type { SubscribeItem, SubscribeNodeItem } from '@/interfaces/subscribe';
import SubscribeApi from '@/apis/SubscribeApi';

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
