<template>
  <div ref="logViewer" class="log-viewer">
    <div v-for="item in logs" :key="item" class="log-line">{{ item }}</div>
  </div>
</template>

<script setup lang="ts">
import { watch, ref, nextTick, onMounted } from 'vue';

export interface LogViewerProps {
  logs: string[];
}

const props = withDefaults(defineProps<LogViewerProps>(), {
  logs: () => [],
});

const logViewer = ref<HTMLInputElement | null>(null);

watch(
  () => props.logs.length,
  () => {
    if (logViewer.value === null) return;
    const element = logViewer.value;
    // console.log(element.scrollTop, element.clientHeight, element.scrollHeight);
    // 如果当前正处于底部，则自动滚动
    if (element.scrollTop + element.clientHeight >= element.scrollHeight - 20) {
      nextTick(() => {
        if (logViewer.value === null) return;
        logViewer.value.scrollTo(0, element?.scrollHeight);
      });
    }
  }
);

onMounted(() => {
  logViewer.value?.scrollTo(0, logViewer.value?.scrollHeight);
});
</script>

<style lang="scss" scoped>
.log-viewer {
  overflow: auto;
  width: 100%;
  height: 100%;
}
.log-line {
  white-space: nowrap;
}
</style>
