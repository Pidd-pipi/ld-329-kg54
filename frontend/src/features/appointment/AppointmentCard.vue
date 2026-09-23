<template>
  <article class="feature-card appointment-card">
    <div class="feature-card__title">
      <strong>{{ appointment.pair }}</strong>
      <el-tag :type="statusTagType">{{ statusLabel }}</el-tag>
    </div>
    <p>{{ appointment.time }} · {{ appointment.place }}</p>
    <p class="muted">{{ appointment.agenda }}</p>
    <div class="tag-row">
      <el-tag v-if="appointment.revision > 1" size="small" effect="plain">
        第 {{ appointment.revision }} 版
      </el-tag>
      <el-tag v-for="name in appointment.confirmedBy" :key="name" size="small" type="success" effect="plain">
        {{ name }} 已确认
      </el-tag>
      <el-tag v-if="waitingFor" size="small" type="warning" effect="plain">
        等待 {{ waitingFor }} 确认
      </el-tag>
    </div>
    <div class="appointment-card__actions">
      <template v-if="isResponder">
        <el-button size="small" type="primary" @click="emit('confirm', appointment)">确认预约</el-button>
        <el-button size="small" type="danger" plain @click="emit('reject', appointment)">拒绝</el-button>
      </template>
      <template v-else-if="isInitiator">
        <el-button size="small" @click="emit('revise', appointment)">修改时间 / 议程</el-button>
        <small class="muted">对方确认前可修改，修改后需对方重新确认</small>
      </template>
      <small v-else class="muted">切换到 {{ appointment.initiator }} 或 {{ appointment.responder }} 视角可操作</small>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Appointment } from '../../types/domain';
import { APPOINTMENT_STATUS_LABELS, APPOINTMENT_STATUS_TAG_TYPES } from '../../constants/appointment.constants';

const props = defineProps<{ appointment: Appointment; currentUser: string }>();
const emit = defineEmits<{
  confirm: [appointment: Appointment];
  reject: [appointment: Appointment];
  revise: [appointment: Appointment];
}>();

const isInitiator = computed(() => props.currentUser === props.appointment.initiator);
const isResponder = computed(() => props.currentUser === props.appointment.responder);
const statusLabel = computed(() => APPOINTMENT_STATUS_LABELS[props.appointment.status]);
const statusTagType = computed(() => APPOINTMENT_STATUS_TAG_TYPES[props.appointment.status]);
const waitingFor = computed(() => {
  if (props.appointment.status !== 'pending') {
    return '';
  }
  return props.appointment.confirmedBy.includes(props.appointment.responder)
    ? ''
    : props.appointment.responder;
});
</script>
