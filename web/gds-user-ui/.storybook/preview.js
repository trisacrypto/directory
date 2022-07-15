import theme from "../src/theme";
import { LanguageProvider } from '../src/contexts/LanguageContext'
import { MemoryRouter } from "react-router";
import { Provider } from 'react-redux';
import store from '../src/application/store'


export const parameters = {
  chakra: {
    theme,
  },
};

export const decorators = [
  (Story) => {
<<<<<<< HEAD
  return (
  <MemoryRouter initialEntries={['/']}>
    <Provider store={store}>
        <LanguageProvider>
          <Story />
        </LanguageProvider>
    </Provider>
  </MemoryRouter>
  )
=======
    return (
      <MemoryRouter initialEntries={['/']}>
        <Provider store={store}>
          <LanguageProvider>
            <Story />
          </LanguageProvider>
        </Provider>
      </MemoryRouter>
    )
>>>>>>> origin/main
  },
];