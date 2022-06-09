import React, { FormEvent, useEffect, useState } from 'react';

import LandingLayout from 'layouts/LandingLayout';
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
      const response = request.results;
      setIsLoading(false);
      if (!response.error) {
        setError('');
        setResult(response);
        setSearch(searchQuery);
      }
    } catch (e: any) {
      setIsLoading(false);
      console.error('error', e);
      if (!e?.response?.data.success) {
        setResult(false);
        setError(e.response.data.error);
        setResult(false);
      } else {
        console.error('error 2', e);
        Sentry.captureMessage('Something wrong happen when we trying to call lookup api');
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
