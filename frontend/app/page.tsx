import Link from 'next/link';
import { APP_ROUTES } from '@/constants/appRoutes';
import { useTranslations } from 'next-intl';

export const metadata = {
  title: 'Flashcards - Learn Languages Fast',
  description: 'Master languages with interactive flashcards. Study English, Spanish, French, and more with our free, easy-to-use app.',
};

export default function HomePage() {
  const t = useTranslations('home_page');

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 text-center">
      <section className="max-w-2xl">
        <h1 className="text-5xl font-extrabold mb-6 leading-tight">
          {t('hero.title')}
        </h1>
        <p className="text-lg mb-8 text-gray-600">
          {t('hero.subtitle')}
        </p>
        <Link
          href={APP_ROUTES.LOGIN}
          className="inline-block rounded-2xl bg-blue-600 px-6 py-3 text-white text-lg font-semibold hover:bg-blue-700 transition-colors shadow-lg"
        >
          {t('hero.button')}
        </Link>
      </section>
      <section className="mt-12 max-w-3xl">
        <h2 className="text-3xl font-bold mb-4">{t('benefits.title')}</h2>
        <ul className="space-y-3 text-left text-gray-700">
          <li>✅ {t('benefits.items.personalized')}</li>
          <li>✅ {t('benefits.items.multi_language')}</li>
          <li>✅ {t('benefits.items.responsive')}</li>
          <li>✅ {t('benefits.items.free')}</li>
        </ul>
      </section>
    </main>
  );
}
