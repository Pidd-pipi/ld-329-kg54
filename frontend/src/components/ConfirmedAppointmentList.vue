<template>
  <el-timeline v-if="appointments.length">
    <el-timeline-item v-for="item in appointments" :key="item.id" :timestamp="item.time" type="success">
      <strong>{{ item.pair }}</strong>
      <el-tag size="small" type="success" class="status-tag">{{ statusLabel }}</el-tag>
      <p>{{ item.place }}</p>
      <p class="muted">{{ item.agenda }}</p>
      <small>{{ item.initiator }} 发起，{{ item.responder }} 已确认</small>
    </el-timeline-item>
  </el-timeline>
  <el-empty v-else description="暂无双方确认的预约" :image-size="72" />
</template>

<script setup lang="ts">
import type { Appointment } from '../types/domain';
import { APPOINTMENT_STATUS_LABEL } from '../constants/appointment.constants';

defineProps<{ appointments: Appointment[] }>();

const statusLabel = APPOINTMENT_STATUS_LABEL.confirmed;
</script>

<style scoped>
.status-tag { margin-left: 8px; }
</style>
