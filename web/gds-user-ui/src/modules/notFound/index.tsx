import LandingLayout from 'layouts/LandingLayout';
import NotFound from 'components/NotFound';
const StartPage: React.FC = () => {
  return (
    <LandingLayout>
      <NotFound />
    </LandingLayout>
  );
};

export default StartPage;
