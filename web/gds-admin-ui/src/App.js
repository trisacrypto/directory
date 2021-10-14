import { Toaster } from 'react-hot-toast';
import './App.css';
import './assets/scss/Creative.scss';

import Routes from './routes/Routes';

const toastOptions = {
  style: {
    borderRadius: 0,
    padding: '1rem',
  }
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
