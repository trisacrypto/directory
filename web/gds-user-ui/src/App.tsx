import LandingLayout from './layouts/LandingLayout';
import { BrowserRouter } from 'react-router-dom';
import AppRouter from 'application/routes';
const App: React.FC = () => {
  return (
    <BrowserRouter>
      <div className="App">
        <AppRouter />
      </div>
    </BrowserRouter>
  );
};

export default App;
