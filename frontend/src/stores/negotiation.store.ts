import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type { Appointment, Match } from '../types/domain';
import type { AppointmentDraft } from '../types/appointment';
import { APPOINTMENT_STATUS_CONFIRMED, APPOINTMENT_STATUS_PENDING, DEFAULT_VIEWER } from '../constants/appointment.constants';
import {
  confirmAppointment,
  createAppointment,
  rejectAppointment,
  reviseAppointment,
} from '../services/appointment.service';

export const useNegotiationStore = defineStore('negotiation', () => {
  const matches = ref<Match[]>([]);
  const appointments = ref<Appointment[]>([]);
  const currentUser = ref(DEFAULT_VIEWER);
  const busy = ref(false);

  const isParty = (match: Match): boolean =>
    match.provider === currentUser.value || match.learner === currentUser.value;

  const otherParty = (match: Match): string =>
    match.provider === currentUser.value ? match.learner : match.provider;

  const appointmentByMatch = computed(() => {
    const map = new Map<number, Appointment>();
    for (const appointment of appointments.value) {
      map.set(appointment.matchId, appointment);
    }
    return map;
  });

  const pendingAppointments = computed(() =>
    appointments.value.filter(
      (appointment) =>
        appointment.status === APPOINTMENT_STATUS_PENDING &&
        (appointment.initiator === currentUser.value ||
          appointment.responder === currentUser.value),
    ),
  );

  const confirmedAppointments = computed(() =>
    appointments.value.filter(
      (appointment) =>
        appointment.status === APPOINTMENT_STATUS_CONFIRMED &&
        (appointment.initiator === currentUser.value ||
          appointment.responder === currentUser.value),
    ),
  );

  function hydrate(nextMatches: Match[], nextAppointments: Appointment[]): void {
    matches.value = nextMatches;
    appointments.value = nextAppointments;
  }

  function switchViewer(user: string): void {
    currentUser.value = user;
  }

  async function initiate(draft: AppointmentDraft): Promise<void> {
    busy.value = true;
    try {
      const appointment = await createAppointment(draft);
      upsert(appointment);
    } finally {
      busy.value = false;
    }
  }

  async function revise(id: number, draft: AppointmentDraft): Promise<void> {
    busy.value = true;
    try {
      const appointment = await reviseAppointment(id, draft);
      upsert(appointment);
    } finally {
      busy.value = false;
    }
  }

  async function confirm(id: number): Promise<void> {
    busy.value = true;
    try {
      const appointment = await confirmAppointment(id, currentUser.value);
      upsert(appointment);
    } finally {
      busy.value = false;
    }
  }

  async function reject(id: number): Promise<void> {
    busy.value = true;
    try {
      await rejectAppointment(id, currentUser.value);
      appointments.value = appointments.value.filter((item) => item.id !== id);
    } finally {
      busy.value = false;
    }
  }

  function upsert(appointment: Appointment): void {
    const index = appointments.value.findIndex((item) => item.id === appointment.id);
    if (index >= 0) {
      appointments.value[index] = appointment;
    } else {
      appointments.value = [...appointments.value, appointment];
    }
  }

  return {
    matches,
    appointments,
    currentUser,
    busy,
    isParty,
    otherParty,
    appointmentByMatch,
    pendingAppointments,
    confirmedAppointments,
    hydrate,
    switchViewer,
    initiate,
    revise,
    confirm,
    reject,
  };
});
