<script lang="ts" setup>
import Fallback from '#/components/Fallback.vue';

defineOptions({ name: 'FallbackOffline' });
</script>

<template>
  <Fallback title="离线" description="网络连接已断开，请检查您的网络设置" />
</template>
