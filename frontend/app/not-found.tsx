import { APP_ROUTES } from "@/constants/appRoutes";
import { useTranslations } from 'next-intl';

export default function NotFound() {
  const t = useTranslations('not_found_page');

  return (
    <div className="flex min-h-screen flex-col items-center justify-center">
      <h1 className="text-4xl font-bold">{t('hero.title')}</h1>
      <p className="mt-4 text-lg">{t('hero.subtitle')}</p>
      <a href={APP_ROUTES.HOME} className="mt-6 text-blue-500 hover:underline">
        {t('hero.button')}
      </a>
    </div>
  );
}