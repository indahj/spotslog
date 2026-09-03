export type PlaceCategory =
  | "restaurant"
  | "cafe"
  | "museum"
  | "library"
  | "attraction";

export type PlaceSource = "curated" | "user";
export type PlaceVisibility = "public" | "private";

export interface User {
  id: number;
  name: string;
  email: string;
  avatar_url?: string;
  role: "user" | "admin";
  created_at: string;
}

export interface Place {
  id: number;
  name: string;
  category: PlaceCategory;
  address: string;
  district?: string;
  lat: number;
  lng: number;
  description?: string;
  price_range?: string;
  opening_hours?: Record<string, unknown>;
  menu?: Record<string, unknown>;
  source: PlaceSource;
  visibility: PlaceVisibility;
  created_by?: number;
  created_at: string;
  updated_at: string;
  cover_photo_url?: string;
}

export interface PlacePhoto {
  id: number;
  place_id: number;
  url: string;
  uploaded_by?: number;
  caption?: string;
  created_at: string;
}

export interface VisitPhoto {
  id: number;
  visit_id: number;
  url: string;
  caption?: string;
  created_at: string;
}

export interface Visit {
  id: number;
  user_id: number;
  place_id: number;
  visited_at: string;
  notes?: string;
  created_at: string;
  place?: Place;
  photos?: VisitPhoto[];
}

export interface SavedPlace {
  id: number;
  user_id: number;
  place_id: number;
  created_at: string;
  place?: Place;
}

export const CATEGORY_LABELS: Record<PlaceCategory, string> = {
  restaurant: "Restaurant",
  cafe: "Café",
  museum: "Museum",
  library: "Library",
  attraction: "Attraction",
};
