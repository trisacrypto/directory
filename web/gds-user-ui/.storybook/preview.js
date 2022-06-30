import customTheme from "../src/utils/theme";
import { LanguageProvider } from '../src/contexts/LanguageContext'
import { MemoryRouter } from "react-router";
import { Provider } from 'react-redux';
import store from '../src/application/store'


export const parameters = {
  chakra: {
    customTheme,
  },
};

export const decorators = [
  (Story) => {
  return (
  <MemoryRouter initialEntries={['/']}>
    <Provider store={store}>
        <LanguageProvider>
          <Story />
        </LanguageProvider>
    </Provider>
  </MemoryRouter>
  )
  },
];