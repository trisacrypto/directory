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

  const handleSearchSubmit = async (evt: FormEvent, searchQuery: string) => {
    setIsLoading(true);
    evt.preventDefault();
    const query = isValidUuid(searchQuery) ? `uuid=${searchQuery}` : `common_name=${searchQuery}`;

    try {
      const response = await lookup(query);
      if (!response.success) {
        setError(response.error);
      } else {
        setResult(response);
      }
      setIsLoading(false);
    } catch (e: any) {
      setIsLoading(false);
      setError('Sorry something went wrong, try again ');
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
      />
    </LandingLayout>
  );
};

export default HomePage;
