export type AppointmentStatus = 'pending' | 'confirmed';

export interface Appointment {
  id: number;
  matchId: number;
  pair: string;
  initiator: string;
  responder: string;
  time: string;
  place: string;
  status: AppointmentStatus;
  agenda: string;
  slots: string[];
  version: number;
}

export interface AppointmentDraft {
  matchId: number;
  actor: string;
  time: string;
  place: string;
  agenda: string;
}
