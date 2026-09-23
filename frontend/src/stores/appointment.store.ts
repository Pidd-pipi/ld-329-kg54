import { defineStore } from 'pinia';
import type { Appointment, AppointmentCreatePayload, AppointmentFormPayload } from '../types/domain';
import {
  confirmAppointment,
  createAppointment,
  fetchAppointments,
  rejectAppointment,
  reviseAppointment,
} from '../services/appointment.service';

export const useAppointmentStore = defineStore('appointment', {
  state: () => ({
    appointments: [] as Appointment[],
    currentUser: '',
    loading: false,
  }),
  getters: {
    negotiating: (state) => state.appointments.filter((item) => item.status === 'pending'),
    confirmed: (state) => state.appointments.filter((item) => item.status === 'confirmed'),
    activeMatchIds: (state) => new Set(state.appointments.map((item) => item.matchId)),
  },
  actions: {
    async load() {
      this.loading = true;
      try {
        this.appointments = await fetchAppointments();
      } finally {
        this.loading = false;
      }
    },
    async create(payload: AppointmentCreatePayload) {
      await createAppointment(payload);
      await this.load();
    },
    async confirm(id: number) {
      await confirmAppointment(id, this.currentUser);
      await this.load();
    },
    async reject(id: number) {
      await rejectAppointment(id, this.currentUser);
      await this.load();
    },
    async revise(id: number, form: AppointmentFormPayload) {
      await reviseAppointment(id, { ...form, actor: this.currentUser });
      await this.load();
    },
  },
});
