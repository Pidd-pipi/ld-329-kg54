<template>
  <el-dialog
    :model-value="visible"
    :title="dialogTitle"
    width="480px"
    @update:model-value="emit('update:visible', $event)"
    @open="resetForm"
  >
    <el-form label-position="top">
      <el-form-item label="共同可约时间">
        <el-select v-model="form.time" placeholder="选择双方都方便的时间" class="full-width">
          <el-option v-for="slot in timeSlots" :key="slot" :label="slot" :value="slot" />
        </el-select>
      </el-form-item>
      <el-form-item label="交换地点">
        <el-input
          v-model="form.place"
          :maxlength="APPOINTMENT_FORM_LIMITS.placeMax"
          show-word-limit
          placeholder="线上会议室 / 校区地点"
        />
        <div class="tag-row">
          <el-tag
            v-for="place in APPOINTMENT_PLACE_SUGGESTIONS"
            :key="place"
            size="small"
            effect="plain"
            class="place-suggestion"
            @click="form.place = place"
          >
            {{ place }}
          </el-tag>
        </div>
      </el-form-item>
      <el-form-item label="协商议程">
        <el-input
          v-model="form.agenda"
          type="textarea"
          :rows="3"
          :maxlength="APPOINTMENT_FORM_LIMITS.agendaMax"
          show-word-limit
          placeholder="本次交换要聊什么、怎么进行"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:visible', false)">取消</el-button>
      <el-button type="primary" :disabled="!canSubmit" :loading="submitting" @click="submit">
        {{ isRevise ? '保存并重新发起确认' : '提交预约' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue';
import type { Appointment, AppointmentFormPayload, Match } from '../../types/domain';
import { APPOINTMENT_FORM_LIMITS, APPOINTMENT_PLACE_SUGGESTIONS } from '../../constants/appointment.constants';

const props = defineProps<{
  visible: boolean;
  match: Match | null;
  appointment: Appointment | null;
  submitting: boolean;
}>();
const emit = defineEmits<{
  'update:visible': [value: boolean];
  submit: [payload: AppointmentFormPayload];
}>();

const form = reactive<AppointmentFormPayload>({ time: '', place: '', agenda: '' });

const isRevise = computed(() => props.appointment !== null);
const dialogTitle = computed(() => {
  if (isRevise.value) {
    return `修改预约 · ${props.appointment?.pair ?? ''}`;
  }
  const pair = props.match ? `${props.match.provider} ↔ ${props.match.learner}` : '';
  return `发起预约 · ${pair}`;
});
const timeSlots = computed(() => props.match?.commonSlots ?? []);
const canSubmit = computed(() => Boolean(form.time && form.place.trim() && form.agenda.trim()));

function resetForm() {
  form.time = props.appointment?.time ?? '';
  form.place = props.appointment?.place ?? '';
  form.agenda = props.appointment?.agenda ?? '';
}

function submit() {
  if (!canSubmit.value) {
    return;
  }
  emit('submit', { time: form.time, place: form.place.trim(), agenda: form.agenda.trim() });
}
</script>
