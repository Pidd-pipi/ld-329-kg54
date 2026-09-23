<template>
  <div class="appointment-actions">
    <template v-if="appointment.initiator === currentUser">
      <el-tag type="info" effect="plain">等待 {{ appointment.responder }} 回应</el-tag>
      <el-button size="small" @click="emit('revise')">改时间/议程</el-button>
    </template>
    <template v-else>
      <el-tag type="warning" effect="plain">{{ appointment.initiator }} 邀请你确认</el-tag>
      <el-button size="small" type="success" @click="emit('confirm')">确认</el-button>
      <el-button size="small" type="danger" plain @click="askReject">拒绝</el-button>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ElMessageBox } from 'element-plus';
import type { Appointment } from '../types/domain';

const props = defineProps<{ appointment: Appointment; currentUser: string }>();
const emit = defineEmits<{
  (e: 'confirm'): void;
  (e: 'reject'): void;
  (e: 'revise'): void;
}>();

async function askReject(): Promise<void> {
  try {
    await ElMessageBox.confirm(
      `拒绝后预约将撤下，${props.appointment.initiator} 可以重新发起，确认拒绝吗？`,
      '拒绝预约',
      { confirmButtonText: '拒绝', cancelButtonText: '再想想', type: 'warning' },
    );
  } catch {
    return;
  }
  emit('reject');
}
</script>

<style scoped>
.appointment-actions { display: flex; align-items: center; flex-wrap: wrap; gap: 8px; }
</style>
