import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';
import JoinUsSection from 'components/Section/JoinUs';
import SearchDirectory from 'components/Section/SearchDirectory';
import AboutTrisaSection from 'components/Section/AboutUs';

const HomePage: React.FC = () => {
  return (
    <LandingLayout>
      <Head hasBtn />
      <AboutTrisaSection />
      <JoinUsSection />
      <SearchDirectory />
    </LandingLayout>
  );
};

export default HomePage;
