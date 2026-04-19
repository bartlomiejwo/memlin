import { API_ROUTES } from '@/constants/apiRoutes';
import backend from './axios';
import { useRouter } from "next/router";
import { APP_ROUTES } from '@/constants/appRoutes';
import { AUTH_KEYS } from '@/constants/auth';


type GoogleAuthResponse = {
  access_token: string;
}

export const handleGoogleLogin = () => {
  const frontendBaseUrl = window.location.origin.replace('127.0.0.1', 'localhost');
  const redirectUri = `${frontendBaseUrl}${APP_ROUTES.GOOGLE_CALLBACK}`;

  // Redirect the user to the backend's Google OAuth endpoint
  window.location.href = `${process.env.NEXT_PUBLIC_BACKEND_URL}${API_ROUTES.AUTH_GOOGLE}?redirect_uri=${redirectUri}`;
};

export const handleGoogleAuthCallback = async (code: string) => {
  try {
    const response = await backend.get<GoogleAuthResponse>(`${API_ROUTES.AUTH_GOOGLE_CALLBACK}?code=${code}`);
    const data = response.data;

    if (data.access_token) {
      localStorage.setItem(AUTH_KEYS.ACCESS_TOKEN, data.access_token);
      return true;
    }
    return false;
  } catch (error) {
    console.error('Error during Google auth callback:', error);
    return false;
  }
};


export function logout() {
  const router = useRouter();

  localStorage.removeItem('access_token');
  router.push(APP_ROUTES.LOGIN); 
}
