import React, { FormEvent, useState } from 'react';

import Head from 'components/Head/LandingHead';
import JoinUsSection from 'components/Section/JoinUs';
import SearchDirectory from 'components/Section/SearchDirectory';
import AboutTrisaSection from 'components/Section/AboutUs';
import * as Sentry from '@sentry/react';
import { lookup } from './service';
import { isValidUuid } from 'utils/utils';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
const HomePage: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState(false);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  const handleSearchSubmit = async (evt: FormEvent, searchQuery: string) => {
    evt.preventDefault();
    // throw new Error('Sentry Frontend Error');
    setIsLoading(true);
    const query = isValidUuid(searchQuery) ? `uuid=${searchQuery}` : `common_name=${searchQuery}`;

    try {
      const request = await lookup(query);

      setIsLoading(false);
      if (request?.mainnet || request?.testnet) {
        setError('');
        setResult(request);
        setSearch(searchQuery);
      } else {
        setResult(false);
        setError('No results found');
      }
    } catch (e: any) {
      setIsLoading(false);
      setResult(false);
      if (!e.response?.data?.success) {
        setError(e.response?.data?.error);
      } else {
        setError('Something went wrong');
        Sentry.captureException(e);
      }
    }
  };
  return (
    <>
      <LandingHeader />
      <Head hasBtn isHomePage />
      <AboutTrisaSection />
      <JoinUsSection />

      <SearchDirectory
        handleSubmit={handleSearchSubmit}
        isLoading={isLoading}
        result={result}
        error={error}
        handleClose={() => {
          setResult(false);
          setError('');
        }}
        query={search}
      />
      <Footer />
    </>
  );
};

export default HomePage;
