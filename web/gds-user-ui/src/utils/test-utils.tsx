import { ChakraProvider } from '@chakra-ui/react';
import { render as rtlRender, RenderOptions as RTLRenderOptions } from '@testing-library/react';
import React, { useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Provider } from 'react-redux';
import customTheme from './theme';
import rootReducer from 'application/store/rootReducer';
import { configureStore, Store } from '@reduxjs/toolkit';
import { I18nProvider } from '@lingui/react';
import { i18n } from '@lingui/core';
import { LanguageContext } from 'contexts/LanguageContext';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

type RenderOptions = {
  rhfDefaultValues?: { [x: string]: any };
  store?: Store;
  preloadedState?: StoreState;
  renderOptions?: Omit<RTLRenderOptions, 'wrapper'>;
  locale?: string;
  route?: string | URL;
};

type StoreState = {
  [x: string]: any;
};

const queryClient = new QueryClient({
  logger: {
    log: console.log,
    warn: console.warn,
    // ✅ no more errors on the console for tests
    error: process.env.NODE_ENV === 'test' ? () => {} : console.error
  }
});

export function render(
  ui: React.ReactElement<any, string | React.JSXElementConstructor<any>>,
  {
    locale = 'en',
    rhfDefaultValues,
    preloadedState,
    route = '/',
    store = configureStore({ reducer: rootReducer, preloadedState }),
    ...renderOptions
  }: RenderOptions = {}
) {
  function Wrapper({ children }: { children: React.ReactNode }) {
    const methods = useForm({ defaultValues: rhfDefaultValues });
    const [language, setLanguage] = useState<string | null>(locale);
    window.history.pushState({}, 'Test page', route);

    return (
      <I18nProvider i18n={i18n}>
        <QueryClientProvider client={queryClient}>
          <LanguageContext.Provider value={[language, setLanguage]}>
            <Provider store={store}>
              <ChakraProvider theme={customTheme}>
                <FormProvider {...methods}>
                  <BrowserRouter>{children}</BrowserRouter>
                </FormProvider>
              </ChakraProvider>
            </Provider>
          </LanguageContext.Provider>
        </QueryClientProvider>
      </I18nProvider>
    );
  }
  return rtlRender(ui, { wrapper: Wrapper, ...renderOptions });
}

export * from '@testing-library/react';
