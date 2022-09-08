<template>
  <div>
    <div class="font-bold">节点列表</div>

    <div v-if="state.nodeList?.length === 0" class="text-gray-500">当前没有可用节点</div>
    <n-list v-else class="border-1 mt-2 overflow-auto" hoverable clickable style="max-height: 50vh">
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
import { getNodeList, subscribeState as state } from '@/composables/use-subscribe';

onMounted(() => {
  getNodeList();
});
</script>
