<template>
  <n-button @click="showModal = true" type="primary" v-bind="$attrs"> 高级设置 </n-button>
  <n-modal v-model:show="showModal">
    <n-card style="width: 600px" title="高级设置" :bordered="false" size="huge" role="dialog" aria-modal="true">
      设置
    </n-card>
  </n-modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue';

const showModal = ref(false);
</script>
