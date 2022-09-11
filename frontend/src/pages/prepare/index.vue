<template>
  <div>
    <div class="text-center font-bold text-4xl">Xray Pool</div>

    <div class="mt-2">
      当前系统架构为：<span class="font-bold">{{ systemState.status?.os }} - {{ systemState.status?.arch }}</span>
    </div>
    <div>
      要运行Xray Pool，请先下载对应平台的以下程序，并解压至
      <span class="font-bold" v-if="systemState.status?.os === 'win32'">程序所在目录的base_things</span>
      <span class="font-bold" v-else-if="systemState.status?.os === 'darwin'"
        >/Users/&lt;user&gt;/.config/xray_pool/base_things</span
      >
      <span class="font-bold" v-else-if="systemState.status?.os === 'linux'">/config/base_things</span>
      文件夹（如果文件夹不存在则手动创建）
    </div>
    <ul>
      <li>
        1、Xray：<a href="https://github.com/XTLS/Xray-core/releases">https://github.com/XTLS/Xray-core/releases</a>
      </li>
      <li>
        2、glider：<a href="https://github.com/nadoo/glider/releases">https://github.com/nadoo/glider/releases</a>
      </li>
    </ul>

    <div>
      更多详细说明请参考：<a href="https://github.com/allanpk716/xray_pool">https://github.com/allanpk716/xray_pool</a>
    </div>

    <div class="mt-2">
      <div v-if="systemState.status?.glider_downloaded" class="text-green-500 flex flex-row items-center gap-x-2">
        <n-icon size="30">
          <checkmark-circle />
        </n-icon>
        <div>已准备好glider</div>
      </div>
      <div v-else class="flex flex-row items-center gap-x-4">
        <n-spin :size="22" />
        <div>正则检测glider状态</div>
      </div>

      <div v-if="systemState.status?.xray_downloaded" class="text-green-500 flex flex-row items-center gap-x-2 mt-2">
        <n-icon>
          <checkmark-circle />
        </n-icon>
        <div>已准备好xray</div>
      </div>
      <div v-else class="flex flex-row items-center gap-x-4">
        <n-spin :size="22" />
        <div>正则检测xray状态</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { watch } from 'vue';
import { CheckmarkCircle } from '@vicons/ionicons5';
import { isSetup, systemState, useSystem } from '@/composables/use-system';

const router = useRouter();

useSystem();

watch(
  () => isSetup.value,
  () => {
    if (isSetup.value) {
      router.push('/');
    }
  }
);
</script>

<style lang="scss">
a {
  color: #3490dc;
}
</style>
