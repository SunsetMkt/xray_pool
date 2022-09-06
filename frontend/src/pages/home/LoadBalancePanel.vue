<template>
  <div>
    <div class="flex row justify-between">
      <div>负载均衡端口(HTTP/Socks5): {{ proxyPoolState?.proxyList?.lib_port || '-' }}</div>
      <div>
        负载均衡策略：<span class="font-bold">{{ settingsState.settings?.glider_strategy }}</span>
      </div>
    </div>

    <n-table class="mt-1" :bordered="false" :single-line="false" size="small">
      <thead>
        <tr>
          <th>HTTP端口</th>
          <th>SOCKS端口</th>
          <th>节点</th>
        </tr>
      </thead>
      <tbody v-if="proxyPoolState.proxyList !== null">
        <tr v-for="item in proxyPoolState.proxyList.open_result_list" :key="item.name">
          <td>{{ item.http_port || '-' }}</td>
          <td>{{ item.socks_port || '-' }}</td>
          <td>{{ item.name }}</td>
        </tr>
      </tbody>
    </n-table>
  </div>
</template>

<script setup lang="ts">
import { proxyPoolState } from '@/composables/use-proxy-pool';
import { settingsState } from '@/composables/use-settings';
</script>
