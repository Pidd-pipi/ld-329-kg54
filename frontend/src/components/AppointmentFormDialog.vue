<template>
  <el-dialog
    :model-value="modelValue"
    :title="existing ? '修改预约方案' : `发起预约 · ${match.provider} × ${match.learner}`"
    width="460px"
    @update:model-value="(value: boolean) => emit('update:modelValue', value)"
    @close="reset"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="84px">
      <el-form-item label="共同时间">
        <el-select v-model="form.time" placeholder="选择双方共同可约时间" class="full-width">
          <el-option
            v-for="slot in match.commonSlots"
            :key="slot"
            :label="slot"
            :value="slot"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="交换地点" prop="place">
        <el-select
          v-model="form.place"
          placeholder="选择或填写线上/线下地点"
          filterable
          allow-create
          default-first-option
          class="full-width"
        >
          <el-option v-for="place in placeSuggestions" :key="place" :label="place" :value="place" />
        </el-select>
      </el-form-item>
      <el-form-item label="议程" prop="agenda">
        <el-input
          v-model="form.agenda"
          type="textarea"
          :rows="3"
          maxlength="120"
          show-word-limit
          placeholder="说明本次交换的内容和安排"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">取消</el-button>
      <el-button type="primary" :loading="loading" @click="submit">
        {{ existing ? '保存并请求重新确认' : '发起预约' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import type { Appointment, Match } from '../types/domain';
import type { AppointmentDraft } from '../types/appointment';
import { PLACE_SUGGESTIONS } from '../constants/appointment.constants';

const props = defineProps<{
  modelValue: boolean;
  match: Match;
  actor: string;
  existing?: Appointment | null;
  loading?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'submit', draft: AppointmentDraft): void;
}>();

const formRef = ref<FormInstance>();
const placeSuggestions = PLACE_SUGGESTIONS;

interface AppointmentForm {
  time: string;
  place: string;
  agenda: string;
}

const form = reactive<AppointmentForm>({ time: '', place: '', agenda: '' });

const rules: FormRules<AppointmentForm> = {
  time: [{ required: true, message: '请选择共同可约时间', trigger: 'change' }],
  place: [{ required: true, message: '请填写交换地点', trigger: 'change' }],
  agenda: [{ required: true, message: '请填写交换议程', trigger: 'blur' }],
};

watch(
  () => props.modelValue,
  (visible) => {
    if (!visible) return;
    form.time = props.existing?.time ?? '';
    form.place = props.existing?.place ?? '';
    form.agenda = props.existing?.agenda ?? '';
  },
);

function reset(): void {
  formRef.value?.clearValidate();
}

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;
  emit('submit', {
    matchId: props.match.id,
    actor: props.actor,
    time: form.time,
    place: form.place.trim(),
    agenda: form.agenda.trim(),
  });
}
</script>

<style scoped>
.full-width { width: 100%; }
</style>
