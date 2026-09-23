<template>
  <header class="app-header">
    <div>
      <p class="eyebrow">Skill Swap Campus</p>
      <h1>{{ APP_NAME }}</h1>
      <p>{{ APP_TAGLINE }}</p>
    </div>
    <div class="app-header__actions">
      <slot name="viewer" />
      <el-badge :value="unread" class="message-badge">
        <el-button type="primary">站内消息</el-button>
      </el-badge>
    </div>
  </header>
</template>

<script setup lang="ts">
import { APP_NAME, APP_TAGLINE } from '../constants/app.constants';

defineProps<{ unread: number }>();
</script>

<style scoped>
.app-header__actions { display: flex; align-items: center; gap: 16px; flex-wrap: wrap; justify-content: flex-end; }
</style>
