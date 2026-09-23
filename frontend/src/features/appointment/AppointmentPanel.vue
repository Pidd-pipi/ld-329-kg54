<template>
  <div class="panel">
    <h2>预约协商</h2>
    <p v-if="!store.negotiating.length" class="muted">暂无等待确认的预约，从智能匹配中发起一个吧。</p>
    <AppointmentCard
      v-for="item in store.negotiating"
      :key="item.id"
      :appointment="item"
      :current-user="currentUser"
      @confirm="emit('confirm', $event)"
      @reject="emit('reject', $event)"
      @revise="emit('revise', $event)"
    />

    <el-divider />

    <h2>预约确认区</h2>
    <p v-if="!store.confirmed.length" class="muted">双方确认后预约才会生效并出现在这里。</p>
    <el-timeline v-else>
      <el-timeline-item v-for="item in store.confirmed" :key="item.id" :timestamp="item.time" type="success">
        <strong>{{ item.pair }}</strong>
        <p>{{ item.place }} · {{ statusLabel(item) }}</p>
        <p class="muted">{{ item.agenda }}</p>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script setup lang="ts">
import AppointmentCard from './AppointmentCard.vue';
import { useAppointmentStore } from '../../stores/appointment.store';
import { APPOINTMENT_STATUS_LABELS } from '../../constants/appointment.constants';
import type { Appointment } from '../../types/domain';

defineProps<{ currentUser: string }>();
const emit = defineEmits<{
  confirm: [appointment: Appointment];
  reject: [appointment: Appointment];
  revise: [appointment: Appointment];
}>();

const store = useAppointmentStore();

function statusLabel(item: Appointment) {
  return APPOINTMENT_STATUS_LABELS[item.status];
}
</script>
