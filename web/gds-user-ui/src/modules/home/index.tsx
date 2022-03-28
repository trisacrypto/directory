import React, { FormEvent, useEffect, useState } from 'react';
import { Heading, Stack } from '@chakra-ui/react';
import CollaboratorsSection from 'components/CollaboratorsSection';

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
      const res = await lookup(query);
      setResult(res.data);
      setIsLoading(false);
    } catch (e) {
      setIsLoading(false);
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
