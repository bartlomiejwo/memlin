"use client";

import { useEffect, useState } from "react";
import { useErrorStore } from "@/stores/errorStore";
import { useRouter } from "next/navigation";
import { APP_ROUTES } from "@/constants/appRoutes";
import LoadingSpinner from "@/components/atoms/LoadingSpinner";
import { useTranslations } from 'next-intl';

export default function ErrorPage() {
  const t = useTranslations('error_page');
  const { error, message, clearError } = useErrorStore();
  const router = useRouter();
  const [isStateReady, setIsStateReady] = useState(false);

  useEffect(() => {
    // Add a slight delay to ensure Zustand state has been initialized
    const timeout = setTimeout(() => {
      if (error !== null) {
        setIsStateReady(true);
      } else {
        router.push(APP_ROUTES.HOME); // Redirect to home if no error is found
      }
    }, 100); // Delay for 100ms to allow Zustand state to load

    return () => clearTimeout(timeout); // Clean up timeout on unmount
  }, [error, router]);

  const handleGoHome = () => {
    router.push(APP_ROUTES.HOME);
    setTimeout(() => {
      clearError();
    }, 500);
  };

  // Don't render page until Zustand state is set
  if (!isStateReady) {
    return <LoadingSpinner />;
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 w-full px-4">
      <div className="max-w-xs sm:max-w-md md:max-w-lg text-center break-words break-all">
        <h1 className="text-2xl sm:text-3xl font-bold text-red-600">
          {t('hero.title')}: {error}
        </h1>
        <p className="text-base sm:text-lg text-gray-700 mt-2 break-words break-all">
          {message}
        </p>
        <button
          className="mt-4 px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-700 transition-all"
          onClick={handleGoHome}
        >
          {t('hero.button')}
        </button>
      </div>
    </div>
  );
}
