'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { handleGoogleAuthCallback } from '@/api/auth';
import { APP_ROUTES } from '@/constants/appRoutes';
import LoadingSpinner from '@/components/atoms/LoadingSpinner';

export default function AuthCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const code = searchParams.get('code');

    if (!code) {
      router.push(APP_ROUTES.LOGIN);
      return;
    }

    handleGoogleAuthCallback(code)
      .then((success) => {
        if (success) {
          router.push(APP_ROUTES.DASHBOARD);
        } else {
          router.push(APP_ROUTES.LOGIN);
        }
      })
      .catch(() => router.push(APP_ROUTES.LOGIN));
  }, [router, searchParams]);

  return <LoadingSpinner />;
}
