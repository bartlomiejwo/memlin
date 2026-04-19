import { AUTH_KEYS } from '@/constants/auth';
import axios from 'axios';
import { logout } from './auth';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useErrorStore } from '@/stores/errorStore';
import { APP_ROUTES } from '@/constants/appRoutes';

const backend = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BACKEND_URL,
  withCredentials: true, // Ensures cookies (refresh_token) are sent
});

// Request Interceptor – Attach `access_token` to every request
backend.interceptors.request.use((config) => {
  const token = localStorage.getItem(AUTH_KEYS.ACCESS_TOKEN); // Get stored token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`; // Attach token
  }
  return config;
});

// Response Interceptor – Update `access_token` when backend provides a new one
backend.interceptors.response.use(
  (response) => {
    const newAccessToken = response.headers['authorization']?.replace('Bearer ', '');
    if (newAccessToken) {
      localStorage.setItem(AUTH_KEYS.ACCESS_TOKEN, newAccessToken); // Store new token
    }
    return response;
  },
  (error) => {
    const setError = useErrorStore.getState().setError;

    if (error.response) {
      const status = error.response.status;

      if (status === 401) {
        // If token is expired and refresh fails, force logout
        logout();
      } else if (status === 500) {
        // Handle 500 - Internal Server Error
        console.error('Server error', error);

        const errorMessage = error.response.data || "Internal Server Error";
        setError("500", errorMessage);

        // Navigate to a custom error page (e.g., /500)
        window.location.href = APP_ROUTES.ERROR;
      } else {
        // Show notification error only for unexpected errors
        toast.error('Something went wrong, please try again!', {
          position: 'top-right',
          autoClose: 5000,
        });
      }
    } else {
      // Handle errors without a response (e.g., network errors)
      toast.error('Network error, please check your connection!', {
        position: 'top-right',
        autoClose: 5000,
      });
    }

    return Promise.reject(error);
  }
);

export default backend;
