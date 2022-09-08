<template>
  <div>
    <div class="flex flex-row items-center">
      <div class="font-bold">节点列表</div>
      <n-button class="ml-2" text type="success" @click="updateNodeList" :loading="updateLoading" title="刷新列表">
        <template #icon>
          <n-icon><refresh /></n-icon>
        </template>
      </n-button>
    </div>
    <n-list class="border-1 mt-2" hoverable clickable style="height: 200px; overflow-y: auto">
      <n-list-item v-for="item in state.nodeList" :key="item.name">
        <div class="flex row">
          <div class="flex-1">{{ item.name }}</div>
          <div class="text-gray-500">{{ item.proto_model }}</div>
        </div>
      </n-list-item>
    </n-list>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { Refresh } from '@vicons/ionicons5';
import { getNodeList, subscribeState as state, updateLoading, updateNodeList } from '@/composables/use-subscribe';

onMounted(() => {
  getNodeList();
});
</script>
