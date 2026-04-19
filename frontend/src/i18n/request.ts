// src/i18n/requests.ts
import { getRequestConfig } from 'next-intl/server';
import { supportedLocales, defaultLocale, Locale } from '../config/locale';
import { cookies, headers } from 'next/headers';

export default getRequestConfig(async () => {
  const tempLocale = (await headers()).get('x-temp-locale'); // injected by middleware
  if (tempLocale && supportedLocales.includes(tempLocale as Locale)) {
    return {
      locale: tempLocale as Locale,
      messages: (await import(`../../locales/${tempLocale}.json`)).default,
    };
  }

  const cookieLocale = (await cookies()).get('NEXT_LOCALE')?.value;
  const selectedLocale: Locale = supportedLocales.includes(cookieLocale as Locale)
    ? (cookieLocale as Locale)
    : defaultLocale;

  const messages = (await import(`../../locales/${selectedLocale}.json`)).default;

  return {
    locale: selectedLocale,
    messages,
  };
});
