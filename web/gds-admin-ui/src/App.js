import AuthProvider from 'contexts/auth/auth-provider';
import { Provider } from 'react-redux';
import { Toaster } from 'react-hot-toast';
import { configureStore } from './redux/store';
import { QueryClientProvider } from '@tanstack/react-query';
import queryClient from 'helpers/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'


import './App.css';
import './assets/scss/Creative.scss';

import Routes from './routes/Routes';

const toastOptions = {
  style: { borderRadius: "0", color: "#FFF" },
  error: {
    style: { background: "#f44336" },
    iconTheme: { primary: "red", secondary: "white" },
  },
  success: {
    style: {
      background: "#4caf50",
      iconTheme: { primary: "green", secondary: "white" },
    },
  },
}

function App() {
  return (
    <div className="App">
      <QueryClientProvider client={queryClient}>
        <Provider store={configureStore({})}>
          <AuthProvider>
            <Routes />
          </AuthProvider>
          <Toaster position='top-right' toastOptions={toastOptions} />
        </Provider>
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </div>
  );
}

export default App;
