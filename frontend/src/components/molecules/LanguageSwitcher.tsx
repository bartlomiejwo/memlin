'use client';

import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { supportedLocales } from '../../config/locale';

export default function LanguageSwitcher() {
  const router = useRouter();
  const currentLocale = useLocale();

  const changeLocale = (newLocale: string) => {
    document.cookie = `NEXT_LOCALE=${newLocale}; path=/; max-age=31536000`;
    router.refresh();
  };

  return (
    <select
      value={currentLocale}
      onChange={(e) => changeLocale(e.target.value)}
      className="border rounded p-2"
    >
      {supportedLocales.map((locale) => (
        <option key={locale} value={locale}>
          {locale.toUpperCase()}
        </option>
      ))}
    </select>
  );
}