import nextIntl from 'next-intl/plugin';

const withNextIntl = nextIntl({
  i18n: {
    locales: ['en', 'de', 'pl'], // remember also about src/config.locale.ts
    defaultLocale: 'en',
    localeDetection: false, // <- Important for no-routing behavior
  }
});

export default withNextIntl({
  // Existing config (e.g., for next-sitemap)
});