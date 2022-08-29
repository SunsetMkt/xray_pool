<template>
  <div>
    <div class="font-bold">节点列表</div>
    <n-list class="border-1 mt-2" hoverable clickable>
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
import { getNodeList, setupState as state } from '@/pages/setup/use-setup';

onMounted(() => {
  getNodeList();
});
</script>
