import { request, uploadFile } from "./client";
import type {
  Place,
  PlaceCategory,
  PlacePhoto,
  SavedPlace,
  User,
  Visit,
  VisitPhoto,
} from "./types";

export const authApi = {
  register: (name: string, email: string, password: string) =>
    request<{ user: User; token: string }>("/auth/register", {
      method: "POST",
      body: { name, email, password },
    }),

  login: (email: string, password: string) =>
    request<{ user: User; token: string }>("/auth/login", {
      method: "POST",
      body: { email, password },
    }),

  me: () => request<User>("/auth/me", { auth: true }),
};

export interface PlaceFilters {
  category?: PlaceCategory;
  area?: string;
  q?: string;
}

export const placesApi = {
  homepage: () => request<Place[]>("/homepage"),

  list: (filters: PlaceFilters = {}) => {
    const params = new URLSearchParams();
    if (filters.category) params.set("category", filters.category);
    if (filters.area) params.set("area", filters.area);
    if (filters.q) params.set("q", filters.q);
    const qs = params.toString();
    return request<Place[]>(`/places${qs ? `?${qs}` : ""}`);
  },

  get: (id: number) =>
    request<{ place: Place; photos: PlacePhoto[] }>(`/places/${id}`),

  create: (place: Partial<Place> & { source: "curated" | "user" }) =>
    request<Place>("/places", { method: "POST", body: place, auth: true }),

  update: (id: number, place: Partial<Place>) =>
    request<Place>(`/places/${id}`, { method: "PUT", body: place, auth: true }),

  remove: (id: number) =>
    request<void>(`/places/${id}`, { method: "DELETE", auth: true }),

  setVisibility: (id: number, visibility: "public" | "private") =>
    request<{ id: number; visibility: string }>(`/places/${id}/visibility`, {
      method: "PATCH",
      body: { visibility },
      auth: true,
    }),

  uploadPhoto: (id: number, file: File, caption?: string) =>
    uploadFile<PlacePhoto>(`/places/${id}/photos`, file, caption),
};

export const visitsApi = {
  list: () => request<Visit[]>("/users/me/visits", { auth: true }),

  create: (placeId: number, notes?: string, visitedAt?: string) =>
    request<Visit>("/visits", {
      method: "POST",
      body: { place_id: placeId, notes, visited_at: visitedAt },
      auth: true,
    }),

  update: (id: number, notes?: string, visitedAt?: string) =>
    request<void>(`/visits/${id}`, {
      method: "PUT",
      body: { notes, visited_at: visitedAt },
      auth: true,
    }),

  remove: (id: number) =>
    request<void>(`/visits/${id}`, { method: "DELETE", auth: true }),

  uploadPhoto: (id: number, file: File, caption?: string) =>
    uploadFile<VisitPhoto>(`/visits/${id}/photos`, file, caption),
};

export const savedApi = {
  list: () => request<SavedPlace[]>("/users/me/saved", { auth: true }),

  add: (placeId: number) =>
    request<SavedPlace>("/saved", {
      method: "POST",
      body: { place_id: placeId },
      auth: true,
    }),

  remove: (placeId: number) =>
    request<void>(`/saved/${placeId}`, { method: "DELETE", auth: true }),
};
