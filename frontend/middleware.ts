// middleware.ts
import { NextResponse } from 'next/server';
import { NextRequest } from 'next/server';
import { supportedLocales, defaultLocale, Locale } from './src/config/locale';

export function middleware(request: NextRequest) {
  const url = request.nextUrl;

  const queryLocale = url.searchParams.get('lang');
  const cookieLocale = request.cookies.get('NEXT_LOCALE')?.value;
  const acceptLang = request.headers.get('accept-language');
  const resolvedFromHeader = acceptLang?.split(',')?.[0]?.split('-')?.[0];
  const headerLocale = supportedLocales.includes(resolvedFromHeader as Locale)
    ? resolvedFromHeader
    : defaultLocale;

  // 1. If ?lang=xx is present — use only for this request
  if (queryLocale && supportedLocales.includes(queryLocale as Locale)) {
    const response = NextResponse.next();
    response.headers.set('x-temp-locale', queryLocale); // key part
    return response;
  }

  // 2. If cookie exists, use it
  if (cookieLocale && supportedLocales.includes(cookieLocale as Locale)) {
    return NextResponse.next();
  }

  // 3. Otherwise, set a new cookie based on Accept-Language
  const response = NextResponse.next();
  response.cookies.set('NEXT_LOCALE', headerLocale, {
    path: '/',
    maxAge: 60 * 60 * 24 * 365,
  });
  return response;
}
