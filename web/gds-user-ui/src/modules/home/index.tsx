import React, { FormEvent, useEffect, useState } from 'react';
import { Heading, Stack } from '@chakra-ui/react';
import LandingLayout from 'layouts/LandingLayout';
import Head from 'components/Head/LandingHead';
import JoinUsSection from 'components/Section/JoinUs';
import SearchDirectory from 'components/Section/SearchDirectory';
import AboutTrisaSection from 'components/Section/AboutUs';

import { lookup } from './service';
import { isValidUuid } from 'utils/utils';
const HomePage: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState(false);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  const handleSearchSubmit = async (evt: FormEvent, searchQuery: string) => {
    evt.preventDefault();
    setIsLoading(true);
    const query = isValidUuid(searchQuery) ? `uuid=${searchQuery}` : `common_name=${searchQuery}`;

    try {
      const request = await lookup(query);
      const response = request.results[0];
      setIsLoading(false);
      if (!response.error) {
        setResult(response);
        setSearch(searchQuery);
        setError('');
      }
    } catch (e: any) {
      setIsLoading(false);
      if (!e.response.data.success) {
        setError(e.response.data.error);
        setResult(false);
      } else {
        console.log('sorry something went wrong , please try again');
      }
    }
  };
  return (
    <LandingLayout>
      <Head hasBtn isHomePage />
      <AboutTrisaSection />
      <JoinUsSection />
      <SearchDirectory
        handleSubmit={handleSearchSubmit}
        isLoading={isLoading}
        result={result}
        error={error}
        query={search}
      />
    </LandingLayout>
  );
};

export default HomePage;
