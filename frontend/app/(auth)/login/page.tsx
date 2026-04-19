'use client';

import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { FcGoogle } from 'react-icons/fc';
import { motion } from 'framer-motion';
import { handleGoogleLogin } from '@/api/auth';
import { APP_ROUTES } from '@/constants/appRoutes';
import { useTranslations } from 'next-intl';

export default function LoginPage() {
  const t = useTranslations('login_page');

  return (
    <main className="flex min-h-screen items-center justify-center bg-gradient-to-br from-indigo-100 to-white p-4">
      <Card className="w-full max-w-md shadow-2xl rounded-2xl p-8">
        <CardContent>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="text-center mb-6"
          >
            <h1 className="text-3xl font-bold mb-2">{t('header.title')}</h1>
            <p className="text-gray-600">{t('header.subtitle')}</p>
          </motion.div>
          <Button
            onClick={handleGoogleLogin}
            className="w-full flex items-center justify-center gap-2 text-base"
            variant="outline"
          >
            <FcGoogle className="text-xl" /> {t('actions.login_google')}
          </Button>
          <p className="text-xs text-center text-gray-500 mt-4">
            {t('messages.more_options')}
          </p>
          <div className="mt-6 text-center">
            <Link href={APP_ROUTES.HOME}>
              <span className="text-indigo-600 hover:underline text-sm">{t('links.back_home')}</span>
            </Link>
          </div>
        </CardContent>
      </Card>
    </main>
  );
}
