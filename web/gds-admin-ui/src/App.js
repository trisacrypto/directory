import { Toaster } from 'react-hot-toast';
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
      <Routes></Routes>
      <Toaster position='top-right' toastOptions={toastOptions} />
    </div>
  );
}

export default App;
