import '../src/styles/globals.css';
import type { Metadata } from 'next';
import { NextIntlClientProvider } from 'next-intl';
import { getLocale, getMessages } from 'next-intl/server';

export const metadata: Metadata = {
  title: 'Flashcards - Learn Languages Fast',
  description: 'Master languages with interactive flashcards. Study English, Spanish, French, and more with our free, easy-to-use app.',
  openGraph: {
    title: 'Flashcards - Learn Languages Fast',
    description: 'Master languages with interactive flashcards. Study English, Spanish, French, and more with our free, easy-to-use app.',
    url: 'https://yourdomain.com',
    siteName: 'Flashcards',
    images: [
      {
        url: '/favicon/android-chrome-512x512.png',
        width: 512,
        height: 512,
        alt: 'Flashcards Logo',
      },
    ],
    locale: 'en_US',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Flashcards - Learn Languages Fast',
    description: 'Master languages with interactive flashcards. Study English, Spanish, French, and more with our free, easy-to-use app.',
    images: ['/favicon/android-chrome-512x512.png'],
  },
  icons: {
    icon: '/favicon/favicon-32x32.png',
    shortcut: '/favicon/favicon-16x16.png',
    apple: '/favicon/apple-touch-icon.png',
  },
};

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const locale = await getLocale();
  const messages = await getMessages();
  
  return (
    <html lang={locale}>
      <body className="min-h-screen bg-white text-gray-900 antialiased">
        <NextIntlClientProvider locale={locale} messages={messages}>
          {children}
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
