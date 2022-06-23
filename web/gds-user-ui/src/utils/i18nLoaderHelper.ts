import { i18n } from '@lingui/core';
import { en, fr, de, zh, ja, ru } from 'make-plural/plurals';
import { Locales } from 'types/type';

export const DEFAULT_LOCALE: Locales = 'en';

export async function dynamicActivate(locale: string) {
  const { messages } = await import(`../locales/${locale}/messages`);

  i18n.load(locale, messages);
  i18n.loadLocaleData({
    en: { plurals: en },
    de: { plurals: de },
    fr: { plurals: fr },
    zh: { plurals: zh },
    ja: { plurals: ja }
  });
  i18n.activate(locale);
}
