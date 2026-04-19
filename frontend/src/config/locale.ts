export const supportedLocales = ['en', 'de', 'pl'] as const; // remember also about next.config.js
export type Locale = (typeof supportedLocales)[number];

export const defaultLocale: Locale = 'en';