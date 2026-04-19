import { create } from "zustand";
import { persist } from "zustand/middleware";

type ErrorState = {
  error: string | null;
  message: string | null;
  setError: (error: string, message: string) => void;
  clearError: () => void;
};

export const useErrorStore = create<ErrorState>()(
  persist(
    (set) => ({
      error: null,
      message: null,
      setError: (error, message) => set({ error, message }),
      clearError: () => set({ error: null, message: null }),
    }),
    { name: "error-storage" } // Saves error state to localStorage
  )
);
